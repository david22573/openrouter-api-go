package cmd

import (
	"fmt"
	"os"

	"github.com/david22573/openrouter-api-go/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openrouter",
	Short: "CLI for interacting with the OpenRouter API",
	Long:  "A Go CLI for sending prompts and managing messages via the OpenRouter API.",
}

// cfg holds the loaded config (API key, etc.)
var cfg *config.Config

// Execute runs the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Load configuration on startup
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Add other subcommands here if needed
	// chatCmd is added inside cmd/chat.go
}
