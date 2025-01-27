package v1

import (
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"
	"autopilot/backends/payment/internal/model"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPayment implements the GetPayment RPC method
func (h *Handler) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.GetPaymentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get payment using the payment service
	payment, err := h.paymentService.GetPayment(ctx, req.Id)
	if err != nil {
		h.Logger.ErrorContext(ctx, "failed to get payment", "error", err)
		return nil, status.Error(codes.Internal, "failed to get payment")
	}

	// Convert domain model to protobuf
	return &paymentv1.GetPaymentResponse{
		Payment: &paymentv1.Payment{
			Id:          payment.ID.String(),
			UserId:      payment.MerchantID.String(),
			Amount:      payment.Amount,
			Currency:    payment.Currency,
			Status:      convertPaymentStatus(payment.Status),
			Description: payment.Description,
			CreatedAt:   payment.CreatedAt.Unix(),
			UpdatedAt:   payment.UpdatedAt.Unix(),
		},
	}, nil
}

// convertPaymentStatus converts domain payment status to protobuf payment status
func convertPaymentStatus(status model.PaymentStatus) paymentv1.PaymentStatus {
	switch status {
	case model.PaymentStatusPending:
		return paymentv1.PaymentStatus_PAYMENT_STATUS_PENDING
	case model.PaymentStatusProcessing:
		return paymentv1.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case model.PaymentStatusSucceeded:
		return paymentv1.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case model.PaymentStatusFailed:
		return paymentv1.PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		return paymentv1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}
