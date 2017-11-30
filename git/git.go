package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Repository holds the repository meta data
type Repository struct {
	Path  string
	root  *git.Repository
	Files []string
}

// Identity is the relevant user information
type Identity struct {
	Name       string
	Email      string
	HomeFolder string
}

// UserID holds the users information
var UserID Identity

func init() {

	u, _ := user.Current()
	UserID.HomeFolder = path.Join("/home", u.Username)

	// TODO: having an os.Exit(1) inside a package, is it considered bad?
	err := UserID.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Load parses the users git config file into Identity
func (u *Identity) Load() error {
	const p string = ".gitconfig"
	var name string
	var email string

	file := path.Join(u.HomeFolder, p)

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

	u.Name = name
	u.Email = email

	return nil
}

// Load a git repository from disk
func (r *Repository) Load() error {

	s, err := git.PlainOpen(r.Path)
	if err != nil {
		return err
	}

	r.root = s
	return nil
}

// Branch creates/switches to a new branch based on the filename of the msg
func (r *Repository) Branch(s string, create bool) error {

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
func (r *Repository) CommitFile(filename string, msg string) error {

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
			Name:  UserID.Name,
			Email: UserID.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("Unable to commit: %s", err)
	}

	return nil
}
