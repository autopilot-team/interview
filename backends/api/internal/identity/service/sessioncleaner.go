package service

import (
	"autopilot/backends/api/pkg/app"
	"context"
	"fmt"

	"github.com/riverqueue/river"
)

// SessionCleanerArgs is the arguments for the session cleaner
type SessionCleanerArgs struct{}

// Kind returns the kind of the worker
func (SessionCleanerArgs) Kind() string {
	return "session_cleaner"
}

// SessionCleaner is a worker that cleans up expired sessions periodically
type SessionCleaner struct {
	*app.Container
	service *Manager
	river.WorkerDefaults[SessionCleanerArgs]
}

// Work is the worker function that cleans up expired sessions
func (s *SessionCleaner) Work(ctx context.Context, job *river.Job[SessionCleanerArgs]) error {
	s.Logger.Info("Starting expired sessions cleanup")

	if err := s.service.Session.CleanUpExpired(ctx); err != nil {
		s.Logger.Error("Failed to clean up expired sessions", "error", err)
		return fmt.Errorf("cleaning up expired sessions: %w", err)
	}

	s.Logger.Info("Successfully cleaned up expired sessions")
	return nil
}
