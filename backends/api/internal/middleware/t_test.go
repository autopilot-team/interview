package middleware

import (
	"autopilot/backends/internal/core"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestWithT(t *testing.T) {
	tests := []struct {
		name           string
		queryLocale    string
		acceptLanguage string
		want           string
	}{
		{
			name: "default locale when no locale specified",
			want: "{\"locale\":\"en\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "locale from query parameter",
			queryLocale: "fr",
			want:        "{\"locale\":\"fr\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:           "locale from Accept-Language header",
			acceptLanguage: "fr-FR,fr;q=0.9,en;q=0.8",
			want:           "{\"locale\":\"fr\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:           "query parameter takes precedence over header",
			queryLocale:    "es",
			acceptLanguage: "fr-FR,fr;q=0.9",
			want:           "{\"locale\":\"es\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:           "only first locale from Accept-Language is used",
			acceptLanguage: "es-ES,fr;q=0.9,en;q=0.8",
			want:           "{\"locale\":\"es\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "normalize en-US to en",
			queryLocale: "en-US",
			want:        "{\"locale\":\"en\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "normalize fr-FR to fr",
			queryLocale: "fr-FR",
			want:        "{\"locale\":\"fr\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "zh-CN uses Simplified Chinese",
			queryLocale: "zh-CN",
			want:        "{\"locale\":\"zh-cn\",\"translated_string\":\"验证电子邮箱\"}",
		},
		{
			name:        "zh-SG falls back to zh-CN",
			queryLocale: "zh-SG",
			want:        "{\"locale\":\"zh-sg\",\"translated_string\":\"验证电子邮箱\"}",
		},
		{
			name:        "zh-HK falls back to zh-TW",
			queryLocale: "zh-HK",
			want:        "{\"locale\":\"zh-hk\",\"translated_string\":\"驗證電子郵箱\"}",
		},
		{
			name:        "zh-MO falls back to zh-TW",
			queryLocale: "zh-MO",
			want:        "{\"locale\":\"zh-mo\",\"translated_string\":\"驗證電子郵箱\"}",
		},
		{
			name:        "normalize mixed case",
			queryLocale: "En-Us",
			want:        "{\"locale\":\"en\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "invalid locale falls back to language code",
			queryLocale: "fr-INVALID",
			want:        "{\"locale\":\"fr\",\"translated_string\":\"Verify Email Address\"}",
		},
		{
			name:        "empty locale uses default",
			queryLocale: "",
			want:        "{\"locale\":\"en\",\"translated_string\":\"Verify Email Address\"}",
		},
	}

	// Initialize the local filesystem
	localFS, err := core.NewLocalFS("./backends/api")
	assert.NoError(t, err)

	// Initialize the i18n bundle
	i18nBundle, err := core.NewI18nBundle(localFS, "locales")
	assert.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			basePath := "/test"
			requestPath := basePath
			if tc.queryLocale != "" {
				requestPath = fmt.Sprintf("%s?locale=%s", basePath, tc.queryLocale)
			}

			_, api := humatest.New(t)
			api.UseMiddleware(
				func(ctx huma.Context, next func(huma.Context)) {
					ctx = huma.WithContext(ctx, AttachT(ctx.Context(), i18nBundle, ctx.Header("Accept-Language"), ctx.Query("locale")))
					next(ctx)
				},
			)

			huma.Register(api, huma.Operation{
				Method: http.MethodPost,
				Path:   basePath,
			}, func(ctx context.Context, _ *struct{}) (*struct {
				Body struct {
					Locale           string `json:"locale"`
					TranslatedString string `json:"translated_string"`
				}
			}, error) {
				translatedString, err := GetT(ctx).LocalizeMessage(&i18n.Message{
					ID: "welcome.verify_button",
				})
				assert.NoError(t, err)

				return &struct {
					Body struct {
						Locale           string `json:"locale"`
						TranslatedString string `json:"translated_string"`
					}
				}{
					Body: struct {
						Locale           string `json:"locale"`
						TranslatedString string `json:"translated_string"`
					}{
						Locale:           GetLocale(ctx),
						TranslatedString: translatedString,
					},
				}, nil
			})

			response := api.Post(requestPath, fmt.Sprintf("Accept-Language: %s", tc.acceptLanguage))
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, tc.want, strings.Trim(response.Body.String(), "\n"))
		})
	}
}
