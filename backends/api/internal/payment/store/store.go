package store

import (
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/internal/types"
	"context"
)

// ModeStore is a collection of stores for a specific operation mode.
type ModeStore struct {
	Payment Paymenter
}

// Manager is a collection of stores used by the services.
type Manager struct {
	*app.Container
	Live *ModeStore
	Test *ModeStore
}

// NewManager creates a new Manager.
func NewManager(container *app.Container) *Manager {
	return &Manager{
		Container: container,
		Live: &ModeStore{
			Payment: NewPayment(container.DB.Payment.Live.Writer()),
		},
		Test: &ModeStore{
			Payment: NewPayment(container.DB.Payment.Test.Writer()),
		},
	}
}

func (m *Manager) WithMode(ctx context.Context) *ModeStore {
	mode := types.GetOperationMode(ctx)
	switch mode {
	case types.OperationModeLive:
		return m.Live
	case types.OperationModeTest:
		return m.Test
	}

	m.Logger.Error("invalid database connection mode", "mode", mode)

	return m.Test
}
