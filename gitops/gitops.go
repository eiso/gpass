package gitops

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var username string
var home string
var gitRepo string

type identity struct {
	name  string
	email string
}

var gitID identity

// Init setups the git system config
func Init() error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("User could not determined: %s", err)
	}
	home = path.Join("/home", u.Username)
	// TO-DO: temporary
	gitRepo = path.Join(home, "temp/gopass")

	parseGitConfig()

	return nil
}

func parseGitConfig() error {
	const p string = ".gitconfig"
	var name string
	var email string

	file := path.Join(home, p)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Git config could not be read: %s", err)
	}

	c := config.NewDecoder(bytes.NewBuffer(f))
	cfg := &config.Config{}

	if err := c.Decode(cfg); err != nil {
		return fmt.Errorf("%s", err)
	}

	for _, section := range cfg.Sections {
		if section.Name != "user" {
			continue
		}
		for _, option := range section.Options {
			switch option.Key {
			case "name":
				name = option.Value
			case "email":
				email = option.Value
			default:
				continue
			}
		}
	}

	gitID.name = name
	gitID.email = email

	return nil
}

// Commit creates a commit for the encrypted message file
func Commit(filename string) error {
	r, err := git.PlainOpen(gitRepo)
	if err != nil {
		return fmt.Errorf(": %s", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("Unable to load the work tree: %s", err)
	}

	_, err = w.Add(filename)
	if err != nil {
		return fmt.Errorf("Unable to git add the file: %s", err)
	}

	msg := fmt.Sprintf("Add: %s", filename)

	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitID.name,
			Email: gitID.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to commit: %s", err)
	}

	return nil
}
