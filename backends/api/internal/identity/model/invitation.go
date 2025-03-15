package model

import (
	"time"
)

// Invitation represents an invitation to join an account/organization/platform
type Invitation struct {
	ID             string    `db:"id"`
	Email          string    `db:"email"`
	ExpiresAt      time.Time `db:"expires_at"`
	AccountID      *string   `db:"account_id"`
	OrganizationID *string   `db:"organization_id"`
	PlatformID     *string   `db:"platform_id"`
	InviterID      string    `db:"inviter_id"`
	Role           *string   `db:"role"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
