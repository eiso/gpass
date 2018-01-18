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

// newTestRepository creates a git repository using the systems `git` used for testing 
func (s *GitSuite) newTestRepository(name string) *Repository {
	dir, err := ioutil.TempDir("", name)
	require.NoError(s.T(), err)

	dotgit := filepath.Join(dir, ".git")

	gitExec(s, dotgit, dir, "init")
	
	repo := &Repository{Path: dir}
	
	// create an master branch with a single file and a single commit
	err = repo.Load()
	require.NoError(s.T(), err)

	w, err := repo.root.Worktree()
	require.NoError(s.T(), err)

	o := &git.CheckoutOptions{}
	o.Branch = plumbing.ReferenceName("refs/heads/master")
	o.Create = true

	err = w.Checkout(o)
	require.NoError(s.T(), err)

	err = ioutil.WriteFile(filepath.Join(repo.Path, ".empty"), []byte(""), 0600)
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

	return repo
}

func (s *GitSuite) TestCreateBranch() {
	gpass := s.newTestRepository("gpass-test")
	defer os.RemoveAll(gpass.Path)

	master := "master"
	n := "testbranch"
	ref := fmt.Sprintf("refs/heads/%s", n)

	err := gpass.Load()
	require.NoError(s.T(), err)
	
	err = gpass.CreateBranch(master, n)
	require.NoError(s.T(), err)

	gpassRef, err := gpass.root.Head()

	s.Equal(string(gpassRef.Name()), ref)
}

func (s *GitSuite) TestCreateOrphanBranch() {
	gpass := s.newTestRepository("gpass-test")
	defer os.RemoveAll(gpass.Path)

	n := "orphanbranch"
	ref := fmt.Sprintf("refs/heads/%s", n)

	err := gpass.Load()
	require.NoError(s.T(), err)

	err = gpass.CreateOrphanBranch(s.user, n)
	require.NoError(s.T(), err)

	_, err = gpass.root.Reference(plumbing.ReferenceName(ref), false)
	require.NoError(s.T(), err)
}

// TODO: should be removed once creating git repositories with go-git is added to this package
// creates an unnecessary dependency for `git` to exist on the system
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
