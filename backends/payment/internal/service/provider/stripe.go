package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/webhook"

	"autopilot/backends/payment/internal/model"
)

var (
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrUnsupportedEventType    = errors.New("unsupported webhook event type")
)

type StripeProvider struct {
	client        *stripe.BackendImplementation
	webhookSecret string
}

func NewStripeProvider(apiKey, webhookSecret string) *StripeProvider {
	stripe.Key = apiKey

	return &StripeProvider{
		client:        stripe.GetBackend(stripe.APIBackend).(*stripe.BackendImplementation),
		webhookSecret: webhookSecret,
	}
}

// CreatePaymentIntent creates a payment intent with Stripe
func (p *StripeProvider) CreatePaymentIntent(ctx context.Context, params CreateIntentParams) (*PaymentIntent, error) {
	// Convert payment method type
	paymentMethod := p.convertPaymentMethodType(params.Method)

	// Create payment intent params
	intentParams := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(params.Amount),
		Currency:           stripe.String(params.Currency),
		PaymentMethodTypes: []*string{stripe.String(paymentMethod)},
		Description:        stripe.String(params.Description),
		Metadata:           p.convertMetadata(params.Metadata),
	}

	// Add return URL if provided
	if params.ReturnURL != "" {
		intentParams.ReturnURL = stripe.String(params.ReturnURL)
	}

	// Create payment intent
	intent, err := paymentintent.New(intentParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe payment intent: %w", err)
	}

	// Convert to provider-agnostic format
	return &PaymentIntent{
		ProviderID:    intent.ID,
		ClientSecret:  intent.ClientSecret,
		Amount:        intent.Amount,
		Currency:      string(intent.Currency),
		Status:        p.convertPaymentStatus(intent.Status),
		ErrorMessage:  p.getErrorMessage(intent.LastPaymentError),
		NextActionURL: p.getNextActionURL(intent.NextAction),
	}, nil
}

// ProcessPayment processes a payment with Stripe
func (p *StripeProvider) ProcessPayment(ctx context.Context, params ProcessPaymentParams) (*Payment, error) {
	// Create payment intent first
	intentParams := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(params.Amount),
		Currency:           stripe.String(params.Currency),
		PaymentMethodTypes: []*string{stripe.String(p.convertPaymentMethodType(params.Method))},
		Description:        stripe.String(params.Description),
		Metadata:           p.convertMetadata(params.Metadata),
		Confirm:            stripe.Bool(true),
	}

	// Attach payment method if provided
	if params.PaymentMethod != nil {
		intentParams.PaymentMethod = stripe.String(*params.PaymentMethod)
	}

	// Create and confirm payment intent
	intent, err := paymentintent.New(intentParams)
	if err != nil {
		return nil, fmt.Errorf("failed to process Stripe payment: %w", err)
	}

	// Convert to provider-agnostic format
	return &Payment{
		ProviderID:    intent.ID,
		Amount:        intent.Amount,
		Currency:      string(intent.Currency),
		Status:        p.convertPaymentStatus(intent.Status),
		ErrorMessage:  p.getErrorMessage(intent.LastPaymentError),
		NextActionURL: p.getNextActionURL(intent.NextAction),
	}, nil
}

// ParseWebhook parses a Stripe webhook event
func (p *StripeProvider) ParseWebhook(payload []byte) (*WebhookEvent, error) {
	// Parse webhook event
	event, err := webhook.ConstructEvent(payload, "", p.webhookSecret)
	if err != nil {
		return nil, ErrInvalidWebhookSignature
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		return p.handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		return p.handlePaymentIntentFailed(event)
	case "payment_intent.canceled":
		return p.handlePaymentIntentCanceled(event)
	default:
		return nil, ErrUnsupportedEventType
	}
}

// SupportsPaymentMethod checks if Stripe supports a payment method
func (p *StripeProvider) SupportsPaymentMethod(method model.PaymentMethodType) bool {
	switch method {
	case model.PaymentMethodTypeCard:
		return true
	case model.PaymentMethodTypeBankTransfer:
		return true
	default:
		return false
	}
}

// Helper functions

func (p *StripeProvider) convertPaymentMethodType(method model.PaymentMethodType) string {
	switch method {
	case model.PaymentMethodTypeCard:
		return "card"
	case model.PaymentMethodTypeBankTransfer:
		return "bank_transfer"
	default:
		return "card" // Default to card
	}
}

func (p *StripeProvider) convertPaymentStatus(status stripe.PaymentIntentStatus) model.PaymentStatus {
	switch status {
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return model.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresConfirmation:
		return model.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresAction:
		return model.PaymentStatusProcessing
	case stripe.PaymentIntentStatusProcessing:
		return model.PaymentStatusProcessing
	case stripe.PaymentIntentStatusSucceeded:
		return model.PaymentStatusSucceeded
	case stripe.PaymentIntentStatusCanceled:
		return model.PaymentStatusCanceled
	default:
		return model.PaymentStatusFailed
	}
}

func (p *StripeProvider) convertMetadata(metadata map[string]any) map[string]string {
	result := make(map[string]string)
	for k, v := range metadata {
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			// Convert non-string values to JSON string
			if bytes, err := json.Marshal(v); err == nil {
				result[k] = string(bytes)
			}
		}
	}
	return result
}

func (p *StripeProvider) getErrorMessage(err *stripe.Error) *string {
	if err == nil {
		return nil
	}

	msg := err.Msg

	return &msg
}

func (p *StripeProvider) getNextActionURL(action *stripe.PaymentIntentNextAction) *string {
	if action == nil || action.RedirectToURL == nil {
		return nil
	}
	return &action.RedirectToURL.URL
}

func (p *StripeProvider) handlePaymentIntentSucceeded(event stripe.Event) (*WebhookEvent, error) {
	var intent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &intent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment intent: %w", err)
	}

	return &WebhookEvent{
		Type:      string(event.Type),
		PaymentID: intent.ID,
		Status:    model.PaymentStatusSucceeded,
		RawData:   event.Data.Object,
	}, nil
}

func (p *StripeProvider) handlePaymentIntentFailed(event stripe.Event) (*WebhookEvent, error) {
	var intent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &intent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment intent: %w", err)
	}

	var errorMessage *string
	if intent.LastPaymentError != nil {
		msg := intent.LastPaymentError.Msg
		errorMessage = &msg
	}

	return &WebhookEvent{
		Type:         string(event.Type),
		PaymentID:    intent.ID,
		Status:       model.PaymentStatusFailed,
		ErrorMessage: errorMessage,
		RawData:      event.Data.Object,
	}, nil
}

func (p *StripeProvider) handlePaymentIntentCanceled(event stripe.Event) (*WebhookEvent, error) {
	var intent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &intent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment intent: %w", err)
	}

	return &WebhookEvent{
		Type:      string(event.Type),
		PaymentID: intent.ID,
		Status:    model.PaymentStatusCanceled,
		RawData:   event.Data.Object,
	}, nil
}
