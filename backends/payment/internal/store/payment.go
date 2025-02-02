package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	"autopilot/backends/payment/internal/app"
	"autopilot/backends/payment/internal/model"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidStatus   = errors.New("invalid payment status transition")
)

// Payment represents the payment store
type Payment struct {
	*app.Container
}

// NewPayment creates a new payment store
func NewPayment(container *app.Container) *Payment {
	return &Payment{container}
}

// CreatePayment creates a new payment record with optimistic locking
func (r *Payment) CreatePayment(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (
			id, merchant_id, amount, currency, status, provider, method,
			description, provider_id, error_message, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	now := time.Now().UTC()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	_, err := r.Live.DB.Primary.Writer().ExecContext(ctx, query,
		payment.ID, payment.MerchantID, payment.Amount, payment.Currency,
		payment.Status, payment.Provider, payment.Method, payment.Description,
		payment.ProviderID, payment.ErrorMessage, payment.Metadata,
		payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

// GetPayment retrieves a payment by ID with strong consistency
func (r *Payment) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, provider_id, error_message, metadata, created_at,
			updated_at, completed_at
		FROM payments
		WHERE id = $1
	`
	payment := &model.Payment{}
	err := r.Live.DB.Primary.Reader().GetContext(ctx, payment, query, id)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentNotFound
	}
	return payment, err
}

// ListPaymentsByMerchant retrieves paginated payments for a merchant with cursor-based pagination
func (r *Payment) ListPaymentsByMerchant(
	ctx context.Context,
	merchantID uuid.UUID,
	limit int,
	cursor *time.Time,
) ([]*model.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, provider_id, error_message, metadata, created_at, updated_at, completed_at
		FROM payments
		WHERE merchant_id = $1
		AND ($2::timestamptz IS NULL OR created_at < $2)
		ORDER BY created_at DESC
		LIMIT $3
	`
	var payments []*model.Payment
	err := r.Live.DB.Primary.Reader().SelectContext(ctx, &payments, query, merchantID, cursor, limit)
	return payments, err
}

// UpdatePaymentStatus updates payment status with optimistic locking and validation
func (r *Payment) UpdatePaymentStatus(
	ctx context.Context,
	id uuid.UUID,
	status model.PaymentStatus,
	errorMsg *string,
) error {
	return r.Live.DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Get current payment status
		var currentStatus model.PaymentStatus
		err := tx.GetContext(ctx, &currentStatus, "SELECT status FROM payments WHERE id = $1 FOR UPDATE", id)
		if err == pgx.ErrNoRows {
			return ErrPaymentNotFound
		}
		if err != nil {
			return err
		}

		// Validate status transition
		if !isValidStatusTransition(currentStatus, status) {
			return ErrInvalidStatus
		}

		now := time.Now().UTC()
		var completedAt *time.Time
		if status == model.PaymentStatusSucceeded || status == model.PaymentStatusFailed {
			completedAt = &now
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE payments
			SET status = $1, error_message = $2, updated_at = $3, completed_at = $4
			WHERE id = $5
		`, status, errorMsg, now, completedAt, id)
		return err
	})
}

// SearchPayments performs an efficient search with composite indexes
func (r *Payment) SearchPayments(ctx context.Context, params SearchParams) ([]*model.Payment, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, provider_id, error_message, metadata, created_at, updated_at, completed_at
		FROM payments
		WHERE merchant_id = $1
		AND ($2::payment_status IS NULL OR status = $2)
		AND ($3::payment_provider IS NULL OR provider = $3)
		AND ($4::timestamptz IS NULL OR created_at >= $4)
		AND ($5::timestamptz IS NULL OR created_at <= $5)
		ORDER BY created_at DESC
		LIMIT $6 OFFSET $7
	`
	var payments []*model.Payment
	err := r.Live.DB.Primary.Reader().SelectContext(ctx, &payments, query,
		params.MerchantID,
		params.Status,
		params.Provider,
		params.StartDate,
		params.EndDate,
		params.Limit,
		params.Offset,
	)
	return payments, err
}

type SearchParams struct {
	MerchantID uuid.UUID
	Status     *model.PaymentStatus
	Provider   *model.PaymentProvider
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

// isValidStatusTransition validates payment status transitions
func isValidStatusTransition(from, to model.PaymentStatus) bool {
	transitions := map[model.PaymentStatus][]model.PaymentStatus{
		model.PaymentStatusPending: {
			model.PaymentStatusProcessing,
			model.PaymentStatusFailed,
			model.PaymentStatusCanceled,
		},
		model.PaymentStatusProcessing: {
			model.PaymentStatusSucceeded,
			model.PaymentStatusFailed,
		},
		model.PaymentStatusSucceeded: {
			model.PaymentStatusRefunded,
		},
	}

	allowed, exists := transitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
