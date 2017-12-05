package cmd

import (
	"fmt"
	"path"

	"github.com/eiso/gpass/encrypt"
	"github.com/eiso/gpass/git"
	"github.com/eiso/gpass/utils"
	"github.com/spf13/cobra"
	"github.com/tucnak/store"
)

type InitCmd struct {
	key string
}

func NewInitCmd() *InitCmd {
	return &InitCmd{}
}

func (c *InitCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init /path/to/git-repository",
		Short: "Initializes gpass for a git repo.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.Execute,
	}

	cmd.Flags().StringVarP(&c.key, "key", "k", "", "Path to your local private key.")
	cmd.MarkFlagRequired("key")

	return cmd
}

func (c *InitCmd) Execute(cmd *cobra.Command, args []string) error {

	u := new(git.User)
	r := new(git.Repository)

	if err := u.Init(); err != nil {
		return err
	}

	r.Path = args[0]

	if err := r.Load(); err != nil {
		return err
	}

	if err := r.Branch("master", true); err != nil {
		return err
	}

	e := ".empty"
	p := path.Join(Cfg.Repository.Path, e)

	if err := utils.TouchFile(p); err != nil {
		return err
	}

	if err := r.CommitFile(Cfg.User, e, "Initial commit."); err != nil {
		return err
	}

	f, err := utils.LoadFile(c.key)
	if err != nil {
		return err
	}

	k := encrypt.NewPGP(f, nil, true)

	if err := k.Keyring(3); err != nil {
		return fmt.Errorf("[exit] only 3 passphrase attempts allowed")
	}

	if err := k.AddPublicKey(); err != nil {
		return fmt.Errorf("Unable to add public key: %s", err)
	}

	Cfg.User = u
	Cfg.Repository = r
	Cfg.PrivateKey = c.key

	if err := store.Save("config.json", Cfg); err != nil {
		return fmt.Errorf("Failed to save the user config: %s", err)
	}

	fmt.Println("Successfully loaded your repository and private key\nConfig file written to your systems config folder as gpass/config.json")
	return nil
}
