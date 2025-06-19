package service

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/internal/payment/store"
	"autopilot/backends/api/pkg/app"
	"context"
)

type Refunder interface {
	Initiate(ctx context.Context, id string) (*model.Refund, error)
}

type Refund struct {
	*app.Container
	store *store.Manager
}

func NewRefund(container *app.Container, store *store.Manager) Refunder {
	return &Refund{
		Container: container,
		store:     store,
	}
}

func (s *Refund) Initiate(ctx context.Context, id string) (*model.Refund, error) {
	// TODO: Implement refund initiation
	return nil, nil
}
