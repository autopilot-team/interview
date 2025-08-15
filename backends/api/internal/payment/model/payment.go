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

type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeBankTransfer PaymentMethodType = "bank_transfer"
)

// Payment represents a payment.
type Payment struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	MerchantID   uuid.UUID         `json:"merchantId" db:"merchant_id"`
	Amount       int64             `json:"amount" db:"amount"`     // Amount in cents
	Currency     string            `json:"currency" db:"currency"` // ISO 4217
	Status       PaymentStatus     `json:"status" db:"status"`
	Provider     string            `json:"provider" db:"provider"`
	Method       PaymentMethodType `json:"method" db:"method"`
	Description  string            `json:"description" db:"description"`
	ErrorMessage *string           `json:"errorMessage,omitempty" db:"error_message"`
	Metadata     map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt    time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time         `json:"updatedAt" db:"updated_at"`
	CompletedAt  *time.Time        `json:"completedAt,omitempty" db:"completed_at"`
}

// PaymentIntent represents a payment intent.
type PaymentIntent struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	MerchantID   uuid.UUID         `json:"merchantId" db:"merchant_id"`
	Amount       int64             `json:"amount" db:"amount"`
	Currency     string            `json:"currency" db:"currency"`
	Status       PaymentStatus     `json:"status" db:"status"`
	Method       PaymentMethodType `json:"method" db:"method"`
	Description  string            `json:"description" db:"description"`
	ClientSecret string            `json:"clientSecret" db:"client_secret"`
	ReturnURL    string            `json:"returnUrl" db:"return_url"`
	WebhookURL   string            `json:"webhookUrl" db:"webhook_url"`
	Metadata     map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt    time.Time         `json:"createdAt" db:"created_at"`
	ExpiresAt    time.Time         `json:"expiresAt" db:"expires_at"`
}

// StoredPaymentMethod represents a stored payment method.
type StoredPaymentMethod struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	MerchantID  uuid.UUID         `json:"merchantId" db:"merchant_id"`
	CustomerID  uuid.UUID         `json:"customerId" db:"customer_id"`
	Type        PaymentMethodType `json:"type" db:"type"`
	ProviderID  string            `json:"providerId" db:"provider_id"`
	Last4       string            `json:"last4" db:"last4"`
	ExpiryMonth int               `json:"expiryMonth" db:"expiry_month"`
	ExpiryYear  int               `json:"expiryYear" db:"expiry_year"`
	Metadata    map[string]any    `json:"metadata" db:"metadata"`
	CreatedAt   time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time         `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time        `json:"deletedAt,omitempty" db:"deleted_at"`
}
