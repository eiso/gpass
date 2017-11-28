package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type pgp struct {
	privateKey []byte
	passphrase string
	message    []byte
}

var entityList openpgp.EntityList

func main() {

	keyPtr := flag.String("key", "", "path to your private key")
	passPtr := flag.String("pass", "", "password to unlock the private key")
	msgPtr := flag.String("msg", "", "path to the encrypted message")
	flag.Parse()

	if *keyPtr == "" || *passPtr == "" || *msgPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	f1, err := loadFile(*keyPtr)
	if err != nil {
		fmt.Println(err)
	}

	f2, err := loadFile(*msgPtr)
	if err != nil {
		fmt.Println(err)
	}

	content := pgp{privateKey: f1,
		passphrase: *passPtr,
		message:    f2,
	}

	// Build the keyring by loading the private key
	if err := content.keyring(); err != nil {
		fmt.Println(err)
	}

	// Decrypt the PGP Message
	msg, err := content.decrypt()
	if err != nil {
		fmt.Println(err)
	}

	decryptedMessage := string(msg)

	fmt.Println(decryptedMessage)
}

func loadFile(filename string) ([]byte, error) {

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Encrypted file could not be read: %s", err)
	}

	return f, err
}

func (f pgp) keyring() error {

	passphraseByte := []byte(f.passphrase)

	s := bytes.NewReader([]byte(f.privateKey))
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

func (f pgp) decrypt() ([]byte, error) {

	block, err := armor.Decode(bytes.NewReader([]byte(f.message)))
	if err != nil {
		return nil, fmt.Errorf("Invalid PGP message or not armor encoded: %s", err)
	}
	if block.Type != "PGP MESSAGE" {
		return nil, fmt.Errorf("This file is not a PGP message: %s", err)
	}

	//c := packet.Config{DefaultCipher: packet.CipherAES256, DefaultCompressionAlgo: packet.CompressionNone, DefaultHash: crypto.SHA256}

	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to decrypt the message: %s", err)
	}

	message, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert the decrypted message to a string: %s", err)
	}

	return message, nil
}
