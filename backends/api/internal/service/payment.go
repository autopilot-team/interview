package service

import (
	"context"
	"fmt"

	"autopilot/backends/api/internal/app"
	"autopilot/backends/internal/grpc/middleware"
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Paymenter is an interface that wraps the Payment methods
type Paymenter interface {
	GetPayment(ctx context.Context, payment *paymentv1.GetPaymentRequest) (*paymentv1.Payment, error)
}

type Payment struct {
	*app.Container
	v1Client paymentv1.PaymentServiceClient
}

// NewPayment creates a new PaymentService instance
func NewPayment(container *app.Container) (*Payment, error) {
	conn, err := grpc.NewClient(
		container.Config.Services.PaymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(middleware.UnaryOperationModeClientInterceptor),
	)
	if err != nil {
		return nil, err
	}

	return &Payment{
		Container: container,
		v1Client:  paymentv1.NewPaymentServiceClient(conn),
	}, nil
}

// GetPayment retrieves payment details using v1 API
func (s *Payment) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.Payment, error) {
	resp, err := s.v1Client.GetPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return resp.Payment, nil
}
