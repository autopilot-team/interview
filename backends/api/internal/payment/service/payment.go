package service

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/internal/payment/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Paymenter defines the interface for payment operations
type Paymenter interface {
	CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*model.Payment, error)
	GetPaymentByIdempotencyKey(ctx context.Context, merchantID uuid.UUID, idempotencyKey string) (*model.Payment, error)
	ListPayments(ctx context.Context, filter *model.PaymentFilter) ([]model.Payment, int64, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, update *model.PaymentStatusUpdate) (*model.Payment, error)
	CreateRefund(ctx context.Context, refund *model.Refund) (*model.Refund, error)
	GetRefundByID(ctx context.Context, id uuid.UUID) (*model.Refund, error)
	GetRefundByIdempotencyKey(ctx context.Context, paymentID uuid.UUID, idempotencyKey string) (*model.Refund, error)
	ListRefunds(ctx context.Context, filter *model.RefundFilter) ([]model.Refund, int64, error)
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

// CreatePayment creates a new payment
func (s *Payment) CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	// Validate merchant exists and is active
	merchant, err := s.store.WithMode(ctx).Payment.GetMerchant(ctx, payment.MerchantID)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}
	if merchant == nil {
		return nil, httpx.ErrUnknown.WithInternal(fmt.Errorf("merchant not found"))
	}
	if !merchant.IsActive {
		return nil, httpx.ErrUnknown.WithInternal(fmt.Errorf("merchant is not active"))
	}

	// Create payment
	createdPayment, err := s.store.WithMode(ctx).Payment.CreatePayment(ctx, payment)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Create payment event
	event := &model.PaymentEvent{
		ID:        uuid.New(),
		PaymentID: createdPayment.ID,
		EventType: "created",
		EventData: map[string]any{
			"amount":   payment.Amount,
			"currency": payment.Currency,
			"method":   payment.PaymentMethod,
		},
	}
	_ = s.store.WithMode(ctx).Payment.CreatePaymentEvent(ctx, event)

	return createdPayment, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *Payment) GetPaymentByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	payment, err := s.store.WithMode(ctx).Payment.GetPaymentByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httpx.ErrPaymentNotFound
		}
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return payment, nil
}

// GetPaymentByIdempotencyKey retrieves a payment by merchant ID and idempotency key
func (s *Payment) GetPaymentByIdempotencyKey(ctx context.Context, merchantID uuid.UUID, idempotencyKey string) (*model.Payment, error) {
	payment, err := s.store.WithMode(ctx).Payment.GetPaymentByIdempotencyKey(ctx, merchantID, idempotencyKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil to indicate not found
		}
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return payment, nil
}

// ListPayments lists payments with filtering
func (s *Payment) ListPayments(ctx context.Context, filter *model.PaymentFilter) ([]model.Payment, int64, error) {
	payments, total, err := s.store.WithMode(ctx).Payment.ListPayments(ctx, filter)
	if err != nil {
		return nil, 0, httpx.ErrUnknown.WithInternal(err)
	}

	return payments, total, nil
}

// UpdatePaymentStatus updates the status of a payment
func (s *Payment) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, update *model.PaymentStatusUpdate) (*model.Payment, error) {
	// Get current payment
	payment, err := s.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update payment
	updatedPayment, err := s.store.WithMode(ctx).Payment.UpdatePaymentStatus(ctx, id, update)
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	// Create payment event
	event := &model.PaymentEvent{
		ID:        uuid.New(),
		PaymentID: id,
		EventType: fmt.Sprintf("status_changed_%s", update.Status),
		EventData: map[string]any{
			"old_status": payment.Status,
			"new_status": update.Status,
		},
	}
	_ = s.store.WithMode(ctx).Payment.CreatePaymentEvent(ctx, event)

	return updatedPayment, nil
}

// CreateRefund creates a new refund
func (s *Payment) CreateRefund(ctx context.Context, refund *model.Refund) (*model.Refund, error) {
	var createdRefund *model.Refund
	
	// Get the appropriate database based on mode
	mode := types.GetOperationMode(ctx)
	var db core.DBer
	switch mode {
	case types.OperationModeLive:
		db = s.Container.DB.Payment.Live
	case types.OperationModeTest:
		db = s.Container.DB.Payment.Test
	default:
		db = s.Container.DB.Payment.Test
	}
	
	// Use database transaction
	err := db.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Create refund
		var err error
		createdRefund, err = s.store.WithMode(ctx).Payment.CreateRefund(ctx, refund, tx)
		if err != nil {
			return err
		}

		// Update payment refund amounts
		err = s.store.WithMode(ctx).Payment.UpdatePaymentRefundAmounts(ctx, refund.PaymentID, refund.Amount, tx)
		if err != nil {
			return err
		}

		// Create payment event
		event := &model.PaymentEvent{
			ID:        uuid.New(),
			PaymentID: refund.PaymentID,
			EventType: "refund_initiated",
			EventData: map[string]any{
				"refund_id": createdRefund.ID,
				"amount":    refund.Amount,
				"reason":    refund.Reason,
			},
		}
		_ = s.store.WithMode(ctx).Payment.WithQuerier(tx).CreatePaymentEvent(ctx, event)

		return nil
	})
	
	if err != nil {
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return createdRefund, nil
}

// GetRefundByID retrieves a refund by ID
func (s *Payment) GetRefundByID(ctx context.Context, id uuid.UUID) (*model.Refund, error) {
	refund, err := s.store.WithMode(ctx).Payment.GetRefundByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httpx.ErrUnknown.WithInternal(fmt.Errorf("refund not found"))
		}
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return refund, nil
}

// GetRefundByIdempotencyKey retrieves a refund by payment ID and idempotency key
func (s *Payment) GetRefundByIdempotencyKey(ctx context.Context, paymentID uuid.UUID, idempotencyKey string) (*model.Refund, error) {
	refund, err := s.store.WithMode(ctx).Payment.GetRefundByIdempotencyKey(ctx, paymentID, idempotencyKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil to indicate not found
		}
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	return refund, nil
}

// ListRefunds lists refunds with filtering
func (s *Payment) ListRefunds(ctx context.Context, filter *model.RefundFilter) ([]model.Refund, int64, error) {
	refunds, total, err := s.store.WithMode(ctx).Payment.ListRefunds(ctx, filter)
	if err != nil {
		return nil, 0, httpx.ErrUnknown.WithInternal(err)
	}

	return refunds, total, nil
}
