package core

import (
	"autopilot/backends/internal/types"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	postgresURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

//go:embed all:testdata
var testMigrationsFS embed.FS

// generateTestDBName generates a unique database name using timestamp and random suffix
func generateTestDBName() string {
	timestamp := time.Now().UTC().Format("20060102150405")
	randomSuffix := rand.Intn(10000)

	return fmt.Sprintf("test_db_%s_%04d", timestamp, randomSuffix)
}

// getTestDBURLs returns writer and reader URLs for the test database
func getTestDBURLs(dbName string) (string, string) {
	return fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName),
		fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
}

func setupTestDB(t *testing.T) (string, string, func()) {
	ctx := context.Background()
	dbName := generateTestDBName()
	writerURL, readerURL := getTestDBURLs(dbName)

	// Connect to postgres database to create/drop test database
	pool, err := pgxpool.New(ctx, postgresURL)
	require.NoError(t, err, "Failed to connect to postgres database")
	defer pool.Close()

	// Create test database
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err, "Failed to create test database")

	// Connect to test database to create schema
	testPool, err := pgxpool.New(ctx, writerURL)
	require.NoError(t, err, "Failed to connect to test database")
	defer testPool.Close()

	// Create test schema and tables if needed
	_, err = testPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err, "Failed to create test schema")

	// Return cleanup function
	return writerURL, readerURL, func() {
		pool, err := pgxpool.New(ctx, postgresURL)
		if err != nil {
			t.Logf("Failed to connect to postgres database during cleanup: %v", err)
			return
		}
		defer pool.Close()

		// Terminate all connections to the database before dropping it
		_, err = pool.Exec(ctx, fmt.Sprintf(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = '%s' AND pid <> pg_backend_pid()
		`, dbName))
		if err != nil {
			t.Logf("Failed to terminate connections to test database: %v", err)
		}

		_, err = pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if err != nil {
			t.Logf("Failed to drop test database during cleanup: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Clean up test-generated migration files
	_, mainFile, _, ok := runtime.Caller(0)
	if ok {
		dir := filepath.Dir(mainFile)
		migrationsDir := filepath.Join(dir, "testdata/migrations/test")
		entries, err := os.ReadDir(migrationsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.Name() != "00000000000000_init.sql" {
					os.Remove(filepath.Join(migrationsDir, entry.Name()))
				}
			}
		}
	}

	os.Exit(code)
}

func TestNewDB(t *testing.T) {
	t.Parallel()
	writerURL, readerURL, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	tests := []struct {
		name    string
		opts    DBOptions
		wantErr bool
	}{
		{
			name: "should initialize successfully with writer only",
			opts: DBOptions{
				Mode:          "debug",
				WriterURL:     writerURL,
				Identifier:    "test",
				Logger:        logger,
				MainFile:      mainFile,
				MigrationsDir: "testdata/migrations",
				MigrationsFS:  testMigrationsFS,
			},
			wantErr: false,
		},
		{
			name: "should initialize successfully with writer and reader",
			opts: DBOptions{
				Mode:          "debug",
				WriterURL:     writerURL,
				ReaderURLs:    []string{readerURL},
				Identifier:    "test",
				Logger:        logger,
				MainFile:      mainFile,
				MigrationsDir: "testdata/migrations",
				MigrationsFS:  testMigrationsFS,
			},
			wantErr: false,
		},
		{
			name: "should reject invalid writer URL",
			opts: DBOptions{
				Mode:          "debug",
				WriterURL:     "invalid://url",
				Identifier:    "test",
				Logger:        logger,
				MainFile:      mainFile,
				MigrationsDir: "testdata/migrations",
				MigrationsFS:  testMigrationsFS,
			},
			wantErr: true,
		},
		{
			name: "should reject invalid reader URL",
			opts: DBOptions{
				Mode:          "debug",
				WriterURL:     writerURL,
				ReaderURLs:    []string{"invalid://url"},
				Identifier:    "test",
				Logger:        logger,
				MainFile:      mainFile,
				MigrationsDir: "testdata/migrations",
				MigrationsFS:  testMigrationsFS,
			},
			wantErr: true,
		},
		{
			name: "should reject missing logger",
			opts: DBOptions{
				Mode:          "debug",
				WriterURL:     writerURL,
				Identifier:    "test",
				MainFile:      mainFile,
				MigrationsDir: "testdata/migrations",
				MigrationsFS:  testMigrationsFS,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(context.Background(), tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					db.Close()
				}
			}
		})
	}
}

func TestDB_Reader(t *testing.T) {
	t.Parallel()
	writerURL, readerURL, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	t.Run("should fallback to writer when no readers available", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)
		defer db.Close()

		reader := db.Reader()
		assert.Equal(t, db.Writer(), reader)
	})

	t.Run("should return reader when available", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			ReaderURLs:    []string{readerURL},
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)
		defer db.Close()

		reader := db.Reader()
		assert.NotNil(t, reader)
	})
}

func TestDB_Writer(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Identifier:    "test",
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
	})
	require.NoError(t, err)
	defer db.Close()

	writer := db.Writer()
	assert.NotNil(t, writer)
}

func TestDB_WithTx(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Identifier:    "test",
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
	})
	require.NoError(t, err)
	defer db.Close()

	t.Run("successful transaction", func(t *testing.T) {
		err := db.WithTx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
			// Insert test data
			_, err := tx.Exec("INSERT INTO test_table (name) VALUES ($1)", "test")
			return err
		})
		assert.NoError(t, err)

		// Verify data was inserted
		var count int
		err = db.Writer().QueryRow("SELECT COUNT(*) FROM test_table WHERE name = $1", "test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("rollback on error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		err := db.WithTx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
			// Insert test data
			_, err := tx.Exec("INSERT INTO test_table (name) VALUES ($1)", "rollback_test")
			if err != nil {
				return err
			}
			return expectedErr
		})
		assert.ErrorIs(t, err, expectedErr)

		// Verify data was not inserted
		var count int
		err = db.Writer().QueryRow("SELECT COUNT(*) FROM test_table WHERE name = $1", "rollback_test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("rollback on panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = db.WithTx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
				// Insert test data
				_, err := tx.Exec("INSERT INTO test_table (name) VALUES ($1)", "panic_test")
				if err != nil {
					return err
				}
				panic("test panic")
			})
		})

		// Verify data was not inserted
		var count int
		err = db.Writer().QueryRow("SELECT COUNT(*) FROM test_table WHERE name = $1", "panic_test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestDB_Close(t *testing.T) {
	t.Parallel()
	writerURL, readerURL, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	t.Run("close writer only", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)
		db.Close()
	})

	t.Run("close writer and readers", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			ReaderURLs:    []string{readerURL},
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)
		db.Close()
	})
}

func TestDB_WithTxTimeout(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	t.Run("successful transaction within timeout", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)

		err = db.WithTxTimeout(context.Background(), 5*time.Second, func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.Exec("INSERT INTO test_table (name) VALUES ($1)", "timeout_test")
			return err
		})
		assert.NoError(t, err)

		var count int
		err = db.Writer().QueryRow("SELECT COUNT(*) FROM test_table WHERE name = $1", "timeout_test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("transaction timeout", func(t *testing.T) {
		db, err := NewDB(context.Background(), DBOptions{
			Mode:          "debug",
			WriterURL:     writerURL,
			Identifier:    "test",
			Logger:        logger,
			MainFile:      mainFile,
			MigrationsDir: "testdata/migrations",
			MigrationsFS:  testMigrationsFS,
		})
		require.NoError(t, err)

		err = db.WithTxTimeout(context.Background(), 50*time.Millisecond, func(ctx context.Context, tx *sqlx.Tx) error {
			// Sleep first to ensure we hit the timeout
			time.Sleep(100 * time.Millisecond)

			_, err := tx.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ($1)", "timeout_fail")
			if err != nil {
				return err
			}

			return ctx.Err()
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded", "expected transaction timeout error")

		var count int
		err = db.Writer().QueryRow("SELECT COUNT(*) FROM test_table WHERE name = $1", "timeout_fail").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestDbLogger_Log(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dbLogger := &DBLogger{
		Logger: logger,
		mode:   types.DebugMode,
	}

	t.Run("logs query in debug mode", func(t *testing.T) {
		data := map[string]any{
			"sql":  "SELECT * FROM users",
			"time": time.Duration(100 * time.Millisecond),
		}
		dbLogger.Log(context.Background(), tracelog.LogLevelInfo, "Query", data)
	})

	t.Run("ignores non-query messages", func(t *testing.T) {
		data := map[string]any{
			"message": "test",
		}
		dbLogger.Log(context.Background(), tracelog.LogLevelInfo, "NotQuery", data)
	})
}

func TestFormatQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		query    string
		contains []string
	}{
		{
			name:     "basic select",
			query:    "SELECT * FROM users WHERE id = 1",
			contains: []string{"SELECT", "FROM", "WHERE"},
		},
		{
			name:     "complex query",
			query:    "SELECT u.*, COALESCE(j.data, '{}') AS job_data FROM users u LEFT JOIN jobs j ON j.user_id = u.id",
			contains: []string{"SELECT", "COALESCE", "LEFT JOIN", "ON"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatQuery(tt.query)
			for _, keyword := range tt.contains {
				assert.Contains(t, result, keyword)
			}
		})
	}
}

func TestFormatArg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		arg  any
	}{
		{
			name: "string",
			arg:  "test",
		},
		{
			name: "integer",
			arg:  42,
		},
		{
			name: "float",
			arg:  3.14,
		},
		{
			name: "bool",
			arg:  true,
		},
		{
			name: "time",
			arg:  time.Now(),
		},
		{
			name: "nil",
			arg:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatArg(tt.arg)
			assert.NotEmpty(t, result)
		})
	}
}

func TestDB_Name(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
		Identifier:    "test_db",
	})
	require.NoError(t, err)
	defer db.Close()

	assert.Equal(t, "test_db", db.Identifier())
}

func TestDB_Migrate(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Identifier:    "test",
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
	})
	require.NoError(t, err)
	defer db.Close()

	err = db.Migrate(context.Background())
	assert.NoError(t, err)
}

func TestDB_GenMigration(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Identifier:    "test",
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
	})
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid migration name",
			input:   "create_users_table",
			wantErr: false,
		},
		{
			name:    "empty migration name",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.GenMigration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDB_MigrateStatus(t *testing.T) {
	t.Parallel()
	writerURL, _, cleanup := setupTestDB(t)
	defer cleanup()

	logger := slog.Default()
	_, mainFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller information")

	db, err := NewDB(context.Background(), DBOptions{
		Mode:          "debug",
		WriterURL:     writerURL,
		Identifier:    "test",
		Logger:        logger,
		MainFile:      mainFile,
		MigrationsDir: "testdata/migrations",
		MigrationsFS:  testMigrationsFS,
	})
	require.NoError(t, err)
	defer db.Close()

	// Check initial status
	pending, err := db.MigrateStatus(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, pending, int64(0))

	// Run migrations
	err = db.Migrate(context.Background())
	assert.NoError(t, err)

	// Check status after migrations
	pending, err = db.MigrateStatus(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), pending)
}
