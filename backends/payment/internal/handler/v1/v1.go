package v1

import (
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"
	"autopilot/backends/payment/internal/app"
	"autopilot/backends/payment/internal/service"
	"autopilot/backends/payment/internal/store"
)

// Handler implements the gRPC PaymentService server
type Handler struct {
	paymentv1.UnimplementedPaymentServiceServer
	*app.Container
	paymentService *service.Payment
}

// NewHandler creates a new Handler instance
func NewHandler(container *app.Container) *Handler {
	// Initialize repositories
	paymentStore := store.NewPayment(container)

	// Initialize payment service
	paymentService := service.NewPayment(
		paymentStore,
	)

	return &Handler{
		Container:      container,
		paymentService: paymentService,
	}
}
