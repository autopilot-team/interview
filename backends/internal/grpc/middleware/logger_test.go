package middleware

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

type testLogBuffer struct {
	logs []string
}

func (b *testLogBuffer) Handle(ctx context.Context, r slog.Record) error {
	var sb strings.Builder
	r.Attrs(func(a slog.Attr) bool {
		sb.WriteString(a.String())
		return true
	})
	b.logs = append(b.logs, sb.String())

	return nil
}

func (b *testLogBuffer) WithAttrs(attrs []slog.Attr) slog.Handler {
	return b
}

func (b *testLogBuffer) WithGroup(name string) slog.Handler {
	return b
}

func (b *testLogBuffer) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name          string
		handler       grpc.UnaryHandler
		expectedError error
		checkLogs     func(t *testing.T, logs []string)
	}{
		{
			name: "successful request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", nil
			},
			expectedError: nil,
			checkLogs: func(t *testing.T, logs []string) {
				if len(logs) != 1 {
					t.Errorf("expected 1 log entry, got %d", len(logs))
				}

				log := logs[0]
				if !strings.Contains(log, "service=test") {
					t.Error("log does not contain service name")
				}

				if !strings.Contains(log, "method=method") {
					t.Error("log does not contain method name")
				}

				if !strings.Contains(log, "latency=") {
					t.Error("log does not contain latency")
				}

				if strings.Contains(log, "error_message=") {
					t.Error("log should not contain error for successful request")
				}
			},
		},
		{
			name: "request with error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, errors.New("test error")
			},
			expectedError: errors.New("test error"),
			checkLogs: func(t *testing.T, logs []string) {
				if len(logs) != 1 {
					t.Errorf("expected 1 log entry, got %d", len(logs))
				}

				log := logs[0]
				if !strings.Contains(log, "service=test") {
					t.Error("log does not contain service name")
				}

				if !strings.Contains(log, "method=method") {
					t.Error("log does not contain method name")
				}

				if !strings.Contains(log, "latency=") {
					t.Error("log does not contain latency")
				}

				if !strings.Contains(log, "error_message=test error") {
					t.Error("log does not contain error message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &testLogBuffer{}
			logger := slog.New(buffer)

			middleware := Logger(logger)
			info := &grpc.UnaryServerInfo{
				FullMethod: "/test/method",
			}

			resp, err := middleware(context.Background(), "request", info, tt.handler)

			if tt.expectedError == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedError != nil && err == nil {
				t.Error("expected error, got nil")
			}

			if tt.expectedError != nil && err != nil && tt.expectedError.Error() != err.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			// Give a small delay to ensure logs are processed
			time.Sleep(time.Millisecond)
			tt.checkLogs(t, buffer.logs)

			if tt.expectedError == nil {
				if resp != "response" {
					t.Errorf("expected response 'response', got %v", resp)
				}
			}
		})
	}
}
