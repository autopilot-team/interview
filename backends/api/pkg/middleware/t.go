package middleware

import (
	"autopilot/backends/internal/core"
	"context"
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	// DefaultLocale is the default locale to use when no locale is specified
	DefaultLocale = "en"

	// LocaleHeader is the HTTP header key for the locale
	LocaleHeader = "Accept-Language"

	// LocaleQueryParam is the query parameter key for the locale
	LocaleQueryParam = "locale"
)

// WithT is a middleware that extracts the locale from the request, initializes
// the i18n translator, and adds both locale and translator to the context
func WithT(i18nBundle *core.I18nBundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := AttachT(r.Context(), i18nBundle, r.Header.Get(LocaleHeader), r.URL.Query().Get(LocaleQueryParam))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AttachT attaches the locale and translator to the request context
func AttachT(ctx context.Context, i18nBundle *core.I18nBundle, acceptLanguageHeader, queryParam string) context.Context {
	var locale string

	// Try to get from query parameter
	locale = queryParam

	// Try to get from header if not in query
	if locale == "" {
		if acceptLanguageHeader != "" {
			// Parse the Accept-Language header which can contain multiple values with q-factors
			// e.g. "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5"
			locales := strings.Split(acceptLanguageHeader, ",")
			if len(locales) > 0 {
				// Take the first (highest priority) locale
				locale = strings.TrimSpace(strings.Split(locales[0], ";")[0])
			}
		}
	}

	// Normalize and set in context
	if locale != "" {
		locale = normalizeLocale(locale)
	} else {
		locale = DefaultLocale
	}

	t := i18n.NewLocalizer(i18nBundle.Bundle, locale)
	ctx = context.WithValue(ctx, LocaleKey, locale)
	ctx = context.WithValue(ctx, TranslatorKey, t)

	return ctx
}

// GetLocale extracts the locale from the context with proper fallback logic.
func GetLocale(ctx context.Context) string {
	if locale, ok := ctx.Value(LocaleKey).(string); ok && locale != "" {
		return locale
	}

	return DefaultLocale
}

// GetT extracts the translator from the context with proper fallback logic.
func GetT(ctx context.Context) *i18n.Localizer {
	if t, ok := ctx.Value(TranslatorKey).(*i18n.Localizer); ok && t != nil {
		return t
	}

	return nil
}

// normalizeLocale normalizes a locale string to a standard format.
// e.g. "en-US" -> "en", "zh-CN" -> "zh-CN"
func normalizeLocale(locale string) string {
	// Convert to lowercase
	locale = strings.ToLower(locale)

	// Special cases for locales that need the region code
	if strings.HasPrefix(locale, "zh-") {
		return locale
	}

	// For other locales, just take the language code
	return strings.Split(locale, "-")[0]
}
