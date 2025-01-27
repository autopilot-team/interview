package cmd

import (
	"autopilot/backends/internal/core"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewGenOpenapiCmd(logger *slog.Logger, httpServer *core.HttpServer) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "gen:openapi",
		Short: "Generate OpenAPI spec(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			for _, api := range httpServer.APIDocs {
				spec, err := api.OpenAPI().MarshalJSON()
				if err != nil {
					return err
				}

				version := api.OpenAPI().Extensions["x-sdk-id"]
				filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", version))
				if err := os.WriteFile(filename, spec, 0644); err != nil {
					return fmt.Errorf("failed to write OpenAPI spec to %s: %w", filename, err)
				}

				logger.Info(fmt.Sprintf("Generated OpenAPI spec: ./%s", filename))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "dir", "d", "./packages/api/src/contracts", "Output directory for OpenAPI specs")
	return cmd
}
