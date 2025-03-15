package payment

import (
	"autopilot/backends/api/internal/payment/service"
	"autopilot/backends/api/internal/payment/store"
	"autopilot/backends/api/pkg/app"
	"context"
)

// Module is the main module for the payment service
type Module struct {
	Service *service.Manager
	Store   *store.Manager
}

// New creates a new payment module
func New(ctx context.Context, container *app.Container) (*Module, error) {
	// Initialize the store manager
	storeManager := store.NewManager(container)

	// Initialize the service manager
	serviceManager := service.NewManager(container, storeManager)

	return &Module{
		Service: serviceManager,
		Store:   storeManager,
	}, nil
}
