package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gpass",
	Short: "gpass is an encrypted account manager built on top of git.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a command or run: gpass --help")
	},
}

// Execute the cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
