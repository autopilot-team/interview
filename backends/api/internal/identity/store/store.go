package store

import (
	"autopilot/backends/internal/core"
)

// Manager is a collection of stores used by the services.
type Manager struct {
	AuditLog     AuditLoger
	Entity       Entityer
	Membership   Membershiper
	Session      Sessioner
	TwoFactor    TwoFactorer
	User         Userer
	Verification Verificationer
}

// NewManager creates a new Manager.
func NewManager(q core.Querier) *Manager {
	return &Manager{
		AuditLog:     NewAuditLog(q),
		Entity:       NewEntity(q),
		Membership:   NewMembership(q),
		Session:      NewSession(q),
		TwoFactor:    NewTwoFactor(q),
		User:         NewUser(q),
		Verification: NewVerification(q),
	}
}
