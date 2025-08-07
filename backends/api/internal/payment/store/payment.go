package store

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
)

// Paymenter is the interface for the transaction store
type Paymenter interface {
	// Create creates a new connection
	Create(ctx context.Context, transaction *model.Payment) (*model.Payment, error)
	Get(ctx context.Context, id string) (*model.Payment, error)
	// WithQuerier returns a new Transaactioner with the given querier
	WithQuerier(core.Querier) Paymenter
}

// Payment is the implementation of the Paymenter interface
type Payment struct {
	core.Querier
}

// NewPayment creates a new transaction store
func NewPayment(q core.Querier) Paymenter {
	return &Payment{q}
}

// WithQuerier returns a new Paymenter with the given querier
func (s *Payment) WithQuerier(q core.Querier) Paymenter {
	return &Payment{q}
}

// Create creates a new transaction
func (s *Payment) Create(ctx context.Context, transaction *model.Payment) (*model.Payment, error) {
	query := `
		INSERT INTO payments (
			merchant_id, amount, currency, status, provider, method, description, error_message, metadata, completed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING
			id, created_at, updated_at
	`

	created := *transaction
	err := s.QueryRowContext(
		ctx,
		query,
		transaction.MerchantID,
		transaction.Amount,
		transaction.Currency,
		transaction.Status,
		transaction.Provider,
		transaction.Method,
		transaction.Description,
		transaction.ErrorMessage,
		transaction.Metadata,
		transaction.CompletedAt,
	).Scan(
		&created.ID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID gets a transaction by transaction ID
func (s *Payment) Get(ctx context.Context, id string) (*model.Payment, error) {
	query := `
		SELECT
		FROM
			payments
		WHERE
			id = $1`

	payment := &model.Payment{}
	err := s.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return payment, nil
}
