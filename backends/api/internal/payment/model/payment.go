package model

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusProcessing        PaymentStatus = "processing"
	PaymentStatusSucceeded         PaymentStatus = "succeeded"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusCancelled         PaymentStatus = "cancelled"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
	PaymentStatusRefunded          PaymentStatus = "refunded"
)

type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeBankTransfer PaymentMethodType = "bank_transfer"
	PaymentMethodTypeWallet       PaymentMethodType = "wallet"
	PaymentMethodTypeOther        PaymentMethodType = "other"
)

// Merchant represents a payment merchant.
type Merchant struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	EntityID           uuid.UUID         `json:"entity_id" db:"entity_id"`
	Name               string            `json:"name" db:"name"`
	Description        *string           `json:"description,omitempty" db:"description"`
	IsActive           bool              `json:"is_active" db:"is_active"`
	PaymentProvider    string            `json:"payment_provider" db:"payment_provider"`
	ProviderMerchantID string            `json:"provider_merchant_id" db:"provider_merchant_id"`
	Settings           map[string]any    `json:"settings" db:"settings"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`
}

// Payment represents a payment.
type Payment struct {
	ID                   uuid.UUID         `json:"id" db:"id"`
	MerchantID           uuid.UUID         `json:"merchant_id" db:"merchant_id"`
	IdempotencyKey       string            `json:"idempotency_key" db:"idempotency_key"`
	ExternalPaymentID    *string           `json:"external_payment_id,omitempty" db:"external_payment_id"`
	Amount               int64             `json:"amount" db:"amount"`                             // Amount in smallest currency unit
	Currency             string            `json:"currency" db:"currency"`                         // ISO 4217
	Status               PaymentStatus     `json:"status" db:"status"`
	PaymentMethod        PaymentMethodType `json:"payment_method" db:"payment_method"`
	PaymentMethodDetails map[string]any    `json:"payment_method_details" db:"payment_method_details"`
	CustomerID           *uuid.UUID        `json:"customer_id,omitempty" db:"customer_id"`
	CustomerEmail        *string           `json:"customer_email,omitempty" db:"customer_email"`
	CustomerName         *string           `json:"customer_name,omitempty" db:"customer_name"`
	Description          *string           `json:"description,omitempty" db:"description"`
	Metadata             map[string]any    `json:"metadata" db:"metadata"`

	// Provider response data
	ProviderResponse     map[string]any    `json:"provider_response" db:"provider_response"`
	ProviderErrorCode    *string           `json:"provider_error_code,omitempty" db:"provider_error_code"`
	ProviderErrorMessage *string           `json:"provider_error_message,omitempty" db:"provider_error_message"`

	// Refund tracking
	RefundedAmount   int64             `json:"refunded_amount" db:"refunded_amount"`
	RefundableAmount int64             `json:"refundable_amount" db:"refundable_amount"`
	RefundCount      int               `json:"refund_count" db:"refund_count"`

	// Timestamps
	InitiatedAt      time.Time         `json:"initiated_at" db:"initiated_at"`
	ProcessedAt      *time.Time        `json:"processed_at,omitempty" db:"processed_at"`
	FailedAt         *time.Time        `json:"failed_at,omitempty" db:"failed_at"`
	CancelledAt      *time.Time        `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
}

// PaymentEvent represents a payment event for audit trail.
type PaymentEvent struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	PaymentID uuid.UUID      `json:"payment_id" db:"payment_id"`
	EventType string         `json:"event_type" db:"event_type"`
	EventData map[string]any `json:"event_data" db:"event_data"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

// PaymentFilter represents filter options for listing payments
type PaymentFilter struct {
	MerchantID    *uuid.UUID
	CustomerID    *uuid.UUID
	CustomerEmail *string
	Status        *PaymentStatus
	Currency      *string
	FromDate      *time.Time
	ToDate        *time.Time
	Limit         int
	Offset        int
}

// RefundFilter represents filter options for listing refunds
type RefundFilter struct {
	PaymentID   *uuid.UUID
	Status      *RefundStatus
	InitiatedBy *uuid.UUID
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
	Offset      int
}

// PaymentStatusUpdate represents data for updating payment status
type PaymentStatusUpdate struct {
	Status               PaymentStatus
	ExternalPaymentID    *string
	ProviderResponse     map[string]any
	ProviderErrorCode    *string
	ProviderErrorMessage *string
}
