package cmd

import (
	"autopilot/backends/internal/core"
	"fmt"
	"log/slog"
	"slices"

	"github.com/spf13/cobra"
)

func NewGenMigrationCmd(logger *slog.Logger, databases []core.DBer) *cobra.Command {
	var dbName string
	var listDBs bool

	cmd := &cobra.Command{
		Use:   "gen:migration [name]",
		Short: "Create a new migration",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if listDBs {
				for _, db := range databases {
					fmt.Println(db.Identifier())
				}
				return nil
			}

			switch {
			case len(databases) == 0:
				return fmt.Errorf("no databases configured")
			case dbName == "":
				return fmt.Errorf("database name not provided (--db)")
			case len(args) == 0:
				return fmt.Errorf("expected migration name")
			}

			i := slices.IndexFunc(databases, func(d core.DBer) bool { return d.Identifier() == dbName })
			if i == -1 {
				return fmt.Errorf("database %q not found", dbName)
			}
			targetDB := databases[i]

			return targetDB.GenMigration(args[0])
		},
	}

	cmd.Flags().BoolVarP(&listDBs, "list", "l", false, "The database name")
	cmd.Flags().StringVar(&dbName, "db", "", "The database name")
	return cmd
}
