package store

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Paymenter is the interface for the payment store
type Paymenter interface {
	// Merchant operations
	GetMerchant(ctx context.Context, id uuid.UUID) (*model.Merchant, error)
	
	// Payment operations
	CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*model.Payment, error)
	GetPaymentByIdempotencyKey(ctx context.Context, merchantID uuid.UUID, idempotencyKey string) (*model.Payment, error)
	ListPayments(ctx context.Context, filter *model.PaymentFilter) ([]model.Payment, int64, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, update *model.PaymentStatusUpdate) (*model.Payment, error)
	UpdatePaymentRefundAmounts(ctx context.Context, paymentID uuid.UUID, refundAmount int64, tx core.Querier) error
	
	// Refund operations
	CreateRefund(ctx context.Context, refund *model.Refund, tx core.Querier) (*model.Refund, error)
	GetRefundByID(ctx context.Context, id uuid.UUID) (*model.Refund, error)
	GetRefundByIdempotencyKey(ctx context.Context, paymentID uuid.UUID, idempotencyKey string) (*model.Refund, error)
	ListRefunds(ctx context.Context, filter *model.RefundFilter) ([]model.Refund, int64, error)
	CancelRefund(ctx context.Context, id uuid.UUID, reason *string) (*model.Refund, error)
	UpdateRefundStatus(ctx context.Context, id uuid.UUID, update *model.RefundStatusUpdate) (*model.Refund, error)
	
	// Payment event operations
	CreatePaymentEvent(ctx context.Context, event *model.PaymentEvent) error
	
	// WithQuerier returns a new Paymenter with the given querier
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

// GetMerchant retrieves a merchant by ID
func (s *Payment) GetMerchant(ctx context.Context, id uuid.UUID) (*model.Merchant, error) {
	query := `
		SELECT id, entity_id, name, description, is_active, payment_provider, 
		       provider_merchant_id, settings, created_at, updated_at
		FROM merchants
		WHERE id = $1`

	merchant := &model.Merchant{}
	err := s.QueryRowContext(ctx, query, id).Scan(
		&merchant.ID,
		&merchant.EntityID,
		&merchant.Name,
		&merchant.Description,
		&merchant.IsActive,
		&merchant.PaymentProvider,
		&merchant.ProviderMerchantID,
		&merchant.Settings,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return merchant, err
}

// CreatePayment creates a new payment
func (s *Payment) CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	query := `
		INSERT INTO payments (
			id, merchant_id, idempotency_key, external_payment_id, amount, currency,
			status, payment_method, payment_method_details, customer_id, customer_email,
			customer_name, description, metadata, provider_response, provider_error_code,
			provider_error_message, refunded_amount, refundable_amount, refund_count,
			initiated_at, processed_at, failed_at, cancelled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24
		)
		RETURNING created_at, updated_at`

	err := s.QueryRowContext(
		ctx,
		query,
		payment.ID,
		payment.MerchantID,
		payment.IdempotencyKey,
		payment.ExternalPaymentID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethod,
		payment.PaymentMethodDetails,
		payment.CustomerID,
		payment.CustomerEmail,
		payment.CustomerName,
		payment.Description,
		payment.Metadata,
		payment.ProviderResponse,
		payment.ProviderErrorCode,
		payment.ProviderErrorMessage,
		payment.RefundedAmount,
		payment.RefundableAmount,
		payment.RefundCount,
		payment.InitiatedAt,
		payment.ProcessedAt,
		payment.FailedAt,
		payment.CancelledAt,
	).Scan(
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return payment, err
}

// GetPaymentByID retrieves a payment by ID
func (s *Payment) GetPaymentByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	query := `
		SELECT id, merchant_id, idempotency_key, external_payment_id, amount, currency,
		       status, payment_method, payment_method_details, customer_id, customer_email,
		       customer_name, description, metadata, provider_response, provider_error_code,
		       provider_error_message, refunded_amount, refundable_amount, refund_count,
		       initiated_at, processed_at, failed_at, cancelled_at, created_at, updated_at
		FROM payments
		WHERE id = $1`

	payment := &model.Payment{}
	err := s.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.MerchantID,
		&payment.IdempotencyKey,
		&payment.ExternalPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.PaymentMethodDetails,
		&payment.CustomerID,
		&payment.CustomerEmail,
		&payment.CustomerName,
		&payment.Description,
		&payment.Metadata,
		&payment.ProviderResponse,
		&payment.ProviderErrorCode,
		&payment.ProviderErrorMessage,
		&payment.RefundedAmount,
		&payment.RefundableAmount,
		&payment.RefundCount,
		&payment.InitiatedAt,
		&payment.ProcessedAt,
		&payment.FailedAt,
		&payment.CancelledAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return payment, err
}

// GetPaymentByIdempotencyKey retrieves a payment by merchant ID and idempotency key
func (s *Payment) GetPaymentByIdempotencyKey(ctx context.Context, merchantID uuid.UUID, idempotencyKey string) (*model.Payment, error) {
	query := `
		SELECT id, merchant_id, idempotency_key, external_payment_id, amount, currency,
		       status, payment_method, payment_method_details, customer_id, customer_email,
		       customer_name, description, metadata, provider_response, provider_error_code,
		       provider_error_message, refunded_amount, refundable_amount, refund_count,
		       initiated_at, processed_at, failed_at, cancelled_at, created_at, updated_at
		FROM payments
		WHERE merchant_id = $1 AND idempotency_key = $2`

	payment := &model.Payment{}
	err := s.QueryRowContext(ctx, query, merchantID, idempotencyKey).Scan(
		&payment.ID,
		&payment.MerchantID,
		&payment.IdempotencyKey,
		&payment.ExternalPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.PaymentMethodDetails,
		&payment.CustomerID,
		&payment.CustomerEmail,
		&payment.CustomerName,
		&payment.Description,
		&payment.Metadata,
		&payment.ProviderResponse,
		&payment.ProviderErrorCode,
		&payment.ProviderErrorMessage,
		&payment.RefundedAmount,
		&payment.RefundableAmount,
		&payment.RefundCount,
		&payment.InitiatedAt,
		&payment.ProcessedAt,
		&payment.FailedAt,
		&payment.CancelledAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return payment, err
}

// ListPayments lists payments with filtering
func (s *Payment) ListPayments(ctx context.Context, filter *model.PaymentFilter) ([]model.Payment, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.MerchantID != nil {
		conditions = append(conditions, fmt.Sprintf("merchant_id = $%d", argIndex))
		args = append(args, *filter.MerchantID)
		argIndex++
	}

	if filter.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argIndex))
		args = append(args, *filter.CustomerID)
		argIndex++
	}

	if filter.CustomerEmail != nil {
		conditions = append(conditions, fmt.Sprintf("customer_email = $%d", argIndex))
		args = append(args, *filter.CustomerEmail)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.Currency != nil {
		conditions = append(conditions, fmt.Sprintf("currency = $%d", argIndex))
		args = append(args, *filter.Currency)
		argIndex++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.ToDate)
		argIndex++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM payments %s", where)
	var total int64
	err := s.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List query
	query := fmt.Sprintf(`
		SELECT id, merchant_id, idempotency_key, external_payment_id, amount, currency,
		       status, payment_method, payment_method_details, customer_id, customer_email,
		       customer_name, description, metadata, provider_response, provider_error_code,
		       provider_error_message, refunded_amount, refundable_amount, refund_count,
		       initiated_at, processed_at, failed_at, cancelled_at, created_at, updated_at
		FROM payments
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.MerchantID,
			&payment.IdempotencyKey,
			&payment.ExternalPaymentID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.PaymentMethodDetails,
			&payment.CustomerID,
			&payment.CustomerEmail,
			&payment.CustomerName,
			&payment.Description,
			&payment.Metadata,
			&payment.ProviderResponse,
			&payment.ProviderErrorCode,
			&payment.ProviderErrorMessage,
			&payment.RefundedAmount,
			&payment.RefundableAmount,
			&payment.RefundCount,
			&payment.InitiatedAt,
			&payment.ProcessedAt,
			&payment.FailedAt,
			&payment.CancelledAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	return payments, total, nil
}

// UpdatePaymentStatus updates payment status
func (s *Payment) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, update *model.PaymentStatusUpdate) (*model.Payment, error) {
	now := time.Now()
	query := `
		UPDATE payments
		SET status = $2,
		    external_payment_id = COALESCE($3, external_payment_id),
		    provider_response = COALESCE($4, provider_response),
		    provider_error_code = COALESCE($5, provider_error_code),
		    provider_error_message = COALESCE($6, provider_error_message),
		    processed_at = CASE WHEN $2 = 'succeeded' THEN $7 ELSE processed_at END,
		    failed_at = CASE WHEN $2 = 'failed' THEN $7 ELSE failed_at END,
		    cancelled_at = CASE WHEN $2 = 'cancelled' THEN $7 ELSE cancelled_at END,
		    updated_at = $7
		WHERE id = $1
		RETURNING id, merchant_id, idempotency_key, external_payment_id, amount, currency,
		          status, payment_method, payment_method_details, customer_id, customer_email,
		          customer_name, description, metadata, provider_response, provider_error_code,
		          provider_error_message, refunded_amount, refundable_amount, refund_count,
		          initiated_at, processed_at, failed_at, cancelled_at, created_at, updated_at`

	payment := &model.Payment{}
	err := s.QueryRowContext(
		ctx,
		query,
		id,
		update.Status,
		update.ExternalPaymentID,
		update.ProviderResponse,
		update.ProviderErrorCode,
		update.ProviderErrorMessage,
		now,
	).Scan(
		&payment.ID,
		&payment.MerchantID,
		&payment.IdempotencyKey,
		&payment.ExternalPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.PaymentMethodDetails,
		&payment.CustomerID,
		&payment.CustomerEmail,
		&payment.CustomerName,
		&payment.Description,
		&payment.Metadata,
		&payment.ProviderResponse,
		&payment.ProviderErrorCode,
		&payment.ProviderErrorMessage,
		&payment.RefundedAmount,
		&payment.RefundableAmount,
		&payment.RefundCount,
		&payment.InitiatedAt,
		&payment.ProcessedAt,
		&payment.FailedAt,
		&payment.CancelledAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return payment, err
}

// UpdatePaymentRefundAmounts updates payment refund amounts
func (s *Payment) UpdatePaymentRefundAmounts(ctx context.Context, paymentID uuid.UUID, refundAmount int64, tx core.Querier) error {
	querier := tx
	if querier == nil {
		querier = s
	}

	query := `
		UPDATE payments
		SET refunded_amount = refunded_amount + $2,
		    refundable_amount = refundable_amount - $2,
		    refund_count = refund_count + 1,
		    status = CASE 
		        WHEN refunded_amount + $2 = amount THEN 'refunded'::payment_status
		        ELSE 'partially_refunded'::payment_status
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := querier.ExecContext(ctx, query, paymentID, refundAmount)
	return err
}

// CreatePaymentEvent creates a payment event
func (s *Payment) CreatePaymentEvent(ctx context.Context, event *model.PaymentEvent) error {
	query := `
		INSERT INTO payment_events (id, payment_id, event_type, event_data)
		VALUES ($1, $2, $3, $4)`

	_, err := s.ExecContext(ctx, query, event.ID, event.PaymentID, event.EventType, event.EventData)
	return err
}

// CreateRefund creates a new refund
func (s *Payment) CreateRefund(ctx context.Context, refund *model.Refund, tx core.Querier) (*model.Refund, error) {
	querier := tx
	if querier == nil {
		querier = s
	}

	query := `
		INSERT INTO refunds (
			id, payment_id, idempotency_key, external_refund_id, amount, currency,
			status, reason, reason_description, metadata, provider_response,
			provider_error_code, provider_error_message, initiated_by,
			initiated_by_email, initiated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING created_at, updated_at`

	err := querier.QueryRowContext(
		ctx,
		query,
		refund.ID,
		refund.PaymentID,
		refund.IdempotencyKey,
		refund.ExternalRefundID,
		refund.Amount,
		refund.Currency,
		refund.Status,
		refund.Reason,
		refund.ReasonDescription,
		refund.Metadata,
		refund.ProviderResponse,
		refund.ProviderErrorCode,
		refund.ProviderErrorMessage,
		refund.InitiatedBy,
		refund.InitiatedByEmail,
		refund.InitiatedAt,
	).Scan(
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	return refund, err
}

// GetRefundByID retrieves a refund by ID
func (s *Payment) GetRefundByID(ctx context.Context, id uuid.UUID) (*model.Refund, error) {
	query := `
		SELECT id, payment_id, idempotency_key, external_refund_id, amount, currency,
		       status, reason, reason_description, metadata, provider_response,
		       provider_error_code, provider_error_message, initiated_by,
		       initiated_by_email, initiated_at, processed_at, failed_at,
		       cancelled_at, created_at, updated_at
		FROM refunds
		WHERE id = $1`

	refund := &model.Refund{}
	err := s.QueryRowContext(ctx, query, id).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.IdempotencyKey,
		&refund.ExternalRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.ReasonDescription,
		&refund.Metadata,
		&refund.ProviderResponse,
		&refund.ProviderErrorCode,
		&refund.ProviderErrorMessage,
		&refund.InitiatedBy,
		&refund.InitiatedByEmail,
		&refund.InitiatedAt,
		&refund.ProcessedAt,
		&refund.FailedAt,
		&refund.CancelledAt,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	return refund, err
}

// GetRefundByIdempotencyKey retrieves a refund by payment ID and idempotency key
func (s *Payment) GetRefundByIdempotencyKey(ctx context.Context, paymentID uuid.UUID, idempotencyKey string) (*model.Refund, error) {
	query := `
		SELECT id, payment_id, idempotency_key, external_refund_id, amount, currency,
		       status, reason, reason_description, metadata, provider_response,
		       provider_error_code, provider_error_message, initiated_by,
		       initiated_by_email, initiated_at, processed_at, failed_at,
		       cancelled_at, created_at, updated_at
		FROM refunds
		WHERE payment_id = $1 AND idempotency_key = $2`

	refund := &model.Refund{}
	err := s.QueryRowContext(ctx, query, paymentID, idempotencyKey).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.IdempotencyKey,
		&refund.ExternalRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.ReasonDescription,
		&refund.Metadata,
		&refund.ProviderResponse,
		&refund.ProviderErrorCode,
		&refund.ProviderErrorMessage,
		&refund.InitiatedBy,
		&refund.InitiatedByEmail,
		&refund.InitiatedAt,
		&refund.ProcessedAt,
		&refund.FailedAt,
		&refund.CancelledAt,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	return refund, err
}

// ListRefunds lists refunds with filtering
func (s *Payment) ListRefunds(ctx context.Context, filter *model.RefundFilter) ([]model.Refund, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Base query parts
	baseSelect := `
		SELECT r.id, r.payment_id, r.idempotency_key, r.external_refund_id, r.amount, r.currency,
		       r.status, r.reason, r.reason_description, r.metadata, r.provider_response,
		       r.provider_error_code, r.provider_error_message, r.initiated_by,
		       r.initiated_by_email, r.initiated_at, r.processed_at, r.failed_at,
		       r.cancelled_at, r.created_at, r.updated_at
		FROM refunds r`
	
	joins := ""
	
	// If filtering by merchant, we need to join with payments table
	if filter.MerchantID != nil {
		joins = " INNER JOIN payments p ON r.payment_id = p.id"
		conditions = append(conditions, fmt.Sprintf("p.merchant_id = $%d", argIndex))
		args = append(args, *filter.MerchantID)
		argIndex++
	}

	if filter.PaymentID != nil {
		conditions = append(conditions, fmt.Sprintf("r.payment_id = $%d", argIndex))
		args = append(args, *filter.PaymentID)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("r.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.InitiatedBy != nil {
		conditions = append(conditions, fmt.Sprintf("r.initiated_by = $%d", argIndex))
		args = append(args, *filter.InitiatedBy)
		argIndex++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("r.created_at >= $%d", argIndex))
		args = append(args, *filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("r.created_at <= $%d", argIndex))
		args = append(args, *filter.ToDate)
		argIndex++
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM refunds r%s%s", joins, where)
	var total int64
	err := s.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List query
	query := fmt.Sprintf(`%s%s%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseSelect, joins, where, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var refunds []model.Refund
	for rows.Next() {
		var refund model.Refund
		err := rows.Scan(
			&refund.ID,
			&refund.PaymentID,
			&refund.IdempotencyKey,
			&refund.ExternalRefundID,
			&refund.Amount,
			&refund.Currency,
			&refund.Status,
			&refund.Reason,
			&refund.ReasonDescription,
			&refund.Metadata,
			&refund.ProviderResponse,
			&refund.ProviderErrorCode,
			&refund.ProviderErrorMessage,
			&refund.InitiatedBy,
			&refund.InitiatedByEmail,
			&refund.InitiatedAt,
			&refund.ProcessedAt,
			&refund.FailedAt,
			&refund.CancelledAt,
			&refund.CreatedAt,
			&refund.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		refunds = append(refunds, refund)
	}

	return refunds, total, nil
}

// CancelRefund cancels a pending or processing refund
func (s *Payment) CancelRefund(ctx context.Context, id uuid.UUID, reason *string) (*model.Refund, error) {
	query := `
		UPDATE refunds
		SET status = $1,
		    cancelled_at = $2,
		    metadata = CASE 
		        WHEN $3::text IS NOT NULL 
		        THEN jsonb_set(COALESCE(metadata, '{}'::jsonb), '{cancellation_reason}', to_jsonb($3::text))
		        ELSE metadata
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND status IN ('pending', 'processing')
		RETURNING id, payment_id, idempotency_key, external_refund_id, amount, currency,
		          status, reason, reason_description, metadata, provider_response,
		          provider_error_code, provider_error_message, initiated_by,
		          initiated_by_email, initiated_at, processed_at, failed_at,
		          cancelled_at, created_at, updated_at`

	refund := &model.Refund{}
	err := s.QueryRowContext(ctx, query, model.RefundStatusCancelled, time.Now(), reason, id).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.IdempotencyKey,
		&refund.ExternalRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.ReasonDescription,
		&refund.Metadata,
		&refund.ProviderResponse,
		&refund.ProviderErrorCode,
		&refund.ProviderErrorMessage,
		&refund.InitiatedBy,
		&refund.InitiatedByEmail,
		&refund.InitiatedAt,
		&refund.ProcessedAt,
		&refund.FailedAt,
		&refund.CancelledAt,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refund not found or cannot be cancelled")
	}

	return refund, err
}

// UpdateRefundStatus updates the status of a refund
func (s *Payment) UpdateRefundStatus(ctx context.Context, id uuid.UUID, update *model.RefundStatusUpdate) (*model.Refund, error) {
	var processedAt, failedAt *time.Time
	now := time.Now()

	switch update.Status {
	case model.RefundStatusSucceeded:
		processedAt = &now
	case model.RefundStatusFailed:
		failedAt = &now
	}

	query := `
		UPDATE refunds
		SET status = $1,
		    external_refund_id = COALESCE($2, external_refund_id),
		    provider_response = COALESCE($3, provider_response),
		    provider_error_code = COALESCE($4, provider_error_code),
		    provider_error_message = COALESCE($5, provider_error_message),
		    processed_at = COALESCE($6, processed_at),
		    failed_at = COALESCE($7, failed_at),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING id, payment_id, idempotency_key, external_refund_id, amount, currency,
		          status, reason, reason_description, metadata, provider_response,
		          provider_error_code, provider_error_message, initiated_by,
		          initiated_by_email, initiated_at, processed_at, failed_at,
		          cancelled_at, created_at, updated_at`

	refund := &model.Refund{}
	err := s.QueryRowContext(
		ctx,
		query,
		update.Status,
		update.ExternalRefundID,
		update.ProviderResponse,
		update.ProviderErrorCode,
		update.ProviderErrorMessage,
		processedAt,
		failedAt,
		id,
	).Scan(
		&refund.ID,
		&refund.PaymentID,
		&refund.IdempotencyKey,
		&refund.ExternalRefundID,
		&refund.Amount,
		&refund.Currency,
		&refund.Status,
		&refund.Reason,
		&refund.ReasonDescription,
		&refund.Metadata,
		&refund.ProviderResponse,
		&refund.ProviderErrorCode,
		&refund.ProviderErrorMessage,
		&refund.InitiatedBy,
		&refund.InitiatedByEmail,
		&refund.InitiatedAt,
		&refund.ProcessedAt,
		&refund.FailedAt,
		&refund.CancelledAt,
		&refund.CreatedAt,
		&refund.UpdatedAt,
	)

	return refund, err
}
