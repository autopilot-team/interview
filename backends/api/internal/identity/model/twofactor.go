package model

import (
	"time"

	"github.com/pquerna/otp/totp"
)

const (
	// TOTPPeriod is the TOTP token validity period in seconds
	TOTPPeriod = 30

	// TOTPDigits is the number of digits in TOTP token
	TOTPDigits = 6

	// BackupCodesCount is the number of backup codes to generate
	BackupCodesCount = 10

	// BackupCodeLength is the length of each backup code
	BackupCodeLength = 10

	// MaxFailedAttempts is the maximum number of failed attempts within the window
	MaxFailedAttempts = 10

	// FailedAttemptsWindow is the time window for tracking failed attempts
	FailedAttemptsWindow = time.Hour
)

// TwoFactor represents a user's two-factor authentication settings
type TwoFactor struct {
	ID                  string     `db:"id"`
	BackupCodes         []string   `db:"backup_codes"` // JSONB array of backup codes
	Secret              string     `db:"secret"`       // Base32 encoded TOTP secret
	UserID              string     `db:"user_id"`
	FailedAttempts      int        `db:"failed_attempts"`
	LastFailedAttemptAt *time.Time `db:"last_failed_attempt_at"`
	LockedUntil         *time.Time `db:"locked_until"`
	EnabledAt           *time.Time `db:"enabled_at"` // When 2FA was successfully enabled
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

// ValidateTOTP validates a TOTP code against the secret
func (t *TwoFactor) ValidateTOTP(code string) bool {
	// Check if it's a valid TOTP format (6 digits)
	if len(code) != TOTPDigits {
		return false
	}

	for _, c := range code {
		if c < '0' || c > '9' {
			return false
		}
	}

	return totp.Validate(code, t.Secret)
}

// ValidateAndConsumeBackupCode validates a backup code and removes it if valid
func (t *TwoFactor) ValidateAndConsumeBackupCode(code string) (bool, error) {
	for i, c := range t.BackupCodes {
		if c == code {
			// Remove the used code
			t.BackupCodes = append(t.BackupCodes[:i], t.BackupCodes[i+1:]...)
			return true, nil
		}
	}

	return false, nil
}

// IsLocked checks if 2FA verification is temporarily locked
func (t *TwoFactor) IsLocked() bool {
	return t.LockedUntil != nil && time.Now().Before(*t.LockedUntil)
}

// ResetFailedAttempts resets the failed attempts counter
func (t *TwoFactor) ResetFailedAttempts() {
	t.FailedAttempts = 0
	t.LastFailedAttemptAt = nil
	t.LockedUntil = nil
	t.UpdatedAt = time.Now()
}

// IncrementFailedAttempts increments failed attempts and handles rate limiting
func (t *TwoFactor) IncrementFailedAttempts() {
	now := time.Now()

	// If last attempt was outside the window, reset counter
	if t.LastFailedAttemptAt == nil || now.Sub(*t.LastFailedAttemptAt) > FailedAttemptsWindow {
		t.FailedAttempts = 1
	} else {
		t.FailedAttempts++
	}

	t.LastFailedAttemptAt = &now
	t.UpdatedAt = now

	// If max attempts reached within window, lock for the remainder of the window
	if t.FailedAttempts >= MaxFailedAttempts {
		// If we have a last failed attempt, lock until 1 hour from that attempt
		// Otherwise, lock for 1 hour from now
		var lockUntil time.Time
		if t.LastFailedAttemptAt != nil {
			lockUntil = t.LastFailedAttemptAt.Add(FailedAttemptsWindow)
		} else {
			lockUntil = now.Add(FailedAttemptsWindow)
		}

		t.LockedUntil = &lockUntil
	}
}
