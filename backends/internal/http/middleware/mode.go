package middleware

import (
	"autopilot/backends/internal/types"
	"context"
	"net/http"
	"net/url"
	"strings"
)

// WithOperationMode adds operation mode to the request context
func WithOperationMode(xModeAllowedUrls []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			mode := types.OperationModeTest // Default to test mode

			// Check if request is from specific allowed URLs
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")
			isXModeReq := isXModeRequest(origin, referer, xModeAllowedUrls)
			apiKey := r.Header.Get("X-Api-Key")

			if strings.HasPrefix(apiKey, "sk_live_") || strings.HasPrefix(apiKey, "pk_live_") {
				mode = types.OperationModeLive
			} else if isXModeReq {
				if xOperationMode := r.Header.Get("X-Operation-Mode"); xOperationMode != "" {
					if xOperationMode == string(types.OperationModeLive) {
						mode = types.OperationModeLive
					}
				}
			}

			ctx = context.WithValue(ctx, types.OperationModeKey, mode)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isXModeRequest checks if the request originated from an URL in the xModeAllowedUrls list
func isXModeRequest(origin, referer string, xModeAllowedUrls []string) bool {
	if len(xModeAllowedUrls) == 0 {
		return false
	}

	// Check origin first
	if origin != "" {
		originURL, err := url.Parse(origin)
		if err == nil {
			for _, domain := range xModeAllowedUrls {
				if originURL.Host != "" && strings.Contains(domain, originURL.Host) {
					return true
				}
			}
		}
	}

	// Fallback to referer if origin is not set
	if referer != "" {
		refererURL, err := url.Parse(referer)
		if err == nil {
			for _, domain := range xModeAllowedUrls {
				if refererURL.Host != "" && strings.Contains(domain, refererURL.Host) {
					return true
				}
			}
		}
	}

	return false
}
