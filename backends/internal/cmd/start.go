package cmd

import (
	"autopilot/backends/internal/core"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func NewStartCmd(ctx context.Context, logger *slog.Logger, httpServer *core.HTTPServer, grpcServer *core.GrpcServer, liveWorker core.Worker, testWorker core.Worker, cleanUp func()) *cobra.Command {
	var runServer, runWorker bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the application server and/or worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Start(ctx, logger, httpServer, grpcServer, liveWorker, testWorker, cleanUp, StartOptions{
				RunServer: runServer,
				RunWorker: runWorker,
			})
		},
	}

	cmd.Flags().BoolVar(&runServer, "server", true, "Run server")
	cmd.Flags().BoolVar(&runWorker, "worker", false, "Run worker")

	return cmd
}

type StartOptions struct {
	RunServer bool
	RunWorker bool
}

func Start(ctx context.Context, logger *slog.Logger, httpServer *core.HTTPServer, grpcServer *core.GrpcServer, liveWorker core.Worker, testWorker core.Worker, cleanUp func(), opts StartOptions) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	serverDone := make(chan struct{})
	liveWorkerDone := make(chan struct{})
	testWorkerDone := make(chan struct{})

	if opts.RunServer {
		if httpServer != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(serverDone)
				logger.Info("Started HTTP server at http://" + httpServer.Addr)
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					errChan <- fmt.Errorf("http server error: %w", err)
				}
			}()
		}

		if grpcServer != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(serverDone)
				logger.Info("Started GRPC server at tcp://" + grpcServer.Addr())
				if err := grpcServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					errChan <- fmt.Errorf("http server error: %w", err)
				}
			}()
		}
	} else {
		close(serverDone)
	}

	if opts.RunWorker {
		if liveWorker != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(liveWorkerDone)
				if err := liveWorker.Start(ctx); err != nil {
					errChan <- fmt.Errorf("live worker error: %w", err)
				}
			}()
		}

		if testWorker != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(testWorkerDone)
				if err := testWorker.Start(ctx); err != nil {
					errChan <- fmt.Errorf("test worker error: %w", err)
				}
			}()
		}
	} else {
		close(liveWorkerDone)
		close(testWorkerDone)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	select {
	case err := <-errChan:
		return err

	case sig := <-shutdown:
		logger.Info("Graceful shutdown signal received", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if opts.RunServer {
			if httpServer != nil {
				if err := httpServer.Shutdown(ctx); err != nil {
					logger.Error("Graceful shutdown failed", slog.String("error", err.Error()))
					if err := httpServer.Close(); err != nil {
						logger.Error("Could not stop HTTP server gracefully: " + err.Error())
					}
				}
			}

			if grpcServer != nil {
				grpcServer.Stop()
			}
		}

		if opts.RunWorker {
			if err := liveWorker.Stop(ctx); err != nil {
				logger.Error("Graceful shutdown failed", slog.String("error", err.Error()))
			}
		}

		if cleanUp != nil {
			cleanUp()
		}

		<-serverDone

		if liveWorker != nil {
			<-liveWorkerDone
		}

		if testWorker != nil {
			<-testWorkerDone
		}
	}

	return nil
}
