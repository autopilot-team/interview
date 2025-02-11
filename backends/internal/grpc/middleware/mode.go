package middleware

import (
	"autopilot/backends/internal/types"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	operationModeMetadataKey = "x-operation-mode"
)

// UnaryOperationMode extracts operation mode from metadata and sets it in context
func UnaryOperationMode(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Get operation mode from metadata
		if modes := md.Get(operationModeMetadataKey); len(modes) > 0 {
			mode := types.OperationMode(modes[0])
			ctx = context.WithValue(ctx, types.OperationModeKey, mode)
		}
	}

	return handler(ctx, req)
}

// UnaryOperationModeClientInterceptor creates a client interceptor that adds operation mode to outgoing requests
func UnaryOperationModeClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add operation mode to outgoing context
	outCtx := AddOperationModeToOutgoingContext(ctx)

	return invoker(outCtx, method, req, reply, cc, opts...)
}

// AddOperationModeToOutgoingContext adds operation mode to gRPC metadata for outgoing requests
func AddOperationModeToOutgoingContext(ctx context.Context) context.Context {
	mode := types.GetOperationMode(ctx)
	md := metadata.New(map[string]string{
		operationModeMetadataKey: string(mode),
	})

	return metadata.NewOutgoingContext(ctx, md)
}
