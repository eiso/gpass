package encrypt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/ssh/terminal"
)

// PGP holds the private key/pass and one message (may be encrypted/decrypted) at a time
type PGP struct {
	PrivateKey []byte
	Message    []byte
	Encrypted  bool
}

var entityList openpgp.EntityList

// NewPGP creates a new instance of PGP struct
func NewPGP(k []byte, m []byte, e bool) *PGP {

	r := new(PGP)

	r.PrivateKey = k
	r.Message = m
	r.Encrypted = e

	return r
}

// LoadKeys loads the private key into entityList (also known as a pgp keyring)
func (f *PGP) LoadKeys() error {
	if len(entityList) > 0 {
		return fmt.Errorf("Keys already loaded")
	}

	s := bytes.NewReader([]byte(f.PrivateKey))
	block, err := armor.Decode(s)
	if err != nil {
		return fmt.Errorf("Not an armor encoded PGP private key: %s", err)
	} else if block.Type != openpgp.PrivateKeyType {
		return fmt.Errorf("Not a OpenPGP private key: %s", err)
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return fmt.Errorf("Unable to read armor decoded key: %s", err)
	}

	entityList = append(entityList, entity)

	return nil
}

func shellPrompt() []byte {
	fmt.Print("Enter passphrase: ")
	passphraseByte, err := terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		fmt.Println("")
	}

	return passphraseByte
}

// WriteFile writes the encrypted message to a new file, fails on existing files
func (f *PGP) WriteFile(repoPath string, filename string) error {
	if len(f.Message) == 0 {
		return fmt.Errorf("The message content has not been loaded")
	}

	if !f.Encrypted {
		return fmt.Errorf("Not allowed to write unencrypted content to a file")
	}

	p := path.Join(repoPath, filename)

	pd := path.Dir(p)
	if pd != "" {
		os.MkdirAll(pd, os.FileMode(0700))
	}

	o, err := os.Open(p)
	if err == nil {
		o.Close()
		return fmt.Errorf("File already exists")
	}
	o.Close()

	d, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("Unable to create the file: %s", err)
	}

	if err := d.Chmod(os.FileMode(0600)); err != nil {
		return fmt.Errorf("Unable to change permissions on file to 0600: %s", err)
	}

	if _, err := d.Write(f.Message); err != nil {
		return fmt.Errorf("Unable to write to file: %s", err)
	}

	return nil
}

//Keyring builds a pgp keyring based upon the users' private key
func (f *PGP) Keyring(attempts int) error {
	entity := entityList[0]

	passphraseByte := shellPrompt()

	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			if attempts > 1 {
				fmt.Println("Sorry, try again.")
				f.Keyring(attempts - 1)
			}
			return fmt.Errorf("Failed to decrypt main private key: %s", err)
		}
	}

	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	entityList = append(entityList, entity)

	return nil
}

// Decrypt a message
func (f *PGP) Decrypt() error {
	if !f.Encrypted {
		return fmt.Errorf("The message is not encrypted")
	}

	block, err := armor.Decode(bytes.NewReader([]byte(f.Message)))
	if err != nil {
		return fmt.Errorf("Invalid PGP message or not armor encoded: %s", err)
	}
	if block.Type != "PGP MESSAGE" {
		return fmt.Errorf("This file is not a PGP message: %s", err)
	}

	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return fmt.Errorf("Unable to decrypt the message: %s", err)
	}

	message, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return fmt.Errorf("Unable to convert the decrypted message to a string: %s", err)
	}

	f.Encrypted = false
	f.Message = message

	return nil
}

// Encrypt a message
func (f *PGP) Encrypt() error {
	if f.Encrypted {
		return fmt.Errorf("The message is encrypted already")
	}

	var w bytes.Buffer

	b, err := armor.Encode(&w, "PGP MESSAGE", nil)
	if err != nil {
		return fmt.Errorf("Unable to armor encode")
	}

	//entity := entityList[0]

	e, err := openpgp.Encrypt(b, entityList, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("Unable to load keyring for encryption: %s", err)
	}

	v, err := e.Write(f.Message)
	if err != nil {
		return fmt.Errorf("%s, ints buffered: %v", err, v)
	}

	if err := e.Close(); err != nil {
		return err
	}

	if err := b.Close(); err != nil {
		return err
	}

	message, err := ioutil.ReadAll(&w)
	if err != nil {
		return err
	}

	f.Encrypted = true
	f.Message = message

	return nil
}
