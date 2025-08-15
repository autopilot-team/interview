package core

import (
	"autopilot/backends/internal/types"
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebugLogHandler(t *testing.T) {
	t.Parallel()
	t.Run("handles standard log message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(LoggerOptions{
			Mode:   types.DebugMode,
			Writer: &buf,
		})

		logger.Info("test message", "password", RedactString("value"))

		output := buf.String()
		assert.Contains(t, output, "INFO  test message password=\"[REDACTED]\"")
	})

	t.Run("should handle HTTP request log with formatted headers", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(LoggerOptions{
			Mode:   types.DebugMode,
			Writer: &buf,
		})
		headers := HTTPHeaders{
			"Accept":     "*/*",
			"User-Agent": "Mozilla/5.0",
			"Referer":    "http://localhost:3000",
		}

		logger.Info("HTTP Request",
			"method", "GET",
			"path", "/test",
			"status", 200,
			"latency", 100*time.Millisecond,
			"ip", "127.0.0.1",
			"headers", headers,
		)

		output := buf.String()
		assert.Contains(t, output, "GET /test 200 100ms 127.0.0.1\n\t\tAccept: */*\n\t\tReferer: http://localhost:3000\n\t\tUser-Agent: Mozilla/5.0")
	})
}

func TestReleaseLogHandler(t *testing.T) {
	t.Parallel()
	t.Run("adds trace information", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(LoggerOptions{
			Mode:   types.ReleaseMode,
			Writer: &buf,
		})

		logger.Info("test message")

		var logEntry map[string]any
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err)
	})

	t.Run("should handle standard log message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(LoggerOptions{
			Mode:   types.ReleaseMode,
			Writer: &buf,
		})

		logger.Info("test message", "password", RedactString("value"))

		output := buf.String()
		assert.Contains(t, output, "\"level\":\"INFO\"")
		assert.Contains(t, output, "\"msg\":\"test message\"")
		assert.Contains(t, output, "\"password\":\"[REDACTED]\"")
	})

	t.Run("should handle HTTP request log with formatted headers", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(LoggerOptions{
			Mode:   types.ReleaseMode,
			Writer: &buf,
		})
		headers := HTTPHeaders{
			"Accept":     "*/*",
			"User-Agent": "Mozilla/5.0",
			"Referer":    "http://localhost:3000",
		}

		logger.Info("HTTP Request",
			"method", "GET",
			"path", "/test",
			"status", 200,
			"latency", 100*time.Millisecond,
			"ip", "127.0.0.1",
			"headers", headers,
		)

		output := buf.String()
		assert.Contains(t, output, "\"level\":\"INFO\"")
		assert.Contains(t, output, "\"msg\":\"HTTP Request\"")
		assert.Contains(t, output, "\"method\":\"GET\"")
		assert.Contains(t, output, "\"path\":\"/test\"")
		assert.Contains(t, output, "\"status\":200")
		assert.Contains(t, output, "\"latency\":100000000")
		assert.Contains(t, output, "\"ip\":\"127.0.0.1\"")
		assert.Contains(t, output, "\"headers\":{")
	})
}

func TestFormatAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		attrs    []slog.Attr
		expected string
	}{
		{
			name:     "should handle empty attributes",
			attrs:    []slog.Attr{},
			expected: "",
		},
		{
			name: "should format string attribute",
			attrs: []slog.Attr{
				slog.String("key", "value"),
			},
			expected: ` key="value"`,
		},
		{
			name: "should format multiple attributes",
			attrs: []slog.Attr{
				slog.Int("count", 42),
				slog.Bool("enabled", true),
			},
			expected: ` count=42 enabled=true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatAttributes(tt.attrs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLevelColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		level    slog.Level
		contains string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
	}

	for _, tt := range tests {
		t.Run("should color "+tt.contains+" level correctly", func(t *testing.T) {
			result := levelColor(tt.level)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestMethodColor(t *testing.T) {
	t.Parallel()
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "CUSTOM"}

	for _, method := range methods {
		t.Run("should color "+method+" method correctly", func(t *testing.T) {
			result := methodColor(method)
			assert.Contains(t, result, method)
		})
	}
}

func TestStatusColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		status string
		code   int
	}{
		{"200", 200},
		{"301", 301},
		{"404", 404},
		{"500", 500},
	}

	for _, tt := range tests {
		t.Run("should color status "+tt.status+" correctly", func(t *testing.T) {
			result := statusColor(tt.status)
			assert.Contains(t, result, tt.status)
		})
	}
}
