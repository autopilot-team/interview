package middleware

import (
	"context"
	"net/http"
	"strings"

	"autopilot/backends/internal/types"
)

// WithOperationMode adds operation mode to the request context
func WithOperationMode() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("Authorization")
			apiKey = strings.TrimPrefix(apiKey, "Bearer ")

			mode := types.OperationModeTest // Default to test mode
			if strings.HasPrefix(apiKey, "sk_live_") || strings.HasPrefix(apiKey, "pk_live_") {
				mode = types.OperationModeLive
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, types.OperationModeKey, mode)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
