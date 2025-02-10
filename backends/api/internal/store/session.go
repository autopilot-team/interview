package store

import (
	"autopilot/backends/api/internal/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
	"fmt"
)

// Sessioner is the store for session operations.
type Sessioner interface {
	Create(ctx context.Context, session *model.Session) error
	GetByToken(ctx context.Context, token string) (*model.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	InvalidateByToken(ctx context.Context, token string) error
	InvalidateByUserID(ctx context.Context, userID string) error
	CleanUpExpired(ctx context.Context) error
	UpdateTwoFactorPending(ctx context.Context, token string, isPending bool) error
}

// Session is the store for session operations.
type Session struct {
	core.DBer
}

// NewSession creates a new Session.
func NewSession(db core.DBer) *Session {
	return &Session{
		db,
	}
}

// Create creates a new session.
func (s *Session) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (
			active_entity_id, expires_at, ip_address, token,
			refresh_token, refresh_expires_at, user_agent, user_id,
			is_two_factor_pending
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9
		)
	`

	_, err := s.Writer().ExecContext(
		ctx,
		query,
		session.ActiveEntityID,
		session.ExpiresAt,
		session.IPAddress,
		session.Token,
		session.RefreshToken,
		session.RefreshExpiresAt,
		session.UserAgent,
		session.UserID,
		session.IsTwoFactorPending,
	)

	if err != nil {
		return fmt.Errorf("inserting session: %w", err)
	}

	return nil
}

// GetByRefreshToken gets a session by refresh token.
func (s *Session) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	var session model.Session
	query := `SELECT * FROM sessions WHERE refresh_token = $1`

	err := s.Writer().QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.ActiveEntityID,
		&session.ExpiresAt,
		&session.IPAddress,
		&session.Token,
		&session.RefreshToken,
		&session.RefreshExpiresAt,
		&session.UserAgent,
		&session.UserID,
		&session.IsTwoFactorPending,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting session by refresh token: %w", err)
	}

	return &session, nil
}

// GetByToken gets a session by token.
func (s *Session) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	query := `SELECT * FROM sessions WHERE token = $1`

	err := s.Writer().QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.ActiveEntityID,
		&session.ExpiresAt,
		&session.IPAddress,
		&session.Token,
		&session.RefreshToken,
		&session.RefreshExpiresAt,
		&session.UserAgent,
		&session.UserID,
		&session.IsTwoFactorPending,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting session by token: %w", err)
	}

	return &session, nil
}

// InvalidateByUserID invalidates all sessions for a user.
func (s *Session) InvalidateByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := s.Writer().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("invalidating sessions by user ID: %w", err)
	}

	return nil
}

// InvalidateByToken invalidates a specific session by token.
func (s *Session) InvalidateByToken(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`

	_, err := s.Writer().ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("invalidating session by token: %w", err)
	}

	return nil
}

// CleanUpExpired removes all expired sessions.
func (s *Session) CleanUpExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`

	_, err := s.Writer().ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("cleaning up expired sessions: %w", err)
	}

	return nil
}

// UpdateTwoFactorPending updates the is_two_factor_pending flag for a session
func (s *Session) UpdateTwoFactorPending(ctx context.Context, token string, isPending bool) error {
	query := `
		UPDATE sessions
		SET is_two_factor_pending = $1,
			updated_at = NOW()
		WHERE token = $2
	`

	_, err := s.Writer().ExecContext(ctx, query, isPending, token)
	if err != nil {
		return fmt.Errorf("updating session two-factor pending status: %w", err)
	}

	return nil
}
