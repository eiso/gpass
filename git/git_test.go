package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

func (s *GitSuite) TestCreateOrphanBranch() {
	gpass := s.newRepository("gpass-test")
	defer os.RemoveAll(gpass.Path)
	fs := s.newRepository("git-test")
	defer os.RemoveAll(fs.Path)

	dotgit := filepath.Join(fs.Path, ".git")

	err := gpass.Load()
	require.NoError(s.T(), err)

	err = fs.Load()
	require.NoError(s.T(), err)

	n := "testbranch"
	err = gpass.CreateOrphanBranch(s.user, n)
	require.NoError(s.T(), err)

	gitExec(s, dotgit, fs.Path, "checkout", "--orphan", n)

	gpassRef, err := gpass.root.Reference(plumbing.ReferenceName(n), false)
	require.NoError(s.T(), err)

	testRef, err := fs.root.Reference(plumbing.ReferenceName(n), false)
	require.NoError(s.T(), err)

	s.Equal(&gpassRef, &testRef)
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
