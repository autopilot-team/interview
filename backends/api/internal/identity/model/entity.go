package model

import (
	"time"
)

type (
	EntityStatus string
	EntityType   string
)

// EntityStatus constants
const (
	EntityStatusPending   EntityStatus = "pending"
	EntityStatusActive    EntityStatus = "active"
	EntityStatusInactive  EntityStatus = "inactive"
	EntityStatusSuspended EntityStatus = "suspended"
)

// EntityType constants
const (
	EntityTypeAccount      EntityType = "account"
	EntityTypeOrganization EntityType = "organization"
	EntityTypePlatform     EntityType = "platform"
)

// Entity represents a entity in the system
type Entity struct {
	ID        string       `db:"id"`
	Domain    *string      `db:"domain"`
	Logo      *string      `db:"logo"`
	Name      string       `db:"name"`
	ParentID  *string      `db:"parent_id"`
	Slug      string       `db:"slug"`
	Status    EntityStatus `db:"status"`
	Type      EntityType   `db:"type"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

// IsActive checks if the entity is active
func (e *Entity) IsActive() bool {
	return e.Status == EntityStatusActive
}
