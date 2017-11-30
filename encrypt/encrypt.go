package encrypt

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type PGP struct {
	PrivateKey []byte
	Passphrase string
	Message    []byte
	Encrypted  bool
}

var entityList openpgp.EntityList

// WriteToFile writes an encrypted message
func (f *PGP) WriteToFile(path string) error {
	if len(f.Message) == 0 {
		return fmt.Errorf("The message content has not been loaded")
	}

	if !f.Encrypted {
		return fmt.Errorf("Not allowed to write unencrypted content to a file")
	}

	if err := ioutil.WriteFile(path, f.Message, 0600); err != nil {
		return fmt.Errorf("Unable to write to file: %s", err)
	}

	return nil
}

func (f *PGP) Keyring() error {
	passphraseByte := []byte(f.Passphrase)

	s := bytes.NewReader([]byte(f.PrivateKey))
	block, err := armor.Decode(s)
	if err != nil {
		return fmt.Errorf("Unable to armor decode: %s", err)
	} else if block.Type != openpgp.PrivateKeyType {
		return fmt.Errorf("Not a OpenPGP private key: %s", err)
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return fmt.Errorf("Unable to read armor decoded key: %s", err)
	}

	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			return fmt.Errorf("Failed to decrypt main private key: %s", err)
		}
	}

	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	entityList = append(entityList, entity)

	return nil
}

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

func (f *PGP) Encrypt() error {
	if f.Encrypted {
		return fmt.Errorf("The message is encrypted already")
	}

	var w bytes.Buffer

	b, err := armor.Encode(&w, "PGP MESSAGE", nil)
	if err != nil {
		return fmt.Errorf("Unable to armor encode")
	}

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
