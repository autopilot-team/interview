package cmd

import (
	"autopilot/backends/internal/core"
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
)

func NewGenMigrationCmd(logger *slog.Logger, databases []core.DBer) *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "gen:migration [name]",
		Short: "Create a new migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(databases) == 0 {
				logger.Info("No databases configured yet.")
				return nil
			}

			var targetDB core.DBer
			for _, db := range databases {
				if db.Identifier() == dbName {
					targetDB = db
					break
				}
			}

			if targetDB == nil {
				return fmt.Errorf("database '%s' not found", dbName)
			}

			// Only run migration for the target database
			if err := targetDB.GenMigration(args[0]); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbName, "db", "", "The database name")
	if err := cmd.MarkFlagRequired("db"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
