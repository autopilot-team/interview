package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithRequestMetadata(t *testing.T) {
	tests := []struct {
		name          string
		remoteAddr    string
		xForwardedFor string
		userAgent     string
		wantIP        string
		wantUserAgent string
	}{
		{
			name:       "use RemoteAddr when X-Forwarded-For is empty",
			remoteAddr: "192.168.1.1:1234",
			wantIP:     "192.168.1.1:1234",
		},
		{
			name:          "prefer X-Forwarded-For over RemoteAddr",
			remoteAddr:    "192.168.1.1:1234",
			xForwardedFor: "10.0.0.1",
			wantIP:        "10.0.0.1",
		},
		{
			name:          "capture User-Agent header",
			remoteAddr:    "192.168.1.1:1234",
			userAgent:     "Mozilla/5.0 Test Browser",
			wantIP:        "192.168.1.1:1234",
			wantUserAgent: "Mozilla/5.0 Test Browser",
		},
		{
			name:          "handle all headers",
			remoteAddr:    "192.168.1.1:1234",
			xForwardedFor: "10.0.0.1",
			userAgent:     "Mozilla/5.0 Test Browser",
			wantIP:        "10.0.0.1",
			wantUserAgent: "Mozilla/5.0 Test Browser",
		},
		{
			name:       "handle empty headers",
			remoteAddr: "192.168.1.1:1234",
			wantIP:     "192.168.1.1:1234",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that will verify the context
			handler := WithRequestMetadata()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				metadata := GetRequestMetadata(r.Context())
				assert.NotNil(t, metadata)
				assert.Equal(t, tc.wantIP, metadata.IPAddress)
				assert.Equal(t, tc.wantUserAgent, metadata.UserAgent)
			}))

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tc.xForwardedFor)
			}
			if tc.userAgent != "" {
				req.Header.Set("User-Agent", tc.userAgent)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Verify response status
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

func TestGetRequestMetadata(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want *RequestMetadata
	}{
		{
			name: "nil context returns nil",
			ctx:  nil,
			want: nil,
		},
		{
			name: "context without metadata returns nil",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "context with metadata returns metadata",
			ctx: context.WithValue(
				context.Background(),
				RequestMetadataKey,
				&RequestMetadata{IPAddress: "127.0.0.1", UserAgent: "test-agent"},
			),
			want: &RequestMetadata{IPAddress: "127.0.0.1", UserAgent: "test-agent"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GetRequestMetadata(tc.ctx)
			assert.Equal(t, tc.want, got)
		})
	}
}
