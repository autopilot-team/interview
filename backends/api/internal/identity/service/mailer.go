package service

import (
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/internal/core"
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
	return "identity.mailer"
}

// Mailer is a background worker that sends emails
type Mailer struct {
	*app.Container
	river.WorkerDefaults[MailerArgs]
	service *Manager
}

// Work is the worker function that sends emails
func (w *Mailer) Work(ctx context.Context, job *river.Job[MailerArgs]) error {
	msg := core.EmailMessage{
		From:    w.Config.App.Support.Email,
		To:      []string{job.Args.Email},
		Data:    job.Args.Data,
		Subject: job.Args.Subject,
	}

	err := w.Mailer.Send(job.Args.Template, msg, &core.RenderOptions{
		Locale: job.Args.Locale,
	})
	if err != nil {
		return err
	}

	return nil
}
