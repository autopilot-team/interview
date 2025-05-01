package store

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
)

// ModeStore is a collection of stores for a specific operation mode.
type ModeStore struct {
	Payment Paymenter
}

// Manager is a collection of stores used by the services.
type Manager struct {
	Live *ModeStore
	Test *ModeStore
}

// NewManager creates a new Manager.
func NewManager(live, test core.Querier) *Manager {
	return &Manager{
		Live: &ModeStore{
			Payment: NewPayment(live),
		},
		Test: &ModeStore{
			Payment: NewPayment(test),
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
	return m.Test
}
