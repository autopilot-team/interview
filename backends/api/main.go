package main

import (
	"autopilot/backends/api/internal"
	"autopilot/backends/api/internal/identity"
	identityhandlerv1 "autopilot/backends/api/internal/identity/handler/v1"
	identitysvc "autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/internal/payment"
	paymenthandlerv1 "autopilot/backends/api/internal/payment/handler/v1"
	paymentsvc "autopilot/backends/api/internal/payment/service"
	"autopilot/backends/api/pkg/app"
	apimdw "autopilot/backends/api/pkg/middleware"
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
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/cors"
	"github.com/riverqueue/river"
	"github.com/spf13/cobra"
	"riverqueue.com/riverui"
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

	// Initialize the application context
	ctx := context.Background()

	// Initialize the application container
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

	mods, err := initAutopilot(ctx, container)
	if err != nil {
		log.Fatalf("Failed to initialize autopilot: %v", err)
	}

	// Initialize the HTTP server
	httpServer, err := initHTTPServer(container, mods)
	if err != nil {
		log.Fatalf("Failed to initialize HTTP server: %v", err)
	}

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

	addCommands(ctx, rootCmd, container, httpServer, mods)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}

func initAutopilot(ctx context.Context, container *app.Container) (*internal.Module, error) {
	identityMod, err := identity.New(ctx, container)
	if err != nil {
		log.Fatalf("Failed to initialize identity module: %v", err)
	}

	paymentMod, err := payment.New(ctx, container)
	if err != nil {
		log.Fatalf("Failed to initialize payment module: %v", err)
	}

	mods := &internal.Module{
		Identity: identityMod,
		Payment:  paymentMod,
	}

	// Initialize background workers
	workers := river.NewWorkers()
	identitysvc.AddWorkers(container, workers, identityMod.Service)
	paymentsvc.AddWorkers(container, workers, paymentMod.Service)

	periodicJobs := slices.Clone(identitysvc.AddPeriodicJobs(container, identityMod.Service))
	periodicJobs = append(periodicJobs, paymentsvc.AddPeriodicJobs(container, paymentMod.Service)...)
	worker, err := core.NewWorker(ctx, core.WorkerOptions{
		Config: &river.Config{
			PeriodicJobs: periodicJobs,
			Queues: map[string]river.QueueConfig{
				river.QueueDefault: {MaxWorkers: 32},
			},
			Workers: workers,
		},
		DBURL:  container.Config.App.Database.Worker,
		Logger: container.Logger,
	})
	if err != nil {
		log.Fatalf("Failed to initialize background worker: %v", err)
	}
	container.Worker = worker

	return mods, nil
}

func initHTTPServer(container *app.Container, mods *internal.Module) (*core.HTTPServer, error) {
	var injectCfg apimdw.InjectCountryConfig
	if container.Mode == types.DebugMode {
		injectCfg = apimdw.InjectCountryConfig{
			Enable:  true,
			Country: "SG",
		}
	}
	httpServer, err := core.NewHTTPServer(core.HTTPServerOptions{
		Host:       container.Config.App.Server.Host,
		Port:       container.Config.App.Server.Port,
		I18nBundle: container.I18nBundle,
		Logger:     container.Logger,
		Mailer:     container.Mailer,
		Mode:       container.Mode,
		Middlewares: []func(http.Handler) http.Handler{
			middleware.WithDebug(container.Mode),
			middleware.Logger(container.Mode, container.Logger),
			middleware.WithOperationMode([]string{
				container.Config.App.DashboardURL,
			}),
			cors.Handler(cors.Options{
				AllowedOrigins: container.Config.App.CORS.AllowedOrigins,
				AllowedMethods: []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
				AllowedHeaders: []string{
					"Accept",
					"Authorization",
					"Content-Type",
					"X-Api-Key",
					"X-Debug",
					"X-File-Name",
					"X-Operation-Mode",
					"X-Requested-With",
					apimdw.ActiveEntityHeader,
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
			apimdw.WithInjectCountry(injectCfg),
			apimdw.WithRateLimit(container, apimdw.DefaultRateLimitConfig()),
			apimdw.WithRequestMetadata(),
			apimdw.WithContainer(container),
			apimdw.WithActiveEntity(),
			apimdw.WithT(container.I18nBundle),
		},
	})
	if err != nil {
		return nil, err
	}

	// Initialize Queue UI server
	if container.Mode == types.DebugMode {
		queueUIServer, err := riverui.NewServer(&riverui.ServerOpts{
			Client: container.Worker.GetClient(),
			DB:     container.Worker.GetDBPool(),
			Logger: container.Logger,
			Prefix: "/queue",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize River UI server: %w", err)
		}

		// Start and mount the queue UI
		if err := queueUIServer.Start(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to start River UI server: %w", err)
		}

		httpServer.Mount("/", queueUIServer)
	}

	httpServer.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name": "%s", "version": "%s"}`, container.Config.App.Name, container.Config.App.Version)
	})

	apiV1 := initAPIV1(container, httpServer)
	authenticator := identity.NewAuthentication(container, apiV1, mods.Identity.Service)

	err = identityhandlerv1.AddRoutes(container, apiV1, mods.Identity.Service, authenticator)
	if err != nil {
		return nil, err
	}

	err = paymenthandlerv1.AddRoutes(container, apiV1, mods.Payment.Service, authenticator)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}

func initAPIV1(container *app.Container, httpServer *core.HTTPServer) huma.API {
	apiConfig := huma.DefaultConfig("Autopilot API", "1.0.0")
	apiConfig.OpenAPI.Info.Description = "The Autopilot Platform API provides a comprehensive suite of endpoints."
	apiConfig.OpenAPI.Security = []map[string][]string{}
	apiConfig.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"API Key Authentication": {
			Type:        "apiKey",
			In:          "header",
			Name:        "X-Api-Key",
			Description: "API key authentication using X-Api-Key header",
		},
	}
	apiConfig.DocsPath = ""
	apiConfig.OpenAPIPath = identityhandlerv1.BasePath("/openapi")
	apiConfig.CreateHooks = []func(huma.Config) huma.Config{
		func(cfg huma.Config) huma.Config {
			if cfg.OpenAPI.Extensions == nil {
				cfg.OpenAPI.Extensions = make(map[string]any)
			}
			cfg.OpenAPI.Extensions["x-sdk-id"] = strings.ReplaceAll(identityhandlerv1.BasePath(""), "/", "")

			return cfg
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

func addCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container, httpServer *core.HTTPServer, _ *internal.Module) {
	databases := []core.DBer{
		container.DB.Identity,
		container.DB.Payment.Live,
		container.DB.Payment.Test,
	}
	rootCmd.AddCommand(cmd.NewDBMigrateCmd(ctx, container.Logger, databases,
		[]core.Worker{
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

func addDebugCommands(ctx context.Context, rootCmd *cobra.Command, container *app.Container, httpServer *core.HTTPServer) {
	databases := []core.DBer{
		container.DB.Identity,
		container.DB.Payment.Live,
		container.DB.Payment.Test,
	}
	rootCmd.AddCommand(cmd.NewDBSeedCmd(ctx, container.Logger, databases))
	rootCmd.AddCommand(cmd.NewGenMigrationCmd(container.Logger, databases))
	rootCmd.AddCommand(cmd.NewGenOpenapiCmd(container.Logger, httpServer))
}
