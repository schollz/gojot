package main

import (
	"errors"
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
	defer timeTrack(time.Now(), "Listed branches for "+folder)
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

func GetLatest(gitfolder string) ([]string, []string, error) {
	defer timeTrack(time.Now(), "Got latest for "+gitfolder)
	var err error
	err = nil
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	addedBranches := []string{}
	deletedBranches := []string{}

	oldBranches, err := ListBranches("./")

	if err != nil {
		return []string{}, []string{}, err
	}

	err = Fetch(gitfolder)
	if err != nil {
		return []string{}, []string{}, err
	}

	newBranches, err := ListBranches("./")
	if err != nil {
		return []string{}, []string{}, err
	}

	oldBranchesMap := make(map[string]bool)
	for _, branch := range oldBranches {
		oldBranchesMap[branch] = true
	}
	for _, branch := range newBranches {
		if _, ok := oldBranchesMap[branch]; !ok {
			addedBranches = append(addedBranches, branch)
		}
	}
	if len(addedBranches) > 0 {
		logger.Debug("Found %d remote branches that have been added: %v", len(addedBranches), addedBranches)
	}

	newBranchesMap := make(map[string]bool)
	for _, branch := range newBranches {
		newBranchesMap[branch] = true
	}
	for _, branch := range oldBranches {
		if _, ok := newBranchesMap[branch]; !ok {
			deletedBranches = append(deletedBranches, branch)
		}
	}
	if len(deletedBranches) > 0 {
		logger.Debug("Found %d remote branches that have been deleted: %v", len(deletedBranches), deletedBranches)
	}

	return addedBranches, deletedBranches, err

}

func Delete(gitfolder string, branch string) error {
	defer timeTrack(time.Now(), "Deleted branch "+branch+" in "+gitfolder)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	// Make sure we aren't on that branch
	cmd := exec.Command("git", "checkout", "master")
	_, err := cmd.Output()
	if err != nil {
		return errors.New("Problem switching to master")
	}

	// Delete branch
	cmd = exec.Command("git", "branch", "-D", branch)
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem deleting branch " + branch)
	}

	// Delete branch remotely
	cmd = exec.Command("git", "push", "origin", "--delete", branch)
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem deleting branch remotely " + branch)
	}
	return nil
}

func Fetch(gitfolder string) error {
	defer timeTrack(time.Now(), "Fetching "+gitfolder)
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

	// Fetch deleted
	cmd = exec.Command("git", "fetch", "-p")
	_, err = cmd.Output()
	if err != nil {
		logger.Error("Problem fetching deleted")
	}

	// Get branchces
	cmd = exec.Command("git", "branch", "-r")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.New("Cannot branch -r")
	}
	branches := []string{}
	allBranches := make(map[string]bool)
	for _, line := range strings.Split(string(stdout), "\n") {
		branchName := strings.TrimSpace(line)
		if strings.Contains(branchName, "->") {
			branchName = strings.TrimSpace(strings.Split(branchName, "->")[1])
		}
		if strings.Contains(branchName, "origin/") {
			branchName = strings.TrimSpace(strings.Split(branchName, "origin/")[1])
		}
		allBranches[branchName] = true
		if len(branchName) == 0 || branchName == "master" {
			continue
		}
		branches = append(branches, branchName)
	}

	// Track each branch
	start := time.Now()
	for _, branch := range branches {
		cmd = exec.Command("git", "branch", "--track", branch, "origin/"+branch)
		cmd.Output()
	}
	logger.Debug("Tracking took " + time.Since(start).String())

	// Find if branches are no longer on remote and delete them locally
	logger.Debug(os.Getwd())
	localBranches, _ := ListBranches("./")
	logger.Debug(os.Getwd())
	for _, localBranch := range localBranches {
		if _, ok := allBranches[localBranch]; !ok {
			logger.Debug("Deleted locally '%s' - branch no longer on remote", localBranch)
			cmd = exec.Command("git", "branch", "-D", localBranch)
			_, err = cmd.Output()
			if err != nil {
				return errors.New("Problem deleting branch " + localBranch)
			}
		}
	}

	return nil
}

func NewDocument(gitfolder string, documentname string, fulltext string, message string, datestring string) (string, error) {
	defer timeTrack(time.Now(), "New document "+documentname+" in "+gitfolder+" created")
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
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot push " + gitfolder)
	}

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
