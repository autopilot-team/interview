package core

import (
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Logger is a wrapper around slog.Logger
type Logger = slog.Logger

// LoggerOptions contains configuration options for the logger
type LoggerOptions struct {
	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// Writer is the writer to write the logs to
	Writer io.Writer
}

// NewLogger creates a new structured logger
func NewLogger(opts LoggerOptions) *Logger {
	if opts.Mode == types.Mode("") {
		opts.Mode = types.DebugMode
	}

	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	var handler slog.Handler
	if opts.Mode == types.DebugMode {
		handler = &DebugLogHandler{out: opts.Writer}
	} else {
		handler = &ReleaseLogHandler{
			handler: slog.NewJSONHandler(opts.Writer, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
		}
	}

	return slog.New(handler)
}

// DebugLogHandler implements slog.Handler interface for debug mode logging.
// It provides colored, human-readable output with special formatting for HTTP requests.
type DebugLogHandler struct {
	out   io.Writer   // Output writer for log messages
	attrs []slog.Attr // Attributes to include with each log entry
	mut   sync.Mutex  // Mutex to synchronize writes
}

// Enabled implements slog.Handler interface.
// In debug mode, only INFO level and above are enabled.
func (h *DebugLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

// Handle implements slog.Handler interface.
// It formats and writes a log record with special handling for HTTP request logs.
func (h *DebugLogHandler) Handle(_ context.Context, r slog.Record) error {
	h.mut.Lock()
	defer h.mut.Unlock()

	timeStr := color.New(color.FgHiBlack).Sprint(r.Time.Format("15:04:05"))
	level := levelColor(r.Level)

	// Extract HTTP-specific attributes
	var method, path, service, status, latency, ip string
	var headers map[string]string
	attrs := append(h.attrs, []slog.Attr{}...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	for _, attr := range attrs {
		switch attr.Key {
		case "headers":
			if valuer, ok := attr.Value.Any().(slog.LogValuer); ok {
				if v := valuer.LogValue(); v.Kind() == slog.KindGroup {
					headers = make(map[string]string)
					for _, a := range v.Group() {
						headers[a.Key] = a.Value.String()
					}
				}
			}
		case "ip":
			ip = attr.Value.String()
		case "latency":
			latency = attr.Value.String()
		case "method":
			method = attr.Value.String()
		case "path":
			path = attr.Value.String()
		case "service":
			service = attr.Value.String()
		case "status":
			// Handle both string and int64 status codes
			switch attr.Value.Kind() {
			case slog.KindInt64:
				status = fmt.Sprintf("%d", attr.Value.Int64())
			case slog.KindString:
				status = attr.Value.String()
			}
		}
	}

	logEntry := ""
	if method != "" {
		if service == "" {
			// Format headers nicely
			var headerStr string
			if len(headers) > 0 {
				headerLines := []string{}
				for k, v := range headers {
					headerLines = append(headerLines, fmt.Sprintf("\t\t%s: %s", k, v))
				}

				sort.Strings(headerLines)
				headerStr = color.BlackString("\n" + strings.Join(headerLines, "\n"))
			}

			// HTTP request log format
			logEntry = fmt.Sprintf("%s %s %s %s %s %s %s%s\n",
				timeStr,
				level,
				methodColor(method),
				path,
				statusColor(status),
				latency,
				ip,
				headerStr,
			)
		} else {
			// gRPC request log format
			var errorCode, errorMsg string
			for _, attr := range attrs {
				if attr.Key == "error_code" {
					errorCode = attr.Value.String()
				}

				if attr.Key == "error_message" {
					errorMsg = attr.Value.String()
				}
			}

			if errorCode != "" && errorMsg != "" {
				logEntry = fmt.Sprintf("%s %s %s %s %s %s\n%s%s\n",
					timeStr,
					level,
					color.HiBlueString(service),
					method,
					color.HiRedString(errorCode),
					color.HiYellowString(latency),
					color.New(color.FgHiBlack).Sprint("\t\t ↳ message = "),
					color.New(color.FgHiWhite).Sprint(errorMsg),
				)
			} else {
				logEntry = fmt.Sprintf("%s %s %s %s %s %s\n",
					timeStr,
					level,
					color.HiBlueString(service),
					method,
					color.HiGreenString("OK"),
					color.HiYellowString(latency),
				)
			}
		}
	} else {
		// Standard log format
		logEntry = fmt.Sprintf("%s %s %s%s\n",
			timeStr,
			level,
			r.Message,
			formatAttributes(attrs),
		)
	}

	_, err := h.out.Write([]byte(logEntry))
	return err
}

// WithAttrs implements slog.Handler interface.
// It returns a new handler with the given attributes added to the existing ones.
func (h *DebugLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DebugLogHandler{
		out:   h.out,
		attrs: append(h.attrs, attrs...),
		mut:   sync.Mutex{},
	}
}

// WithGroup implements slog.Handler interface.
// Groups are not supported in debug mode, so it returns the handler unchanged.
func (h *DebugLogHandler) WithGroup(name string) slog.Handler {
	return h
}

// ReleaseLogHandler implements slog.Handler interface for release mode logging.
// It wraps the standard JSON handler and adds OpenTelemetry integration.
type ReleaseLogHandler struct {
	attrs   []slog.Attr  // Additional attributes to include with each log entry
	handler slog.Handler // Underlying JSON handler
}

// Enabled implements slog.Handler interface.
// Delegates to the underlying handler's Enabled method.
func (h *ReleaseLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler interface.
// Creates an OpenTelemetry span for each log entry and adds trace information.
func (h *ReleaseLogHandler) Handle(ctx context.Context, r slog.Record) error {
	logAttrs := append(h.attrs, []slog.Attr{}...)
	r.Attrs(func(a slog.Attr) bool {
		logAttrs = append(logAttrs, a)
		return true
	})

	span := trace.SpanFromContext(ctx)
	span.AddEvent(r.Level.String(), trace.WithAttributes(
		attribute.String("message", r.Message),
	))

	return h.handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler interface.
// Returns a new handler with the given attributes added to both the wrapper and underlying handler.
func (h *ReleaseLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ReleaseLogHandler{
		handler: h.handler.WithAttrs(attrs),
		attrs:   append(h.attrs, attrs...),
	}
}

// WithGroup implements slog.Handler interface.
// Delegates group creation to the underlying handler.
func (h *ReleaseLogHandler) WithGroup(name string) slog.Handler {
	return &ReleaseLogHandler{
		handler: h.handler.WithGroup(name),
		attrs:   h.attrs,
	}
}

// formatAttributes formats a slice of attributes as a space-separated string.
func formatAttributes(attrs []slog.Attr) string {
	if len(attrs) == 0 {
		return ""
	}

	var parts []string
	for _, attr := range attrs {
		parts = append(parts, fmt.Sprintf("%s=%s", attr.Key, formatAttrValue(attr.Value)))
	}

	return " " + strings.Join(parts, " ")
}

// formatAttrValue formats a slog.Value based on its kind.
func formatAttrValue(v slog.Value) string {
	// Handle LogValuer interface first
	if valuer, ok := v.Any().(slog.LogValuer); ok {
		return formatAttrValue(valuer.LogValue())
	}

	switch v.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%q", v.String())
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindBool:
		return fmt.Sprintf("%t", v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	case slog.KindAny:
		if m, ok := v.Any().(map[string]string); ok {
			parts := make([]string, 0, len(m))
			for k, v := range m {
				parts = append(parts, fmt.Sprintf("%s:%s", k, v))
			}

			return strings.Join(parts, " ")
		}

		return fmt.Sprintf("%v", v.Any())
	default:
		return fmt.Sprintf("%v", v)
	}
}

// levelColor returns a colored string representation of the log level.
func levelColor(level slog.Level) string {
	var bg, fg color.Attribute
	switch level {
	case slog.LevelDebug:
		bg, fg = color.BgMagenta, color.FgWhite
	case slog.LevelInfo:
		bg, fg = color.BgBlue, color.FgWhite
	case slog.LevelWarn:
		bg, fg = color.BgYellow, color.FgBlack
	case slog.LevelError:
		bg, fg = color.BgRed, color.FgWhite
	default:
		bg, fg = color.BgWhite, color.FgBlack
	}

	return color.New(bg, fg, color.Bold).Sprint(" " + strings.ToUpper(level.String()) + " ")
}

// methodColor returns a colored string representation of the HTTP method.
func methodColor(method string) string {
	switch method {
	case "GET":
		return color.BlueString(method)
	case "POST":
		return color.GreenString(method)
	case "PUT":
		return color.YellowString(method)
	case "DELETE":
		return color.RedString(method)
	case "PATCH":
		return color.CyanString(method)
	default:
		return color.MagentaString(method)
	}
}

// statusColor returns a colored string representation of the HTTP status code.
func statusColor(status string) string {
	code := 0
	_, _ = fmt.Sscanf(status, "%d", &code)
	switch {
	case code >= 200 && code < 300:
		return color.GreenString(status)
	case code >= 300 && code < 400:
		return color.CyanString(status)
	case code >= 400 && code < 500:
		return color.YellowString(status)
	default:
		return color.RedString(status)
	}
}
