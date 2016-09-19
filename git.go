package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ListBranches returns a slice of the branch names of the repo
func ListBranches(folder string) ([]string, error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []string{}, errors.New("Cannot chdir into " + folder)
	}
	cmd := exec.Command("git", "branch", "--list")
	stdout, err := cmd.Output()
	if err != nil {
		return []string{}, errors.New("Unable to find branches in " + folder)
	}
	possibleBranches := strings.Split(string(stdout), "\n")
	branches := make([]string, len(possibleBranches))
	for i, name := range possibleBranches {
		branches[i] = strings.TrimSpace(strings.Replace(name, "*", "", -1))
	}
	return branches, nil
}
