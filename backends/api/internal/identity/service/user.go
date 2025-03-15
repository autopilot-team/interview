package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/crypto/bcrypt"
)

// Userer is an interface that wraps the User methods
type Userer interface {
	Create(ctx context.Context, user *model.User, password string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	InitiatePasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
	UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) error
	VerifyEmail(ctx context.Context, token string) error
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
func (s *User) Create(ctx context.Context, user *model.User, password string) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), model.PasswordHashBcryptCost)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Set up user fields
	user.PasswordHash = &[]string{string(hashedPassword)}[0]

	// Check if email already exists
	exists, err := s.store.User.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if exists {
		return nil, httpx.ErrEmailExists
	}

	created, err := s.store.User.Create(ctx, user)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Create verification record and send email
	if err := createEmailVerification(ctx, s.store, s.Container, created.Email, created); err != nil {
		return nil, err
	}

	// Create audit log for user creation
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionCreate, created.ID, created.ID, nil); err != nil {
		return nil, err
	}

	return created, nil
}

// GetByID retrieves user by ID.
func (s *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.store.User.GetByID(ctx, id)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if user == nil {
		return nil, httpx.ErrUserNotFound
	}

	// Create audit log for user profile access
	metadata := map[string]any{
		"accessed_at": time.Now(),
		"email":       user.Email,
		"is_verified": user.IsEmailVerified(),
	}
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionRead, user.ID, user.ID, metadata); err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a field.
func (s *User) Update(ctx context.Context, u *model.User) (*model.User, error) {
	user, err := s.store.User.GetByID(ctx, u.ID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Image != nil {
		user.Image = u.Image
	}
	if err := s.store.User.Update(ctx, user); err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionUpdate, u.ID, u.ID, nil); err != nil {
		return nil, err
	}
	return user, nil
}

// VerifyEmail verifies a user's email address.
func (s *User) VerifyEmail(ctx context.Context, token string) error {
	// Get verification
	verification, err := s.store.User.GetVerification(ctx, model.VerificationContextEmailVerification, token)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}
	if verification == nil || verification.IsExpired() {
		return httpx.ErrInvalidOrExpiredToken
	}

	// Get user by email
	user, err := s.store.User.GetByEmail(ctx, verification.Value)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}
	if user == nil {
		return httpx.ErrUserNotFound
	}

	// Update user and delete verification within transaction
	now := time.Now()
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now

	err = s.DB.Identity.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		u := s.store.User.WithQuerier(tx)
		if err := u.Update(ctx, user); err != nil {
			return err
		}

		if err := u.DeleteVerification(ctx, verification.ID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for email verification
	metadata := map[string]any{
		"email":           user.Email,
		"verification_id": verification.ID,
		"verification_at": now,
	}
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionVerify, user.ID, user.ID, metadata); err != nil {
		return err
	}

	return nil
}

// InitiatePasswordReset starts the password reset process for a user
func (s *User) InitiatePasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.store.User.GetByEmail(ctx, email)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
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
	if _, err := s.store.User.CreateVerification(ctx, verification); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for password reset initiation
	metadata := map[string]any{
		"email":           email,
		"verification_id": verification.ID,
		"expiration_time": verification.ExpiresAt,
	}
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionCreate, verification.ID, user.ID, metadata); err != nil {
		return err
	}

	// Send password reset email
	locale := middleware.GetLocale(ctx)
	t := middleware.GetT(ctx)
	subject, err := t.Localize(&i18n.LocalizeConfig{
		MessageID: "password_reset.title",
		TemplateData: map[string]any{
			"AppName": s.Config.App.Name,
		},
	})
	if err != nil {
		s.Logger.Error("Failed to localize email subject", "error", err)
		subject = fmt.Sprintf("Reset your %s password", s.Config.App.Name)
	}

	if _, err := s.Worker.Insert(ctx, MailerArgs{
		Data: map[string]any{
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
		s.Logger.Error("Failed to queue password reset email", "error", err)
	}

	return nil
}

// ResetPassword completes the password reset process
func (s *User) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// Get verification record
	verification, err := s.store.User.GetVerification(ctx, model.VerificationContextPasswordReset, token)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if verification == nil || verification.IsExpired() {
		return httpx.ErrInvalidOrExpiredToken
	}

	// Get user by email
	user, err := s.store.User.GetByEmail(ctx, verification.Value)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}
	if user == nil {
		return httpx.ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), model.PasswordHashBcryptCost)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	user.PasswordHash = &[]string{string(hashedPassword)}[0]
	user.UpdatedAt = time.Now()
	user.FailedLoginAttempts = 0 // Reset failed login attempts
	user.LockedAt = nil          // Remove account lock

	// Update password and delete verification within transaction
	err = s.DB.Identity.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		u := s.store.User.WithQuerier(tx)
		v := s.store.Verification.WithQuerier(tx)

		if err := u.Update(ctx, user); err != nil {
			return err
		}

		if err := v.Delete(ctx, verification.ID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for password reset completion
	metadata := map[string]any{
		"verification_id": verification.ID,
		"reset_at":        time.Now(),
		"unlocked":        user.LockedAt != nil,         // Whether this reset also unlocked the account
		"attempts_reset":  user.FailedLoginAttempts > 0, // Whether failed attempts were reset
	}
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionResetPassword, user.ID, user.ID, metadata); err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates a user's password after verifying their current password
func (s *User) UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) error {
	// Get user
	user, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if user == nil {
		return httpx.ErrUserNotFound
	}

	// Verify current password
	if !user.VerifyPassword(currentPassword) {
		return httpx.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), model.PasswordHashBcryptCost)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Update user password within transaction
	now := time.Now()
	user.PasswordHash = &[]string{string(hashedPassword)}[0]
	user.UpdatedAt = now

	if err := s.store.User.Update(ctx, user); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Create audit log for password update
	metadata := map[string]any{
		"updated_at": now,
	}
	if err := auditLog(ctx, s.store, types.ResourceUser, types.ActionUpdate, user.ID, user.ID, metadata); err != nil {
		return err
	}

	return nil
}
