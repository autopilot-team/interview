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

// WorkerOptions contains configuration options for the worker
type WorkerOptions struct {
	// DbURL is the database connection string
	DbURL string

	// Logger is used for worker-related logging
	Logger *slog.Logger

	// Config contains River-specific configuration
	Config *river.Config
}

// Worker is a River worker client
type Worker struct {
	*river.Client[pgx.Tx]
	Config *river.Config
	DbURL  string
}

// NewWorker creates a new River worker client
func NewWorker(ctx context.Context, opts WorkerOptions) (*Worker, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if opts.DbURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	dbPool, err := pgxpool.New(ctx, opts.DbURL)
	if err != nil {
		return nil, err
	}

	opts.Config.ErrorHandler = &errorHandler{logger: opts.Logger}
	opts.Config.Logger = opts.Logger

	client, err := river.NewClient(riverpgxv5.New(dbPool), opts.Config)
	if err != nil {
		return nil, err
	}

	return &Worker{
		client,
		opts.Config,
		opts.DbURL,
	}, nil
}

// errorHandler implements river.ErrorHandler to handle job errors and panics.
type errorHandler struct {
	logger *slog.Logger
}

// HandleError processes errors that occur during job execution.
func (h *errorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	return nil
}

// HandlePanic processes panics that occur during job execution.
func (h *errorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
	return nil
}
