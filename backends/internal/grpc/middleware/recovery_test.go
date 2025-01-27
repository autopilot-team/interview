package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecovery(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelError,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	tests := []struct {
		name         string
		handler      func(ctx context.Context, req interface{}) (interface{}, error)
		expectPanic  bool
		expectedCode codes.Code
		expectedResp interface{}
		expectedLog  string
	}{
		{
			name: "successful execution - no panic",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			},
			expectPanic:  false,
			expectedCode: codes.OK,
			expectedResp: "success",
			expectedLog:  "",
		},
		{
			name: "panic with error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				panic("test panic error")
			},
			expectPanic:  true,
			expectedCode: codes.Internal,
			expectedResp: nil,
			expectedLog:  "level=ERROR msg=\"Panic recovered\" error=\"test panic error\"",
		},
		{
			name: "panic with non-error value",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				panic(123)
			},
			expectPanic:  true,
			expectedCode: codes.Internal,
			expectedResp: nil,
			expectedLog:  "level=ERROR msg=\"Panic recovered\" error=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the buffer before each test
			logBuffer.Reset()

			recoveryFunc := Recovery(logger)

			// Create a mock UnaryServerInfo
			info := &grpc.UnaryServerInfo{
				FullMethod: "test.method",
			}

			// Create a wrapper handler that matches the grpc.UnaryHandler signature
			wrappedHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return tt.handler(ctx, req)
			}

			resp, err := recoveryFunc(context.Background(), "test-request", info, wrappedHandler)

			if tt.expectPanic {
				if err == nil {
					t.Error("expected error from panic, got nil")
					return
				}

				st, ok := status.FromError(err)
				if !ok {
					t.Error("expected gRPC status error")
					return
				}

				if st.Code() != tt.expectedCode {
					t.Errorf("expected status code %v, got %v", tt.expectedCode, st.Code())
				}

				// Check log output
				logOutput := strings.TrimSpace(logBuffer.String())
				if tt.expectedLog != "" && logOutput != tt.expectedLog {
					t.Errorf("expected log:\n%q\ngot:\n%q", tt.expectedLog, logOutput)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if resp != tt.expectedResp {
					t.Errorf("expected response %v, got %v", tt.expectedResp, resp)
				}

				// Verify no logs were written for successful cases
				if logBuffer.Len() > 0 {
					t.Errorf("expected no logs for successful case, got: %s", logBuffer.String())
				}
			}
		})
	}
}
