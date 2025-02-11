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
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrPaymentMethodDeleted  = errors.New("payment method has been deleted")
)

// PaymentMethod represents the store for payment methods
type PaymentMethod struct {
	BaseStore
}

// NewPaymentMethod creates a new instance of PaymentMethod
func NewPaymentMethod(container *app.Container) *PaymentMethod {
	return &PaymentMethod{BaseStore: BaseStore{container: container}}
}

// CreatePaymentMethod stores a new payment method with proper validation
func (r *PaymentMethod) CreatePaymentMethod(ctx context.Context, method *model.StoredPaymentMethod) error {
	query := `
		INSERT INTO payment_methods (
			id, merchant_id, customer_id, type, provider, provider_id,
			last4, expiry_month, expiry_year, metadata, created_at,
			updated_at, deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	now := time.Now().UTC()
	method.CreatedAt = now
	method.UpdatedAt = now

	_, err := r.Infra(ctx).DB.Primary.Writer().ExecContext(ctx, query,
		method.ID, method.MerchantID, method.CustomerID, method.Type,
		method.Provider, method.ProviderID, method.Last4,
		method.ExpiryMonth, method.ExpiryYear, method.Metadata,
		method.CreatedAt, method.UpdatedAt, method.DeletedAt,
	)
	return err
}

// GetPaymentMethod retrieves a payment method by ID with soft delete check
func (r *PaymentMethod) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*model.StoredPaymentMethod, error) {
	query := `
		SELECT id, merchant_id, customer_id, type, provider, provider_id,
			last4, expiry_month, expiry_year, metadata, created_at,
			updated_at, deleted_at
		FROM payment_methods
		WHERE id = $1
	`
	method := &model.StoredPaymentMethod{}
	err := r.Infra(ctx).DB.Primary.Writer().GetContext(ctx, method, query, id)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentMethodNotFound
	}
	if err != nil {
		return nil, err
	}

	if method.DeletedAt != nil {
		return nil, ErrPaymentMethodDeleted
	}

	return method, nil
}

// ListPaymentMethodsByCustomer retrieves all active payment methods for a customer
func (r *PaymentMethod) ListPaymentMethodsByCustomer(
	ctx context.Context,
	customerID uuid.UUID,
	limit int,
	offset int,
) ([]*model.StoredPaymentMethod, error) {
	query := `
		SELECT id, merchant_id, customer_id, type, provider, provider_id,
			last4, expiry_month, expiry_year, metadata, created_at,
			updated_at, deleted_at
		FROM payment_methods
		WHERE customer_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var methods []*model.StoredPaymentMethod
	err := r.Infra(ctx).DB.Primary.Reader().SelectContext(ctx, &methods, query, customerID, limit, offset)
	return methods, err
}

// SoftDeletePaymentMethod marks a payment method as deleted
func (r *PaymentMethod) SoftDeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	return r.Infra(ctx).DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Check if payment method exists and is not already deleted
		var deletedAt *time.Time
		err := tx.GetContext(ctx, &deletedAt, "SELECT deleted_at FROM payment_methods WHERE id = $1 FOR UPDATE", id)
		if err == pgx.ErrNoRows {
			return ErrPaymentMethodNotFound
		}
		if err != nil {
			return err
		}

		if deletedAt != nil {
			return ErrPaymentMethodDeleted
		}

		now := time.Now().UTC()
		_, err = tx.ExecContext(ctx, `
			UPDATE payment_methods
			SET deleted_at = $1, updated_at = $1
			WHERE id = $2
		`, now, id)
		return err
	})
}

// GetPaymentMethodByProviderID retrieves a payment method by provider ID
func (r *PaymentMethod) GetPaymentMethodByProviderID(ctx context.Context, providerID string) (*model.StoredPaymentMethod, error) {
	query := `
		SELECT id, merchant_id, customer_id, type, provider, provider_id,
			last4, expiry_month, expiry_year, metadata, created_at,
			updated_at, deleted_at
		FROM payment_methods
		WHERE provider_id = $1
		AND deleted_at IS NULL
	`
	method := &model.StoredPaymentMethod{}
	err := r.Infra(ctx).DB.Primary.Writer().GetContext(ctx, method, query, providerID)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentMethodNotFound
	}
	return method, err
}

// UpdatePaymentMethodMetadata updates the metadata of a payment method
func (r *PaymentMethod) UpdatePaymentMethodMetadata(
	ctx context.Context,
	id uuid.UUID,
	metadata map[string]any,
) error {
	return r.Infra(ctx).DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Check if payment method exists and is not deleted
		var deletedAt *time.Time
		err := tx.GetContext(ctx, &deletedAt, "SELECT deleted_at FROM payment_methods WHERE id = $1 FOR UPDATE", id)
		if err == pgx.ErrNoRows {
			return ErrPaymentMethodNotFound
		}
		if err != nil {
			return err
		}

		if deletedAt != nil {
			return ErrPaymentMethodDeleted
		}

		now := time.Now().UTC()
		_, err = tx.ExecContext(ctx, `
			UPDATE payment_methods
			SET metadata = $1, updated_at = $2
			WHERE id = $3
		`, metadata, now, id)
		return err
	})
}
