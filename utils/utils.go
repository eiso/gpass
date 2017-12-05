package utils

import (
	"fmt"
	"io/ioutil"
	"os"
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
