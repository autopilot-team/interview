package middleware

import (
	"context"
	"net/http"
)

// RequestMetadataKey is the context key for request metadata
const RequestMetadataKey contextKey = "request_metadata"

// RequestMetadata contains metadata about the HTTP request
type RequestMetadata struct {
	IPAddress string // Client's IP address from X-Forwarded-For or RemoteAddr
	UserAgent string // Client's User-Agent header
}

// WithRequestMetadata is a middleware that adds request metadata to the
// context. It extracts client IP from X-Forwarded-For header if available,
// falling back to RemoteAddr.
func WithRequestMetadata() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := AttachRequestMetadata(r.Context(), r.Header.Get("X-Forwarded-For"), r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AttachRequestMetadata attaches request metadata to the context
func AttachRequestMetadata(ctx context.Context, clientIP, remoteAddr, userAgent string) context.Context {
	// Extract client IP, preferring X-Forwarded-For header
	if clientIP == "" {
		clientIP = remoteAddr
	}

	metadata := &RequestMetadata{
		IPAddress: clientIP,
		UserAgent: userAgent,
	}

	return context.WithValue(ctx, RequestMetadataKey, metadata)
}

// GetRequestMetadata retrieves request metadata from the context
// Returns nil if metadata is not found in context
func GetRequestMetadata(ctx context.Context) *RequestMetadata {
	if ctx == nil {
		return nil
	}
	if metadata, ok := ctx.Value(RequestMetadataKey).(*RequestMetadata); ok {
		return metadata
	}

	return nil
}
