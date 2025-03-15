package middleware

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"net/http"
	"strings"
)

// WithDebug adds debug flag to the request context based on x-debug header
func WithDebug(mode types.Mode) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if mode == types.DebugMode {
				debug := false
				xDebug := r.Header.Get("X-Debug")
				if xDebug != "" {
					debug = strings.EqualFold(xDebug, "true") || strings.EqualFold(xDebug, "1")
				}

				ctx = context.WithValue(ctx, core.DebugContextKey, debug)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
