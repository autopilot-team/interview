package middleware

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithDebug(t *testing.T) {
	tests := []struct {
		name         string
		mode         types.Mode
		headerValue  string
		expectedFlag bool
	}{
		{
			name:         "should return false when header is not present",
			mode:         types.DebugMode,
			headerValue:  "",
			expectedFlag: false,
		},
		{
			name:         "should return true when header value is 'true'",
			mode:         types.DebugMode,
			headerValue:  "true",
			expectedFlag: true,
		},
		{
			name:         "should be case insensitive and return true for 'TRUE'",
			mode:         types.DebugMode,
			headerValue:  "TRUE",
			expectedFlag: true,
		},
		{
			name:         "should return true when header value is '1'",
			mode:         types.DebugMode,
			headerValue:  "1",
			expectedFlag: true,
		},
		{
			name:         "should return false when header value is 'false'",
			mode:         types.DebugMode,
			headerValue:  "false",
			expectedFlag: false,
		},
		{
			name:         "should return false when header value is invalid",
			mode:         types.DebugMode,
			headerValue:  "invalid",
			expectedFlag: false,
		},
		{
			name:         "should not set debug flag when mode is not DebugMode",
			mode:         types.ReleaseMode,
			headerValue:  "true",
			expectedFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flagFromContext bool

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				flagFromContext = core.GetDebugContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			middleware := WithDebug(tt.mode)
			middlewareHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Debug", tt.headerValue)
			}

			rr := httptest.NewRecorder()
			middlewareHandler.ServeHTTP(rr, req)

			if flagFromContext != tt.expectedFlag {
				t.Errorf("Expected debug flag to be %v, got %v", tt.expectedFlag, flagFromContext)
			}
		})
	}
}

func TestGetDebugContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "should return true when context has debug=true",
			ctx:      context.WithValue(context.Background(), core.DebugContextKey, true),
			expected: true,
		},
		{
			name:     "should return false when context has debug=false",
			ctx:      context.WithValue(context.Background(), core.DebugContextKey, false),
			expected: false,
		},
		{
			name:     "should return false when context does not have debug value",
			ctx:      context.Background(),
			expected: false,
		},
		{
			name:     "should return false when context has invalid debug type",
			ctx:      context.WithValue(context.Background(), core.DebugContextKey, "not-a-bool"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.GetDebugContext(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
