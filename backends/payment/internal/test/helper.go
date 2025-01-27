package test

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"autopilot/backends/payment/internal/app"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// NewMockContainer creates a new mock container for testing
func NewMockContainer(t *testing.T) *app.Container {
	ctx := context.Background()
	cleanUp := make([]func() error, 0)
	mode := types.DebugMode

	// Initialize the configuration
	config, err := app.NewConfig()
	require.NoError(t, err)

	// Initialize the logger
	logger := core.NewLogger(core.LoggerOptions{Mode: mode})

	// Initialize the local filesystem
	localFS, err := core.NewLocalFS("./backends/api")
	require.NoError(t, err)

	projectRoot, err := core.FindProjectRoot()
	require.NoError(t, err)
	mainFile := filepath.Join(projectRoot, "backends/api/internal/app/main.go")

	// Initialize the live primary database
	livePrimaryDB, livePrimaryDBCleanUp, err := NewDB(ctx, "payment_live", core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     mainFile,
		MigrationsFS: localFS,
	})
	require.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		livePrimaryDBCleanUp()
		return nil
	})

	// Initialize the test primary database
	testPrimaryDB, testPrimaryDBCleanUp, err := NewDB(ctx, "payment_test", core.DBOptions{
		Identifier:   "primary",
		Logger:       logger,
		MainFile:     mainFile,
		MigrationsFS: localFS,
	})
	require.NoError(t, err)
	cleanUp = append(cleanUp, func() error {
		testPrimaryDBCleanUp()
		return nil
	})

	// Initialize the container
	container := &app.Container{
		Config: config,
		FS: app.ContainerFS{
			Migrations: localFS,
		},
		Logger: logger,
		Mode:   mode,
		Live: &app.ContainerInfra{
			DB: app.ContainerDB{
				Primary: livePrimaryDB,
			},
		},
		Test: &app.ContainerInfra{
			DB: app.ContainerDB{
				Primary: testPrimaryDB,
			},
		},
		CleanUp: cleanUp,
	}

	return container
}
