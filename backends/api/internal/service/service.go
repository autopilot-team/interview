package service

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/middleware"
	"autopilot/backends/api/internal/model"
	"autopilot/backends/api/internal/store"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Manager is a collection of services used by the handlers/workers.
type Manager struct {
	Payment Paymenter
	Session Sessioner
	User    Userer
}

// NewManager creates a new service manager
func NewManager(container *app.Container, store *store.Manager) (*Manager, error) {
	paymentService, err := NewPayment(container)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Payment: paymentService,
		Session: NewSession(container, store),
		User:    NewUser(container, store),
	}, nil
}

var (
	// ErrAccountLocked indicates that the account is temporarily locked due to too many failed attempts
	ErrAccountLocked = NewError("identity.account_locked", "Account is temporarily locked", nil)

	// ErrEmailNotVerified indicates that the email address has not been verified
	ErrEmailNotVerified = NewError("identity.email_not_verified", "Email verification required", nil)

	// ErrInvalidCredentials indicates invalid login credentials
	ErrInvalidCredentials = NewError("identity.invalid_credentials", "Invalid credentials", nil)

	// ErrInvalidRefreshToken indicates an invalid or expired refresh token
	ErrInvalidRefreshToken = NewError("identity.invalid_refresh_token", "Invalid refresh token", nil)

	// ErrInvalidSession indicates an invalid or expired session
	ErrInvalidSession = NewError("identity.invalid_session", "Invalid or expired session", nil)

	// ErrUserNotFound indicates the requested user was not found
	ErrUserNotFound = NewError("identity.user_not_found", "User not found", nil)

	// ErrEmailExists indicates that the email is already registered
	ErrEmailExists = NewError("identity.email_exists", "Email already exists", nil)

	// ErrVerificationNotFound indicates that the verification token was not found
	ErrVerificationNotFound = NewError("identity.verification_not_found", "Verification token not found", nil)

	// ErrVerificationExpired indicates that the verification token has expired
	ErrVerificationExpired = NewError("identity.verification_expired", "Verification token has expired", nil)

	// ErrInvalidOrExpiredToken is returned when a verification token is invalid or expired
	ErrInvalidOrExpiredToken = &Error{
		Code:    "invalid_or_expired_token",
		Message: "The verification token is invalid or has expired",
		Status:  http.StatusBadRequest,
	}
)

// Error is a custom error type that includes a code and an error
type Error struct {
	// Code is the error code
	Code string

	// Err is the underlying error that should only be used for logging and
	// debugging, not for user-facing messages
	Err error

	// Message is the error message
	Message string

	// Status is the HTTP status code associated with the error
	Status int
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new error with the given code and error
func NewError(code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// NewUnknownError creates a new unknown error with the given error
func NewUnknownError(err error) *Error {
	return NewError("unknown_error", "Unknown error", err)
}

// createEmailVerification creates a new email verification record and sends the verification email
func createEmailVerification(ctx context.Context, store *store.Manager, container *app.Container, email string, user *model.User) *Error {
	now := time.Now()
	verification := &model.Verification{
		Context:   model.VerificationContextEmailVerification,
		Value:     email,
		ExpiresAt: now.Add(model.EmailVerificationDuration),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := store.User.CreateVerification(ctx, verification); err != nil {
		return NewUnknownError(err)
	}

	// Queue verification email
	locale := middleware.GetLocale(ctx)
	t := middleware.GetT(ctx)
	if t == nil {
		t = i18n.NewLocalizer(container.I18nBundle.Bundle, locale)
	}

	subject, err := t.Localize(&i18n.LocalizeConfig{
		MessageID: "welcome.title",
		TemplateData: map[string]interface{}{
			"AppName": container.Config.App.Name,
		},
	})
	if err != nil {
		container.Logger.Error("Failed to localize email subject", "error", err.Error())
		subject = fmt.Sprintf("Welcome to %s", container.Config.App.Name)
	}

	if _, err := container.Worker.Insert(ctx, MailerArgs{
		Data: map[string]interface{}{
			"AssetsURL":       container.Config.App.AssetsURL,
			"AppName":         container.Config.App.Name,
			"Duration":        model.EmailVerificationDuration.Hours(),
			"Email":           user.Email,
			"Name":            user.Name,
			"VerificationURL": fmt.Sprintf("%s/verify-email?token=%s", container.Config.App.DashboardURL, verification.ID),
		},
		Email:    user.Email,
		Locale:   locale,
		Subject:  subject,
		Template: "welcome",
	}, nil); err != nil {
		container.Logger.Error("Failed to queue verification email", "error", err.Error())
	}

	return nil
}

// generateSecureToken generates a secure random token of the specified length
func generateSecureToken(length int) (string, error) {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
