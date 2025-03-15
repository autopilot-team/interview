package model

import (
	"autopilot/backends/internal/types"
	"time"
)

// Session represents a user session
type Session struct {
	ID                 string        `db:"id"`
	Token              string        `db:"token"`
	RefreshToken       string        `db:"refresh_token"`
	UserID             string        `db:"user_id"`
	Memberships        []*Membership `db:"-"` // All memberships for the user
	ExpiresAt          time.Time     `db:"expires_at"`
	RefreshExpiresAt   time.Time     `db:"refresh_expires_at"`
	IsTwoFactorPending bool          `db:"is_two_factor_pending"`
	IPAddress          *string       `db:"ip_address"`
	Country            *string       `db:"country"`
	UserAgent          *string       `db:"user_agent"`
	CreatedAt          time.Time     `db:"created_at"`
	UpdatedAt          time.Time     `db:"updated_at"`
}

// HasPermission checks if the session's active member has the given permission
func (s *Session) HasPermission(entityID string, resource types.Resource, action types.Action) bool {
	for _, m := range s.Memberships {
		if m.EntityID != nil && *m.EntityID == entityID {
			return m.Role.HasPermission(resource, action)
		}
	}
	return false
}

func (s *Session) Role(entityID string) types.Role {
	for _, m := range s.Memberships {
		if m.EntityID != nil && *m.EntityID == entityID {
			return m.Role
		}
	}
	return types.RoleNone
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
