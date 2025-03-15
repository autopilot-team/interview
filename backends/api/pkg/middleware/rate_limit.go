package middleware

import (
	"autopilot/backends/api/pkg/app"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/httprate"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig holds rate limit configuration for different endpoints
type RateLimitConfig struct {
	// Rate limits for unauthenticated endpoints (sign-in, sign-up, etc)
	Public struct {
		Requests int
		Window   time.Duration
	}

	// Rate limits for authenticated endpoints (all operations requiring auth)
	Private struct {
		Requests int
		Window   time.Duration
	}
}

// DefaultRateLimitConfig returns the default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	cfg := RateLimitConfig{}

	// Public endpoints (unauthenticated) - Stricter limits to prevent abuse
	cfg.Public = struct {
		Requests int
		Window   time.Duration
	}{
		Requests: 30,              // 30 requests
		Window:   5 * time.Minute, // per 5 minutes
	}

	// Private endpoints (authenticated) - More generous limits for legitimate operations
	cfg.Private = struct {
		Requests int
		Window   time.Duration
	}{
		Requests: 300,         // 300 requests
		Window:   time.Minute, // per minute
	}

	return cfg
}

// ClusterRateLimit implements a rate limiting middleware that works with Valkey cluster
func ClusterRateLimit(container *app.Container, client redis.UniversalClient, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a key from IP and path
			ip := r.RemoteAddr
			if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" {
				ip = fwdFor
			}
			endpoint := r.URL.Path
			key := fmt.Sprintf("ratelimit:%s:%s", ip, endpoint)

			// Use Valkey for counting
			val, err := client.Incr(r.Context(), key).Result()
			if err != nil {
				// Fall back to allowing the request if Valkey fails
				next.ServeHTTP(w, r)
				return
			}

			// Set expiry on the key
			if val == 1 {
				client.Expire(r.Context(), key, window)
			}

			// Get TTL for the remaining time header
			ttl, err := client.TTL(r.Context(), key).Result()
			if err != nil {
				ttl = window
			}

			// If we've exceeded the limit, return 429
			if val > int64(limit) {
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(ttl.Seconds())))
				w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(int64(limit)-val)))
			w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(ttl.Seconds())))

			next.ServeHTTP(w, r)
		})
	}
}

// isPublicEndpoint determines if the given path is a public authentication endpoint
func isPublicEndpoint(path string) bool {
	return path == "/v1/identity/sign-in" ||
		path == "/v1/identity/sign-up" ||
		path == "/v1/identity/forgot-password" ||
		path == "/v1/identity/reset-password" ||
		path == "/v1/identity/verify-email" ||
		path == "/v1/identity/verify-two-factor"
}

// WithMemoryRateLimit creates a rate limiter using in-memory storage via httprate
func WithMemoryRateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var limiter func(http.Handler) http.Handler
			path := r.URL.Path

			// Apply appropriate rate limits based on endpoint type
			if isPublicEndpoint(path) {
				limiter = httprate.Limit(
					config.Public.Requests,
					config.Public.Window,
					httprate.WithKeyFuncs(httprate.KeyByIP),
					httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, "Rate limit exceeded for authentication operations", http.StatusTooManyRequests)
					}),
				)
			} else {
				// All other endpoints use private rate limits
				limiter = httprate.Limit(
					config.Private.Requests,
					config.Private.Window,
					httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
					httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, "Rate limit exceeded for API operations", http.StatusTooManyRequests)
					}),
				)
			}

			// Apply the rate limiter
			limiter(next).ServeHTTP(w, r)
		})
	}
}

// WithClusterRateLimit creates a rate limiter using Redis/Valkey for distributed rate limiting
func WithClusterRateLimit(container *app.Container, client redis.UniversalClient, config RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// Apply appropriate rate limits based on endpoint type
			if isPublicEndpoint(path) {
				ClusterRateLimit(
					container,
					client,
					config.Public.Requests,
					config.Public.Window,
				)(next).ServeHTTP(w, r)
			} else {
				// All other endpoints use private rate limits
				ClusterRateLimit(
					container,
					client,
					config.Private.Requests,
					config.Private.Window,
				)(next).ServeHTTP(w, r)
			}
		})
	}
}

// WithRateLimit adds rate limiting middleware to the handler
func WithRateLimit(container *app.Container, config RateLimitConfig) func(http.Handler) http.Handler {
	// Use RateLimiter Valkey client from container
	valkeyClient := container.RateLimiter

	// Choose appropriate rate limiting strategy based on available resources
	if valkeyClient == nil {
		// Log that we're falling back to memory-based rate limiting
		container.Logger.Warn("Valkey client unavailable for rate limiting, falling back to in-memory implementation")
		return WithMemoryRateLimit(config)
	}

	// Use our cluster-compatible implementation
	return WithClusterRateLimit(container, valkeyClient, config)
}
