package cmd

import (
	"fmt"

	"github.com/eiso/gpass/encrypt"
	"github.com/eiso/gpass/utils"
	"github.com/spf13/cobra"
)

type InsertCmd struct{}

func NewInsertCmd() *InsertCmd {
	return &InsertCmd{}
}

func (c *InsertCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a new encrypted account.",
		RunE:  c.Execute,
	}

	return cmd
}

func (c *InsertCmd) Execute(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a name for the account you are inserting")
	}

	path := args[0]
	filename := path + ".gpg"

	var prompts []string

	prompts = append(prompts, "Enter password for "+path+": ")
	prompts = append(prompts, "Retype password for "+path+": ")

	r := Cfg.Repository

	if err := r.Load(); err != nil {
		return err
	}

	if r.BranchExists(path) {
		return fmt.Errorf("the account already exists")
	}

	if err := r.Branch("master", false); err != nil {
		return err
	}

	if err := r.Branch(path, true); err != nil {
		return err
	}

	f, err := utils.PassShellPrompt(prompts)
	if err != nil {
		return err
	}

	pk, err := utils.LoadFile(Cfg.PrivateKey)
	if err != nil {
		return err
	}

	p := encrypt.NewPGP(pk, f, false)

	if err := p.LoadKeys(); err != nil {
		return err
	}

	if err := p.Encrypt(); err != nil {
		return err
	}

	if err := p.WriteFile(r.Path, filename); err != nil {
		return err
	}

	msg := fmt.Sprintf("Add: %s", path)
	if err := r.CommitFile(Cfg.User, filename, msg); err != nil {
		return err
	}

	return nil
}
