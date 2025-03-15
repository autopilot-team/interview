package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/internal/types"
	"context"
)

// Entityer defines the interface for entity operations
type Entityer interface {
	Get(ctx context.Context, id string) (*model.Entity, error)
	GetByID(ctx context.Context, id string) (*model.Entity, error)
	GetBySlug(ctx context.Context, mode types.OperationMode, slug string) (*model.Entity, error)
}

// Entity implements the Entityer interface
type Entity struct {
	*app.Container
	store *store.Manager
}

// NewEntity creates a new Entity service
func NewEntity(container *app.Container, store *store.Manager) Entityer {
	return &Entity{
		Container: container,
		store:     store,
	}
}

// GetByID retrieves an entity by ID or slug
func (s *Entity) Get(ctx context.Context, id string) (*model.Entity, error) {
	entity, err := s.store.Entity.Get(ctx, id)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if entity == nil {
		return nil, httpx.ErrEntityNotFound
	}

	return entity, nil
}

// Get retrieves an entity by ID
func (s *Entity) GetByID(ctx context.Context, id string) (*model.Entity, error) {
	entity, err := s.store.Entity.GetByID(ctx, id)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if entity == nil {
		return nil, httpx.ErrEntityNotFound
	}

	return entity, nil
}

// GetBySlug retrieves an entity by slug
func (s *Entity) GetBySlug(ctx context.Context, mode types.OperationMode, slug string) (*model.Entity, error) {
	entity, err := s.store.Entity.GetBySlug(ctx, mode, slug)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if entity == nil {
		return nil, httpx.ErrEntityNotFound
	}

	return entity, nil
}
