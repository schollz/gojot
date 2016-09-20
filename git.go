package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ListBranches returns a slice of the branch names of the repo
// excluding the master branch
func ListBranches(folder string) ([]string, error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []string{}, errors.New("Cannot chdir into " + folder)
	}

	// run git command
	cmd := exec.Command("git", "branch", "--list")
	stdout, err := cmd.Output()
	if err != nil {
		return []string{}, errors.New("Unable to find branches in " + folder)
	}
	possibleBranches := strings.Split(string(stdout), "\n")

	// clean names for branches
	branches := []string{}
	for _, name := range possibleBranches {
		possibleName := strings.TrimSpace(strings.Replace(name, "*", "", -1))
		if possibleName != "master" && len(possibleName) > 0 {
			branches = append(branches, possibleName)
		}
	}

	return branches, nil
}

func ParseFetch(message string) {

}

func NewDocumentOnNewBranch(fulltext string, name string, gitfolder string) {

}

func CloneRepo(folder string, remote string) error {
	logger.Debug("Cloning %s into directory at %s", remote, folder)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return errors.New("Cannot chdir into " + folder)
	}
	cmd := exec.Command("git", "clone", "git@github.com:schollz/test.git")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot clone into " + folder)
	} else {
		return nil
	}
}
