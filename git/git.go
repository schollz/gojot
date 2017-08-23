package git

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// GitRepo is the basic store object.
type GitRepo struct {
	repo   string
	root   string
	folder string
	log    *logrus.Logger
}

// New returns a new GPGStore that can then needs to be initialized with Init()
// The repo is cloned into the `rootDir`.
func New(repo, rootDir string) (*GitRepo, error) {
	gr := new(GitRepo)
	gr.repo = repo
	gr.root = rootDir
	gr.folder = parseRepoFolder(repo)
	if !exists(rootDir) {
		os.MkdirAll(rootDir, 0644)
	}
	gr.log = logrus.New()
	gr.log.SetLevel(logrus.WarnLevel)
	return gr, nil
}

func (gr *GitRepo) Debug(on bool) {
	if on {
		gr.log.SetLevel(logrus.InfoLevel)
	} else {
		gr.log.SetLevel(logrus.WarnLevel)
	}
}

// Update will clone a repo if it doesn't exist or pull a repo, if it does.
func (gr *GitRepo) Update() (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(cwd)
	err = os.Chdir(gr.root)
	if err != nil {
		return
	}
	var cmd *exec.Cmd
	var stdoutStderr []byte
	if !exists(gr.folder) {
		gr.log.Infof("Running: git pull %s %s", gr.repo, gr.folder)
		cmd = exec.Command("git", "clone", gr.repo, gr.folder)
	} else {
		os.Chdir(gr.folder)
		gr.log.Info("Running: git pull --rebase origin master")
		cmd = exec.Command("git", "pull", "--rebase", "origin", "master")
	}
	stdoutStderr, err = cmd.CombinedOutput()
	gr.log.Infof("Output: [%s]\n", stdoutStderr)
	if bytes.Contains(stdoutStderr, []byte("fatal")) {
		err = errors.New("Could not clone repo")
	}
	return
}

func (gr *GitRepo) AddData(data []byte, filepath string) (err error) {
	return
}

func GetRemoteOriginURL(repoFolder string) (repo string, err error) {
	// move to repo
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(cwd)
	err = os.Chdir(repoFolder)
	if err != nil {
		return
	}

	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	stdoutStderr, err := cmd.CombinedOutput()
	repo = strings.TrimSpace(string(stdoutStderr))
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
