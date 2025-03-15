package service

import (
	"autopilot/backends/api/pkg/app"
	"time"

	"github.com/riverqueue/river"
)

// AddWorkers returns the background workers
func AddWorkers(container *app.Container, workers *river.Workers, serviceManager *Manager) {
	river.AddWorker(workers, &Mailer{Container: container, service: serviceManager})
	river.AddWorker(workers, &SessionCleaner{Container: container, service: serviceManager})
}

// AddPeriodicJobs returns the periodic jobs
func AddPeriodicJobs(container *app.Container, serviceManager *Manager) []*river.PeriodicJob {
	jobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour*12),
			func() (river.JobArgs, *river.InsertOpts) {
				return SessionCleanerArgs{}, nil
			},
			&river.PeriodicJobOpts{
				RunOnStart: false,
			},
		),
	}

	return jobs
}
