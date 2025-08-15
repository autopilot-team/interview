package cmd

import (
	"autopilot/backends/internal/core"
	"context"
	"fmt"
	"log/slog"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/spf13/cobra"
)

func NewDBSeedCmd(ctx context.Context, logger *slog.Logger, databases []core.DBer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db:seed",
		Short: "Seed databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(databases) == 0 {
				logger.Info("No databases configured yet.")
				return nil
			}

			for _, db := range databases {
				msg := fmt.Sprintf("Seeding data for '%s' database...", db.Name())
				logger.Info(msg)

				if err := db.Seed(ctx); err != nil {
					return err
				}

				logger.Info(fmt.Sprintf("%s DONE\n", msg))
			}

			return nil
		},
	}

	return cmd
}
