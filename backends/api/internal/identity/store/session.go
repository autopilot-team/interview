package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
	"fmt"
)

// Sessioner is the store for session operations.
type Sessioner interface {
	CleanUpExpired(ctx context.Context) error
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
	GetByToken(ctx context.Context, token string) (*model.Session, error)
	ListByUser(ctx context.Context, userID string) ([]*model.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	InvalidateByToken(ctx context.Context, token string) error
	InvalidateByUserID(ctx context.Context, userID string, token string) error
	InvalidateByID(ctx context.Context, id, userID string) error
	UpdateTwoFactorPending(ctx context.Context, token string, isPending bool) error
	WithQuerier(q core.Querier) Sessioner
}

// Session is the store for session operations.
type Session struct {
	core.Querier
}

func (s *Session) WithQuerier(q core.Querier) Sessioner {
	return &Session{q}
}

// NewSession creates a new Session.
func NewSession(db core.Querier) Sessioner {
	return &Session{db}
}

// Create creates a new session.
func (s *Session) Create(ctx context.Context, session *model.Session) (*model.Session, error) {
	query := `
		INSERT INTO sessions (
			expires_at, ip_address, token, country,
			refresh_token, refresh_expires_at, user_agent,
			user_id, is_two_factor_pending
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9
		) RETURNING
			id, expires_at, ip_address, country, token,
			refresh_token, refresh_expires_at, user_agent,
			user_id, is_two_factor_pending, created_at, updated_at
	`

	var created model.Session
	err := s.QueryRowContext(
		ctx,
		query,
		session.ExpiresAt,
		session.IPAddress,
		session.Token,
		session.Country,
		session.RefreshToken,
		session.RefreshExpiresAt,
		session.UserAgent,
		session.UserID,
		session.IsTwoFactorPending,
	).Scan(
		&created.ID,
		&created.ExpiresAt,
		&created.IPAddress,
		&created.Country,
		&created.Token,
		&created.RefreshToken,
		&created.RefreshExpiresAt,
		&created.UserAgent,
		&created.UserID,
		&created.IsTwoFactorPending,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// CleanUpExpired removes all expired sessions.
func (s *Session) CleanUpExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`

	_, err := s.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("cleaning up expired sessions: %w", err)
	}

	return nil
}

// GetByRefreshToken gets a session by refresh token.
func (s *Session) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	var session model.Session
	query := `SELECT * FROM sessions WHERE refresh_token = $1`

	err := s.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.ExpiresAt,
		&session.IPAddress,
		&session.Country,
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
		return nil, err
	}
	return &session, nil
}

// GetByToken gets a session by token.
func (s *Session) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	query := `SELECT * FROM sessions WHERE token = $1`

	err := s.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.ExpiresAt,
		&session.IPAddress,
		&session.Country,
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
		return nil, err
	}

	return &session, nil
}

// ListByUser lists all sessions for that user.
func (s *Session) ListByUser(ctx context.Context, userID string) ([]*model.Session, error) {
	query := `SELECT * FROM sessions WHERE user_id = $1`

	rows, err := s.QueryContext(ctx, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.Session
	for rows.Next() {
		var session model.Session
		if err := rows.Scan(
			&session.ID,
			&session.ExpiresAt,
			&session.IPAddress,
			&session.Country,
			&session.Token,
			&session.RefreshToken,
			&session.RefreshExpiresAt,
			&session.UserAgent,
			&session.UserID,
			&session.IsTwoFactorPending,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// InvalidateByUserID invalidates all sessions for a user, except for the provided session.
func (s *Session) InvalidateByUserID(ctx context.Context, userID string, token string) error {
	query := `DELETE FROM sessions WHERE user_id = $1 AND token <> $2`

	_, err := s.ExecContext(ctx, query, userID, token)
	if err != nil {
		return err
	}

	return nil
}

// InvalidateByToken invalidates a specific session by token.
func (s *Session) InvalidateByToken(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`

	_, err := s.ExecContext(ctx, query, token)
	return err
}

// InvalidateByID invalidates a specific session by session ID.
func (s *Session) InvalidateByID(ctx context.Context, id, userID string) error {
	query := `DELETE FROM sessions WHERE id = $1 AND user_id = $2`

	_, err := s.ExecContext(ctx, query, id, userID)
	return err
}

// UpdateTwoFactorPending updates the is_two_factor_pending flag for a session
func (s *Session) UpdateTwoFactorPending(ctx context.Context, token string, isPending bool) error {
	query := `
		UPDATE sessions
		SET is_two_factor_pending = $1,
			updated_at = NOW()
		WHERE token = $2
	`

	_, err := s.ExecContext(ctx, query, isPending, token)
	return err
}
