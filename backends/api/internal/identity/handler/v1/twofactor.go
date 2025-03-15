package v1

import (
	"autopilot/backends/api/pkg/httpx"
	"context"
	"errors"
	"net/http"
	"time"
)

// DisableTwoFactorRequest is the request body for the disable two-factor endpoint.
type DisableTwoFactorRequest struct{}

// DisableTwoFactorResponse is the response body for the disable two-factor endpoint.
type DisableTwoFactorResponse struct {
	Body struct {
		Enabled bool `json:"enabled" doc:"Whether two-factor authentication is enabled"`
	}
}

// DisableTwoFactor disables two-factor authentication
func (v *V1) DisableTwoFactor(ctx context.Context, input *DisableTwoFactorRequest) (*DisableTwoFactorResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	if err := v.identity.TwoFactor.Disable(ctx, auth.UserID); err != nil {
		v.Logger.Error("Failed to disable two-factor", "error", err)
		return nil, err
	}

	response := &DisableTwoFactorResponse{}
	response.Body.Enabled = false

	return response, nil
}

// EnableTwoFactorRequest is the request body for the enable two-factor endpoint.
type EnableTwoFactorRequest struct {
	Body struct {
		Code string `json:"code" required:"true" doc:"The verification code to confirm setup" example:"123456"`
	}
}

// EnableTwoFactorResponse is the response body for the enable two-factor endpoint.
type EnableTwoFactorResponse struct {
	Body struct {
		Enabled bool `json:"enabled" doc:"Whether two-factor authentication is enabled"`
	}
}

// EnableTwoFactor enables two-factor authentication after setup
func (v *V1) EnableTwoFactor(ctx context.Context, input *EnableTwoFactorRequest) (*EnableTwoFactorResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	if err := v.identity.TwoFactor.Enable(ctx, auth.UserID, input.Body.Code); err != nil {
		v.Logger.Error("Failed to enable two-factor", "error", err)
		return nil, err
	}

	response := &EnableTwoFactorResponse{}
	response.Body.Enabled = true

	return response, nil
}

// RegenerateQRCodeRequest is the request body for the regenerate QR code endpoint.
type RegenerateQRCodeRequest struct{}

// RegenerateQRCodeResponse is the response body for the regenerate QR code endpoint.
type RegenerateQRCodeResponse struct {
	Body struct {
		QRCode string `json:"qr_code" doc:"The QR code for scanning with authenticator apps"`
	}
}

// RegenerateQRCode regenerates the QR code for an existing 2FA setup
func (v *V1) RegenerateQRCode(ctx context.Context, input *RegenerateQRCodeRequest) (*RegenerateQRCodeResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	qrCode, err := v.identity.TwoFactor.RegenerateQRCode(ctx, auth.UserID)
	if err != nil {
		v.Logger.Error("Failed to regenerate QR code", "error", err)
		return nil, err
	}

	response := &RegenerateQRCodeResponse{}
	response.Body.QRCode = qrCode

	return response, nil
}

// SetupTwoFactorRequest is the request body for the setup two-factor endpoint.
type SetupTwoFactorRequest struct{}

// SetupTwoFactorResponse is the response body for the setup two-factor endpoint.
type SetupTwoFactorResponse struct {
	Body struct {
		Secret      string   `json:"secret" doc:"The TOTP secret key"`
		BackupCodes []string `json:"backupCodes" doc:"The backup codes for account recovery"`
		QRCode      string   `json:"qrCode" doc:"The QR code for scanning with authenticator apps"`
	}
}

// SetupTwoFactor initiates two-factor authentication setup
func (v *V1) SetupTwoFactor(ctx context.Context, input *SetupTwoFactorRequest) (*SetupTwoFactorResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	setupData, err := v.identity.TwoFactor.Setup(ctx, auth.UserID)
	if err != nil {
		v.Logger.Error("Failed to setup two-factor", "error", err)
		return nil, err
	}

	response := &SetupTwoFactorResponse{}
	response.Body.Secret = setupData.TwoFactor.Secret
	response.Body.BackupCodes = setupData.BackupCodes
	response.Body.QRCode = setupData.QRCode

	return response, nil
}

// VerifyTwoFactorRequest is the request body for the verify two-factor endpoint.
type VerifyTwoFactorRequest struct {
	Session http.Cookie `cookie:"session" doc:"The session cookie"`
	Body    struct {
		Code string `json:"code" required:"true" doc:"The two-factor authentication code" example:"123456"`
	}
}

// VerifyTwoFactorResponse is the response body for the verify two-factor endpoint.
type VerifyTwoFactorResponse struct {
	SetCookies []http.Cookie `header:"Set-Cookie"`
}

// VerifyTwoFactor verifies a two-factor authentication code during sign-in
func (v *V1) VerifyTwoFactor(ctx context.Context, input *VerifyTwoFactorRequest) (*VerifyTwoFactorResponse, error) {
	// Get session to get user ID
	session, err := v.identity.Session.GetByToken(ctx, input.Session.Value)
	if err != nil && !errors.Is(err, httpx.ErrTwoFactorPending) {
		v.Logger.Error("Failed to get session", "error", err)
		return nil, err
	}

	// Verify the 2FA code with the user ID from the session
	if err := v.identity.TwoFactor.Verify(ctx, session.UserID, input.Body.Code); err != nil {
		v.Logger.Error("Failed to verify two-factor code", "error", err)
		return nil, err
	}

	// Update session to mark 2FA as completed
	if err := v.identity.Session.UpdateTwoFactorStatus(ctx, session.Token, false); err != nil {
		v.Logger.Error("Failed to update session two-factor status", "error", err)
		return nil, err
	}

	// Create new session after successful 2FA by refreshing the current session
	newSession, err := v.identity.Session.Refresh(ctx, session.RefreshToken)
	if err != nil {
		v.Logger.Error("Failed to create new session after successful 2FA", "error", err)
		return nil, err
	}

	// Invalidate the old pending 2FA session
	// Don't return error here as new session is already created
	_ = v.identity.Session.Invalidate(ctx, session.Token)

	response := &VerifyTwoFactorResponse{
		SetCookies: []http.Cookie{
			v.newSessionCookie(
				newSession.Token,
				int(time.Until(newSession.ExpiresAt).Seconds()),
				newSession.ExpiresAt,
			),
			v.newRefreshCookie(
				newSession.RefreshToken,
				int(time.Until(newSession.RefreshExpiresAt).Seconds()),
				newSession.RefreshExpiresAt,
			),
		},
	}
	return response, nil
}
