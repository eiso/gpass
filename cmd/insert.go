package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Inserts a new account.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("INSERT")
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)
}
