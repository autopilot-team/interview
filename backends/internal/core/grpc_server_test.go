package core

import (
	"log/slog"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewGrpcServer(t *testing.T) {
	tests := []struct {
		name    string
		opts    GrpcServerOptions
		wantErr bool
	}{
		{
			name: "valid configuration",
			opts: GrpcServerOptions{
				Logger: slog.Default(),
				Port:   "50051",
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			opts: GrpcServerOptions{
				Logger: slog.Default(),
				Port:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewGrpcServer(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGrpcServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && server == nil {
				t.Error("NewGrpcServer() returned nil server when no error was expected")
			}
			if server != nil {
				server.Stop()
			}
		})
	}
}

func TestGrpcServer_Addr(t *testing.T) {
	opts := GrpcServerOptions{
		Logger: slog.Default(),
		Port:   "0", // Use port 0 to get a random available port
	}

	server, err := NewGrpcServer(opts)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Stop()

	addr := server.Addr()
	if addr == "" {
		t.Error("Addr() returned empty string")
	}
}

func TestGrpcServer_ListenAndServe(t *testing.T) {
	// Create a bufconn listener for testing
	listener := bufconn.Listen(1024 * 1024)
	server := &GrpcServer{
		Server: grpc.NewServer(),
		lis:    listener,
	}

	// Use WaitGroup to ensure the goroutine completes
	var wg sync.WaitGroup
	wg.Add(1)

	// Start server in a goroutine
	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil && err.Error() != "closed" {
			t.Logf("ListenAndServe returned error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server
	server.Stop()

	// Wait for the goroutine to complete
	wg.Wait()
}

func TestGrpcServer_Stop(t *testing.T) {
	opts := GrpcServerOptions{
		Logger: slog.Default(),
		Port:   "0",
	}

	server, err := NewGrpcServer(opts)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Server should stop gracefully
	server.Stop()
}
