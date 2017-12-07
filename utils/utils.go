package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// LoadFile loads a file
func LoadFile(filename string) ([]byte, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("File could not be read: %s", err)
	}

	return f, err
}

// Touch file
func TouchFile(filename string) error {
	_, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	return nil
}

func PassShellPrompt(prompts []string) ([]byte, error) {

	if len(prompts) != 2 {
		return nil, fmt.Errorf("Two prompt phrases are required")
	}

	fmt.Print(prompts[0])
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println("")

	fmt.Print(prompts[1])
	p2, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println("")

	if string(p) != string(p2) {
		return nil, fmt.Errorf("the entered passwords do not match")
	}

	return p, nil
}
