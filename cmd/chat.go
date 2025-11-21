package cmd

import (
	"fmt"

	"github.com/david22573/openrouter-api-go/internal/api"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a prompt to OpenRouter and get a response",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.APIKey)
		resp, err := client.Chat(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Response:", resp)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
