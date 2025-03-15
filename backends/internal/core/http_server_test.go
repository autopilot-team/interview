package core

import (
	"autopilot/backends/internal/types"
	"embed"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed all:testdata/dist
var testFS embed.FS

func TestHttpServer_ServeStaticFiles(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedType   string
		expectedBody   string
	}{
		{
			name:           "should serve index.html for root path",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html",
			expectedBody:   "<h1>Index</h1>",
		},
		{
			name:           "should serve 404.html for non-existent path",
			path:           "/nonexistent",
			expectedStatus: http.StatusOK,
			expectedType:   "text/html",
			expectedBody:   "<h1>404</h1>",
		},
		{
			name:           "should serve static file correctly",
			path:           "/style.css",
			expectedStatus: http.StatusOK,
			expectedType:   "text/css; charset=utf-8",
			expectedBody:   "body {\n	color: black;\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new test server for each test case
			server, err := NewHttpServer(HttpServerOptions{
				Mode:       types.DebugMode,
				I18nBundle: &I18nBundle{},
				Logger:     logger,
				Host:       "localhost",
				Port:       "3000",
				SpaFS:      testFS,
				SpaDir:     "testdata/dist",
			})
			if err != nil {
				t.Fatal(err)
			}

			// Create a test request with a valid path
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			// Serve the request
			server.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, tt.expectedType) {
				t.Errorf("expected content type %s, got %s", tt.expectedType, contentType)
			}

			// Check body content
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
