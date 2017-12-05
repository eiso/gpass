package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type InsertCmd struct {
	accountName string
}

func NewInsertCmd() *InsertCmd {
	return &InsertCmd{}
}

func (c *InsertCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a new encrypted account.",
		RunE:  c.Execute,
	}

	cmd.Flags().StringVar(&c.accountName, "account-name", "", "Insert a new account.")

	return cmd
}

func (c *InsertCmd) Execute(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a name for the account you are inserting")
	}

	return nil
}
