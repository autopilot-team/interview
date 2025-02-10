package model

import "time"

// Session represents a user authentication session
type Session struct {
	ID                 string    `db:"id"`
	ActiveEntityID     *string   `db:"active_entity_id"`
	ExpiresAt          time.Time `db:"expires_at"`
	IPAddress          *string   `db:"ip_address"`
	Token              string    `db:"token"`
	RefreshToken       string    `db:"refresh_token"`
	RefreshExpiresAt   time.Time `db:"refresh_expires_at"`
	UserAgent          *string   `db:"user_agent"`
	UserID             string    `db:"user_id"`
	IsTwoFactorPending bool      `db:"is_two_factor_pending"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
