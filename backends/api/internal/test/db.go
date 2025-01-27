package test

import (
	"autopilot/backends/internal/core"
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDB creates a new test database from the template database
func NewDB(ctx context.Context, dbName string, dbOpts core.DBOptions) (core.DBer, func(), error) {
	postgresDbUrl := fmt.Sprintf(core.DbUrlTemplate, "postgres")
	pool, err := pgxpool.New(ctx, postgresDbUrl)
	if err != nil {
		return nil, nil, err
	}
	defer pool.Close()

	testDbName := generateTestDbName(dbName)
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE template_%s", testDbName, dbName))
	if err != nil {
		return nil, nil, err
	}

	// Return cleanup function
	deleteTestDB := func() {
		cleanupCtx := context.Background()
		pool, err := pgxpool.New(cleanupCtx, postgresDbUrl)
		if err != nil {
			return
		}
		defer pool.Close()

		// Terminate all connections to the database
		_, _ = pool.Exec(cleanupCtx, fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, testDbName))

		// Drop the test database
		_, _ = pool.Exec(cleanupCtx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDbName))
	}

	dbUrl := fmt.Sprintf(core.DbUrlTemplate, testDbName)
	dbOpts.WriterURL = dbUrl
	dbOpts.ReaderURLs = []string{dbUrl}
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
func generateTestDbName(dbName string) string {
	randomSuffix := rand.Intn(1_000_000_000_000)

	return fmt.Sprintf("%s_test_%010d", dbName, randomSuffix)
}
