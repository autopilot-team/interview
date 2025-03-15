package middleware

import (
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"context"
	"net/http"
)

// WithContainer is a middleware that adds the container to the request context.
func WithContainer(container *app.Container) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := AttachContainer(r.Context(), container)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AttachContainer attaches the container to the request context
func AttachContainer(ctx context.Context, container *app.Container) context.Context {
	ctx = context.WithValue(ctx, ContainerKey, container)
	ctx = context.WithValue(ctx, httpx.TurnstileKey, container.Turnstile)
	return ctx
}
