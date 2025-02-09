package store

import (
	"autopilot/backends/internal/core"
)

// Manager is a collection of stores used by the services.
type Manager struct {
}

// NewManager creates a new Manager.
func NewManager(primaryDB core.DBer) *Manager {
	return &Manager{}
}
