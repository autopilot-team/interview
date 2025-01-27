package middleware

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(mode types.Mode, logger *core.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			// Define important headers to log
			importantHeaders := []string{
				"Accept",
				"Accept-Language",
				"Authorization",
				"Content-Type",
				"Referer",
				"User-Agent",
				"X-Forwarded-For",
				"X-Real-IP",
				"X-Request-ID",
			}

			headers := core.HttpHeaders{}
			for _, name := range importantHeaders {
				if value := r.Header.Get(name); value != "" {
					headers[name] = value
				}
			}

			defer func() {
				if mode == types.DebugMode {
					duration := time.Since(start)
					status := ww.Status()

					logger.InfoContext(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL),
						"method", r.Method,
						"path", r.URL.Path,
						"status", status,
						"latency", duration,
						"ip", r.RemoteAddr,
						"headers", headers,
					)
				}
			}()

			next.ServeHTTP(ww, r)
		})

		return otelhttp.NewHandler(
			handler,
			"http_request",
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				route := chi.RouteContext(r.Context()).RoutePattern()
				if route == "" {
					route = r.URL.Path
				}

				return fmt.Sprintf("%s %s", r.Method, route)
			}),
		)
	}
}
