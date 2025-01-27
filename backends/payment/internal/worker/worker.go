package worker

import (
	"autopilot/backends/payment/internal/app"

	"github.com/riverqueue/river"
)

// Register returns the background workers
func Register(container *app.Container) *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &Mailer{Container: container})

	return workers
}
