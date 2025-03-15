package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"context"
)

// Membershiper defines the interface for membership operations
type Membershiper interface {
	GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error)
	GetByEntityID(ctx context.Context, entityID string) ([]*model.Membership, error)
	GetByUserIDWithInheritance(ctx context.Context, userID string) ([]*model.Membership, error)
}

// Membership implements the Memberer interface
type Membership struct {
	*app.Container
	store *store.Manager
}

// NewMembership creates a new Membership service
func NewMembership(container *app.Container, store *store.Manager) Membershiper {
	return &Membership{
		Container: container,
		store:     store,
	}
}

// GetByUserID retrieves all memberships for a user
func (s *Membership) GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	memberships, err := s.store.Membership.GetByUserID(ctx, userID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return memberships, nil
}

// GetByEntityID retrieves all memberships for an entity
func (s *Membership) GetByEntityID(ctx context.Context, entityID string) ([]*model.Membership, error) {
	members, err := s.store.Membership.GetByEntityID(ctx, entityID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return members, nil
}

// GetByUserIDWithInheritance retrieves all memberships for a user including inherited ones
func (s *Membership) GetByUserIDWithInheritance(ctx context.Context, userID string) ([]*model.Membership, error) {
	memberships, err := s.store.Membership.GetByUserIDWithInheritance(ctx, userID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return memberships, nil
}
