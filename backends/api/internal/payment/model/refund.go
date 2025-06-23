package model

import (
	"time"

	"github.com/google/uuid"
)

type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSucceeded  RefundStatus = "succeeded"
	RefundStatusFailed     RefundStatus = "failed"
	RefundStatusCancelled  RefundStatus = "cancelled"
)

type RefundReason string

const (
	RefundReasonDuplicate       RefundReason = "duplicate"
	RefundReasonFraudulent      RefundReason = "fraudulent"
	RefundReasonCustomerRequest RefundReason = "customer_request"
	RefundReasonProductIssue    RefundReason = "product_issue"
	RefundReasonOther           RefundReason = "other"
)

type Refund struct {
	ID                uuid.UUID      `json:"id" db:"id"`
	PaymentID         uuid.UUID      `json:"payment_id" db:"payment_id"`
	IdempotencyKey    string         `json:"idempotency_key" db:"idempotency_key"`
	ExternalRefundID  *string        `json:"external_refund_id,omitempty" db:"external_refund_id"`
	Amount            int64          `json:"amount" db:"amount"`                                     // Amount in smallest currency unit
	Currency          string         `json:"currency" db:"currency"`                                 // ISO 4217
	Status            RefundStatus   `json:"status" db:"status"`
	Reason            RefundReason   `json:"reason" db:"reason"`
	ReasonDescription *string        `json:"reason_description,omitempty" db:"reason_description"`
	Metadata          map[string]any `json:"metadata" db:"metadata"`

	// Provider response data
	ProviderResponse     map[string]any `json:"provider_response" db:"provider_response"`
	ProviderErrorCode    *string        `json:"provider_error_code,omitempty" db:"provider_error_code"`
	ProviderErrorMessage *string        `json:"provider_error_message,omitempty" db:"provider_error_message"`

	// User who initiated the refund
	InitiatedBy      *uuid.UUID `json:"initiated_by,omitempty" db:"initiated_by"`
	InitiatedByEmail *string    `json:"initiated_by_email,omitempty" db:"initiated_by_email"`

	// Timestamps
	InitiatedAt time.Time  `json:"initiated_at" db:"initiated_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	FailedAt    *time.Time `json:"failed_at,omitempty" db:"failed_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// RefundStatusUpdate represents data for updating refund status
type RefundStatusUpdate struct {
	Status               RefundStatus
	ExternalRefundID     *string
	ProviderResponse     map[string]any
	ProviderErrorCode    *string
	ProviderErrorMessage *string
}
