package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/eiso/gpass/utils"
	"github.com/spf13/cobra"
)

type RmCmd struct{}

func NewRmCmd() *RmCmd {
	return &RmCmd{}
}

func (c *RmCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Removes an encrypted account and all its history.",
		RunE:  c.Execute,
	}

	return cmd
}

func (c *RmCmd) Execute(cmd *cobra.Command, args []string) error {
	if err := InitCheck(); err != nil {
		return err
	}

	r := Cfg.Repository
	filename := args[0] + ".gpg"
	d := strings.Split(args[0], string(os.PathSeparator))
	path := path.Join(r.Path, d[0])

	if err := r.Load(); err != nil {
		return err
	}

	if !r.BranchExists("gpass") {
		return fmt.Errorf("gpass has not been initialized yet, please run: gpass init")
	}

	if !r.BranchExists(filename) {
		return fmt.Errorf("%s does not exist", args[0])
	}

	m := fmt.Sprintf("Are you sure you would like to delete %s?", args[0])
	p := utils.ConfirmShellPrompt(m)

	if !p {
		return nil
	}

	if err := r.Branch(filename, false); err != nil {
		return err
	}

	if err := utils.DeletePath(path); err != nil {
		return err
	}

	msg := fmt.Sprintf("Remove: %s", args[0])
	if err := r.Commit(Cfg.User, filename, msg); err != nil {
		return err
	}
	/*
		if err := r.RemoveBranch(args[0]); err != nil {
			return err
		}
		fmt.Println("Successfully removed (if you add the account again, all version history will be back)")
	*/
	return nil
}
