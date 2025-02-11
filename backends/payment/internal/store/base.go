package store

import (
	"context"

	"autopilot/backends/internal/types"
	"autopilot/backends/payment/internal/app"
)

// BaseStore provides common store functionality
type BaseStore struct {
	container *app.Container
}

// Infra returns the correct infrastructure based on operation mode
func (s *BaseStore) Infra(ctx context.Context) *app.ContainerInfra {
	if types.GetOperationMode(ctx) == types.OperationModeLive {
		return s.container.Live
	}

	return s.container.Test
}
