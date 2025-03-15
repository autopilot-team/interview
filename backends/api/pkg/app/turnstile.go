package app

import (
	"context"

	"github.com/9ssi7/turnstile"
)

// Turnstiler is an interface that wraps the Verify method
type Turnstiler interface {
	// Verify verifies a Turnstile token
	Verify(ctx context.Context, token string, action string) (bool, error)
}

// NewTurnstile creates a new Turnstiler
func NewTurnstile(secret string) turnstile.Service {
	return turnstile.New(turnstile.Config{
		Secret: secret,
	})
}
