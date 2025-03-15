package model

import (
	"autopilot/backends/internal/types"
	"time"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           string         `db:"id"`
	Action       types.Action   `db:"action"`
	ResourceID   string         `db:"resource_id"`
	ResourceType types.Resource `db:"resource_type"`
	IPAddress    *string        `db:"ip_address"`
	Metadata     []byte         `db:"metadata"`
	UserAgent    *string        `db:"user_agent"`
	UserID       string         `db:"user_id"`
	CreatedAt    time.Time      `db:"created_at"`
}
