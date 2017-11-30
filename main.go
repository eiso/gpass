package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	e "github.com/eiso/gpass/encrypt"
	"github.com/eiso/gpass/git"
	"github.com/eiso/gpass/utils"
)

func main() {
	keyPtr := flag.String("key", "", "path to your private key")
	passPtr := flag.String("pass", "", "password to unlock the private key")
	msgPtr := flag.String("msg", "", "path to the encrypted message")
	flag.Parse()

	if *keyPtr == "" || *passPtr == "" || *msgPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	f1, err := utils.LoadFile(*keyPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f2, err := utils.LoadFile(*msgPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var r git.Repository
	var c e.PGP

	r.Path = path.Join(git.UserID.HomeFolder, "temp/gopass/")

	if err := r.Load(); err != nil {
		fmt.Println(err)
	}

	c = e.PGP{PrivateKey: f1,
		Passphrase: *passPtr,
		Message:    f2,
		Encrypted:  true,
	}

	if err := c.Keyring(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.Decrypt(); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	if err := c.Encrypt(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := r.Branch("msg", true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.WriteFile(r.Path, "msg.gpg"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg := fmt.Sprintf("Add: %s", "msg")
	if err := r.CommitFile("msg.gpg", msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := r.Branch("msg2", true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.WriteFile(r.Path, "msg2.gpg"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg = fmt.Sprintf("Add: %s", "msg2")
	if err = r.CommitFile("msg2.gpg", msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := r.Branch("msg", false); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := c.WriteFile(r.Path, "msg1-2.gpg"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg = fmt.Sprintf("Add: %s", "msg1-2")
	if err = r.CommitFile("msg1-2.gpg", msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
