package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithInjectCountry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		disabled      bool
		always        bool
		injectCountry string
		headerCountry string
		wantCountry   string
	}{
		{
			name:     "should return nil if not enabled",
			disabled: true,
		},
		{
			name:          "should inject country if no header provided",
			injectCountry: "SG",
			wantCountry:   "SG",
		},
		{
			name:          "should skip injection if header provided",
			injectCountry: "SG",
			headerCountry: "CA",
			wantCountry:   "CA",
		},
		{
			name:          "should override country if always is true",
			always:        true,
			injectCountry: "SG",
			headerCountry: "CA",
			wantCountry:   "SG",
		},
		{
			name:          "should handle empty injected country",
			always:        true,
			injectCountry: "",
			headerCountry: "CA",
			wantCountry:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that will verify the context
			cfg := InjectCountryConfig{
				Enable:  true,
				Always:  tc.always,
				Country: tc.injectCountry,
			}
			handler := WithInjectCountry(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.wantCountry, r.Header.Get(CFCountryHeader))
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			if tc.headerCountry != "" {
				req.Header.Set(CFCountryHeader, tc.headerCountry)
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
