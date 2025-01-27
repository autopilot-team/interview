package worker

import (
	"autopilot/backends/payment/internal/app"
	"context"

	"github.com/riverqueue/river"
)

type MailerArgs struct {
	Email string
}

func (MailerArgs) Kind() string {
	return "mailer"
}

type Mailer struct {
	*app.Container
	river.WorkerDefaults[MailerArgs]
}

func (w *Mailer) Work(ctx context.Context, job *river.Job[MailerArgs]) error {
	return nil
}
