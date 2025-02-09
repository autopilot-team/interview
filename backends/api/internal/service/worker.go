package service

import (
	"autopilot/backends/api/internal/app"

	"github.com/riverqueue/river"
)

// AddWorkers returns the background workers
func AddWorkers(container *app.Container, serviceManager *Manager) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &Mailer{Container: container, service: serviceManager})

	return workers
}
