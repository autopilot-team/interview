package service

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/model"
	"autopilot/backends/api/internal/store"
	"context"
	"time"
)

const (
	// AccountLockoutDuration is the duration for which an account remains locked after too many failed attempts
	AccountLockoutDuration = 30 * time.Minute // minutes

	// MaxFailedLoginAttempts is the maximum number of failed login attempts before account lockout
	MaxFailedLoginAttempts = 5

	// RefreshTokenDuration is the duration for which a refresh token remains valid
	RefreshTokenDuration = 30 * 24 * time.Hour

	// SessionDuration is the duration for which a session token remains valid
	SessionDuration = 24 * time.Hour

	// TempTokenDuration is the duration for which a temporary token remains valid during 2FA flow
	TempTokenDuration = 5 * time.Minute
)

// Sessioner is an interface that wraps the Session methods
type Sessioner interface {
	Create(ctx context.Context, email, password string) (*model.Session, *Error)
	Refresh(ctx context.Context, refreshToken string) (*model.Session, *Error)
	Validate(ctx context.Context, token string) (*model.Session, *Error)
	InvalidateAllSessions(ctx context.Context, userID string) *Error
	GetByToken(ctx context.Context, token string) (*model.Session, *Error)
	Invalidate(ctx context.Context, token string) *Error
	UpdateTwoFactorStatus(ctx context.Context, token string, isPending bool) *Error
}

// Session is the service for session operations.
type Session struct {
	*app.Container
	store *store.Manager
}

// NewSession creates a new Session service.
func NewSession(container *app.Container, store *store.Manager) *Session {
	return &Session{
		Container: container,
		store:     store,
	}
}

// CreateSession creates a new session for a user
func (s *Session) Create(ctx context.Context, email, password string) (*model.Session, *Error) {
	user, err := s.store.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if user.EmailVerifiedAt == nil {
		// Check if there's an existing verification
		verification, err := s.store.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, email)
		if err != nil {
			return nil, NewUnknownError(err)
		}

		// If no verification exists or it's expired, create a new one
		if verification == nil || verification.IsExpired() {
			// Delete any existing verification
			if verification != nil {
				if err := s.store.User.DeleteVerification(ctx, verification.ID); err != nil {
					return nil, NewUnknownError(err)
				}
			}

			// Create new verification and send email
			if err := createEmailVerification(ctx, s.store, s.Container, email, user); err != nil {
				return nil, err
			}
		}

		return nil, ErrEmailNotVerified
	}

	// Check for too many failed login attempts
	if user.FailedLoginAttempts >= MaxFailedLoginAttempts &&
		user.LockedAt != nil &&
		time.Now().Before(user.LockedAt.Add(AccountLockoutDuration)) {
		return nil, ErrAccountLocked
	}

	if !user.VerifyPassword(password) {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= MaxFailedLoginAttempts {
			now := time.Now()
			user.LockedAt = &now
		}
		user.UpdatedAt = time.Now()

		// Update user in database
		if err := s.store.User.Update(ctx, user); err != nil {
			s.Logger.Error("Failed to update failed login attempts", "error", err.Error())
		}

		return nil, ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	if user.FailedLoginAttempts > 0 {
		user.FailedLoginAttempts = 0
		user.LockedAt = nil
		user.UpdatedAt = time.Now()
		if err := s.store.User.Update(ctx, user); err != nil {
			s.Logger.Error("Failed to reset failed login attempts", "error", err.Error())
		}
	}

	// Generate tokens
	accessToken, err := generateSecureToken(32)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	refreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	now := time.Now()

	// Create a new session
	session := &model.Session{
		Token:            accessToken,
		RefreshToken:     refreshToken,
		UserID:           user.ID,
		ExpiresAt:        now.Add(SessionDuration),
		RefreshExpiresAt: now.Add(RefreshTokenDuration),
	}

	// Store the session
	if err := s.store.Session.Create(ctx, session); err != nil {
		return nil, NewUnknownError(err)
	}

	return session, nil
}

// GetByToken retrieves a session by token
func (s *Session) GetByToken(ctx context.Context, token string) (*model.Session, *Error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if session == nil {
		return nil, ErrInvalidSession
	}

	return session, nil
}

// Invalidate invalidates a specific session (sign out)
func (s *Session) Invalidate(ctx context.Context, token string) *Error {
	// Get session first to log the user ID
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return NewUnknownError(err)
	}

	if session == nil {
		return ErrInvalidSession
	}

	if err := s.store.Session.InvalidateByToken(ctx, token); err != nil {
		return NewUnknownError(err)
	}

	return nil
}

// Refresh creates a new session using a refresh token
func (s *Session) Refresh(ctx context.Context, refreshToken string) (*model.Session, *Error) {
	// Get existing session by refresh token
	oldSession, err := s.store.Session.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if oldSession == nil || oldSession.IsTwoFactorPending || time.Now().After(oldSession.RefreshExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	// Invalidate the old session before creating a new one
	if err := s.store.Session.InvalidateByToken(ctx, oldSession.Token); err != nil {
		return nil, NewUnknownError(err)
	}

	// Generate new tokens
	newAccessToken, err := generateSecureToken(32)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	newRefreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	now := time.Now()
	// Create new session
	newSession := &model.Session{
		Token:            newAccessToken,
		RefreshToken:     newRefreshToken,
		UserID:           oldSession.UserID,
		ActiveEntityID:   oldSession.ActiveEntityID,
		IPAddress:        oldSession.IPAddress,
		UserAgent:        oldSession.UserAgent,
		ExpiresAt:        now.Add(SessionDuration),
		RefreshExpiresAt: now.Add(RefreshTokenDuration),
	}

	// Store the new session
	if err := s.store.Session.Create(ctx, newSession); err != nil {
		return nil, NewUnknownError(err)
	}

	return newSession, nil
}

// Validate validates a session token
func (s *Session) Validate(ctx context.Context, token string) (*model.Session, *Error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, NewUnknownError(err)
	}

	if session == nil || time.Now().After(session.ExpiresAt) {
		return nil, ErrInvalidSession
	}

	return session, nil
}

// InvalidateAllSessions invalidates all sessions for a user (sign out from all devices)
func (s *Session) InvalidateAllSessions(ctx context.Context, userID string) *Error {
	if err := s.store.Session.InvalidateByUserID(ctx, userID); err != nil {
		return NewUnknownError(err)
	}

	return nil
}

// UpdateTwoFactorStatus updates the two-factor pending status of a session
func (s *Session) UpdateTwoFactorStatus(ctx context.Context, token string, isPending bool) *Error {
	if err := s.store.Session.UpdateTwoFactorPending(ctx, token, isPending); err != nil {
		return NewUnknownError(err)
	}

	return nil
}
