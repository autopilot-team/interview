package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/types"
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
	CleanUpExpired(ctx context.Context) error
	Create(ctx context.Context, email, password string) (*model.Session, error)
	GetByToken(ctx context.Context, token string) (*model.Session, error)
	GetByTokenFull(ctx context.Context, token string) (*model.Session, error)
	ListByToken(ctx context.Context, userID string) ([]*model.Session, error)
	Invalidate(ctx context.Context, token string) error
	InvalidateByID(ctx context.Context, token string, sessionID string) error
	InvalidateAllSessions(ctx context.Context, userID string, token string) error
	Refresh(ctx context.Context, refreshToken string) (*model.Session, error)
	UpdateTwoFactorStatus(ctx context.Context, token string, isPending bool) error
	Validate(ctx context.Context, token string) (*model.Session, error)
}

// Session is the service for session operations.
type Session struct {
	*app.Container
	store     *store.Manager
	TwoFactor TwoFactorer
}

// NewSession creates a new Session service.
func NewSession(container *app.Container, store *store.Manager, twoFactor TwoFactorer) Sessioner {
	return &Session{
		Container: container,
		store:     store,
		TwoFactor: twoFactor,
	}
}

// CleanUpExpired removes all expired sessions from the database
func (s *Session) CleanUpExpired(ctx context.Context) error {
	s.Logger.Info("Starting expired sessions cleanup")

	if err := s.store.Session.CleanUpExpired(ctx); err != nil {
		s.Logger.Error("Failed to clean up expired sessions", "error", err)
		return httpx.ErrUnknown.WithInternal(err)
	}

	s.Logger.Info("Successfully cleaned up expired sessions")
	return nil
}

// CreateSession creates a new session for a user
func (s *Session) Create(ctx context.Context, email, password string) (*model.Session, error) {
	user, err := s.store.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if user == nil {
		return nil, httpx.ErrInvalidCredentials
	}

	if user.EmailVerifiedAt == nil {
		// Check if there's an existing verification
		verification, err := s.store.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, email)
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}

		// If no verification exists or it's expired, create a new one
		if verification == nil || verification.IsExpired() {
			// Delete any existing verification
			if verification != nil {
				if err := s.store.User.DeleteVerification(ctx, verification.ID); err != nil {
					return nil, httpx.ErrUnknown.WithInternal(err)
				}
			}

			// Create new verification and send email
			if err := createEmailVerification(ctx, s.store, s.Container, email, user); err != nil {
				return nil, err
			}
		}

		return nil, httpx.ErrEmailNotVerified
	}

	// Check for too many failed login attempts
	if user.FailedLoginAttempts >= MaxFailedLoginAttempts &&
		user.LockedAt != nil &&
		time.Now().Before(user.LockedAt.Add(AccountLockoutDuration)) {
		return nil, httpx.ErrAccountLocked
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
			s.Logger.Error("Failed to update failed login attempts", "error", err)
		}

		return nil, httpx.ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	if user.FailedLoginAttempts > 0 {
		user.FailedLoginAttempts = 0
		user.LockedAt = nil
		user.UpdatedAt = time.Now()
		if err := s.store.User.Update(ctx, user); err != nil {
			s.Logger.Error("Failed to reset failed login attempts", "error", err)
		}
	}

	// Get user's memberships
	memberships, err := s.store.Membership.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Generate tokens
	accessToken, err := generateSecureToken(32)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	refreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	var (
		ipAddress *string
		userAgent *string
		country   *string
	)

	reqMetadata := middleware.GetRequestMetadata(ctx)
	if reqMetadata != nil {
		ipAddress = &reqMetadata.IPAddress
		userAgent = &reqMetadata.UserAgent
		country = &reqMetadata.Country
	}

	now := time.Now()
	// TODO: Check if any of the user's entities require 2FA
	// This should:
	// 1. Get all entities the user belongs to (through members table)
	// 2. Check if any of those entities have enforced 2FA settings
	// 3. If yes, treat the user as having 2FA enabled regardless of their personal setting
	// 4. This ensures organization-wide security policies are enforced
	// Check if 2FA is required by checking if user has 2FA setup
	twoFactor, err := s.store.TwoFactor.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if twoFactor != nil {
		session := &model.Session{
			ExpiresAt:          now.Add(TempTokenDuration),
			IPAddress:          ipAddress,
			Country:            country,
			IsTwoFactorPending: true,
			RefreshExpiresAt:   now.Add(TempTokenDuration),
			RefreshToken:       refreshToken,
			Token:              accessToken,
			UserID:             user.ID,
			UserAgent:          userAgent,
			Memberships:        memberships,
		}

		// Store the temporary session
		created, err := s.store.Session.Create(ctx, session)
		if err != nil {
			return nil, httpx.ErrUnknown.WithInternal(err)
		}
		session = created

		return session, httpx.ErrTwoFactorPending
	}

	// Create a new session
	session := &model.Session{
		ExpiresAt:        now.Add(SessionDuration),
		IPAddress:        ipAddress,
		Country:          country,
		RefreshExpiresAt: now.Add(RefreshTokenDuration),
		RefreshToken:     refreshToken,
		Token:            accessToken,
		UserID:           user.ID,
		UserAgent:        userAgent,
		Memberships:      memberships,
	}

	// Store the session
	created, err := s.store.Session.Create(ctx, session)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}
	session = created

	// Log session creation
	if err := auditLog(ctx, s.store, types.ResourceSession, types.ActionCreate, session.ID, session.UserID, nil); err != nil {
		return nil, err
	}

	return session, nil
}

// GetByToken retrieves a session by token
func (s *Session) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil {
		return nil, httpx.ErrUnauthenticated
	}

	if session.IsTwoFactorPending {
		return session, httpx.ErrTwoFactorPending
	}

	// Get user's current membership.
	if activeEntity := middleware.GetActiveEntity(ctx); activeEntity != "" {
		m, err := s.store.Membership.GetByEntityIDWithInheritance(ctx, session.UserID, activeEntity)
		if err != nil {
			s.Logger.Warn("invalid entity header found", "userID", session.UserID, "entity", activeEntity)
		} else {
			session.Memberships = append(session.Memberships, m...)
		}
	}

	return session, nil
}

// GetByToken retrieves a session by token
func (s *Session) GetByTokenFull(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil {
		return nil, httpx.ErrUnauthenticated
	}

	if session.IsTwoFactorPending {
		return session, httpx.ErrTwoFactorPending
	}

	// Get user's memberships
	session.Memberships, err = s.store.Membership.GetByUserIDWithInheritance(ctx, session.UserID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return session, nil
}

// ListByToken retrieves all user sessions by token.
func (s *Session) ListByToken(ctx context.Context, token string) ([]*model.Session, error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil {
		return nil, httpx.ErrUnauthenticated
	}

	if session.IsTwoFactorPending {
		return nil, httpx.ErrTwoFactorPending
	}

	sessions, err := s.store.Session.ListByUser(ctx, session.UserID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return sessions, nil
}

// Invalidate invalidates a specific session (sign out)
func (s *Session) Invalidate(ctx context.Context, token string) error {
	// Get session first to log the user ID
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil {
		return httpx.ErrUnauthenticated
	}

	if err := s.store.Session.InvalidateByToken(ctx, token); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Log session invalidation
	if err := auditLog(ctx, s.store, types.ResourceSession, types.ActionDelete, session.ID, session.UserID, nil); err != nil {
		return err
	}

	return nil
}

// InvalidateByID invalidates a specific session (revoke session)
func (s *Session) InvalidateByID(ctx context.Context, token, sessionID string) error {
	// Get session first to validate session owner.
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil {
		return httpx.ErrUnauthenticated
	}

	if err := s.store.Session.InvalidateByID(ctx, sessionID, session.UserID); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Log session invalidation
	if err := auditLog(ctx, s.store, types.ResourceSession, types.ActionDelete, sessionID, session.UserID, nil); err != nil {
		return err
	}

	return nil
}

// InvalidateAllSessions invalidates all sessions for a user (revoke all other sessions)
func (s *Session) InvalidateAllSessions(ctx context.Context, userID string, token string) error {
	if err := s.store.Session.InvalidateByUserID(ctx, userID, token); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Log session invalidation
	if err := auditLog(ctx, s.store, types.ResourceSession, types.ActionDelete, userID, userID, nil); err != nil {
		return err
	}

	return nil
}

// Refresh creates a new session using a refresh token
func (s *Session) Refresh(ctx context.Context, refreshToken string) (*model.Session, error) {
	// Get existing session by refresh token
	oldSession, err := s.store.Session.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if oldSession == nil || oldSession.IsTwoFactorPending || time.Now().After(oldSession.RefreshExpiresAt) {
		return nil, httpx.ErrInvalidRefreshToken
	}

	// Invalidate the old session before creating a new one
	if err := s.store.Session.InvalidateByToken(ctx, oldSession.Token); err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Generate new tokens
	newAccessToken, err := generateSecureToken(32)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	newRefreshToken, err := generateSecureToken(32)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Get user's memberships
	memberships, err := s.store.Membership.GetByUserID(ctx, oldSession.UserID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	now := time.Now()
	// Create new session
	newSession := &model.Session{
		ExpiresAt:        now.Add(SessionDuration),
		IPAddress:        oldSession.IPAddress,
		Country:          oldSession.Country,
		RefreshExpiresAt: now.Add(RefreshTokenDuration),
		RefreshToken:     newRefreshToken,
		Token:            newAccessToken,
		UserAgent:        oldSession.UserAgent,
		UserID:           oldSession.UserID,
		Memberships:      memberships,
	}

	// Store the new session
	created, err := s.store.Session.Create(ctx, newSession)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}
	newSession = created

	// Log session refresh
	if err := auditLog(ctx, s.store, types.ResourceSession, types.ActionUpdate, newSession.ID, newSession.UserID, nil); err != nil {
		return nil, err
	}

	return newSession, nil
}

// Validate validates a session token
func (s *Session) Validate(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.store.Session.GetByToken(ctx, token)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if session == nil || time.Now().After(session.ExpiresAt) {
		return nil, httpx.ErrUnauthenticated
	}

	// Get user's memberships
	session.Memberships, err = s.store.Membership.GetByUserID(ctx, session.UserID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return session, nil
}

// UpdateTwoFactorStatus updates the two-factor pending status of a session
func (s *Session) UpdateTwoFactorStatus(ctx context.Context, token string, isPending bool) error {
	if err := s.store.Session.UpdateTwoFactorPending(ctx, token, isPending); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	return nil
}
