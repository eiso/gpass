package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// TouchFile creates an empty file
func TouchFile(filename string) error {
	_, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	return nil
}

// DeletePath removes everything in a path (incl. dirs)
func DeletePath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

// DeleteEmptyFolders recursively removes all empty folders in a path
func DeleteEmptyFolders(dir string) error {

	dir = filepath.Clean(dir)

	err := filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		//log.Printf("\n------\ndir: %s\nfpath: %s", dir, fpath)
		
		if dir == fpath {
			return nil
		}

		if info.IsDir() {
			fd, err := os.Open(fpath)	
			if err != nil {
				return err
			}

			names, err := fd.Readdirnames(100)

			if len(names) < 1 {
				//log.Printf("\nDeleting: %s\n", fpath)
				err = os.RemoveAll(fpath)
				if err != nil {
					return err
				}
			}

			log.Printf("\nDirnames: %s", names)
			for _, name := range names {	
				//log.Printf("\nFolder content:%s\n", fpath + string(os.PathSeparator) + name)		
				f, err := os.Stat(fpath + string(os.PathSeparator) + name)	
				if err != nil {
					return err
				}
				
				if !f.IsDir() {
					log.Printf("YOOOO")
					fd.Close()
					return nil				
				}
			}

			fd.Close()
			DeleteEmptyFolders(fpath)
		} else {
			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// RenamePath renames a path from old to new
func RenamePath(old string, new string) error {
	pd := path.Dir(new)
	if pd != "" {
		if err := os.MkdirAll(pd, os.FileMode(0700)); err != nil {
			return err
		}
	}

	err := os.Rename(old, new)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmShellPrompt load a prompt for a [y/n] confirmation
// source: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func ConfirmShellPrompt(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// PassShellPrompt loads a shell prompt for entering and confirming a passphrase
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
