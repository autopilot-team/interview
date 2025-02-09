package testutil

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/handler"
	"autopilot/backends/api/internal/middleware"
	"autopilot/backends/api/internal/service"
	"autopilot/backends/api/internal/store"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTurnstile is a mock implementation of app.Turnstiler using testify/mock
type MockTurnstile struct {
	mock.Mock
}

// Verify implements app.Turnstiler
func (m *MockTurnstile) Verify(ctx context.Context, token string, action string) (bool, error) {
	args := m.Called(ctx, token, action)
	return args.Bool(0), args.Error(1)
}

// NewMockTurnstile creates a new MockTurnstile instance
func NewMockTurnstile() *MockTurnstile {
	return &MockTurnstile{}
}

// NewMocks creates a new mock container for testing
func NewMocks(t *testing.T) (humatest.TestAPI, *app.Container, *service.Manager, *store.Manager) {
	ctx := context.Background()
	cleanUp := make([]func() error, 0)
	mode := types.DebugMode

	// Initialize the configuration
	config, err := app.NewConfig()
	assert.NoError(t, err)

	// Initialize the logger
	logger := core.NewLogger(core.LoggerOptions{Mode: mode})

	// Initialize the local filesystem
	localFS, err := core.NewLocalFS("./backends/api")
	assert.NoError(t, err)

	// Initialize the i18n bundle
	i18nBundle, err := core.NewI18nBundle(localFS, "locales")
	assert.NoError(t, err)

	// Initialize the mailer
	mailer, err := core.NewMail(core.MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       mode,
		PreviewData: map[string]map[string]interface{}{
			"welcome": {
				"AppName":         config.App.Name,
				"AssetsURL":       config.App.AssetsURL,
				"Duration":        (24 * time.Hour).Hours(),
				"Name":            "John Doe",
				"VerificationURL": fmt.Sprintf("%s/verify-email?token=01948450-988e-7976-a454-7163b6f1c6c6", config.App.DashboardURL),
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
			FS:     localFS,
			Layout: "layouts/transactional",
		},
	})
	assert.NoError(t, err)

	projectRoot, err := core.FindProjectRoot()
	assert.NoError(t, err)
	mainFile := filepath.Join(projectRoot, "backends/api/internal/app/main.go")

	// Initialize the primary database
	primaryDB, primaryDBCleanUp, err := NewDB(ctx, "api", core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     mainFile,
		MigrationsFS: localFS,
	})
	assert.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		primaryDBCleanUp()
		return nil
	})

	// Initialize the container
	container := &app.Container{
		CleanUp: cleanUp,
		Config:  config,
		DB: app.ContainerDB{
			Primary: primaryDB,
		},
		FS: app.ContainerFS{
			Locales:    localFS,
			Migrations: localFS,
			Templates:  localFS,
		},
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       mode,
		Mailer:     mailer,
		Turnstile:  &MockTurnstile{},
	}

	// Initialize the service and store manager
	storeManager := store.NewManager(container.DB.Primary)
	serviceManager, err := service.NewManager(container, storeManager)
	assert.NoError(t, err)

	// Initialize the background worker
	worker, err := core.NewWorker(ctx, core.WorkerOptions{
		Config: &river.Config{
			Workers: service.AddWorkers(container, serviceManager),
		},
		DbURL:  primaryDB.Options().WriterURL,
		Logger: logger,
	})
	assert.NoError(t, err)
	container.Worker = worker

	_, api := humatest.New(t)
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachRequestMetadata(ctx.Context(), ctx.Header("X-Forwarded-For"), ctx.RemoteAddr(), ctx.Header("User-Agent")))
			next(ctx)
		},
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachContainer(ctx.Context(), container))
			next(ctx)
		},
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithContext(ctx, middleware.AttachT(ctx.Context(), container.I18nBundle, ctx.Header("Accept-Language"), ctx.Query("locale")))
			next(ctx)
		},
	)
	huma.NewError = handler.NewCustomStatusError

	return api, container, serviceManager, storeManager
}
