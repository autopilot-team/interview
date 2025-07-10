package core

import (
	"autopilot/backends/internal/types"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
)

const (
	DbUrlTemplate = "postgres://postgres:postgres@localhost:5432/%s?sslmode=disable"
)

// Querier is an interface for database queries
type Querier interface {
	// Exec executes a query without returning any rows
	Exec(query string, args ...any) (sql.Result, error)

	// ExecContext executes a query without returning any rows
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// PrepareContext creates a prepared statement for later queries or executions
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Query executes a query that returns rows, typically a SELECT
	Query(query string, args ...any) (*sql.Rows, error)

	// QueryContext executes a query that returns rows, typically a SELECT
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row
	QueryRow(query string, args ...any) *sql.Row

	// QueryRowContext executes a query that is expected to return at most one row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DBer is a database connection pool
type DBer interface {
	// Close closes the database connections
	Close()

	// GenMigration creates a new migration file
	GenMigration(name string) error

	// HealthCheck checks the health of the database connections
	HealthCheck(ctx context.Context) error

	// Identifier is the identifier of the database which is used for displaying and identifying the migrations directory
	Identifier() string

	// Migrate runs all pending migrations
	Migrate(ctx context.Context) error

	// Status returns the number of pending migrations
	MigrateStatus(ctx context.Context) (int64, error)

	// DBName returns the name of the database
	Name() string

	// Options returns the options of the database
	Options() DBOptions

	// Reader returns a connection from the read replica pool
	// Falls back to writer if no replicas are available
	Reader() *sqlx.DB

	// Seed runs the seeders
	Seed(ctx context.Context) error

	// Writer returns the writer (primary) database connection
	Writer() *sqlx.DB

	// WithTx starts a transaction on the writer (primary) database
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error

	// WithTxTimeout starts a transaction on the writer (primary) database with a timeout
	WithTxTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context, tx *sqlx.Tx) error) error
}

type DB struct {
	opts       DBOptions
	identifier string
	logger     *slog.Logger
	mainFile   string
	mu         sync.RWMutex
	migrator   *dbmate.DB
	name       string
	readers    []*sqlx.DB
	seeder     func(ctx context.Context, db DBer) error
	writer     *sqlx.DB
}

// DBOptions contains configuration options for database connections
type DBOptions struct {
	// Logger is used for database-related logging
	Logger *slog.Logger

	// DisableLogQueries determines whether queries are logged.
	DisableLogQueries bool

	// MainFile is the path to the main application file
	MainFile string

	// MigrationsFS is the embedded filesystem containing migrations
	MigrationsFS FS

	// MigrationsDir is the directory path within MigrationsFS containing migrations
	// Defaults to "migrations"
	MigrationsDir string

	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// Identifier is the identifier of the database which is used for displaying and identifying the migrations directory
	Identifier string

	// WriterURL is the URL for the primary (writer) database connection
	WriterURL string

	// ReaderURLs is a list of URLs for read replica database connections
	ReaderURLs []string

	// Seeder is a function that runs the seeders
	Seeder func(ctx context.Context, db DBer) error
}

// NewDB creates a new database connection pool
func NewDB(ctx context.Context, opts DBOptions) (DBer, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	level := tracelog.LogLevelInfo
	if opts.DisableLogQueries {
		level = tracelog.LogLevelWarn
	}

	if opts.MainFile == "" {
		return nil, fmt.Errorf("main file is required")
	}

	if opts.Mode == "" {
		return nil, fmt.Errorf("mode is required")
	}

	if opts.Identifier == "" {
		return nil, fmt.Errorf("identifier is required")
	}

	if opts.WriterURL == "" {
		return nil, fmt.Errorf("writer URL is required")
	}

	writerPoolConfig, err := pgxpool.ParseConfig(opts.WriterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse writer URL: %w", err)
	}
	writerPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		LogLevel: level,
		Logger: &DbLogger{
			opts.Logger,
			opts.Mode,
		},
	}

	writerPool, err := pgxpool.NewWithConfig(ctx, writerPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer pool: %w", err)
	}
	writer := sqlx.NewDb(stdlib.OpenDBFromPool(writerPool), "pgx")

	var readers []*sqlx.DB
	for _, readerURL := range opts.ReaderURLs {
		readerPoolConfig, err := pgxpool.ParseConfig(readerURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reader URL: %w", err)
		}
		readerPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			LogLevel: level,
			Logger: &DbLogger{
				opts.Logger,
				opts.Mode,
			},
		}

		readerPool, err := pgxpool.NewWithConfig(ctx, readerPoolConfig)
		if err != nil {
			opts.Logger.WarnContext(ctx, "failed to create reader pool",
				"url", readerURL,
				"error", err)
			continue
		}
		readers = append(readers, sqlx.NewDb(stdlib.OpenDBFromPool(readerPool), "pgx"))
	}

	var zeroFS embed.FS
	if opts.MigrationsFS == zeroFS {
		return nil, fmt.Errorf("migrations filesystem is required")
	}

	if opts.MigrationsDir == "" {
		opts.MigrationsDir = "migrations"
	}

	if _, err := opts.MigrationsFS.ReadDir(opts.MigrationsDir); err != nil {
		return nil, fmt.Errorf("migrations directory '%s' not found in filesystem: %w", opts.MigrationsDir, err)
	}

	parsedUrl, err := url.Parse(sanitizeDBURLForMigrator(opts.WriterURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse writer URL: %w", err)
	}

	migrator := dbmate.New(parsedUrl)
	migrator.AutoDumpSchema = false
	migrator.FS = opts.MigrationsFS
	migrator.MigrationsDir = []string{filepath.Join(opts.MigrationsDir, opts.Identifier)}
	migrator.Log = NewDbMigrateLogger(opts.Logger)

	db := &DB{
		opts:       opts,
		logger:     opts.Logger,
		identifier: opts.Identifier,
		mainFile:   opts.MainFile,
		migrator:   migrator,
		name:       strings.TrimPrefix(parsedUrl.Path, "/"),
		writer:     writer,
		readers:    readers,
		seeder:     opts.Seeder,
	}

	return db, nil
}

// Close closes the database connections
func (d *DB) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.writer.Close()
	for _, reader := range d.readers {
		reader.Close()
	}
}

// GenMigration creates a new migration file
func (d *DB) GenMigration(name string) error {
	newMigrator := dbmate.New(&url.URL{})
	newMigrator.Log = NewDbMigrateLogger(d.logger)
	newMigrator.MigrationsDir = []string{
		filepath.Join(filepath.Dir(d.mainFile), d.migrator.MigrationsDir[0]),
	}

	err := newMigrator.NewMigration(name)
	if err != nil {
		return err
	}

	return nil
}

// HealthCheck checks the health of the database connections
func (d *DB) HealthCheck(ctx context.Context) error {
	if err := d.writer.PingContext(ctx); err != nil {
		return fmt.Errorf("writer database health check failed: %w", err)
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	for i, reader := range d.readers {
		if err := reader.PingContext(ctx); err != nil {
			return fmt.Errorf("reader %d database health check failed: %w", i, err)
		}
	}

	return nil
}

// Identifier returns the identifier of the database
func (d *DB) Identifier() string {
	return d.identifier
}

// Migrate runs all pending migrations
func (d *DB) Migrate(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.migrator.Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// MigrateStatus returns the number of pending migrations
func (d *DB) MigrateStatus(ctx context.Context) (int64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	pending, err := d.migrator.Status(true)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration status: %w", err)
	}

	return int64(pending), nil
}

// Name returns the name of the database
func (d *DB) Name() string {
	return d.name
}

// Options returns the options of the database
func (d *DB) Options() DBOptions {
	return d.opts
}

// Reader returns a connection from the read replica pool
func (d *DB) Reader() *sqlx.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.readers) == 0 {
		return d.writer
	}
	return d.readers[rand.Intn(len(d.readers))]
}

// Seed runs the seeders
func (d *DB) Seed(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.seeder == nil {
		return fmt.Errorf("seeder function not initialized")
	}

	if err := d.seeder(ctx, d); err != nil {
		return fmt.Errorf("seeder failed: %w", err)
	}

	return nil
}

// Writer returns the writer (primary) database connection
func (d *DB) Writer() *sqlx.DB {
	return d.writer
}

// WithTx starts a transaction on the writer (primary) database
func (d *DB) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := d.writer.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %v", err, rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTxTimeout starts a transaction on the writer (primary) database with a timeout
func (d *DB) WithTxTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return d.WithTx(ctx, fn)
}

// DbLogger implements tracelog.Logger interface
type DbLogger struct {
	*slog.Logger
	mode types.Mode
}

// Log logs a database query
func (l *DbLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	isDebuggingRequest := GetDebugContext(ctx)
	if !isDebuggingRequest {
		return
	}

	// Extract and format the query execution time
	var executionTime string
	duration, ok := data["time"].(time.Duration)
	if ok {
		executionTime = color.MagentaString(fmt.Sprintf(" (%s)", duration.String()))
	}

	if l.mode == types.ReleaseMode {
		l.InfoContext(ctx, "SQL Query", "query", data["sql"], "duration", duration.String())
		return
	}

	query, ok := data["sql"].(string)
	if !ok {
		return
	}

	// Filter out type maps from args if it's a slice
	var cleanArgs any
	if argSlice, ok := data["args"].([]any); ok && len(argSlice) > 0 {
		// Check if first element is the type information
		if _, ok := argSlice[0].(map[uint32]int); ok {
			cleanArgs = argSlice[1:]
		} else {
			// Check the type name as a fallback
			typeName := fmt.Sprintf("%T", argSlice[0])
			if strings.Contains(typeName, "QueryResultFormats") {
				cleanArgs = argSlice[1:]
			} else {
				cleanArgs = argSlice
			}
		}
	} else {
		cleanArgs = data["args"]
	}

	l.Logger.InfoContext(ctx, color.CyanString("SQL")+executionTime+"\n\n\t\t"+formatQuery(query)+"\n"+formatArgs(cleanArgs))
}

func formatArgs(args any) string {
	if args == nil {
		return ""
	}

	switch v := args.(type) {
	case []any:
		parts := make([]string, len(v))

		if len(parts) == 0 {
			return ""
		}

		for i, arg := range v {
			parts[i] = formatArg(arg)
		}

		return "\n\t\t" + color.YellowString("[ ") + strings.Join(parts, color.YellowString(", ")) + color.YellowString(" ]") + "\n"
	default:
		return "\n\t\t" + formatArg(v)
	}
}

func formatArg(arg any) string {
	switch v := arg.(type) {
	case string:
		return color.GreenString("%q", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return color.BlueString("%d", v)
	case float32, float64:
		return color.BlueString("%f", v)
	case bool:
		return color.MagentaString("%t", v)
	case time.Time:
		return color.CyanString("%s", v.Format("2006-01-02 15:04:05"))
	case nil:
		return color.RedString("NULL")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatQuery(query string) string {
	keywords := []string{
		"COALESCE", "AS", "FILTER", "ON", "JSONB_AGG", "JSONB_BUILD_OBJECT", "SELECT", "FROM",
		"WHERE", "LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "JOIN",
		"ORDER BY", "GROUP BY", "HAVING", "LIMIT", "OFFSET", "INSERT", "UPDATE", "RETURNING", "SET",
		"DELETE", "CREATE", "ALTER", "DROP", "TABLE", "INDEX", "VIEW", "INTO", "VALUES",
		"AND", "OR", "IS NOT", "NOT", "IN", "LIKE", "IS NULL",
	}

	for _, keyword := range keywords {
		// Create a regular expression that matches the keyword,
		// ignoring case and ensuring it's not part of another word
		re := regexp.MustCompile(`(?i)\b(` + regexp.QuoteMeta(keyword) + `)\b`)

		// Replace all occurrences of the keyword with its colored version
		query = re.ReplaceAllString(query, color.BlueString("$1"))
	}

	query = strings.TrimSpace(strings.ReplaceAll(query, "\n", "\n   "))
	query = strings.ReplaceAll(query, "begin", color.BlueString("BEGIN"))
	query = strings.ReplaceAll(query, "commit", color.GreenString("COMMIT"))
	query = strings.ReplaceAll(query, "rollback", color.RedString("ROLLBACK"))

	return query
}

// DbMigrateLogger implements io.Writer interface
type DbMigrateLogger struct {
	logger *slog.Logger
}

// NewDbMigrateLogger creates a new DbMigrateLogger
func NewDbMigrateLogger(logger *slog.Logger) *DbMigrateLogger {
	return &DbMigrateLogger{logger: logger}
}

// Write writes a message to the logger
func (l *DbMigrateLogger) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		l.logger.InfoContext(context.Background(), msg)
	}

	return len(p), nil
}

func sanitizeDBURLForMigrator(dbURL string) string {
	// List of pgx v5 specific parameters to remove
	pgxParams := []string{
		"pool_max_conns",
		"pool_min_conns",
		"pool_max_conn_lifetime",
		"pool_max_conn_idle_time",
		"pool_health_check_period",
	}

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return dbURL
	}

	query := parsedURL.Query()
	for _, param := range pgxParams {
		query.Del(param)
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}
