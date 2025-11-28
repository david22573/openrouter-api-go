package cmd

import (
	"fmt"

	"github.com/david22573/openrouter-api-go/internal/app"
	"github.com/david22573/openrouter-api-go/internal/config"
	"github.com/david22573/openrouter-api-go/pkg/openrouter"
	"github.com/spf13/cobra"
)

// rootCmd is the base command.
var rootCmd = &cobra.Command{
	Use:   "openrouter",
	Short: "CLI for interacting with the OpenRouter API",
	Long:  "A command-line interface for sending prompts and managing messages via the OpenRouter API.",

	// This runs before *every* command.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		app.A.Config = cfg

		c, err := openrouter.NewClient(cfg.APIKey)
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}
		app.A.Client = c

		return nil
	},
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
