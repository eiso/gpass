package cmd

import (
	"fmt"
	"path"

	"github.com/eiso/gpass/encrypt"
	"github.com/eiso/gpass/utils"
	"github.com/spf13/cobra"
)

type ShowCmd struct{}

func NewShowCmd() *ShowCmd {
	return &ShowCmd{}
}

func (c *ShowCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Args:  cobra.ExactArgs(1),
		Short: "Decrypts and shows the contents of an encrypted account.",
		RunE:  c.Execute,
	}

	return cmd
}

func (c *ShowCmd) Execute(cmd *cobra.Command, args []string) error {
	if err := Load(); err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf("please provide a name for the account you are inserting")
	}

	r := Cfg.Repository
	arg := args[0]
	file := path.Join(r.Path, arg+".gpg")

	pk, err := utils.LoadFile(Cfg.PrivateKey)
	if err != nil {
		return err
	}

	if err := r.Load(); err != nil {
		return err
	}

	if !r.BranchExists("gpass") {
		return fmt.Errorf("gpass has not been initialized yet, please run: gpass init")
	}

	if !r.BranchExists(arg) {
		return fmt.Errorf("the account does not exist")
	}

	if err := r.Branch(arg, false); err != nil {
		return err
	}

	f, err := utils.LoadFile(file)
	if err != nil {
		return err
	}

	p := encrypt.NewPGP(pk, f, true)

	if err := p.LoadKeys(); err != nil {
		return err
	}

	if err := p.Keyring(3); err != nil {
		return fmt.Errorf("[exit] only 3 passphrase attempts allowed")
	}

	if err := p.Decrypt(); err != nil {
		return err
	}

	fmt.Println(string(p.Message))

	return nil
}
