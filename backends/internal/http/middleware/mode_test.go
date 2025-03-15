package middleware

import (
	"autopilot/backends/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOperationMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		xModeAllowedUrls []string
		headers          map[string]string
		expectedMode     types.OperationMode
	}{
		{
			name:             "should default to test mode",
			xModeAllowedUrls: []string{},
			headers:          map[string]string{},
			expectedMode:     types.OperationModeTest,
		},
		{
			name:             "should use live mode for dashboard request with live mode header",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			headers: map[string]string{
				"Origin":           "https://dashboard.autopilot.com",
				"X-Operation-Mode": string(types.OperationModeLive),
			},
			expectedMode: types.OperationModeLive,
		},
		{
			name:             "should use test mode for dashboard request without X-Operation-Mode",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			headers: map[string]string{
				"Origin": "https://dashboard.autopilot.com",
			},
			expectedMode: types.OperationModeTest,
		},
		{
			name:             "should use test mode for dashboard request with empty X-Operation-Mode",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			headers: map[string]string{
				"Origin":           "https://dashboard.autopilot.com",
				"X-Operation-Mode": "",
			},
			expectedMode: types.OperationModeTest,
		},
		{
			name:             "should use live mode for API request with live secret key",
			xModeAllowedUrls: []string{},
			headers: map[string]string{
				"X-API-Key": "sk_live_123",
			},
			expectedMode: types.OperationModeLive,
		},
		{
			name:             "should use live mode for API request with live publishable key",
			xModeAllowedUrls: []string{},
			headers: map[string]string{
				"X-API-Key": "pk_live_123",
			},
			expectedMode: types.OperationModeLive,
		},
		{
			name:             "should use test mode for API request with test secret key",
			xModeAllowedUrls: []string{},
			headers: map[string]string{
				"X-API-Key": "sk_test_123",
			},
			expectedMode: types.OperationModeTest,
		},
		{
			name:             "should use test mode for API request with test publishable key",
			xModeAllowedUrls: []string{},
			headers: map[string]string{
				"X-API-Key": "pk_test_123",
			},
			expectedMode: types.OperationModeTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedMode types.OperationMode

			// Create a test handler that captures the mode
			handler := WithOperationMode(tt.xModeAllowedUrls)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedMode = r.Context().Value(types.OperationModeKey).(types.OperationMode)
			}))

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// Serve the request
			handler.ServeHTTP(httptest.NewRecorder(), req)

			assert.Equal(t, tt.expectedMode, capturedMode)
		})
	}
}

func TestIsXModeRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		origin           string
		referer          string
		xModeAllowedUrls []string
		expected         bool
	}{
		{
			name:             "should reject empty allowed URLs",
			origin:           "https://dashboard.autopilot.com",
			referer:          "",
			xModeAllowedUrls: []string{},
			expected:         false,
		},
		{
			name:             "should accept matching origin",
			origin:           "https://dashboard.autopilot.com",
			referer:          "",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         true,
		},
		{
			name:             "should reject non-matching origin",
			origin:           "https://other.com",
			referer:          "",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         false,
		},
		{
			name:             "should accept matching referer when origin empty",
			origin:           "",
			referer:          "https://dashboard.autopilot.com/settings",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         true,
		},
		{
			name:             "should reject non-matching referer when origin empty",
			origin:           "",
			referer:          "https://other.com/settings",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         false,
		},
		{
			name:             "should reject invalid origin URL",
			origin:           "not-a-url",
			referer:          "",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         false,
		},
		{
			name:             "should reject invalid referer URL",
			origin:           "",
			referer:          "not-a-url",
			xModeAllowedUrls: []string{"https://dashboard.autopilot.com"},
			expected:         false,
		},
		{
			name:             "should accept matching origin from multiple allowed domains",
			origin:           "https://dashboard.autopilot.com",
			referer:          "",
			xModeAllowedUrls: []string{"https://api.autopilot.com", "https://dashboard.autopilot.com"},
			expected:         true,
		},
		{
			name:             "should reject non-matching origin from multiple allowed domains",
			origin:           "https://other.com",
			referer:          "",
			xModeAllowedUrls: []string{"https://api.autopilot.com", "https://dashboard.autopilot.com"},
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isXModeRequest(tt.origin, tt.referer, tt.xModeAllowedUrls)
			assert.Equal(t, tt.expected, result)
		})
	}
}
