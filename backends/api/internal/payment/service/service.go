package service

import (
	"autopilot/backends/api/internal/payment/store"
	"autopilot/backends/api/pkg/app"
)

// Manager is a collection of services used by the handlers/workers.
type Manager struct {
	Payment Paymenter
}

// New creates a new service manager
func NewManager(container *app.Container, store *store.Manager) *Manager {
	return &Manager{
		Payment: NewPayment(container, store),
	}
}
