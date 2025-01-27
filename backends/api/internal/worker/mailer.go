package worker

import (
	"autopilot/backends/api/internal/app"
	"context"

	"autopilot/backends/internal/core"

	"github.com/riverqueue/river"
)

// MailerArgs is the arguments for the mailer worker
type MailerArgs struct {
	Data     map[string]interface{}
	Email    string
	Locale   string
	Subject  string
	Template string
}

// Kind returns the kind of the worker
func (MailerArgs) Kind() string {
	return "mailer"
}

// Mailer is a worker that sends emails
type Mailer struct {
	*app.Container
	river.WorkerDefaults[MailerArgs]
}

// Work is the worker function that sends emails
func (w *Mailer) Work(ctx context.Context, job *river.Job[MailerArgs]) error {
	msg := core.EmailMessage{
		From:    w.Config.App.Support.Email,
		To:      []string{job.Args.Email},
		Data:    job.Args.Data,
		Subject: job.Args.Subject,
	}

	return w.Mailer.Send(job.Args.Template, msg, &core.RenderOptions{
		Locale: job.Args.Locale,
	})
}
