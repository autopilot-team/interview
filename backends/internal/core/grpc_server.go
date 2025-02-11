package core

import (
	"autopilot/backends/internal/grpc/middleware"
	"fmt"
	"log/slog"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// GrpcServerOptions contains configuration options for creating a new gRPC server.
type GrpcServerOptions struct {
	// Logger is used for logging server events
	Logger *slog.Logger
	// Host is the host the server will listen on
	Host string
	// Port is the port number the server will listen on
	Port string
	// ServerOptions are additional gRPC server options
	ServerOptions []grpc.ServerOption
}

// GrpcServer wraps a gRPC server with additional functionality.
type GrpcServer struct {
	*grpc.Server
	lis net.Listener
}

// NewGrpcServer creates and returns a new GrpcServer instance.
// It initializes the server with the provided options and sets up middleware
// for logging and recovery.
func NewGrpcServer(opts GrpcServerOptions) (*GrpcServer, error) {
	lis, err := net.Listen("tcp", opts.Host+":"+opts.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	serverOpts := append(opts.ServerOptions,
		grpc.UnaryInterceptor(middleware.UnaryOperationMode),
		grpc.ChainUnaryInterceptor(
			middleware.Logger(opts.Logger),
			middleware.Recovery(opts.Logger),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return &GrpcServer{
		grpc.NewServer(serverOpts...),
		lis,
	}, nil
}

// Addr returns the address the server is listening on.
func (s *GrpcServer) Addr() string {
	return s.lis.Addr().String()
}

// ListenAndServe starts the gRPC server and begins accepting connections.
// It blocks until the server is stopped or encounters an error.
func (s *GrpcServer) ListenAndServe() error {
	return s.Serve(s.lis)
}

// Stop gracefully stops the gRPC server, waiting for all RPCs to complete.
func (s *GrpcServer) Stop() {
	s.GracefulStop()
}
