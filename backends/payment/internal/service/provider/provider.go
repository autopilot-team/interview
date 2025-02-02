package provider

import (
	"autopilot/backends/payment/internal/model"
	"context"
)

// PaymentProvider defines the interface that all payment providers must implement
type PaymentProvider interface {
	// CreatePaymentIntent initializes a payment with the provider
	CreatePaymentIntent(ctx context.Context, params CreateIntentParams) (*PaymentIntent, error)

	// ProcessPayment processes a payment with the provider
	ProcessPayment(ctx context.Context, params ProcessPaymentParams) (*Payment, error)

	// ParseWebhook parses a webhook payload from the provider
	ParseWebhook(payload []byte) (*WebhookEvent, error)

	// SupportsPaymentMethod checks if the provider supports a payment method
	SupportsPaymentMethod(method model.PaymentMethodType) bool
}

// CreateIntentParams contains parameters for creating a payment intent
type CreateIntentParams struct {
	Amount      int64
	Currency    string
	Method      model.PaymentMethodType
	ReturnURL   string
	WebhookURL  string
	Description string
	Metadata    map[string]any
}

// ProcessPaymentParams contains parameters for processing a payment
type ProcessPaymentParams struct {
	Amount        int64
	Currency      string
	Method        model.PaymentMethodType
	PaymentMethod *string
	Description   string
	Metadata      map[string]any
}

// PaymentIntent represents a payment intent from a provider
type PaymentIntent struct {
	ProviderID    string
	ClientSecret  string
	Amount        int64
	Currency      string
	Status        model.PaymentStatus
	ErrorMessage  *string
	NextActionURL *string
}

// Payment represents a payment from a provider
type Payment struct {
	ProviderID    string
	Amount        int64
	Currency      string
	Status        model.PaymentStatus
	ErrorMessage  *string
	NextActionURL *string
}

// WebhookEvent represents a webhook event from a provider
type WebhookEvent struct {
	Type         string
	PaymentID    string
	Status       model.PaymentStatus
	ErrorMessage *string
	RawData      map[string]any
}
