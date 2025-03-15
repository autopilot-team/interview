package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
)

// AuditLoger is the store for audit log operations.
type AuditLoger interface {
	Create(ctx context.Context, log *model.AuditLog) (*model.AuditLog, error)
	WithQuerier(q core.Querier) AuditLoger
}

// AuditLog is the store for audit log operations.
type AuditLog struct {
	core.Querier
}

func (s *AuditLog) WithQuerier(q core.Querier) AuditLoger {
	return &AuditLog{q}
}

// NewAuditLog creates a new AuditLog.
func NewAuditLog(db core.Querier) *AuditLog {
	return &AuditLog{db}
}

// Create creates a new audit log entry.
func (s *AuditLog) Create(ctx context.Context, log *model.AuditLog) (*model.AuditLog, error) {
	query := `
		INSERT INTO audit_logs (
			action, resource_type, resource_id, ip_address, metadata, user_agent, user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING
			id, action, resource_type, resource_id, ip_address, metadata, user_agent, user_id, created_at
	`

	var created model.AuditLog
	err := s.QueryRowContext(
		ctx,
		query,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.IPAddress,
		log.Metadata,
		log.UserAgent,
		log.UserID,
	).Scan(
		&created.ID,
		&created.Action,
		&created.ResourceType,
		&created.ResourceID,
		&created.IPAddress,
		&created.Metadata,
		&created.UserAgent,
		&created.UserID,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}
