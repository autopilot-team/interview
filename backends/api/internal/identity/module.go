package identity

import (
	"autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"context"
)

// Module is the main module for the identity service
type Module struct {
	Service *service.Manager
	Store   *store.Manager
}

// New creates a new identity module
func New(ctx context.Context, container *app.Container) (*Module, error) {
	// Initialize the store manager
	storeManager := store.NewManager(container.DB.Identity.Writer())

	// Initialize the service manager
	serviceManager := service.NewManager(container, storeManager)

	return &Module{
		Service: serviceManager,
		Store:   storeManager,
	}, nil
}
