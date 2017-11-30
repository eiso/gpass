package gitops

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type repo struct {
	path string
	root *git.Repository
}

type identity struct {
	name  string
	email string
	home  string
}

var gitID identity
var r repo

// Init setups the git system config & adds .gitignore
func Init() error {

	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("User could not determined: %s", err)
	}

	// TODO: temporary
	gitID.home = path.Join("/home", u.Username)
	r.path = path.Join(gitID.home, "temp/gopass")

	if err := parseGitConfig(); err != nil {
		fmt.Println(err)
	}

	if err := r.load(); err != nil {
		fmt.Println(err)
	}

	return nil
}

func parseGitConfig() error {
	const p string = ".gitconfig"
	var name string
	var email string

	file := path.Join(gitID.home, p)

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

func (r *repo) load() error {

	s, err := git.PlainOpen(r.path)
	if err != nil {
		return err
	}

	r.root = s
	return nil
}

// Branch creates/switches to a new branch based on the filename of the msg
func Branch(s string, create bool) error {

	name := fmt.Sprintf("refs/heads/%s", s)

	w, err := r.root.Worktree()
	if err != nil {
		return fmt.Errorf("Unable to load the work tree: %s", err)
	}

	o := &git.CheckoutOptions{}

	if err = o.Validate(); err != nil {
		return err
	}

	o.Branch = plumbing.ReferenceName(name)
	o.Create = create

	if err = w.Checkout(o); err != nil {
		return fmt.Errorf("Unable to create a new branch: %s", err)
	}

	return nil
}

// CommitFile adds the file & commits it
func CommitFile(filename string, msg string) error {

	w, err := r.root.Worktree()
	if err != nil {
		return fmt.Errorf("Unable to load the work tree: %s", err)
	}

	_, err = w.Add(filename)
	if err != nil {
		return fmt.Errorf("Unable to git add the file: %s", err)
	}

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
