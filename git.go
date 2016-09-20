package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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

func GetLatest(gitfolder string) ([]string, error) {
	var err error
	err = nil
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	addedBranches := []string{}

	oldBranches, err := ListBranches(gitfolder)
	if err != nil {
		return []string{}, err
	}
	fmt.Println(oldBranches)

	err = Fetch(gitfolder)
	if err != nil {
		return []string{}, err
	}

	newBranches, err := ListBranches(gitfolder)
	if err != nil {
		return []string{}, err
	}

	fmt.Println(oldBranches)
	fmt.Println(newBranches)

	return addedBranches, err

}

func Fetch(gitfolder string) error {
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	// Fetch all
	cmd := exec.Command("git", "fetch", "--all")
	_, err = cmd.Output()
	if err != nil {
		logger.Error("Problem fetching all")
	}

	// Get branchces
	cmd = exec.Command("git", "branch", "-r")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.New("Cannot branch -r")
	}
	branches := []string{}
	for _, line := range strings.Split(string(stdout), "\n") {
		branchName := strings.TrimSpace(line)
		if strings.Contains(branchName, "->") {
			branchName = strings.TrimSpace(strings.Split(branchName, "->")[1])
		}
		if strings.Contains(branchName, "origin/") {
			branchName = strings.TrimSpace(strings.Split(branchName, "origin/")[1])
		}
		if len(branchName) == 0 || branchName == "master" {
			continue
		}
		branches = append(branches, branchName)
	}

	// Track each branch
	for _, branch := range branches {
		cmd = exec.Command("git", "branch", "--track", branch, "origin/"+branch)
		cmd.Output()
	}

	// Fetch all
	cmd = exec.Command("git", "fetch", "--all")
	stdout, err = cmd.Output()
	if err != nil {
		logger.Error("Problem fetching all")
	}

	return nil
}

func NewDocument(gitfolder string, documentname string, fulltext string, message string, datestring string) (string, error) {
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	newBranch := RandStringBytesMaskImprSrc(6, time.Now().UnixNano())
	cmd := exec.Command("git", "checkout", "--orphan", newBranch)
	_, err = cmd.Output()
	if err != nil {
		log.Println(err)
		return newBranch, errors.New("Cannot checkout branch " + newBranch)
	}

	err = ioutil.WriteFile(documentname, []byte(fulltext), 0644)
	if err != nil {
		return newBranch, errors.New("Cannot write file " + documentname)
	}

	cmd = exec.Command("git", "add", documentname)
	_, err = cmd.Output()
	if err != nil {
		return newBranch, errors.New("Cannot add file " + documentname)
	}

	cmd = exec.Command("git", "commit", "--date", datestring, "-m", message, documentname)
	_, err = cmd.Output()
	if err != nil {
		return newBranch, errors.New("Cannot commit " + documentname)
	}

	return newBranch, err
}

func Push(gitfolder string) error {
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	cmd := exec.Command("git", "push", "--all", "origin")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.New("Cannot push " + gitfolder)
	}

	logger.Debug(string(stdout))
	return nil
}

func Clone(folder string, remote string) error {
	var err error
	logger.Debug("Cloning %s into directory at %s", remote, folder)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	cmd := exec.Command("git", "clone", "git@github.com:schollz/test.git", folder)
	_, err = cmd.Output()
	if err != nil {
		log.Println(err)
		return errors.New("Cannot clone into " + folder)
	}

	Fetch(folder)

	return nil
}
