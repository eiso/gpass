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
		Args:  cobra.ExactArgs(1),
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
	f := []byte("this is my test message")
	r := Cfg.Repository

	pk, err := utils.LoadFile(Cfg.PrivateKey)
	if err != nil {
		return err
	}

	p := encrypt.NewPGP(pk, f, false)

	/*	if err := p.Keyring(3); err != nil {
			return fmt.Errorf("[exit] only 3 passphrase attempts allowed")
		}
	*/
	if err := p.Encrypt(); err != nil {
		return err
	}

	if err := r.Load(); err != nil {
		return err
	}

	if err := r.Branch("master", false); err != nil {
		return err
	}

	if err := r.Branch(path, true); err != nil {
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
