package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("openrouter-api-go v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
