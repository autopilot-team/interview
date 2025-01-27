package app

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"

	"github.com/redis/go-redis/v9"
)

// ContainerDB holds the database connections for the container
type ContainerDB struct {
	// Primary is the primary database connection
	Primary core.DBer
}

// ContainerFS holds the filesystems for the container
type ContainerFS struct {
	// Migrations is the filesystem for the migrations
	Migrations core.FS
}

// ContainerInfra holds the infrastructure environment
type ContainerInfra struct {
	// DB is the database connection
	DB ContainerDB

	// Worker is the background worker
	Worker *core.Worker
}

// Container holds the application container
type Container struct {
	// Cache is the Redis client
	Cache *redis.Client

	// Config is the application configuration
	Config *Config

	// FS holds the filesystems for the container
	FS ContainerFS

	// I18nBundle is the i18n bundle
	I18nBundle *core.I18nBundle

	// Logger is a structured logger
	Logger *core.Logger

	// Mailer is the mailer
	Mailer *core.Mailer

	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// Live holds the live infrastructure environment
	Live *ContainerInfra

	// Test holds the test infrastructure environment
	Test *ContainerInfra

	// CleanUp is a list of functions to clean up the container
	CleanUp []func() error
}

// ContainerOpts holds the options for the container
type ContainerOpts struct {
	// FS holds the filesystems for the container
	FS ContainerFS

	// MainFile is the main file for the application
	MainFile string

	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// LiveWorker is the background worker for the live environment
	LiveWorker *core.Worker

	// TestWorker is the background worker for the test environment
	TestWorker *core.Worker
}

// NewContainer creates a new Container instance
func NewContainer(ctx context.Context, opts ContainerOpts) (*Container, error) {
	// Initialize the clean up functions
	cleanUp := make([]func() error, 0)

	// Initialize the configuration
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	// Initialize the tracer
	tracerShutdown, err := core.NewTracer(
		ctx,
		config.App.Environment,
		config.Observability.AxiomApiToken,
		config.App.Service,
		config.App.Version,
	)
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		return tracerShutdown(ctx)
	})

	// Initialize the logger
	logger := core.NewLogger(core.LoggerOptions{Mode: opts.Mode})

	// Initialize the primary databases
	livePrimaryDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     opts.MainFile,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Database.LivePrimaryWriter,
		ReaderURLs:   config.Database.LivePrimaryReaders,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		livePrimaryDB.Close()
		return nil
	})

	testPrimaryDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     opts.MainFile,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Database.TestPrimaryWriter,
		ReaderURLs:   config.Database.TestPrimaryReaders,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		livePrimaryDB.Close()
		return nil
	})

	return &Container{
		Config: config,
		FS: ContainerFS{
			Migrations: opts.FS.Migrations,
		},
		Logger: logger,
		Mode:   opts.Mode,
		Live: &ContainerInfra{
			DB: ContainerDB{
				Primary: livePrimaryDB,
			},
		},
		Test: &ContainerInfra{
			DB: ContainerDB{
				Primary: testPrimaryDB,
			},
		},
		CleanUp: cleanUp,
	}, nil
}

// Close cleans up the container
func (c *Container) Close() []error {
	var errs []error
	for i := len(c.CleanUp) - 1; i >= 0; i-- {
		if err := c.CleanUp[i](); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
