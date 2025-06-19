package v1

import (
	"autopilot/backends/api/internal/payment/model"
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	Headers struct {
		IdempotencyKey string `header:"Idempotency-Key" required:"true" minLength:"1" maxLength:"255" doc:"Unique key to prevent duplicate payments"`
	}
	Body struct {
		MerchantID           uuid.UUID               `json:"merchant_id" required:"true" doc:"ID of the merchant receiving payment"`
		Amount               int64                   `json:"amount" required:"true" minimum:"1" doc:"Amount in smallest currency unit (e.g., cents)"`
		Currency             string                  `json:"currency" required:"true" pattern:"^[A-Z]{3}$" doc:"ISO 4217 currency code"`
		PaymentMethod        model.PaymentMethodType `json:"payment_method" required:"true" doc:"Payment method type"`
		PaymentMethodDetails map[string]any          `json:"payment_method_details,omitempty" doc:"Payment method specific details"`
		CustomerID           *uuid.UUID              `json:"customer_id,omitempty" doc:"Customer ID from identity system"`
		CustomerEmail        *string                 `json:"customer_email,omitempty" format:"email" doc:"Customer email address"`
		CustomerName         *string                 `json:"customer_name,omitempty" doc:"Customer name"`
		Description          *string                 `json:"description,omitempty" maxLength:"500" doc:"Payment description"`
		Metadata             map[string]any          `json:"metadata,omitempty" doc:"Custom metadata for the payment"`
	}
}

// CreatePaymentResponse represents the response for creating a payment
type CreatePaymentResponse struct {
	Body model.Payment
}

// CreatePayment creates a new payment with idempotency support
func (v1 *V1) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {

	// Create payment object
	payment := &model.Payment{
		ID:                   uuid.New(),
		MerchantID:           req.Body.MerchantID,
		IdempotencyKey:       req.Headers.IdempotencyKey,
		Amount:               req.Body.Amount,
		Currency:             req.Body.Currency,
		Status:               model.PaymentStatusPending,
		PaymentMethod:        req.Body.PaymentMethod,
		PaymentMethodDetails: req.Body.PaymentMethodDetails,
		CustomerID:           req.Body.CustomerID,
		CustomerEmail:        req.Body.CustomerEmail,
		CustomerName:         req.Body.CustomerName,
		Description:          req.Body.Description,
		Metadata:             req.Body.Metadata,
		RefundableAmount:     req.Body.Amount, // Initially, full amount is refundable
		InitiatedAt:          time.Now(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Check for existing payment with same idempotency key
	existingPayment, err := v1.payment.Payment.GetPaymentByIdempotencyKey(ctx, req.Body.MerchantID, req.Headers.IdempotencyKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check for existing payment")
	}

	// If payment exists with same idempotency key, return it (idempotent behavior)
	if existingPayment != nil {
		return &CreatePaymentResponse{
			Body: *existingPayment,
		}, nil
	}

	// Create the payment
	createdPayment, err := v1.payment.Payment.CreatePayment(ctx, payment)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create payment")
	}

	return &CreatePaymentResponse{
		Body: *createdPayment,
	}, nil
}

// GetPaymentRequest represents the request to get a payment
type GetPaymentRequest struct {
	PaymentID uuid.UUID `path:"payment_id" required:"true" doc:"Payment ID"`
}

// GetPaymentResponse represents the response for getting a payment
type GetPaymentResponse struct {
	Body model.Payment
}

// GetPayment retrieves a payment by ID
func (v1 *V1) GetPayment(ctx context.Context, req *GetPaymentRequest) (*GetPaymentResponse, error) {
	payment, err := v1.payment.Payment.GetPaymentByID(ctx, req.PaymentID)
	if err != nil {
		return nil, huma.Error404NotFound("Payment not found")
	}

	return &GetPaymentResponse{
		Body: *payment,
	}, nil
}

// ListPaymentsRequest represents the request to list payments
type ListPaymentsRequest struct {
	Query struct {
		MerchantID    *uuid.UUID             `query:"merchant_id" doc:"Filter by merchant ID"`
		CustomerID    *uuid.UUID             `query:"customer_id" doc:"Filter by customer ID"`
		CustomerEmail *string                `query:"customer_email" doc:"Filter by customer email"`
		Status        *model.PaymentStatus   `query:"status" doc:"Filter by payment status"`
		Currency      *string                `query:"currency" doc:"Filter by currency"`
		FromDate      *time.Time             `query:"from_date" doc:"Filter payments created after this date"`
		ToDate        *time.Time             `query:"to_date" doc:"Filter payments created before this date"`
		Limit         int                    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Number of results to return"`
		Offset        int                    `query:"offset" default:"0" minimum:"0" doc:"Number of results to skip"`
	}
}

// ListPaymentsResponse represents the response for listing payments
type ListPaymentsResponse struct {
	Body struct {
		Payments []model.Payment `json:"payments"`
		Total    int64           `json:"total"`
		Limit    int             `json:"limit"`
		Offset   int             `json:"offset"`
	}
}

// ListPayments lists payments with filtering options
func (v1 *V1) ListPayments(ctx context.Context, req *ListPaymentsRequest) (*ListPaymentsResponse, error) {
	filter := &model.PaymentFilter{
		MerchantID:    req.Query.MerchantID,
		CustomerID:    req.Query.CustomerID,
		CustomerEmail: req.Query.CustomerEmail,
		Status:        req.Query.Status,
		Currency:      req.Query.Currency,
		FromDate:      req.Query.FromDate,
		ToDate:        req.Query.ToDate,
		Limit:         req.Query.Limit,
		Offset:        req.Query.Offset,
	}

	payments, total, err := v1.payment.Payment.ListPayments(ctx, filter)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list payments")
	}

	return &ListPaymentsResponse{
		Body: struct {
			Payments []model.Payment `json:"payments"`
			Total    int64           `json:"total"`
			Limit    int             `json:"limit"`
			Offset   int             `json:"offset"`
		}{
			Payments: payments,
			Total:    total,
			Limit:    req.Query.Limit,
			Offset:   req.Query.Offset,
		},
	}, nil
}

// UpdatePaymentStatusRequest represents the request to update payment status
type UpdatePaymentStatusRequest struct {
	PaymentID uuid.UUID `path:"payment_id" required:"true" doc:"Payment ID"`
	Body struct {
		Status               model.PaymentStatus    `json:"status" required:"true" doc:"New payment status"`
		ExternalPaymentID    *string                `json:"external_payment_id,omitempty" doc:"External payment ID from provider"`
		ProviderResponse     map[string]any         `json:"provider_response,omitempty" doc:"Provider response data"`
		ProviderErrorCode    *string                `json:"provider_error_code,omitempty" doc:"Provider error code"`
		ProviderErrorMessage *string                `json:"provider_error_message,omitempty" doc:"Provider error message"`
	}
}

// UpdatePaymentStatusResponse represents the response for updating payment status
type UpdatePaymentStatusResponse struct {
	Body model.Payment
}

// UpdatePaymentStatus updates the status of a payment (webhook endpoint)
func (v1 *V1) UpdatePaymentStatus(ctx context.Context, req *UpdatePaymentStatusRequest) (*UpdatePaymentStatusResponse, error) {
	// Update payment status
	updateData := &model.PaymentStatusUpdate{
		Status:               req.Body.Status,
		ExternalPaymentID:    req.Body.ExternalPaymentID,
		ProviderResponse:     req.Body.ProviderResponse,
		ProviderErrorCode:    req.Body.ProviderErrorCode,
		ProviderErrorMessage: req.Body.ProviderErrorMessage,
	}

	payment, err := v1.payment.Payment.UpdatePaymentStatus(ctx, req.PaymentID, updateData)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update payment status")
	}

	return &UpdatePaymentStatusResponse{
		Body: *payment,
	}, nil
}
