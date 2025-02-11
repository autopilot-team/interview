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
	ErrPaymentIntentNotFound = errors.New("payment intent not found")
	ErrPaymentIntentExpired  = errors.New("payment intent expired")
)

// PaymentIntent is a store for payment intents
type PaymentIntent struct {
	BaseStore
}

// NewPaymentIntent creates a new PaymentIntent instance
func NewPaymentIntent(container *app.Container) *PaymentIntent {
	return &PaymentIntent{BaseStore: BaseStore{container: container}}
}

// CreatePaymentIntent creates a new payment intent with proper expiration
func (r *PaymentIntent) CreatePaymentIntent(ctx context.Context, intent *model.PaymentIntent) error {
	query := `
		INSERT INTO payment_intents (
			id, merchant_id, amount, currency, status, provider, method,
			description, client_secret, return_url, webhook_url, metadata,
			created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	now := time.Now().UTC()
	intent.CreatedAt = now
	// Set expiration to 1 hour by default
	intent.ExpiresAt = now.Add(1 * time.Hour)

	_, err := r.Infra(ctx).DB.Primary.Writer().ExecContext(ctx, query,
		intent.ID, intent.MerchantID, intent.Amount, intent.Currency,
		intent.Status, intent.Provider, intent.Method, intent.Description,
		intent.ClientSecret, intent.ReturnURL, intent.WebhookURL, intent.Metadata,
		intent.CreatedAt, intent.ExpiresAt,
	)
	return err
}

// GetPaymentIntent retrieves a payment intent by ID and validates expiration
func (r *PaymentIntent) GetPaymentIntent(ctx context.Context, id uuid.UUID) (*model.PaymentIntent, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, client_secret, return_url, webhook_url, metadata,
			created_at, expires_at
		FROM payment_intents
		WHERE id = $1
	`
	intent := &model.PaymentIntent{}
	err := r.Infra(ctx).DB.Primary.Writer().GetContext(ctx, intent, query, id)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentIntentNotFound
	}
	if err != nil {
		return nil, err
	}

	// Check if the payment intent has expired
	if time.Now().UTC().After(intent.ExpiresAt) {
		return nil, ErrPaymentIntentExpired
	}

	return intent, nil
}

// GetPaymentIntentByClientSecret retrieves a payment intent by client secret
func (r *PaymentIntent) GetPaymentIntentByClientSecret(ctx context.Context, clientSecret string) (*model.PaymentIntent, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, client_secret, return_url, webhook_url, metadata,
			created_at, expires_at
		FROM payment_intents
		WHERE client_secret = $1
	`
	intent := &model.PaymentIntent{}
	err := r.Infra(ctx).DB.Primary.Writer().GetContext(ctx, intent, query, clientSecret)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentIntentNotFound
	}
	if err != nil {
		return nil, err
	}

	if time.Now().UTC().After(intent.ExpiresAt) {
		return nil, ErrPaymentIntentExpired
	}

	return intent, nil
}

// UpdatePaymentIntentStatus updates the status of a payment intent with validation
func (r *PaymentIntent) UpdatePaymentIntentStatus(ctx context.Context, id uuid.UUID, status model.PaymentStatus) error {
	return r.Infra(ctx).DB.Primary.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Get current payment intent
		var expiresAt time.Time
		err := tx.GetContext(ctx, &expiresAt, "SELECT expires_at FROM payment_intents WHERE id = $1 FOR UPDATE", id)
		if err == pgx.ErrNoRows {
			return ErrPaymentIntentNotFound
		}
		if err != nil {
			return err
		}

		// Check expiration
		if time.Now().UTC().After(expiresAt) {
			return ErrPaymentIntentExpired
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE payment_intents
			SET status = $1
			WHERE id = $2
		`, status, id)
		return err
	})
}

// ListPendingPaymentIntents retrieves all pending payment intents for cleanup
func (r *PaymentIntent) ListPendingPaymentIntents(ctx context.Context, batchSize int) ([]*model.PaymentIntent, error) {
	query := `
		SELECT id, merchant_id, amount, currency, status, provider, method,
			description, client_secret, return_url, webhook_url, metadata,
			created_at, expires_at
		FROM payment_intents
		WHERE status = $1
		AND expires_at < $2
		LIMIT $3
	`
	var intents []*model.PaymentIntent
	err := r.Infra(ctx).DB.Primary.Reader().SelectContext(ctx, &intents, query,
		model.PaymentStatusPending,
		time.Now().UTC(),
		batchSize,
	)
	return intents, err
}

// CleanupExpiredPaymentIntents removes expired payment intents in batches
func (r *PaymentIntent) CleanupExpiredPaymentIntents(ctx context.Context, batchSize int) error {
	_, err := r.Infra(ctx).DB.Primary.Writer().ExecContext(ctx, `
		DELETE FROM payment_intents
		WHERE status = $1
		AND expires_at < $2
		LIMIT $3
	`,
		model.PaymentStatusPending,
		time.Now().UTC(),
		batchSize,
	)

	return err
}
