package cmd

import (
	"fmt"
	"os"

	"github.com/eiso/gpass/git"
	"github.com/spf13/cobra"
	"github.com/tucnak/store"
)

// Config holds the gpass/config.json file information
type Config struct {
	User       *git.User       `json:"user"`
	Repository *git.Repository `json:"repository"`
	PrivateKey string          `json:"private-key"`
}

// Cfg is a package variable initialized with the init command
var Cfg Config

var rootCmd = &cobra.Command{
	Use:   "gpass",
	Short: "gpass is an encrypted account manager built on top of git.",
	Run: func(cmd *cobra.Command, args []string) {
		c := NewListCmd()
		c.Execute(cmd, args)
		//fmt.Println("Please specify a command or run: gpass --help")
	},
}

func InitCheck() error {
	if Cfg.Repository == nil {
		return fmt.Errorf("gpass has not been initialized yet, please run: gpass init")
	}

	return nil
}

func init() {
	store.Init("gpass")

	store.Load("config.json", &Cfg)

	rootCmd.AddCommand(NewInitCmd().Cmd())
	rootCmd.AddCommand(NewInsertCmd().Cmd())
	rootCmd.AddCommand(NewShowCmd().Cmd())
	rootCmd.AddCommand(NewListCmd().Cmd())
	rootCmd.AddCommand(NewRmCmd().Cmd())
	rootCmd.AddCommand(NewMvCmd().Cmd())
	rootCmd.AddCommand(NewCpCmd().Cmd())
}

// Execute the cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
