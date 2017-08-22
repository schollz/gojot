package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// GitRepo is the basic store object.
type GitRepo struct {
	repo   string
	root   string
	folder string
}

// New returns a new GPGStore that can then needs to be initialized with Init()
func New(repo, rootDir string) (*GitRepo, error) {
	gr := new(GitRepo)
	gr.repo = repo
	gr.root = rootDir
	gr.folder = parseRepoFolder(repo)
	return gr, nil
}

// Update will clone a repo if it doesn't exist or pull a repo, if it does.
func (gr *GitRepo) Update() (err error) {
	cwd, err := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gr.root)
	var cmd *exec.Cmd
	var stdoutStderr []byte
	if !exists(path.Join(gr.root, gr.folder)) {
		// git clone repo
		cmd = exec.Command("git", "clone", gr.repo)
	} else {
		// git pull --rebase origin master
		os.Chdir(gr.folder)
		cmd = exec.Command("git", "pull", "--rebase", "origin", "master")
	}
	stdoutStderr, err = cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if bytes.Contains(stdoutStderr, []byte("fatal")) {
		err = errors.New("Could not clone repo")
	}
	return
}

func parseRepoFolder(repo string) (folder string) {
	firstPart := strings.Split(repo, ".git")[0]
	firstPartSplit := strings.Split(firstPart, "/")
	folder = strings.TrimSpace(firstPartSplit[len(firstPartSplit)-1])
	return
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
