package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestGitSuite(t *testing.T) {
	suite.Run(t, new(GitSuite))
}

type GitSuite struct {
	suite.Suite
	user *User
}

func (s *GitSuite) SetupSuite() {
	s.user = s.newUser()
}

func (s *GitSuite) newUser() *User {
	return &User{Name: "John Doe",
		Email:      "john@doe.org",
		HomeFolder: "/home/john-doe/"}
}

func (s *GitSuite) newRepository(name string) *Repository {
	dir, err := ioutil.TempDir("", name)
	require.NoError(s.T(), err)

	dotgit := filepath.Join(dir, ".git")

	gitExec(s, dotgit, dir, "init")

	return &Repository{Path: dir}
}

func (s *GitSuite) TestCreateBranch() {
	gpass := s.newRepository("gpass-test")
	defer os.RemoveAll(gpass.Path)
	fs := s.newRepository("git-test")
	defer os.RemoveAll(fs.Path)

	dotgit := filepath.Join(fs.Path, ".git")
	origin := "master"
	n := "testbranch"
	ref1 := fmt.Sprintf("refs/heads/%s", origin)
	ref2 := fmt.Sprintf("refs/heads/%s", n)

	// gpass case:
	//   create an origin branch so we can test a checkout
	//   based on an exiting branch
	err := gpass.Load()
	require.NoError(s.T(), err)

	w, err := gpass.root.Worktree()
	require.NoError(s.T(), err)

	o := &git.CheckoutOptions{}
	o.Branch = plumbing.ReferenceName(ref1)
	o.Create = true

	err = w.Checkout(o)
	require.NoError(s.T(), err)

	err = ioutil.WriteFile(filepath.Join(gpass.Path, ".empty"), []byte(""), 0600)
	require.NoError(s.T(), err)

	_, err = w.Add(".empty")
	require.NoError(s.T(), err)

	_, err = w.Commit("add .empty", &git.CommitOptions{
		Author: &object.Signature{
			Name:  s.user.Name,
			Email: s.user.Email,
			When:  time.Now(),
		},
	})
	require.NoError(s.T(), err)

	//   create a new branch based on origin
	gpass.CreateBranch(origin, n)
	require.NoError(s.T(), err)

	// git case:
	err = fs.Load()
	require.NoError(s.T(), err)

	gitExec(s, dotgit, fs.Path, "checkout", "-b", origin)

	err = ioutil.WriteFile(filepath.Join(fs.Path, ".empty"), []byte(""), 0600)
	require.NoError(s.T(), err)

	gitExec(s, dotgit, fs.Path, "add", ".")
	gitExec(s, dotgit, fs.Path, "commit", "-m", "add .empty")

	gitExec(s, dotgit, fs.Path, "checkout", "-b", n, origin)

	// verify both git & gpass case
	gpassRef, err := gpass.root.Head()
	gitRef, err := gpass.root.Head()

	s.Equal(string(gpassRef.Name()), ref2)
	s.Equal(string(gitRef.Name()), ref2)
}

func (s *GitSuite) TestCreateOrphanBranch() {
	gpass := s.newRepository("gpass-test")
	defer os.RemoveAll(gpass.Path)
	fs := s.newRepository("git-test")
	defer os.RemoveAll(fs.Path)

	dotgit := filepath.Join(fs.Path, ".git")
	n := "testbranch"
	ref := fmt.Sprintf("refs/heads/%s", n)

	// gpass case:
	err := gpass.Load()
	require.NoError(s.T(), err)

	err = gpass.CreateOrphanBranch(s.user, n)
	require.NoError(s.T(), err)

	// git case:
	err = fs.Load()
	require.NoError(s.T(), err)

	gitExec(s, dotgit, fs.Path, "checkout", "--orphan", n)

	err = ioutil.WriteFile(filepath.Join(fs.Path, ".empty"), []byte(""), 0600)
	require.NoError(s.T(), err)

	gitExec(s, dotgit, fs.Path, "add", ".")
	gitExec(s, dotgit, fs.Path, "commit", "-m", "creating branch for "+n)

	// verify both git & gpass case:
	_, err = gpass.root.Reference(plumbing.ReferenceName(ref), false)
	require.NoError(s.T(), err)

	_, err = fs.root.Reference(plumbing.ReferenceName(ref), false)
	require.NoError(s.T(), err)
}

func gitExec(s *GitSuite, dir string, worktree string, command ...string) {
	cmd := exec.Command("git",
		append([]string{"--git-dir", dir, "--work-tree", worktree},
			command...)...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	require.NoError(s.T(), err)
}
