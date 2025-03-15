package httpx

import (
	"autopilot/backends/internal/types"
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type AuthInfo struct {
	Authenticated bool
	EntityID      string
	UserID        string
	Mode          types.OperationMode
	APIKeyUsed    bool

	EntityRole types.Role
}

type Authenticator interface {
	// RequireUserSession enforces authentication using a session cookie.
	RequireUserSession(ctx huma.Context, next func(huma.Context))

	// RequireSecretKey enforces authentication using a secret API key.
	RequireSecretKey(ctx huma.Context, next func(huma.Context))

	// RequirePublishableKey enforces authentication using a publishable API key.
	RequirePublishableKey(ctx huma.Context, next func(huma.Context))

	// RequireAuthenticated allows authentication via either a session cookie or a secret API key.
	RequireAuthenticated(ctx huma.Context, next func(huma.Context))
}

func WithAuthInfo(ctx huma.Context, info AuthInfo) huma.Context {
	return huma.WithValue(ctx, types.AuthKey, info)
}

func GetAuthInfo(ctx context.Context) AuthInfo {
	auth, _ := ctx.Value(types.AuthKey).(AuthInfo)
	return auth
}

// Turnstiler is an interface that wraps the Verify method
type Turnstiler interface {
	// Verify verifies a Turnstile token
	Verify(ctx context.Context, token string, action string) (bool, error)
}

type contextKey string

const (
	// ContainerKey is the key used to store the container in the context.
	TurnstileKey contextKey = "turnstile"
)
