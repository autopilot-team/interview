package internal

import (
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/payment"
)

type Module struct {
	Identity *identity.Module
	Payment  *payment.Module
}
