package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"autopilot/backends/payment/internal/app"
	"autopilot/backends/payment/internal/model"
	"autopilot/backends/payment/internal/service/provider"
	"autopilot/backends/payment/internal/store"
)

var (
	ErrInvalidAmount      = errors.New("invalid payment amount")
	ErrInvalidCurrency    = errors.New("invalid currency code")
	ErrProviderNotFound   = errors.New("payment provider not found")
	ErrMethodNotSupported = errors.New("payment method not supported by provider")
)

// Payment represents the payment service
type Payment struct {
	paymentStore       *store.Payment
	paymentIntentStore *store.PaymentIntent
	paymentMethodStore *store.PaymentMethod
	providers          map[model.PaymentProvider]provider.PaymentProvider
}

// NewPayment creates a new payment service
func NewPayment(
	container *app.Container,
) *Payment {
	// Initialize payment service
	// TODO: Initialize payment providers
	providers := make(map[model.PaymentProvider]provider.PaymentProvider)

	return &Payment{
		paymentStore:       store.NewPayment(container),
		paymentIntentStore: store.NewPaymentIntent(container),
		paymentMethodStore: store.NewPaymentMethod(container),
		providers:          providers,
	}
}

// CreatePaymentIntent initializes a new payment flow
func (s *Payment) CreatePaymentIntent(ctx context.Context, params CreatePaymentIntentParams) (*model.PaymentIntent, error) {
	// Validate amount and currency
	if params.Amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if !isValidCurrency(params.Currency) {
		return nil, ErrInvalidCurrency
	}

	// Get provider implementation
	providerImpl, exists := s.providers[params.Provider]
	if !exists {
		return nil, ErrProviderNotFound
	}

	// Validate payment method support
	if !providerImpl.SupportsPaymentMethod(params.Method) {
		return nil, ErrMethodNotSupported
	}

	// Create payment intent with provider
	providerIntent, err := providerImpl.CreatePaymentIntent(ctx, provider.CreateIntentParams{
		Amount:      params.Amount,
		Currency:    params.Currency,
		Method:      params.Method,
		ReturnURL:   params.ReturnURL,
		WebhookURL:  params.WebhookURL,
		Description: params.Description,
		Metadata:    params.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create provider payment intent: %w", err)
	}

	// Store payment intent
	intent := &model.PaymentIntent{
		ID:           uuid.New(),
		MerchantID:   params.MerchantID,
		Amount:       params.Amount,
		Currency:     params.Currency,
		Status:       model.PaymentStatusPending,
		Provider:     params.Provider,
		Method:       params.Method,
		Description:  params.Description,
		ClientSecret: providerIntent.ClientSecret,
		ReturnURL:    params.ReturnURL,
		WebhookURL:   params.WebhookURL,
		Metadata:     params.Metadata,
	}

	if err := s.paymentIntentStore.CreatePaymentIntent(ctx, intent); err != nil {
		return nil, fmt.Errorf("failed to store payment intent: %w", err)
	}

	return intent, nil
}

// ProcessPayment handles the payment processing flow
func (s *Payment) ProcessPayment(ctx context.Context, params ProcessPaymentParams) (*model.Payment, error) {
	// Get provider implementation
	providerImpl, exists := s.providers[params.Provider]
	if !exists {
		return nil, ErrProviderNotFound
	}

	// Process payment with provider
	providerPayment, err := providerImpl.ProcessPayment(ctx, provider.ProcessPaymentParams{
		Amount:        params.Amount,
		Currency:      params.Currency,
		Method:        params.Method,
		PaymentMethod: params.PaymentMethodID,
		Description:   params.Description,
		Metadata:      params.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment with provider: %w", err)
	}

	// Create payment record
	payment := &model.Payment{
		ID:          uuid.New(),
		MerchantID:  params.MerchantID,
		Amount:      params.Amount,
		Currency:    params.Currency,
		Status:      model.PaymentStatusProcessing,
		Provider:    params.Provider,
		Method:      params.Method,
		Description: params.Description,
		ProviderID:  providerPayment.ProviderID,
		Metadata:    params.Metadata,
	}

	if err := s.paymentStore.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to store payment: %w", err)
	}

	return payment, nil
}

// HandleWebhook processes provider webhooks
func (s *Payment) HandleWebhook(ctx context.Context, provider model.PaymentProvider, payload []byte) error {
	providerImpl, exists := s.providers[provider]
	if !exists {
		return ErrProviderNotFound
	}

	event, err := providerImpl.ParseWebhook(payload)
	if err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Update payment status based on webhook
	if event.PaymentID != "" {
		payment, err := s.paymentStore.GetPayment(ctx, event.PaymentID)
		if err != nil {
			return err
		}

		if err := s.paymentStore.UpdatePaymentStatus(ctx, payment.ID, event.Status, event.ErrorMessage); err != nil {
			return err
		}
	}

	return nil
}

// GetPayment retrieves a payment by ID
func (s *Payment) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	payment, err := s.paymentStore.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// ListPayments retrieves paginated payments for a merchant
func (s *Payment) ListPayments(ctx context.Context, params ListPaymentsParams) ([]*model.Payment, error) {
	payments, err := s.paymentStore.SearchPayments(ctx, store.SearchParams{
		MerchantID: params.MerchantID,
		Status:     params.Status,
		Provider:   params.Provider,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search payments: %w", err)
	}

	return payments, nil
}

// CreatePaymentIntentParams represents the parameters for creating a payment intent
type CreatePaymentIntentParams struct {
	MerchantID  uuid.UUID
	Amount      int64
	Currency    string
	Provider    model.PaymentProvider
	Method      model.PaymentMethodType
	Description string
	ReturnURL   string
	WebhookURL  string
	Metadata    map[string]any
}

// ProcessPaymentParams represents the parameters for processing a payment
type ProcessPaymentParams struct {
	MerchantID      uuid.UUID
	Amount          int64
	Currency        string
	Provider        model.PaymentProvider
	Method          model.PaymentMethodType
	PaymentMethodID *string
	Description     string
	Metadata        map[string]any
}

// ListPaymentsParams represents the parameters for listing payments
type ListPaymentsParams struct {
	MerchantID uuid.UUID
	Status     *model.PaymentStatus
	Provider   *model.PaymentProvider
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

// isValidCurrency validates ISO 4217 currency codes
func isValidCurrency(currency string) bool {
	// Add currency validation logic here
	// For now, just check if it's 3 characters
	return len(currency) == 3
}
