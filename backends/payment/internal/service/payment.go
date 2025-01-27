package service

import (
	"context"
	"errors"

	"autopilot/backends/payment/internal/model"
	"autopilot/backends/payment/internal/store"
)

var (
	ErrInvalidAmount      = errors.New("invalid payment amount")
	ErrInvalidCurrency    = errors.New("invalid currency code")
	ErrProviderNotFound   = errors.New("payment provider not found")
	ErrMethodNotSupported = errors.New("payment method not supported by provider")
)

type Payment struct {
	paymentStore *store.Payment
}

func NewPayment(
	paymentStore *store.Payment,
) *Payment {
	return &Payment{
		paymentStore: paymentStore,
	}
}

// GetPayment retrieves a payment by ID
func (s *Payment) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	payment, err := s.paymentStore.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
