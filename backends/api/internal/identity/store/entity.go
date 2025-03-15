package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"database/sql"
)

// Entityer defines the interface for entity store operations
type Entityer interface {
	Create(ctx context.Context, entity *model.Entity) (*model.Entity, error)
	Get(ctx context.Context, id string) (*model.Entity, error)
	GetByID(ctx context.Context, id string) (*model.Entity, error)
	GetBySlug(ctx context.Context, mode types.OperationMode, slug string) (*model.Entity, error)
	WithQuerier(core.Querier) Entityer
}

// Entity implements Entityer interface
type Entity struct {
	core.Querier
}

func (s *Entity) WithQuerier(q core.Querier) Entityer {
	return &Entity{q}
}

// NewEntity creates a new entity store
func NewEntity(db core.Querier) *Entity {
	return &Entity{db}
}

// Create creates a new entity
func (s *Entity) Create(ctx context.Context, entity *model.Entity) (*model.Entity, error) {
	query := `
		INSERT INTO entities (
			domain,
			logo,
			name,
			parent_id,
			slug,
			status,
			type
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING
			id, domain, logo, name, parent_id, slug, status, type, created_at, updated_at
	`

	var created model.Entity
	err := s.QueryRowContext(
		ctx,
		query,
		entity.Domain,
		entity.Logo,
		entity.Name,
		entity.ParentID,
		entity.Slug,
		entity.Status,
		entity.Type,
	).Scan(
		&created.ID,
		&created.Domain,
		&created.Logo,
		&created.Name,
		&created.ParentID,
		&created.Slug,
		&created.Status,
		&created.Type,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// Get retrieves an entity by ID or slug
func (s *Entity) Get(ctx context.Context, id string) (*model.Entity, error) {
	query := `
		SELECT
			id,
			domain,
			logo,
			name,
			parent_id,
			slug,
			status,
			type,
			created_at,
			updated_at
		FROM entities
		WHERE id = $1
		OR slug = $2
	`

	var entity model.Entity
	err := s.QueryRowContext(ctx, query, id, id).Scan(
		&entity.ID,
		&entity.Domain,
		&entity.Logo,
		&entity.Name,
		&entity.ParentID,
		&entity.Slug,
		&entity.Status,
		&entity.Type,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// GetByID retrieves an entity by ID
func (s *Entity) GetByID(ctx context.Context, id string) (*model.Entity, error) {
	query := `
		SELECT
			id,
			domain,
			logo,
			name,
			parent_id,
			slug,
			status,
			type,
			created_at,
			updated_at
		FROM entities
		WHERE id = $1
	`

	var entity model.Entity
	err := s.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		&entity.Domain,
		&entity.Logo,
		&entity.Name,
		&entity.ParentID,
		&entity.Slug,
		&entity.Status,
		&entity.Type,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// GetBySlug retrieves an entity by slug
func (s *Entity) GetBySlug(ctx context.Context, mode types.OperationMode, slug string) (*model.Entity, error) {
	query := `
		SELECT
			id,
			domain,
			logo,
			name,
			parent_id,
			slug,
			status,
			type,
			created_at,
			updated_at
		FROM entities
		WHERE mode = $1 AND slug = $2
	`

	var entity model.Entity
	err := s.QueryRowContext(ctx, query, mode, slug).Scan(
		&entity.ID,
		&entity.Domain,
		&entity.Logo,
		&entity.Name,
		&entity.ParentID,
		&entity.Slug,
		&entity.Status,
		&entity.Type,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &entity, nil
}
