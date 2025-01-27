package validator

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
)

// TurnstileToken is a huma validation type that validates a Cloudflare Turnstile token.
type TurnstileToken string

const (
	ErrUnableToVerifyTurnstileToken = "unable to verify Turnstile token"
	ErrFailedToVerifyTurnstileToken = "failed to verify Turnstile token"
)

// Resolve implements the huma.ResolverWithPath interface.
func (t TurnstileToken) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	var errors []error

	ctxx := ctx.Context()
	container := ctxx.Value(middleware.ContainerKey).(*app.Container)

	// Verify Cloudflare Turnstile token
	ok, err := container.Turnstile.Verify(ctxx, string(t), "")
	if err != nil {
		container.Logger.Error("Unable to verify Turnstile token", "error", err)
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  ErrUnableToVerifyTurnstileToken,
		})
	}

	if !ok {
		container.Logger.Error("Failed to verify Turnstile token")
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  ErrFailedToVerifyTurnstileToken,
		})
	}

	return errors
}

// Ensure our resolver meets the expected interface
var _ huma.ResolverWithPath = (*TurnstileToken)(nil)
