package app

import (
	"autopilot/backends/api/seeders"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"html/template"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ContainerDB holds the database connections for the container
type ContainerDB struct {
	// Primary is the primary database connection
	Primary core.DBer
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

// Container holds the application container
type Container struct {
	// Cache is the Redis client
	Cache *redis.Client

	// Config is the application configuration
	Config *Config

	// DB is the database connection
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

	// Turnstile is the Turnstile client
	Turnstile Turnstiler

	// Worker is the background worker
	Worker *core.Worker

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
	Worker *core.Worker
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
		PreviewData: map[string]map[string]interface{}{
			"welcome": {
				"AppName":         config.App.Name,
				"Duration":        (24 * time.Hour).Hours(),
				"Name":            "John Doe",
				"VerificationURL": "http://localhost:3000/verify-email?token=01948450-988e-7976-a454-7163b6f1c6c6",
			},
		},
		SmtpUrl: config.Mailer.SmtpUrl,
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

	// Initialize the primary database
	primaryDB, err := core.NewDB(ctx, core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     opts.MainFile,
		MigrationsFS: opts.FS.Migrations,
		WriterURL:    config.Database.PrimaryWriter,
		ReaderURLs:   config.Database.PrimaryReaders,
		Seeder:       seeders.Primary,
	})
	if err != nil {
		return nil, err
	}
	cleanUp = append(cleanUp, func() error {
		primaryDB.Close()
		return nil
	})

	return &Container{
		CleanUp: cleanUp,
		Config:  config,
		DB: ContainerDB{
			Primary: primaryDB,
		},
		FS: ContainerFS{
			Locales:    opts.FS.Locales,
			Migrations: opts.FS.Migrations,
			Templates:  opts.FS.Templates,
		},
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mailer:     mailer,
		Mode:       opts.Mode,
		Turnstile:  NewTurnstile(config.Cloudflare.TurnstileSecretKey),
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
