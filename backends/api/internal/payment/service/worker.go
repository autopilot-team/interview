package service

import (
	"autopilot/backends/api/pkg/app"

	"github.com/riverqueue/river"
)

// AddWorkers returns the background workers
func AddWorkers(container *app.Container, workers *river.Workers, serviceManager *Manager) {
	river.AddWorker(workers, &Mailer{Container: container, service: serviceManager})
}

// AddPeriodicJobs returns the periodic jobs
func AddPeriodicJobs(container *app.Container, serviceManager *Manager) []*river.PeriodicJob {
	jobs := []*river.PeriodicJob{}

	return jobs
}
