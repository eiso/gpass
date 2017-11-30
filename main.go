package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	e "github.com/eiso/go-pass/encrypt"
	"github.com/eiso/go-pass/git"
	"github.com/eiso/go-pass/utils"
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
	var content e.PGP

	folder := path.Join(git.UserID.Home, "temp/gopass")

	if err := r.Load(folder); err != nil {
		fmt.Println(err)
	}

	// Decryption

	content = e.PGP{PrivateKey: f1,
		Passphrase: *passPtr,
		Message:    f2,
		Encrypted:  true,
	}

	if err := content.Keyring(); err != nil {
		fmt.Println(err)
	}

	if err := content.Decrypt(); err != nil {
		fmt.Println(err)
	}

	if err := content.Encrypt(); err != nil {
		fmt.Println(err)
	}

	if err := r.Branch("msg", true); err != nil {
		fmt.Println(err)
	}

	// TODO
	if err := content.WriteToFile("/home/mthek/temp/gopass/msg.gpg"); err != nil {
		fmt.Println(err)
	}

	msg := fmt.Sprintf("Add: %s", "msg")
	if err := r.CommitFile("msg.gpg", msg); err != nil {
		fmt.Println(err)
	}

	if err := r.Branch("msg2", true); err != nil {
		fmt.Println(err)
	}

	if err := content.WriteToFile("/home/mthek/temp/gopass/msg2.gpg"); err != nil {
		fmt.Println(err)
	}

	msg = fmt.Sprintf("Add: %s", "msg2")
	if err = r.CommitFile("msg2.gpg", msg); err != nil {
		fmt.Println(err)
	}

	if err := r.Branch("msg", false); err != nil {
		fmt.Println(err)
	}

	if err := content.WriteToFile("/home/mthek/temp/gopass/msg1-2.gpg"); err != nil {
		fmt.Println(err)
	}

	msg = fmt.Sprintf("Add: %s", "msg1-2")
	if err = r.CommitFile("msg1-2.gpg", msg); err != nil {
		fmt.Println(err)
	}

}
