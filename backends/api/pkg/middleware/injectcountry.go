package middleware

import (
	"net/http"
	"net/textproto"
)

type InjectCountryConfig struct {
	// Enable indicates whether to inject the country.
	Enable bool
	// Always indicates the header should be overwritten even if already present.
	Always bool
	// Country to set the header value to.
	Country string
}

// WithInjectCountry overrides the request IP country with the specific value.
// if override, the header will always be replaced.
func WithInjectCountry(cfg InjectCountryConfig) func(http.Handler) http.Handler {
	if !cfg.Enable {
		return nil
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := textproto.CanonicalMIMEHeaderKey(CFCountryHeader)
			if _, ok := r.Header[key]; !ok || cfg.Always {
				r.Header.Set(CFCountryHeader, cfg.Country)
			}
			next.ServeHTTP(w, r)
		})
	}
}
