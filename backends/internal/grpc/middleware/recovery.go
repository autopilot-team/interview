package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recovery is a middleware that recovers from panics in gRPC handlers and logs them.
// It converts panics into gRPC errors with Internal status code to ensure the service
// remains stable even when unexpected errors occur.
//
// The middleware accepts a slog.Logger instance and returns a gRPC UnaryServerInterceptor
// that will catch any panics, log them using the provided logger, and return an appropriate
// gRPC error response to the client.
func Recovery(logger *slog.Logger) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, "Panic recovered", "error", r)
				err = status.Error(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
