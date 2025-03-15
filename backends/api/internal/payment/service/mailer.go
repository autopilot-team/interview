package service

import (
	"autopilot/backends/api/pkg/app"
	"context"

	"github.com/riverqueue/river"
)

// MailerArgs is the arguments for the mailer worker
type MailerArgs struct {
	Data     map[string]any
	Email    string
	Locale   string
	Subject  string
	Template string
}

// Kind returns the kind of the worker
func (MailerArgs) Kind() string {
	return "payment.mailer"
}

// Mailer is a background worker that sends emails
type Mailer struct {
	*app.Container
	river.WorkerDefaults[MailerArgs]
	service *Manager
}

// Work is the worker function that sends emails
func (w *Mailer) Work(ctx context.Context, job *river.Job[MailerArgs]) error {
	return nil
}
