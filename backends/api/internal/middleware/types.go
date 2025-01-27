// Package middleware provides HTTP middleware for the API.
package middleware

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// Context keys for middleware
const (
	// ContainerKey is the key used to store the container in the context
	ContainerKey contextKey = "container"

	// LocaleKey is the context key for the locale
	LocaleKey contextKey = "locale"

	// TranslatorKey is the context key for the translator
	TranslatorKey contextKey = "translator"
)
