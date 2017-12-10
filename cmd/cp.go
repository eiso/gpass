package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/eiso/gpass/utils"
	"github.com/spf13/cobra"
)

type CpCmd struct{}

func NewCpCmd() *CpCmd {
	return &CpCmd{}
}

func (c *CpCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cp",
		Short: "Copies an encrypted account and all its history to a new path.",
		RunE:  c.Execute,
	}

	return cmd
}

func (c *CpCmd) Execute(cmd *cobra.Command, args []string) error {
	if err := InitCheck(); err != nil {
		return err
	}

	if len(args) != 2 {
		return fmt.Errorf("please provide an old and a new path")
	}

	r := Cfg.Repository
	filename := args[0] + ".gpg"
	new := args[1] + ".gpg"
	pathPrev := path.Join(r.Path, filename)
	pathNew := path.Join(r.Path, new)
	d := strings.Split(filename, string(os.PathSeparator))
	root := path.Join(r.Path, d[0])

	if err := r.Load(); err != nil {
		return err
	}

	if !r.BranchExists("gpass") {
		return fmt.Errorf("gpass has not been initialized yet, please run: gpass init")
	}

	if !r.BranchExists(filename) {
		return fmt.Errorf("%s does not exist", args[0])
	}

	if r.BranchExists(new) {
		return fmt.Errorf("%s already exists", args[1])
	}

	if err := r.CreateBranch(filename, new); err != nil {
		return err
	}

	if err := utils.RenamePath(pathPrev, pathNew); err != nil {
		return err
	}

	if err := utils.DeletePath(root); err != nil {
		return err
	}

	msg := fmt.Sprintf("Copied: %s to %s", args[0], args[1])
	if err := r.CommitFile(Cfg.User, new, msg); err != nil {
		return err
	}

	fmt.Println("Successfully moved the account to:", args[1])

	return nil
}
