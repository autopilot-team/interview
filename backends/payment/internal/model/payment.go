package model

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded  PaymentStatus = "succeeded"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCanceled   PaymentStatus = "canceled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

// PaymentProvider represents the payment provider.
type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderAdyen  PaymentProvider = "adyen"
)

type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeBankTransfer PaymentMethodType = "bank_transfer"
	PaymentMethodTypeCrypto       PaymentMethodType = "crypto"
)

// Payment represents a payment.
type Payment struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	MerchantID   uuid.UUID         `json:"merchant_id" db:"merchant_id"`
	Amount       int64             `json:"amount" db:"amount"`     // Amount in cents
	Currency     string            `json:"currency" db:"currency"` // ISO 4217
	Status       PaymentStatus     `json:"status" db:"status"`
	Provider     PaymentProvider   `json:"provider" db:"provider"`
	Method       PaymentMethodType `json:"method" db:"method"`
	Description  string            `json:"description" db:"description"`
	ProviderID   string            `json:"provider_id" db:"provider_id"`
	ErrorMessage *string           `json:"error_message,omitempty" db:"error_message"`
	Metadata     map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
	CompletedAt  *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
}

// PaymentIntent represents a payment intent.
type PaymentIntent struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	MerchantID   uuid.UUID         `json:"merchant_id" db:"merchant_id"`
	Amount       int64             `json:"amount" db:"amount"`
	Currency     string            `json:"currency" db:"currency"`
	Status       PaymentStatus     `json:"status" db:"status"`
	Provider     PaymentProvider   `json:"provider" db:"provider"`
	Method       PaymentMethodType `json:"method" db:"method"`
	Description  string            `json:"description" db:"description"`
	ClientSecret string            `json:"client_secret" db:"client_secret"`
	ReturnURL    string            `json:"return_url" db:"return_url"`
	WebhookURL   string            `json:"webhook_url" db:"webhook_url"`
	Metadata     map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	ExpiresAt    time.Time         `json:"expires_at" db:"expires_at"`
}

// StoredPaymentMethod represents a stored payment method.
type StoredPaymentMethod struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	MerchantID  uuid.UUID         `json:"merchant_id" db:"merchant_id"`
	CustomerID  uuid.UUID         `json:"customer_id" db:"customer_id"`
	Type        PaymentMethodType `json:"type" db:"type"`
	Provider    PaymentProvider   `json:"provider" db:"provider"`
	ProviderID  string            `json:"provider_id" db:"provider_id"`
	Last4       string            `json:"last4" db:"last4"`
	ExpiryMonth int               `json:"expiry_month" db:"expiry_month"`
	ExpiryYear  int               `json:"expiry_year" db:"expiry_year"`
	Metadata    map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
}
