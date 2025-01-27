package middleware

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logger returns a gRPC UnaryServerInterceptor for logging
func Logger(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		// Extract useful information from the method
		service := path.Dir(info.FullMethod)[1:] // Remove leading slash
		method := path.Base(info.FullMethod)

		// Handle the request
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		// Build log attributes
		attrs := []any{
			"service", service,
			"method", method,
			"latency", duration,
		}

		// Add error details if present
		if err != nil {
			st, _ := status.FromError(err)
			attrs = append(attrs,
				"error_code", st.Code().String(),
				"error_message", st.Message(),
			)

			// Log errors at error level
			logger.ErrorContext(ctx, "gRPC request failed", attrs...)
			return resp, err
		}

		// Log success at info level
		logger.InfoContext(ctx, "gRPC request completed", attrs...)
		return resp, nil
	}
}
