package utils

import (
	"fmt"
	"io/ioutil"
)

// LoadFile loads a file
func LoadFile(filename string) ([]byte, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("File could not be read: %s", err)
	}

	return f, err
}
