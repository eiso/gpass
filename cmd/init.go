package cmd

import (
	"fmt"

	"github.com/eiso/gpass/git"
	"github.com/spf13/cobra"
	"github.com/tucnak/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gpass for a git repo.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	store.Init("gpass")
	rootCmd.AddCommand(initCmd)
}

func execute() error {
	u := new(git.User)

	if err := u.Init(); err != nil {
		return err
	}

	if err := store.Save("config.json", u); err != nil {
		return fmt.Errorf("Failed to save the user config: ", err)
	}

	return nil
}
