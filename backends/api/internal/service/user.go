package service

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/middleware"
	"autopilot/backends/api/internal/model"
	"autopilot/backends/api/internal/store"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/crypto/bcrypt"
)

// Userer is an interface that wraps the User methods
type Userer interface {
	Create(ctx context.Context, user *model.User, password string) (*model.User, *Error)
	GetByID(ctx context.Context, id string) (*model.User, *Error)
	VerifyEmail(ctx context.Context, token string) *Error
	InitiatePasswordReset(ctx context.Context, email string) *Error
	ResetPassword(ctx context.Context, token string, newPassword string) *Error
}

// User is the service for user operations.
type User struct {
	*app.Container
	store *store.Manager
}

// NewUser creates a new User service.
func NewUser(container *app.Container, store *store.Manager) *User {
	return &User{
		container,
		store,
	}
}

// Create creates a new user.
func (s *User) Create(ctx context.Context, user *model.User, password string) (*model.User, *Error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), model.PasswordHashBcryptCostBcryptCost)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	// Set up user fields
	user.PasswordHash = &[]string{string(hashedPassword)}[0]

	// Check if email already exists
	exists, err := s.store.User.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if exists {
		return nil, ErrEmailExists
	}

	// Create verification record and send email
	if err := createEmailVerification(ctx, s.store, s.Container, user.Email, user); err != nil {
		return nil, err
	}

	// Create user within transaction
	err = s.DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.store.User.Create(ctx, user, tx.Tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, NewUnknownError(err)
	}

	return user, nil
}

// GetByID retrieves user by ID.
func (s *User) GetByID(ctx context.Context, id string) (*model.User, *Error) {
	user, err := s.store.User.GetByID(ctx, id)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// VerifyEmail verifies a user's email address.
func (s *User) VerifyEmail(ctx context.Context, token string) *Error {
	// Get verification
	verification, err := s.store.User.GetVerification(ctx, model.VerificationContextEmailVerification, token)
	if err != nil {
		return NewUnknownError(err)
	}

	if verification == nil {
		return ErrVerificationNotFound
	}

	if verification.IsExpired() {
		return ErrVerificationExpired
	}

	// Get user by email
	user, err := s.store.User.GetByEmail(ctx, verification.Value)
	if err != nil {
		return NewUnknownError(err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Update user and delete verification within transaction
	now := time.Now()
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now

	err = s.DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.store.User.Update(ctx, user, tx.Tx); err != nil {
			return err
		}

		if err := s.store.User.DeleteVerification(ctx, verification.ID, tx.Tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return NewUnknownError(err)
	}

	return nil
}

// InitiatePasswordReset starts the password reset process for a user
func (s *User) InitiatePasswordReset(ctx context.Context, email string) *Error {
	// Check if user exists
	user, err := s.store.User.GetByEmail(ctx, email)
	if err != nil {
		return NewUnknownError(err)
	}

	if user == nil {
		// Return success even if user doesn't exist to prevent email enumeration
		return nil
	}

	// Create verification record
	now := time.Now()
	verification := &model.Verification{
		Context:   model.VerificationContextPasswordReset,
		Value:     email,
		ExpiresAt: now.Add(model.PasswordResetDuration),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create verification within transaction
	err = s.DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.store.User.CreateVerification(ctx, verification, tx.Tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return NewUnknownError(err)
	}

	// Send password reset email
	locale := middleware.GetLocale(ctx)
	t := middleware.GetT(ctx)
	subject, err := t.Localize(&i18n.LocalizeConfig{
		MessageID: "password_reset.title",
		TemplateData: map[string]interface{}{
			"AppName": s.Config.App.Name,
		},
	})
	if err != nil {
		s.Logger.Error("Failed to localize email subject", "error", err.Error())
		subject = fmt.Sprintf("Reset your %s password", s.Config.App.Name)
	}

	if _, err := s.Worker.Insert(ctx, MailerArgs{
		Data: map[string]interface{}{
			"AssetsURL": s.Config.App.AssetsURL,
			"AppName":   s.Config.App.Name,
			"Duration":  model.PasswordResetDuration.Hours(),
			"Email":     user.Email,
			"Name":      user.Name,
			"ResetURL":  fmt.Sprintf("%s/reset-password?token=%s", s.Config.App.DashboardURL, verification.ID),
		},
		Email:    user.Email,
		Locale:   locale,
		Subject:  subject,
		Template: "password_reset",
	}, nil); err != nil {
		s.Logger.Error("Failed to queue password reset email", "error", err.Error())
	}

	return nil
}

// ResetPassword completes the password reset process
func (s *User) ResetPassword(ctx context.Context, token string, newPassword string) *Error {
	// Get verification record
	verification, err := s.store.User.GetVerification(ctx, model.VerificationContextPasswordReset, token)
	if err != nil {
		return NewUnknownError(err)
	}

	if verification == nil || verification.IsExpired() {
		return ErrInvalidOrExpiredToken
	}

	// Get user by email
	user, err := s.store.User.GetByEmail(ctx, verification.Value)
	if err != nil {
		return NewUnknownError(err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), model.PasswordHashBcryptCostBcryptCost)
	if err != nil {
		return NewUnknownError(err)
	}

	user.PasswordHash = &[]string{string(hashedPassword)}[0]
	user.UpdatedAt = time.Now()
	user.FailedLoginAttempts = 0 // Reset failed login attempts
	user.LockedAt = nil          // Remove account lock

	// Update password and delete verification within transaction
	err = s.DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.store.User.Update(ctx, user, tx.Tx); err != nil {
			return err
		}

		if err := s.store.Verification.Delete(ctx, verification.ID, tx.Tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return NewUnknownError(err)
	}

	return nil
}
