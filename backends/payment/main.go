package main

import (
	"autopilot/backends/internal/cmd"
	"autopilot/backends/internal/core"
	paymentv1 "autopilot/backends/internal/pbgen/payment/v1"
	"autopilot/backends/internal/types"
	"autopilot/backends/payment/internal/app"
	v1 "autopilot/backends/payment/internal/handler/v1"
	"autopilot/backends/payment/internal/worker"
	"context"
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/riverqueue/river"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	//go:embed all:migrations
	migrationsFS embed.FS
)

func main() {
	if !types.Mode(mode).IsValid() {
		log.Fatalf("Invalid mode: %s. Must be either 'debug' or 'release'", mode)
	}

	_, mainFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get caller information")
	}

	// Initialize the context
	ctx := context.Background()

	// Initialize the container
	container, err := app.NewContainer(ctx, app.ContainerOpts{
		FS: app.ContainerFS{
			Migrations: migrationsFS,
		},
		MainFile: mainFile,
		Mode:     mode,
	})
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize the GRPC server
	grpcServer, err := initGrpcServer(container)
	if err != nil {
		log.Fatalf("Failed to initialize GRPC server: %v", err)
	}

	// Initialize the background workers
	workers := worker.Register(container)

	liveWorker, err := core.NewWorker(ctx, core.WorkerOptions{
		Config: &river.Config{
			Queues: map[string]river.QueueConfig{
				"default": {
					MaxWorkers: 100,
				},
			},
			Workers: workers,
		},
		DbURL:  container.Config.Database.LiveWorker,
		Logger: container.Logger,
	})
	if err != nil {
		log.Fatalf("Failed to initialize live worker: %v", err)
	}
	container.Live.Worker = liveWorker

	testWorker, err := core.NewWorker(ctx, core.WorkerOptions{
		Config: &river.Config{
			Queues: map[string]river.QueueConfig{
				"default": {
					MaxWorkers: 100,
				},
			},
			Workers: workers,
		},
		DbURL:  container.Config.Database.TestWorker,
		Logger: container.Logger,
	})
	if err != nil {
		log.Fatalf("Failed to initialize test worker: %v", err)
	}
	container.Test.Worker = testWorker

	// Initialize the root command
	rootCmd := &cobra.Command{
		Use:   "payment",
		Short: fmt.Sprintf("%s service", container.Config.App.Name),
	}
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Add debug commands if in debug mode
	if mode == types.DebugMode {
		addDebugCommands(ctx, rootCmd, container)
	}

	addCommands(ctx, rootCmd, container, grpcServer)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}

func addDebugCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container) {
	rootCmd.AddCommand(cmd.NewDbSeedCmd(ctx, container.Logger,
		[]core.DBer{
			container.Live.DB.Primary,
			container.Test.DB.Primary,
		},
	))
	rootCmd.AddCommand(cmd.NewGenMigrationCmd(container.Logger, []core.DBer{container.Live.DB.Primary}))
}

func addCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container, grpcServer *core.GrpcServer) {
	rootCmd.AddCommand(cmd.NewDbMigrateCmd(ctx, container.Logger,
		[]core.DBer{
			container.Live.DB.Primary,
			container.Test.DB.Primary,
		},
		[]*core.Worker{
			container.Live.Worker,
			container.Test.Worker,
		},
	))
	rootCmd.AddCommand(cmd.NewStartCmd(ctx, container.Logger, nil, grpcServer, container.Live.Worker, container.Test.Worker, func() {
		errs := container.Close()
		if len(errs) > 0 {
			log.Fatalf("Failed to close application: %v", errs)
		}
	}))
}

func initGrpcServer(container *app.Container) (*core.GrpcServer, error) {
	grpcServer, err := core.NewGrpcServer(core.GrpcServerOptions{
		Host:          container.Config.Server.Host,
		Port:          container.Config.Server.Port,
		Logger:        container.Logger,
		ServerOptions: []grpc.ServerOption{},
	})
	if err != nil {
		return nil, err
	}

	paymentv1.RegisterPaymentServiceServer(
		grpcServer,
		v1.New(container),
	)
	reflection.Register(grpcServer)

	return grpcServer, nil
}
