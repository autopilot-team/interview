package service

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/internal/payment/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"context"
)

// Paymenter defines the interface for payment operations
type Paymenter interface {
	Create(ctx context.Context, id string) (*model.Payment, error)
	Get(ctx context.Context, id string) (*model.Payment, error)
}

// Payment implements the Paymenter interface
type Payment struct {
	*app.Container
	store *store.Manager
}

// NewPayment creates a new Payment service
func NewPayment(container *app.Container, store *store.Manager) Paymenter {
	return &Payment{
		Container: container,
		store:     store,
	}
}

// GetByID retrieves a payment by ID or slug
func (s *Payment) Create(ctx context.Context, id string) (*model.Payment, error) {
	entity, err := s.store.WithMode(ctx).Payment.Get(ctx, id)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if entity == nil {
		return nil, httpx.ErrPaymentNotFound
	}

	return entity, nil
}

// GetByID retrieves a payment by the payment ID.
func (s *Payment) Get(ctx context.Context, id string) (*model.Payment, error) {
	entity, err := s.store.WithMode(ctx).Payment.Get(ctx, id)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if entity == nil {
		return nil, httpx.ErrPaymentNotFound
	}

	return entity, nil
}
