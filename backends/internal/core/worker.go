package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
)

// Worker implements the River worker client
type Worker interface {
	GetClient() *river.Client[pgx.Tx]
	GetDBPool() *pgxpool.Pool
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	Queues() *river.QueueBundle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	StopAndCancel(ctx context.Context) error
	Stopped() <-chan struct{}
}

// WorkerOptions contains configuration options for the worker
type WorkerOptions struct {
	// DBURL is the database connection string
	DBURL string

	// Logger is used for worker-related logging
	Logger *slog.Logger

	// Config contains River-specific configuration
	Config *river.Config
}

// BackgroundWorker is a River worker client
type BackgroundWorker struct {
	*river.Client[pgx.Tx]
	dbPool *pgxpool.Pool
	Config *river.Config
	DBURL  string
}

// NewWorker creates a new River worker client
func NewWorker(ctx context.Context, opts WorkerOptions) (*BackgroundWorker, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if opts.DBURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	dbPool, err := pgxpool.New(ctx, opts.DBURL)
	if err != nil {
		return nil, err
	}

	opts.Config.ErrorHandler = &errorHandler{logger: opts.Logger}
	opts.Config.Logger = opts.Logger

	client, err := river.NewClient(riverpgxv5.New(dbPool), opts.Config)
	if err != nil {
		return nil, err
	}

	return &BackgroundWorker{
		client,
		dbPool,
		opts.Config,
		opts.DBURL,
	}, nil
}

// errorHandler implements river.ErrorHandler to handle job errors and panics.
type errorHandler struct {
	logger *slog.Logger
}

// HandleError processes errors that occur during job execution.
func (h *errorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	h.logger.Error("Job failed", "job", job, "error", err)

	return nil
}

// HandlePanic processes panics that occur during job execution.
func (h *errorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
	h.logger.Error("Job panic", "job", job, "panic", panicVal, "trace", trace)

	return nil
}

// GetClient returns the underlying River client
func (w *BackgroundWorker) GetClient() *river.Client[pgx.Tx] {
	return w.Client
}

// GetDbPool returns the underlying database pool
func (w *BackgroundWorker) GetDBPool() *pgxpool.Pool {
	return w.dbPool
}
