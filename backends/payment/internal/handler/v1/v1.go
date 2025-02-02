package v1

import (
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"
	"autopilot/backends/payment/internal/app"
	"autopilot/backends/payment/internal/model"
	"autopilot/backends/payment/internal/service"
)

// V1 implements the gRPC PaymentService server
type V1 struct {
	paymentv1.UnimplementedPaymentServiceServer
	*app.Container
	paymentService *service.Payment
}

// New creates a new V1 instance
func New(container *app.Container) *V1 {
	return &V1{
		Container:      container,
		paymentService: service.NewPayment(container),
	}
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
