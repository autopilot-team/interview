package model

import (
	"autopilot/backends/internal/types"
	"time"
)

// Membership represents a membership of an account/organization/platform
type Membership struct {
	ID        string     `db:"id"`
	EntityID  *string    `db:"entity_id"`
	Role      types.Role `db:"role"`
	UserID    string     `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
