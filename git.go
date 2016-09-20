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
