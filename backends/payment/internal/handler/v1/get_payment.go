package v1

import (
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPayment implements the GetPayment RPC method
func (h *V1) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.GetPaymentResponse, error) {
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
