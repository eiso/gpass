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
	encrypted  bool
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

	// Decryption

	content := pgp{privateKey: f1,
		passphrase: *passPtr,
		message:    f2,
		encrypted:  true,
	}

	// Build the keyring by loading the private key
	if err := content.keyring(); err != nil {
		fmt.Println(err)
	}

	// Decrypt the PGP Message
	if err := content.decrypt(); err != nil {
		fmt.Println(err)
	}

	decryptedMessage := string(content.message)

	fmt.Println(decryptedMessage)

	// Encryption

	content = pgp{
		message:   content.message,
		encrypted: false,
	}

	// Build the keyring by loading the private key
	if err := content.keyring(); err != nil {
		fmt.Println(err)
	}

	// Decrypt the PGP Message
	if err := content.encrypt(); err != nil {
		fmt.Println(err)
	}

	encryptedMessage := string(content.message)

	fmt.Println(encryptedMessage)

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

func (f *pgp) decrypt() error {
	if !f.encrypted {
		return fmt.Errorf("The message is not encrypted")
	}

	block, err := armor.Decode(bytes.NewReader([]byte(f.message)))
	if err != nil {
		return fmt.Errorf("Invalid PGP message or not armor encoded: %s", err)
	}
	if block.Type != "PGP MESSAGE" {
		return fmt.Errorf("This file is not a PGP message: %s", err)
	}

	//c := packet.Config{DefaultCipher: packet.CipherAES256, DefaultCompressionAlgo: packet.CompressionNone, DefaultHash: crypto.SHA256}

	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return fmt.Errorf("Unable to decrypt the message: %s", err)
	}

	message, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return fmt.Errorf("Unable to convert the decrypted message to a string: %s", err)
	}

	f.encrypted = false
	f.message = message

	return nil
}

func (f *pgp) encrypt() error {
	if f.encrypted {
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

	v, err := e.Write(f.message)
	if err != nil {
		return fmt.Errorf("%s, ints buffered: %v", err, v)
	}

	if err := e.Close(); err != nil {
		return fmt.Errorf("%s", err)
	}

	if err := b.Close(); err != nil {
		return fmt.Errorf("%s", err)
	}

	message, err := ioutil.ReadAll(&w)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	f.encrypted = true
	f.message = message

	return nil
}
