package middleware

import (
	"context"
	"net/http"
)

const ActiveEntityHeader = "X-Entity-Id"

// WithEntity is a middleware that adds the active entity from the header to the request context.
func WithActiveEntity() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if entity := r.Header.Get(ActiveEntityHeader); entity != "" {
				ctx = context.WithValue(ctx, EntityKey, entity)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetActiveEntity retrieves request request active entity from the context.
func GetActiveEntity(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	entity, _ := ctx.Value(EntityKey).(string)
	return entity
}
