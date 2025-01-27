package cmd

import (
	"autopilot/backends/internal/core"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/spf13/cobra"
)

func NewDbMigrateCmd(ctx context.Context, logger *slog.Logger, databases []core.DBer, workers []*core.Worker) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db:migrate",
		Short: "Migrate database(s) to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(databases) == 0 && len(workers) == 0 {
				logger.Info("No databases or workers configured yet.")
				return nil
			}

			for _, worker := range workers {
				parsedUrl, err := url.Parse(worker.DbURL)
				if err != nil {
					return fmt.Errorf("failed to parse worker database URL: %w", err)
				}

				msg := fmt.Sprintf("Migrating worker schema for '%s' database...", strings.TrimPrefix(parsedUrl.Path, "/"))
				logger.Info(msg)

				dbPool, err := pgxpool.New(ctx, worker.DbURL)
				if err != nil {
					return fmt.Errorf("failed to create 'worker' database pool: %w", err)
				}
				defer dbPool.Close()

				migrator, err := rivermigrate.New(riverpgxv5.New(dbPool), nil)
				if err != nil {
					return fmt.Errorf("failed to create 'worker' database migrator: %w", err)
				}

				_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
				if err != nil {
					return fmt.Errorf("failed to migrate 'worker' database: %w", err)
				}

				logger.Info(fmt.Sprintf("%s DONE\n", msg))
			}

			for _, db := range databases {
				msg := fmt.Sprintf("Migrating application schema for '%s' database...", db.Name())
				logger.Info(msg)

				if err := db.Migrate(ctx); err != nil {
					return err
				}

				logger.Info(fmt.Sprintf("%s DONE\n", msg))
			}

			return nil
		},
	}

	return cmd
}
