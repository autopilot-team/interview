package middleware

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		mode           types.Mode
		method         string
		path           string
		headers        map[string]string
		expectedLogs   []string
		unexpectedLogs []string
	}{
		{
			name:   "should log headers in debug mode",
			mode:   types.DebugMode,
			method: "GET",
			path:   "/api/v1/test",
			headers: map[string]string{
				"Accept":          "*/*",
				"Accept-Language": "en-US",
				"User-Agent":      "test-agent",
				"Referer":         "http://localhost:3000/docs",
			},
			expectedLogs: []string{
				"INFO  GET /api/v1/test",
				"Accept: */*",
				"Accept-Language: en-US",
				"Referer: http://localhost:3000/docs",
				"User-Agent: test-agent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture logs
			var buf bytes.Buffer

			logger := core.NewLogger(core.LoggerOptions{
				Mode:   tt.mode,
				Writer: &buf,
			})

			// Create a test handler that returns 200 OK
			handler := Logger(tt.mode, logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				time.Sleep(10 * time.Millisecond) // Add a small delay to test duration logging
			}))

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = "[::1]:53150" // Set a consistent remote address
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check response status
			assert.Equal(t, http.StatusOK, rr.Code)

			// Get the log output
			logOutput := buf.String()

			// Check expected logs
			for _, expected := range tt.expectedLogs {
				assert.Contains(t, logOutput, expected)
			}
		})
	}
}

func TestLoggerWithPanic(t *testing.T) {
	t.Parallel()
	// Create a buffer to capture logs
	var buf bytes.Buffer

	logger := core.NewLogger(core.LoggerOptions{
		Mode:   types.DebugMode,
		Writer: &buf,
	})

	// Create a new router with the middleware chain
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(Logger(types.DebugMode, logger))

	// Add a handler that panics
	r.Get("/api/v1/test", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	req.RemoteAddr = "[::1]:53150"
	rr := httptest.NewRecorder()

	// Serve the request (should recover from panic)
	r.ServeHTTP(rr, req)

	// Check that we got a 500 status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify the panic was logged in the correct format
	logOutput := buf.String()
	assert.Contains(t, logOutput, "INFO  GET /api/v1/test")
	// The log format is: "INFO  GET /path status duration [ip]"
	assert.Regexp(t, `INFO\s+GET /api/v1/test \d+ \d+(\.\d+)?(µs|ms|s) \[::1\]:53150`, logOutput)
}

func TestLoggerWithCustomStatus(t *testing.T) {
	t.Parallel()
	// Create a buffer to capture logs
	var buf bytes.Buffer

	logger := core.NewLogger(core.LoggerOptions{
		Mode:   types.DebugMode,
		Writer: &buf,
	})

	// Create handlers with different status codes
	statuses := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}

	for _, status := range statuses {
		t.Run(fmt.Sprintf("%d", status), func(t *testing.T) {
			buf.Reset() // Clear the buffer for each test

			handler := Logger(types.DebugMode, logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))

			// Create a test request
			req := httptest.NewRequest("GET", "/api/v1/test", nil)
			req.RemoteAddr = "[::1]:53150"
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Verify the status was logged in the correct format
			logOutput := buf.String()
			assert.Contains(t, logOutput, "INFO  GET /api/v1/test")
			assert.Contains(t, logOutput, fmt.Sprintf(" %d ", status)) // Status should be surrounded by spaces
		})
	}
}
