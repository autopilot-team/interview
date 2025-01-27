package core

import (
	"autopilot/backends/internal/types"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed all:testdata
	mailerPreviewFS embed.FS
)

func setupTestMailerPreview(t *testing.T) (*MailerPreview, Mailer) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerPreviewFS, "testdata/locales")
	require.NoError(t, err)

	mailer, err := NewMail(MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       types.DebugMode,
		SmtpUrl:    smtpUrl,
		TemplateOptions: &MailTemplateOptions{
			FS:     mailerPreviewFS,
			Dir:    "testdata/templates",
			Layout: "layouts/test",
			ExtraFuncs: []template.FuncMap{
				{
					"fail": func(msg string) (string, error) {
						return "", fmt.Errorf("%s", msg)
					},
					"upper": strings.ToUpper,
				},
			},
		},
	})
	require.NoError(t, err)

	preview := NewMailerPreview(mailer, i18nBundle, logger)
	return preview, mailer
}

func TestNewMailerPreview(t *testing.T) {
	preview, mailer := setupTestMailerPreview(t)

	assert.NotNil(t, preview)
	assert.Equal(t, mailer, preview.mailer)
	assert.NotNil(t, preview.data)
	assert.Empty(t, preview.data)
}

func TestSetPreviewData(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)

	testData := map[string]map[string]interface{}{
		"welcome": {
			"Name":    "John",
			"AppName": "TestApp",
		},
	}

	preview.SetPreviewData(testData)
	assert.Equal(t, testData, preview.data)

	// Test GetPreviewData
	data := preview.GetPreviewData("welcome")
	assert.Equal(t, testData["welcome"], data)

	// Test non-existent template
	data = preview.GetPreviewData("nonexistent")
	assert.Empty(t, data)
}

func TestHandleListEmailTemplates(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)
	router := chi.NewRouter()
	preview.setupEmailPreviewRoutes(router)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "list templates without query params",
			path:           "/mailer/preview",
			expectedStatus: http.StatusFound, // Should redirect to first template
		},
		{
			name:           "list templates with template and locale",
			path:           "/mailer/preview?template=welcome&locale=en",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "Mailer Preview")
				assert.Contains(t, w.Body.String(), "welcome")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestHandlePreviewEmailTemplate(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)
	router := chi.NewRouter()
	preview.setupEmailPreviewRoutes(router)

	// Set test data
	preview.SetPreviewData(map[string]map[string]interface{}{
		"basic": {
			"Name": "John",
		},
		"with-functions": {
			"Text": "hello",
		},
	})

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "preview basic HTML template",
			path:           "/mailer/preview/basic?format=html&locale=en",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
				assert.Contains(t, w.Body.String(), "Hello, John!")
			},
		},
		{
			name:           "preview template with functions",
			path:           "/mailer/preview/with-functions?format=html&locale=en",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
				assert.Contains(t, w.Body.String(), "HELLO")
			},
		},
		{
			name:           "preview basic text template",
			path:           "/mailer/preview/basic?format=text&locale=en",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
				assert.Contains(t, w.Body.String(), "Hello, John!")
			},
		},
		{
			name:           "preview non-existent template",
			path:           "/mailer/preview/nonexistent?format=html&locale=en",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestHandleSendTestEmail(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)
	router := chi.NewRouter()
	preview.setupEmailPreviewRoutes(router)

	// Set test data
	preview.SetPreviewData(map[string]map[string]interface{}{
		"basic": {
			"Name": "John",
		},
	})

	tests := []struct {
		name           string
		templateName   string
		requestBody    map[string]string
		expectedStatus int
	}{
		{
			name:         "send test email success",
			templateName: "basic",
			requestBody: map[string]string{
				"email":  "test@example.com",
				"locale": "en",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			templateName:   "basic",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "non-existent template",
			templateName: "nonexistent",
			requestBody: map[string]string{
				"email":  "test@example.com",
				"locale": "en",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.requestBody != nil {
				jsonBody, err := json.Marshal(tt.requestBody)
				require.NoError(t, err)
				body = strings.NewReader(string(jsonBody))
			}

			req := httptest.NewRequest("POST", "/mailer/preview/"+tt.templateName+"/send", body)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTemplateManager(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)
	templateMgr := newTemplateManager(preview.mailer.TemplateOptions().FS, preview.mailer.TemplateOptions().Dir)

	templates, err := templateMgr.getTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, templates)

	// Verify that layout templates are not included
	for _, tmpl := range templates {
		assert.False(t, strings.HasPrefix(tmpl, "layouts/"))
	}
}

func TestPreviewRenderer(t *testing.T) {
	preview, _ := setupTestMailerPreview(t)
	templateMgr := newTemplateManager(preview.mailer.TemplateOptions().FS, preview.mailer.TemplateOptions().Dir)
	renderer := newPreviewRenderer(templateMgr)

	assert.NotNil(t, renderer.cssProvider)
	assert.NotNil(t, renderer.jsProvider)
	assert.NotNil(t, renderer.templateMgr)

	// Test CSS and JS providers
	css := renderer.cssProvider.getCSS()
	assert.NotEmpty(t, css)
	assert.Contains(t, css, "body")

	js := renderer.jsProvider.getJS()
	assert.NotEmpty(t, js)
	assert.Contains(t, js, "function")
}
