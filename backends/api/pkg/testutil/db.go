package testutil

import (
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDB creates a new test database from the template database
func NewDB(ctx context.Context, dbName string, dbOpts core.DBOptions) (core.DBer, func(), error) {
	postgresDBURL := fmt.Sprintf(core.DBURLTemplate, "postgres")
	pool, err := pgxpool.New(ctx, postgresDBURL)
	if err != nil {
		return nil, nil, err
	}
	defer pool.Close()

	testDBName := generateTestDBName(dbName)
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE template_%s", testDBName, dbName))
	if err != nil {
		return nil, nil, err
	}

	// Return cleanup function
	deleteTestDB := func() {
		cleanupCtx := context.Background()
		pool, err := pgxpool.New(cleanupCtx, postgresDBURL)
		if err != nil {
			return
		}
		defer pool.Close()

		// Terminate all connections to the database
		_, _ = pool.Exec(cleanupCtx, fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, testDBName))

		// Drop the test database
		_, _ = pool.Exec(cleanupCtx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	}

	dbURL := fmt.Sprintf(core.DBURLTemplate, testDBName)
	dbOpts.WriterURL = dbURL
	dbOpts.ReaderURLs = []string{dbURL}
	dbOpts.Mode = types.DebugMode
	db, err := core.NewDB(ctx, dbOpts)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
		deleteTestDB()
	}, nil
}

// generateTestDbName generates a unique database name using timestamp and random suffix
func generateTestDBName(dbName string) string {
	randomSuffix := rand.Intn(1_000_000_000_000)

	return fmt.Sprintf("%s_test_%010d", dbName, randomSuffix)
}
