package main

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/handler"
	v1 "autopilot/backends/api/internal/handler/v1"
	apimdw "autopilot/backends/api/internal/middleware"
	"autopilot/backends/api/internal/service"
	"autopilot/backends/api/internal/store"
	"autopilot/backends/internal/cmd"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/http/middleware"
	"autopilot/backends/internal/types"
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/cors"
	"github.com/riverqueue/river"
	"github.com/spf13/cobra"
)

var (
	//go:embed all:locales
	localesFS embed.FS

	//go:embed all:migrations
	migrationsFS embed.FS

	//go:embed all:templates
	templatesFS embed.FS
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
			Locales:    localesFS,
			Migrations: migrationsFS,
			Templates:  templatesFS,
		},
		MainFile: mainFile,
		Mode:     mode,
	})
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize the store manager
	storeManager := store.NewManager(container.DB.Primary)

	// Initialize the service manager
	serviceManager, err := service.NewManager(container, storeManager)
	if err != nil {
		log.Fatalf("Failed to initialize service manager: %v", err)
	}

	// Initialize the HTTP server
	httpServer, err := initHttpServer(container, serviceManager)
	if err != nil {
		log.Fatalf("Failed to initialize HTTP server: %v", err)
	}

	// Initialize the background worker
	worker, err := core.NewWorker(ctx, core.WorkerOptions{
		Config: &river.Config{
			Queues: map[string]river.QueueConfig{
				"default": {
					MaxWorkers: 100,
				},
			},
			Workers: service.AddWorkers(container, serviceManager),
		},
		DbURL:  container.Config.Database.Worker,
		Logger: container.Logger,
	})
	if err != nil {
		log.Fatalf("Failed to initialize background worker: %v", err)
	}
	container.Worker = worker

	// Initialize the root command
	rootCmd := &cobra.Command{
		Use:   "api",
		Short: fmt.Sprintf("%s API", container.Config.App.Name),
	}
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Add debug commands if in debug mode
	if mode == types.DebugMode {
		addDebugCommands(ctx, rootCmd, container, httpServer)
	}

	addCommands(ctx, rootCmd, container, httpServer)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}

func addDebugCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container, httpServer *core.HttpServer) {
	rootCmd.AddCommand(cmd.NewDbSeedCmd(ctx, container.Logger,
		[]core.DBer{
			container.DB.Primary,
		},
	))
	rootCmd.AddCommand(cmd.NewGenMigrationCmd(container.Logger, []core.DBer{container.DB.Primary}))
	rootCmd.AddCommand(cmd.NewGenOpenapiCmd(container.Logger, httpServer))
}

func addCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container, httpServer *core.HttpServer) {
	rootCmd.AddCommand(cmd.NewDbMigrateCmd(ctx, container.Logger,
		[]core.DBer{
			container.DB.Primary,
		},
		[]*core.Worker{
			container.Worker,
		},
	))
	rootCmd.AddCommand(cmd.NewStartCmd(ctx, container.Logger, httpServer, nil, container.Worker, nil, func() {
		errs := container.Close()
		if len(errs) > 0 {
			log.Fatalf("Failed to close application: %v", errs)
		}
	}))
}

func initHttpServer(container *app.Container, serviceManager *service.Manager) (*core.HttpServer, error) {
	httpServer, err := core.NewHttpServer(core.HttpServerOptions{
		Host:       container.Config.Server.Host,
		Port:       container.Config.Server.Port,
		I18nBundle: container.I18nBundle,
		Logger:     container.Logger,
		Mailer:     container.Mailer,
		Mode:       container.Mode,
		Middlewares: []func(http.Handler) http.Handler{
			middleware.Logger(container.Mode, container.Logger),
			cors.Handler(cors.Options{
				AllowedOrigins: container.Config.Server.CORS.AllowedOrigins,
				AllowedMethods: []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
				AllowedHeaders: []string{
					"Accept",
					"Authorization",
					"Content-Type",
					"X-Requested-With",
				},
				ExposedHeaders: []string{
					"Link",
					"X-Total-Count",
					"X-Request-ID",
					"X-RateLimit-Limit",
					"X-RateLimit-Remaining",
					"X-RateLimit-Reset",
				},
				AllowCredentials: true,
				MaxAge:           30,
			}),
			apimdw.WithRequestMetadata(),
			apimdw.WithContainer(container),
			apimdw.WithT(container.I18nBundle),
		},
	})
	if err != nil {
		return nil, err
	}

	huma.NewError = handler.NewCustomStatusError
	apiV1 := initApiV1(container, httpServer)
	err = v1.AddRoutes(container, apiV1, serviceManager)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}

func initApiV1(container *app.Container, httpServer *core.HttpServer) huma.API {
	apiConfig := huma.DefaultConfig("Autopilot API", "1.0.0")
	apiConfig.OpenAPI.Info.Description = "The Autopilot Platform API provides a comprehensive suite of endpoints."
	apiConfig.OpenAPI.Security = []map[string][]string{}
	apiConfig.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "API key authentication",
		},
	}
	apiConfig.DocsPath = ""
	apiConfig.OpenAPIPath = "/v1/openapi"
	apiConfig.CreateHooks = []func(huma.Config) huma.Config{
		func(cfg huma.Config) huma.Config {
			if cfg.OpenAPI.Extensions == nil {
				cfg.OpenAPI.Extensions = make(map[string]interface{})
			}
			cfg.OpenAPI.Extensions["x-sdk-id"] = "v1"

			return cfg
		},
	}
	apiConfig.OpenAPI.Tags = []*huma.Tag{
		{
			Name:        v1.TagIdentity,
			Description: "Identity management",
		},
	}

	scheme := "http://"
	if !strings.HasPrefix(container.Config.App.Domain, "localhost:") {
		scheme = "https://"
	}

	apiConfig.OpenAPI.Servers = []*huma.Server{
		{
			URL:         fmt.Sprintf("%s%s", scheme, container.Config.App.Domain),
			Description: "Production API",
		},
	}

	api := humachi.New(httpServer.Mux, apiConfig)
	httpServer.APIDocs = append(httpServer.APIDocs, api)

	return api
}
