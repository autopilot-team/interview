package app

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/seeders"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ContainerCache holds the cache connections for the container
type ContainerCache struct {
	// Identity is the identity cache connection
	Identity redis.UniversalClient

	// Payment is the payment cache connection
	Payment redis.UniversalClient
}

// PaymentDB holds the payment database connections
type PaymentDB struct {
	// Live is the live payment database connection
	Live core.DBer

	// Test is the test payment database connection
	Test core.DBer
}

// ContainerDB holds the database connections for the container
type ContainerDB struct {
	// Identity is the identity database connection
	Identity core.DBer

	// Payment is the payment database connection
	Payment PaymentDB
}

// ContainerFS holds the filesystems for the container
type ContainerFS struct {
	// Locales is the filesystem for the locales
	Locales core.FS

	// Migrations is the filesystem for the migrations
	Migrations core.FS

	// Templates is the filesystem for the templates
	Templates core.FS
}

// ContainerS3 holds the S3 connections for the container
type ContainerS3 struct {
	// Identity is the identity S3 connection
	Identity core.Storage

	// Payment is the payment S3 connection
	Payment core.Storage
}

// Container holds the application container
type Container struct {
	// Cache holds the cache connections for the container
	Cache ContainerCache

	// Config is the application configuration
	Config *Config

	// DB holds the database connections for the container
	DB ContainerDB

	// FS holds the filesystems for the container
	FS ContainerFS

	// I18nBundle is the i18n bundle
	I18nBundle *core.I18nBundle

	// Logger is a structured logger
	Logger *core.Logger

	// Mailer is the mailer
	Mailer core.Mailer

	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// RateLimiter is the Valkey client for rate limiting
	RateLimiter redis.UniversalClient

	// Storage is the S3 storage
	Storage ContainerS3

	// Turnstile is the Turnstile client
	Turnstile Turnstiler

	// Worker is the background worker
	Worker core.Worker

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

	// Worker is the background worker
	Worker core.Worker
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
		config.App.Observability.ApiEndpoint,
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

	// Initialize the i18n bundle
	i18nBundle, err := core.NewI18nBundle(opts.FS.Locales, "locales")
	if err != nil {
		return nil, err
	}

	// Initialize the mailer
	mailer, err := core.NewMail(core.MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       opts.Mode,
		PreviewData: map[string]map[string]any{
			"welcome": {
				"AppName":         config.App.Name,
				"AssetsURL":       config.App.AssetsURL,
				"Duration":        model.EmailVerificationDuration.Hours(),
				"Name":            "John Doe",
				"VerificationURL": fmt.Sprintf("%s/verify-email?token=01948450-988e-7976-a454-7163b6f1c6c6", config.App.DashboardURL),
			},
			"password_reset": {
				"AppName":   config.App.Name,
				"AssetsURL": config.App.AssetsURL,
				"Duration":  model.PasswordResetDuration.Hours(),
				"Name":      "John Doe",
				"ResetURL":  fmt.Sprintf("%s/reset-password?token=01948450-988e-7976-a454-7163b6f1c6c6", config.App.DashboardURL),
			},
		},
		SmtpUrl: config.App.Mailer.SmtpUrl,
		TemplateOptions: &core.MailTemplateOptions{
			Dir: "templates",
			ExtraFuncs: []template.FuncMap{
				{
					"currentYear": func() string {
						return strconv.Itoa(time.Now().Year())
					},
				},
			},
			FS:     opts.FS.Templates,
			Layout: "layouts/transactional",
		},
	})
	if err != nil {
		return nil, err
	}

	// Initialize RateLimiter (Valkey client)
	rateLimiter, err := core.NewRedis(ctx, core.RedisOptions{
		URL:       config.App.RateLimiter.ValkeyURLs,
		IsCluster: true,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to api rate limiter", "error", err)
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		return rateLimiter.Close()
	})

	// Initialize the identity cache
	identityCache, err := core.NewRedis(ctx, core.RedisOptions{
		URL:       config.Identity.Cache.ValkeyURLs,
		IsCluster: true,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to identity cache", "error", err)
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		return identityCache.Close()
	})

	// Initialize the identity database
	identityDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "identity",
		Logger:       logger,
		MainFile:     opts.MainFile,
		Mode:         opts.Mode,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Identity.Database.PrimaryWriter,
		ReaderURLs:   config.Identity.Database.PrimaryReaders,
		Seeder:       seeders.Identity,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		identityDB.Close()
		return nil
	})

	// Initialize the identity S3 storage
	identityStorage, err := core.NewStorage(ctx, core.StorageOptions{
		Logger:          logger,
		Endpoint:        config.Identity.Storage.Endpoint,
		Region:          config.Identity.Storage.Region,
		AccessKeyID:     config.Identity.Storage.AccessKeyID,
		SecretAccessKey: config.Identity.Storage.SecretAccessKey,
		Bucket:          config.Identity.Storage.Bucket,
		UsePathStyle:    config.Identity.Storage.UsePathStyle,
	})
	if err != nil {
		return nil, err
	}

	// Initialize the payment databases
	livePaymentDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "payment",
		Logger:       logger,
		MainFile:     opts.MainFile,
		Mode:         opts.Mode,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Payment.Database.LivePrimaryWriter,
		ReaderURLs:   config.Payment.Database.LivePrimaryReaders,
		Seeder:       seeders.Payment,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		livePaymentDB.Close()
		return nil
	})

	testPaymentDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "payment",
		Logger:       logger,
		MainFile:     opts.MainFile,
		Mode:         opts.Mode,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Payment.Database.TestPrimaryWriter,
		ReaderURLs:   config.Payment.Database.TestPrimaryReaders,
		Seeder:       seeders.Payment,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		testPaymentDB.Close()
		return nil
	})

	return &Container{
		Cache: ContainerCache{
			Identity: identityCache,
		},
		CleanUp: cleanUp,
		Config:  config,
		DB: ContainerDB{
			Identity: identityDB,
			Payment: PaymentDB{
				Live: livePaymentDB,
				Test: testPaymentDB,
			},
		},
		FS: ContainerFS{
			Locales:    opts.FS.Locales,
			Migrations: opts.FS.Migrations,
			Templates:  opts.FS.Templates,
		},
		I18nBundle:  i18nBundle,
		Logger:      logger,
		Mailer:      mailer,
		Mode:        opts.Mode,
		RateLimiter: rateLimiter,
		Storage: ContainerS3{
			Identity: identityStorage,
		},
		Turnstile: NewTurnstile(config.App.Cloudflare.TurnstileSecretKey),
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
