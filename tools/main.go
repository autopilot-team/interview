package main

import (
	"autopilot/tools/cmd"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Use:           "tools",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	addCommands(ctx, rootCmd)
	return rootCmd.Execute()
}

func addCommands(ctx context.Context, rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmd.NewSyncLocaleCmd(ctx))
	rootCmd.AddCommand(cmd.NewStringerCmd(ctx))
}
