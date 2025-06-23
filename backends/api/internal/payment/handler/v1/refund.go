package v1

import (
	"autopilot/backends/api/internal/payment/model"
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// InitiateRefundRequest represents the request to create a refund
type InitiateRefundRequest struct {
	Headers struct {
		IdempotencyKey string `header:"Idempotency-Key" required:"true" minLength:"1" maxLength:"255" doc:"Unique key to prevent duplicate refunds"`
	}
	Body struct {
		PaymentID         uuid.UUID            `json:"payment_id" required:"true" doc:"ID of the payment to refund"`
		Amount            int64                `json:"amount" required:"true" minimum:"1" doc:"Amount to refund in smallest currency unit"`
		Reason            model.RefundReason   `json:"reason" required:"true" doc:"Reason for the refund"`
		ReasonDescription *string              `json:"reason_description,omitempty" maxLength:"500" doc:"Additional details about the refund reason"`
		Metadata          map[string]any       `json:"metadata,omitempty" doc:"Custom metadata for the refund"`
	}
}

// InitiateRefundResponse represents the response for creating a refund
type InitiateRefundResponse struct {
	Body model.Refund
}

// InitiateRefund creates a new refund with validation and idempotency support
func (v1 *V1) InitiateRefund(ctx context.Context, req *InitiateRefundRequest) (*InitiateRefundResponse, error) {
	// For now, we'll leave user info as nil since we don't have session context
	var initiatedBy *uuid.UUID
	var initiatedByEmail *string

	// Fetch the payment to validate refund
	payment, err := v1.payment.Payment.GetPaymentByID(ctx, req.Body.PaymentID)
	if err != nil {
		return nil, huma.Error404NotFound("Payment not found")
	}

	// Validate payment status allows refunds
	if payment.Status != model.PaymentStatusSucceeded && payment.Status != model.PaymentStatusPartiallyRefunded {
		return nil, huma.Error400BadRequest("Payment must be succeeded or partially refunded to issue a refund")
	}

	// Validate refund amount
	if req.Body.Amount > payment.RefundableAmount {
		return nil, huma.Error400BadRequest("Refund amount exceeds refundable amount")
	}

	// Check for existing refund with same idempotency key
	existingRefund, err := v1.payment.Payment.GetRefundByIdempotencyKey(ctx, req.Body.PaymentID, req.Headers.IdempotencyKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check for existing refund")
	}

	// If refund exists with same idempotency key, return it (idempotent behavior)
	if existingRefund != nil {
		return &InitiateRefundResponse{
			Body: *existingRefund,
		}, nil
	}

	// Check for duplicate refund attempts - list all pending/processing refunds for this payment
	refundFilter := &model.RefundFilter{
		PaymentID: &req.Body.PaymentID,
		Limit:     100,
	}
	existingRefunds, _, err := v1.payment.Payment.ListRefunds(ctx, refundFilter)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check for existing refunds")
	}

	// Check if there's already a pending/processing refund for the same amount
	for _, refund := range existingRefunds {
		if (refund.Status == model.RefundStatusPending || refund.Status == model.RefundStatusProcessing) &&
			refund.Amount == req.Body.Amount {
			return nil, huma.Error409Conflict("A refund for this amount is already being processed")
		}
	}

	// Create refund object
	refund := &model.Refund{
		ID:                uuid.New(),
		PaymentID:         req.Body.PaymentID,
		IdempotencyKey:    req.Headers.IdempotencyKey,
		Amount:            req.Body.Amount,
		Currency:          payment.Currency, // Use same currency as payment
		Status:            model.RefundStatusPending,
		Reason:            req.Body.Reason,
		ReasonDescription: req.Body.ReasonDescription,
		Metadata:          req.Body.Metadata,
		InitiatedBy:       initiatedBy,
		InitiatedByEmail:  initiatedByEmail,
		InitiatedAt:       time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Create the refund
	createdRefund, err := v1.payment.Payment.CreateRefund(ctx, refund)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create refund")
	}

	return &InitiateRefundResponse{
		Body: *createdRefund,
	}, nil
}

// GetRefundRequest represents the request to get a refund
type GetRefundRequest struct {
	RefundID uuid.UUID `path:"refund_id" required:"true" doc:"Refund ID"`
}

// GetRefundResponse represents the response for getting a refund
type GetRefundResponse struct {
	Body model.Refund
}

// GetRefund retrieves a refund by ID
func (v1 *V1) GetRefund(ctx context.Context, req *GetRefundRequest) (*GetRefundResponse, error) {
	refund, err := v1.payment.Payment.GetRefundByID(ctx, req.RefundID)
	if err != nil {
		return nil, huma.Error404NotFound("Refund not found")
	}

	return &GetRefundResponse{
		Body: *refund,
	}, nil
}

// ListRefundsRequest represents the request to list refunds
type ListRefundsRequest struct {
	Query struct {
		PaymentID    *uuid.UUID           `query:"payment_id" doc:"Filter by payment ID"`
		Status       *model.RefundStatus  `query:"status" doc:"Filter by refund status"`
		InitiatedBy  *uuid.UUID           `query:"initiated_by" doc:"Filter by user who initiated"`
		FromDate     *time.Time           `query:"from_date" doc:"Filter refunds created after this date"`
		ToDate       *time.Time           `query:"to_date" doc:"Filter refunds created before this date"`
		Limit        int                  `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Number of results to return"`
		Offset       int                  `query:"offset" default:"0" minimum:"0" doc:"Number of results to skip"`
	}
}

// ListRefundsResponse represents the response for listing refunds
type ListRefundsResponse struct {
	Body struct {
		Refunds []model.Refund `json:"refunds"`
		Total   int64          `json:"total"`
		Limit   int            `json:"limit"`
		Offset  int            `json:"offset"`
	}
}

// ListRefunds lists refunds with filtering options
func (v1 *V1) ListRefunds(ctx context.Context, req *ListRefundsRequest) (*ListRefundsResponse, error) {
	filter := &model.RefundFilter{
		PaymentID:   req.Query.PaymentID,
		Status:      req.Query.Status,
		InitiatedBy: req.Query.InitiatedBy,
		FromDate:    req.Query.FromDate,
		ToDate:      req.Query.ToDate,
		Limit:       req.Query.Limit,
		Offset:      req.Query.Offset,
	}

	refunds, total, err := v1.payment.Payment.ListRefunds(ctx, filter)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list refunds")
	}

	return &ListRefundsResponse{
		Body: struct {
			Refunds []model.Refund `json:"refunds"`
			Total   int64          `json:"total"`
			Limit   int            `json:"limit"`
			Offset  int            `json:"offset"`
		}{
			Refunds: refunds,
			Total:   total,
			Limit:   req.Query.Limit,
			Offset:  req.Query.Offset,
		},
	}, nil
}

// ListMerchantRefundsRequest represents the request to list refunds for a specific merchant
type ListMerchantRefundsRequest struct {
	MerchantID uuid.UUID `path:"merchant_id" required:"true" doc:"Merchant ID"`
	Query struct {
		Status       *model.RefundStatus  `query:"status" doc:"Filter by refund status"`
		FromDate     *time.Time           `query:"from_date" doc:"Filter refunds created after this date"`
		ToDate       *time.Time           `query:"to_date" doc:"Filter refunds created before this date"`
		Limit        int                  `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Number of results to return"`
		Offset       int                  `query:"offset" default:"0" minimum:"0" doc:"Number of results to skip"`
	}
}

// ListMerchantRefundsResponse represents the response for listing merchant refunds
type ListMerchantRefundsResponse struct {
	Body struct {
		Refunds []model.Refund `json:"refunds"`
		Total   int64          `json:"total"`
		Limit   int            `json:"limit"`
		Offset  int            `json:"offset"`
	}
}

// ListMerchantRefunds lists all refunds associated with a particular merchant
func (v1 *V1) ListMerchantRefunds(ctx context.Context, req *ListMerchantRefundsRequest) (*ListMerchantRefundsResponse, error) {
	filter := &model.RefundFilter{
		MerchantID: &req.MerchantID,
		Status:     req.Query.Status,
		FromDate:   req.Query.FromDate,
		ToDate:     req.Query.ToDate,
		Limit:      req.Query.Limit,
		Offset:     req.Query.Offset,
	}

	refunds, total, err := v1.payment.Payment.ListRefunds(ctx, filter)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list merchant refunds")
	}

	return &ListMerchantRefundsResponse{
		Body: struct {
			Refunds []model.Refund `json:"refunds"`
			Total   int64          `json:"total"`
			Limit   int            `json:"limit"`
			Offset  int            `json:"offset"`
		}{
			Refunds: refunds,
			Total:   total,
			Limit:   req.Query.Limit,
			Offset:  req.Query.Offset,
		},
	}, nil
}

// CancelRefundRequest represents the request to cancel a refund
type CancelRefundRequest struct {
	RefundID uuid.UUID `path:"refund_id" required:"true" doc:"Refund ID"`
	Body struct {
		Reason *string `json:"reason,omitempty" maxLength:"500" doc:"Reason for cancellation"`
	}
}

// CancelRefundResponse represents the response for cancelling a refund
type CancelRefundResponse struct {
	Body model.Refund
}

// CancelRefund cancels a pending refund
func (v1 *V1) CancelRefund(ctx context.Context, req *CancelRefundRequest) (*CancelRefundResponse, error) {
	// Get the refund
	refund, err := v1.payment.Payment.GetRefundByID(ctx, req.RefundID)
	if err != nil {
		return nil, huma.Error404NotFound("Refund not found")
	}

	// Validate refund can be cancelled
	if refund.Status != model.RefundStatusPending && refund.Status != model.RefundStatusProcessing {
		return nil, huma.Error400BadRequest("Only pending or processing refunds can be cancelled")
	}

	// Cancel the refund
	cancelledRefund, err := v1.payment.Payment.CancelRefund(ctx, req.RefundID, req.Body.Reason)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to cancel refund")
	}

	return &CancelRefundResponse{
		Body: *cancelledRefund,
	}, nil
}

// UpdateRefundStatusRequest represents the request to update refund status (for async processing)
type UpdateRefundStatusRequest struct {
	RefundID uuid.UUID `path:"refund_id" required:"true" doc:"Refund ID"`
	Body struct {
		Status               model.RefundStatus     `json:"status" required:"true" doc:"New refund status"`
		ExternalRefundID     *string                `json:"external_refund_id,omitempty" doc:"External refund ID from provider"`
		ProviderResponse     map[string]any         `json:"provider_response,omitempty" doc:"Provider response data"`
		ProviderErrorCode    *string                `json:"provider_error_code,omitempty" doc:"Provider error code"`
		ProviderErrorMessage *string                `json:"provider_error_message,omitempty" doc:"Provider error message"`
	}
}

// UpdateRefundStatusResponse represents the response for updating refund status
type UpdateRefundStatusResponse struct {
	Body model.Refund
}

// UpdateRefundStatus updates the status of a refund (webhook endpoint for async processing)
func (v1 *V1) UpdateRefundStatus(ctx context.Context, req *UpdateRefundStatusRequest) (*UpdateRefundStatusResponse, error) {
	// Update refund status
	updateData := &model.RefundStatusUpdate{
		Status:               req.Body.Status,
		ExternalRefundID:     req.Body.ExternalRefundID,
		ProviderResponse:     req.Body.ProviderResponse,
		ProviderErrorCode:    req.Body.ProviderErrorCode,
		ProviderErrorMessage: req.Body.ProviderErrorMessage,
	}

	refund, err := v1.payment.Payment.UpdateRefundStatus(ctx, req.RefundID, updateData)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update refund status")
	}

	return &UpdateRefundStatusResponse{
		Body: *refund,
	}, nil
}
