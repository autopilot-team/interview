package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/internal/types"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/skip2/go-qrcode"
)

// TwoFactorSetupData contains the data needed for 2FA setup
type TwoFactorSetupData struct {
	TwoFactor   *model.TwoFactor
	BackupCodes []string
	QRCode      string // Base64 encoded QR code image
}

// TwoFactorer is an interface that wraps the two-factor authentication methods
type TwoFactorer interface {
	Disable(ctx context.Context, userID string) error
	Enable(ctx context.Context, userID string, code string) error
	GetByUserID(ctx context.Context, userID string) (*model.TwoFactor, error)
	RegenerateQRCode(ctx context.Context, userID string) (string, error)
	Setup(ctx context.Context, userID string) (*TwoFactorSetupData, error)
	Verify(ctx context.Context, userID string, code string) error
}

// TwoFactor implements TwoFactorer interface
type TwoFactor struct {
	*app.Container
	store *store.Manager
}

// NewTwoFactor creates a new TwoFactor service
func NewTwoFactor(container *app.Container, store *store.Manager) *TwoFactor {
	return &TwoFactor{
		Container: container,
		store:     store,
	}
}

// Disable disables 2FA for a user
func (s *TwoFactor) Disable(ctx context.Context, userID string) error {
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor == nil || twoFactor.EnabledAt == nil {
		return httpx.ErrTwoFactorNotEnabled
	}

	// Create audit log before deletion
	metadata := map[string]any{
		"disabled_at":      time.Now(),
		"had_backup_codes": len(twoFactor.BackupCodes),
		"was_enabled_at":   twoFactor.EnabledAt,
	}
	if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionDisable, twoFactor.ID, userID, metadata); err != nil {
		return err
	}

	// Delete TwoFactor record - this is all we need to do to disable 2FA
	if err := s.store.TwoFactor.Delete(ctx, twoFactor.ID); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	return nil
}

// Enable enables 2FA for a user after verifying the initial setup
func (s *TwoFactor) Enable(ctx context.Context, userID string, code string) error {
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor == nil {
		return httpx.ErrTwoFactorNotEnabled
	}

	// Verify the provided code
	if !twoFactor.ValidateTOTP(code) {
		// Create audit log for failed enable attempt
		metadata := map[string]any{
			"attempt_at": time.Now(),
			"success":    false,
			"reason":     "invalid_code",
		}
		if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionEnable, twoFactor.ID, userID, metadata); err != nil {
			return err
		}

		return httpx.ErrInvalidTwoFactorCode
	}

	// Set enabled_at timestamp
	now := time.Now()
	twoFactor.EnabledAt = &now
	if err := s.store.TwoFactor.Update(ctx, twoFactor); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for successful enable
	metadata := map[string]any{
		"enabled_at":         now,
		"success":            true,
		"backup_codes_count": len(twoFactor.BackupCodes),
	}
	if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionEnable, twoFactor.ID, userID, metadata); err != nil {
		return err
	}

	return nil
}

// GetByUserID retrieves 2FA settings for a user
func (s *TwoFactor) GetByUserID(ctx context.Context, userID string) (*model.TwoFactor, error) {
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor == nil || twoFactor.EnabledAt == nil {
		return nil, httpx.ErrTwoFactorNotEnabled
	}

	return twoFactor, nil
}

// RegenerateQRCode regenerates the QR code for an existing 2FA setup
func (s *TwoFactor) RegenerateQRCode(ctx context.Context, userID string) (string, error) {
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return "", httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor == nil {
		return "", httpx.ErrTwoFactorNotEnabled
	}

	qrCode, err := s.generateQRCode(twoFactor.Secret)
	if err != nil {
		return "", httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for QR code regeneration
	metadata := map[string]any{
		"regenerated_at": time.Now(),
	}
	if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionCreate, twoFactor.ID, userID, metadata); err != nil {
		return "", err
	}

	return qrCode, nil
}

// Setup initiates 2FA setup for a user
func (s *TwoFactor) Setup(ctx context.Context, userID string) (*TwoFactorSetupData, error) {
	// Check if 2FA is already enabled
	existing, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if existing != nil && existing.EnabledAt != nil {
		return nil, httpx.ErrTwoFactorAlreadyEnabled
	}

	var (
		now         time.Time
		twoFactor   *model.TwoFactor
		backupCodes []string
		qrCode      string
	)

	if existing != nil {
		twoFactor = existing
		backupCodes = existing.BackupCodes
		qrCode, err = s.generateQRCode(existing.Secret)
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}
	} else {
		// Generate TOTP secret
		secret := make([]byte, 20)
		if _, err := rand.Read(secret); err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}
		secretBase32 := base32.StdEncoding.EncodeToString(secret)

		// Generate backup codes
		backupCodes, err = generateBackupCodes()
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}

		// Create TwoFactor record
		now = time.Now()
		twoFactor = &model.TwoFactor{
			Secret:      secretBase32,
			BackupCodes: backupCodes,
			UserID:      userID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Generate QR code
		qrCode, err = s.generateQRCode(secretBase32)
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}

		// Store in database
		created, err := s.store.TwoFactor.Create(ctx, twoFactor)
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}
		twoFactor = created
	}

	// Create audit log for 2FA setup initiation
	metadata := map[string]any{
		"setup_at":           now,
		"backup_codes_count": len(backupCodes),
	}
	if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionCreate, twoFactor.ID, userID, metadata); err != nil {
		return nil, err
	}

	return &TwoFactorSetupData{
		TwoFactor:   twoFactor,
		BackupCodes: backupCodes,
		QRCode:      qrCode,
	}, nil
}

// Verify validates a 2FA code for a user
func (s *TwoFactor) Verify(ctx context.Context, userID string, code string) error {
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, userID)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor == nil {
		return httpx.ErrTwoFactorNotEnabled
	}

	// Check if 2FA is locked
	if twoFactor.IsLocked() {
		// Create audit log for locked verification attempt
		metadata := map[string]any{
			"attempt_at":   time.Now(),
			"success":      false,
			"reason":       "locked",
			"locked_until": twoFactor.LockedUntil,
		}
		if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionVerify, twoFactor.ID, userID, metadata); err != nil {
			return err
		}

		return httpx.ErrTwoFactorLocked
	}

	// First try TOTP code
	if twoFactor.ValidateTOTP(code) {
		// Reset failed attempts on successful verification
		twoFactor.ResetFailedAttempts()
		if err := s.store.TwoFactor.Update(ctx, twoFactor); err != nil {
			return httpx.ErrUnknown.WithInternal(err)
		}

		// Create audit log for successful TOTP verification
		metadata := map[string]any{
			"verified_at": time.Now(),
			"success":     true,
			"method":      "totp",
		}
		if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionVerify, twoFactor.ID, userID, metadata); err != nil {
			return err
		}

		return nil
	}

	// Then try backup code
	valid, err := twoFactor.ValidateAndConsumeBackupCode(code)
	if err != nil {
		return httpx.ErrBackupCodeValidation
	}

	if !valid {
		// Increment failed attempts
		twoFactor.IncrementFailedAttempts()
		if err := s.store.TwoFactor.Update(ctx, twoFactor); err != nil {
			return httpx.ErrUnknown.WithInternal(err)
		}

		// Create audit log for failed verification
		metadata := map[string]any{
			"attempt_at":      time.Now(),
			"success":         false,
			"reason":          "invalid_code",
			"failed_attempts": twoFactor.FailedAttempts,
			"locked_until":    twoFactor.LockedUntil,
		}
		if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionVerify, twoFactor.ID, userID, metadata); err != nil {
			return err
		}

		return httpx.ErrInvalidTwoFactorCode
	}

	// Update the backup codes in database
	if err := s.store.TwoFactor.Update(ctx, twoFactor); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for successful backup code verification
	metadata := map[string]any{
		"verified_at":            time.Now(),
		"success":                true,
		"method":                 "backup_code",
		"remaining_backup_codes": len(twoFactor.BackupCodes),
	}
	if err := auditLog(ctx, s.store, types.ResourceTwoFactor, types.ActionVerify, twoFactor.ID, userID, metadata); err != nil {
		return err
	}

	return nil
}

// generateQRCode generates a QR code for the TOTP URI
func (s *TwoFactor) generateQRCode(secret string) (string, error) {
	// Generate the otpauth URI
	uri := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s",
		url.QueryEscape(s.Config.App.Name),
		url.QueryEscape(secret),
		url.QueryEscape(s.Config.App.Name),
	)

	// Generate QR code
	qr, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return "", err
	}

	// Get PNG data
	png, err := qr.PNG(256) // 256x256 pixels
	if err != nil {
		return "", err
	}

	// Convert to base64 data URL
	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png)), nil
}

// generateBackupCodes generates a set of backup codes
func generateBackupCodes() ([]string, error) {
	codes := make([]string, model.BackupCodesCount)
	for i := 0; i < model.BackupCodesCount; i++ {
		bytes := make([]byte, model.BackupCodeLength/2)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}

		codes[i] = fmt.Sprintf("%x", bytes)
	}

	return codes, nil
}
