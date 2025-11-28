package cmd

import (
	"fmt"

	"github.com/david22573/openrouter-api-go/internal/app"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Send a message to OpenRouter",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing chat prompt")
		}

		prompt := args[0]

		resp, err := app.A.Client.Chat(cmd.Context(), prompt)
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
