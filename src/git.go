package sdees

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ListBranches returns a slice of the branch names of the repo
// excluding the master branch
func ListBranches(folder string) ([]string, error) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Listing branches %s", id, folder)
	defer timeTrack(time.Now(), "["+id+"]Listing branches "+folder)
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

// GetLatest assumes you already have a cloned repo in `gitfolder`
// and it will fetch the latest and compare the before and after
// to update the cache incrementally
func GetLatest(gitfolder string) ([]string, []string, error) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Getting latest for %s", id, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Getting latest "+gitfolder)
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

// Delete will permanetly delete by deleting and creating a new orphan
// branch. This will erase history on all copies!
func Delete(gitfolder string, branch string) error {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Deleting branch %s in %s", id, branch, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Deleting "+branch)

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

	// Create empty file and commit to same branch
	NewDocument(gitfolder, ".deleted", "", "deleted", "Thu, 07 Apr 2005 22:13:13 +0200", branch)

	return nil
}

func DeleteBranch(branch string) error {
	err := Delete(RemoteFolder, branch)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	err = DeleteCache()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	err = Push(RemoteFolder)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	return nil
}

// Fetch will force fetch and update tracking and rebase all branches so
// that it matches the remote origin. It will not destroy local copies of things.
func Fetch(gitfolder string) error {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Fetching %s", id, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Fetching "+gitfolder)
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	// Fetch all
	cmd := exec.Command("git", "fetch", "--all", "--force", "--prune")
	_, err2 := cmd.Output()
	if err2 != nil {
		logger.Debug("Problem fetching all")
	}

	// Get branchces
	cmd = exec.Command("git", "branch", "-r")
	stdout, err := cmd.Output()
	if err != nil {
		return errors.New("Cannot branch -r")
	}
	branches := []string{}
	remotelyTrackedBranches := make(map[string]bool)
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
		remotelyTrackedBranches[branchName] = true
		branches = append(branches, branchName)
	}

	// Find all locally tracked branches with
	//		git branch -vv
	cmd = exec.Command("git", "branch", "-vv")
	stdout, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot branch -vv")
	}
	branchesToReset := []string{}
	locallyTrackedBranches := make(map[string]bool)
	for _, line := range strings.Split(string(stdout), "\n") {
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			split1 := strings.Split(line, "[")
			split2 := strings.Split(split1[1], "]")
			insides := split2[0]
			split1 = strings.Split(insides, ":")
			branch := split1[0]
			if strings.Contains(branch, "/") {
				split2 = strings.Split(branch, "/")
				branch = split2[len(split2)-1]
			}
			branch = strings.TrimSpace(branch)
			if len(branch) == 0 {
				continue
			}
			if strings.Contains(insides, "ahead ") || strings.Contains(insides, "behind ") {
				branchesToReset = append(branchesToReset, branch)
			}
			locallyTrackedBranches[branch] = true
		}
	}

	// Track branches not being tracked.
	//       BRANCHES NOT BEING TRACKED
	//                 =
	//   SET OF BRANCHES FROM git branch -r
	//                  -
	//   SET OF BRANCHES FROM git branch -vv
	start := time.Now()
	for branch := range remotelyTrackedBranches {
		if _, ok := locallyTrackedBranches[branch]; !ok {
			cmd = exec.Command("git", "branch", "--track", branch, "origin/"+branch)
			cmd.Output()
		}
	}
	logger.Debug("Tracking took " + time.Since(start).String())

	// Find ANY that have "ahead" or "behind", and do
	//      git checkout branch
	//      git reset --hard HEAD
	// 			git rebase
	for _, branch := range branchesToReset {
		logger.Debug("Resetting branch %s", branch)
		cmd = exec.Command("git", "checkout", branch)
		stdout, err = cmd.Output()
		if err != nil {
			return errors.New("Cannot checkout" + branch)
		}
		cmd = exec.Command("git", "reset", "--hard", "HEAD")
		stdout, err = cmd.Output()
		if err != nil {
			return errors.New("Cannot reset --hard HEAD of " + branch)
		}
		cmd = exec.Command("git", "rebase")
		stdout, err = cmd.Output()
		if err != nil {
			return errors.New("Cannot rebase " + branch)
		}
		cmd = exec.Command("git", "checkout", "master")
		stdout, err = cmd.Output()
		if err != nil {
			return errors.New("Cannot checkout master")
		}
	}

	// NOT NEEDED - THIS IS TAKEN CARE OF WITH FETCH --FORCE
	// // Find if branches are no longer on remote and delete them locally
	// localBranches, _ := ListBranches("./")
	// for _, localBranch := range localBranches {
	// 	if _, ok := allBranches[localBranch]; !ok {
	// 		logger.Debug("Deleted locally '%s' - branch no longer on remote", localBranch)
	// 		cmd = exec.Command("git", "branch", "-D", localBranch)
	// 		_, err = cmd.Output()
	// 		if err != nil {
	// 			return errors.New("Problem deleting branch " + localBranch)
	// 		}
	// 	}
	// }

	return err2
}

func NewDocument(gitfolder string, documentname string, fulltext string, message string, datestring string, branchNameOverride string) (string, error) {
	defer timeTrack(time.Now(), "NewDocument")
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	newBranch := branchNameOverride
	if len(branchNameOverride) == 0 {
		newBranch = MakeAlliteration()
	}
	cmd := exec.Command("git", "checkout", "--orphan", newBranch)
	_, err = cmd.Output()
	if err != nil {
		cmd2 := exec.Command("git", "checkout", newBranch)
		_, err2 := cmd2.Output()
		if err2 != nil {
			return newBranch, errors.New("Cannot checkout branch " + newBranch)
		}
	}

	err = ioutil.WriteFile(documentname, []byte(fulltext), 0644)
	if err != nil {
		return newBranch, errors.New("Cannot write file " + documentname)
	}
	if Encrypt {
		err = EncryptFile(documentname, Passphrase)
		if err != nil {
			return newBranch, err
		}
		documentname += ".gpg"
		message = EncryptString(message, Passphrase)
	}

	cmd = exec.Command("git", "add", documentname)
	_, err = cmd.Output()
	if err != nil {
		return newBranch, errors.New("Cannot add file " + documentname)
	}

	cmd = exec.Command("git", "commit", "--date", datestring, "-m", message, documentname)
	_, err = cmd.Output()
	if err != nil {
		return newBranch, errors.New("Cannot commit " + documentname + " error: " + err.Error())
	}

	logger.Debug("Updated document %s in branch %s", documentname, newBranch)

	// Check if its a new document
	_, errExistence := GetTextOfOne(gitfolder, "master", documentname)
	if errExistence != nil {
		logger.Debug("It seems %s doesn't exist yet, making a index file for it in master", documentname)
		cmd2 := exec.Command("git", "checkout", "master")
		_, err2 := cmd2.Output()
		if err2 != nil {
			logger.Warn("Something wrong checking out master")
		}
		text := "Yay some text!"
		if Encrypt {
			text = EncryptString(text, Passphrase)
		}
		err2 = ioutil.WriteFile(documentname, []byte(text), 0644)
		if err != nil {
			logger.Warn("Something wrong with writing " + documentname)
		}
		cmd2 = exec.Command("git", "add", documentname)
		_, err2 = cmd2.Output()
		if err2 != nil {
			logger.Warn("Something wrong checking out master")
		}
		cmd2 = exec.Command("git", "commit", "--date", datestring, "-m", message, documentname)
		_, err2 = cmd2.Output()
		if err2 != nil {
			logger.Warn("Something wrong checking out master")
		}
	}

	return newBranch, err
}

func Push(gitfolder string) error {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Pushing %s", id, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Pushing")
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	cmd := exec.Command("git", "push", "--force", "--all", "origin")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot push " + gitfolder)
	}

	return nil
}

func Clone(folder string, remote string) error {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Cloning %s", id, remote)
	defer timeTrack(time.Now(), "["+id+"]Cloning")
	var err error
	logger.Debug("Cloning %s into directory at %s", remote, folder)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	cmd := exec.Command("git", "clone", remote, folder)
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Cloning command failed: '" + strings.Join([]string{"git", "clone", remote, folder}, " ") + "'")
	}

	Fetch(folder)

	return nil
}
