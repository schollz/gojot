package sdees

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func WriteToMaster(gitfolder string, filename string, text string) {
	dir, _ := os.Getwd()
	logger.Debug("Current DIR: %s", dir)
	os.Chdir(gitfolder)
	cmd2 := exec.Command("git", "checkout", "-f", "master")
	_, err2 := cmd2.Output()
	if err2 != nil {
		cmd2 = exec.Command("git", "checkout", "--orphan", "master")
		_, err2 = cmd2.Output()
		if err2 != nil {
			logger.Error("Couldn't checkout master")
		}
	}
	text = EncryptString(text, Passphrase)
	err2 = ioutil.WriteFile(filename, []byte(text), 0644)
	if err2 != nil {
		logger.Warn("Something wrong with writing " + filename)
	}
	cmd2 = exec.Command("git", "add", filename)
	_, err2 = cmd2.Output()
	if err2 != nil {
		logger.Warn("Something wrong adding " + filename)
	}
	cmd2 = exec.Command("git", "commit", "-m", "'added key'", filename)
	_, err2 = cmd2.Output()
	if err2 != nil {
		logger.Warn("Something wrong: %s ", strings.Join([]string{"git", "commit", "-m", "'added key'", filename}, " "))
	}
	os.Chdir(dir)
}

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

func trackBranchInParallel(gitfolder string, branch string, wg *sync.WaitGroup) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)
	cmd := exec.Command("git", "branch", "--track", branch, "origin/"+branch)
	cmd.Output()
	cmd = exec.Command("git", "branch", "--set-upstream-to=origin/"+branch, branch)
	cmd.Output()
	wg.Done()
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

	// Get password asyncrhonsly
	gettingPassword := false
	wg := sync.WaitGroup{}
	if Passphrase == "alskdfjalskdjfalskjdflajsdfljasd" {
		wg.Add(1)
		go func() { Passphrase = PromptPassword(RemoteFolder); wg.Done(); fmt.Println("Fetching...") }()
		gettingPassword = true
	}

	// Get clean DIR
	cmd := exec.Command("git", "reset", "--hard", "HEAD")
	out2, _ := cmd.StderrPipe()
	cmd.Start()
	out2b, _ := ioutil.ReadAll(out2)
	cmd.Wait()
	if strings.Contains(string(out2b), "fatal:") {
		return errors.New(strings.TrimSpace(string(out2b)))
	}

	// Fetch all
	cmd = exec.Command("git", "fetch", "--all", "--force", "--prune")
	out2, _ = cmd.StderrPipe()
	cmd.Start()
	out2b, _ = ioutil.ReadAll(out2)
	cmd.Wait()
	if strings.Contains(string(out2b), "fatal:") {
		return errors.New(strings.TrimSpace(string(out2b)))
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

	logger.Debug("Found %d branches", len(branches))

	// Find all locally tracked branches with
	//		git branch -vv
	logger.Debug("git branch -vv")
	cmd = exec.Command("git", "branch", "-vv")
	stdout, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot branch -vv")
	}
	locallyTrackedBranches := make(map[string]bool)
	for _, line := range strings.Split(string(stdout), "\n") {
		branch := strings.Split(strings.TrimSpace(line), " ")[0]
		branch = strings.TrimSpace(branch)
		if len(branch) < 2 || strings.Contains(branch, "Found local") || !strings.Contains(line, "[") {
			continue
		}
		locallyTrackedBranches[branch] = true
	}

	// Track branches not being tracked.
	//       BRANCHES NOT BEING TRACKED
	//                 =
	//   SET OF BRANCHES FROM git branch -r
	//                  -
	//   SET OF BRANCHES FROM git branch -vv
	start := time.Now()
	numTracked := 0
	// wg2 := sync.WaitGroup{}
	for branch := range remotelyTrackedBranches {
		if _, ok := locallyTrackedBranches[branch]; !ok {
			logger.Debug("remote '%s' not in local", branch)
			// wg2.Add(1)
			// trackBranchInParallel(gitfolder, branch, &wg2)

			cmd = exec.Command("git", "branch", "--track", branch, "origin/"+branch)
			cmd.Output()
			cmd = exec.Command("git", "branch", "--set-upstream-to=origin/"+branch, branch)
			cmd.Output()
			numTracked++
		}
	}
	// wg2.Wait()
	logger.Debug("Tracking took " + time.Since(start).String())

	// Find ANY that have "ahead" or "behind", and do
	branchesToReset := []string{}
	logger.Debug(`git for-each-ref --format="%%(refname)-=-%%(push:track)" refs/heads`)
	cmd = exec.Command("git", "for-each-ref", `--format="%(refname)-=-%(push:track)"`, "refs/heads")
	stdout, err = cmd.Output()
	if err != nil {
		return errors.New("Cannot for-each-ref")
	}
	for _, line := range strings.Split(string(stdout), "\n") {
		line = strings.Replace(line, `"`, "", -1)
		// logger.Debug(line)
		items := strings.Split(line, "-=-")
		if len(items) < 2 {
			continue
		}
		branch := strings.Replace(items[0], "refs/heads/", "", -1)
		if strings.Contains(items[1], "ahead") || strings.Contains(items[1], "behind") {
			branchesToReset = append(branchesToReset, branch)
		}
	}

	if len(branchesToReset) > 0 {
		DeleteCache()
	}
	if gettingPassword {
		// Wait until password is recieved
		wg.Wait()
	}
	var gotErr error
	for _, branch := range branchesToReset {
		if branch == "master" {
			continue
		}
		logger.Debug("Resetting branch %s", branch)

		cmd = exec.Command("git", "reset", "--hard", "HEAD")
		_, err = cmd.Output()
		if err != nil {
			logger.Debug("Cannot reset" + branch)
			continue
		}

		cmd = exec.Command("git", "checkout", branch)
		stdout, err = cmd.Output()
		if err != nil {
			logger.Debug("Cannot checkout" + branch)
			continue
		}

		cmd = exec.Command("git", "reset", "--hard", "HEAD")
		_, err = cmd.Output()
		if err != nil {
			logger.Debug("Cannot reset" + branch)
			continue
		}

		cmd = exec.Command("git", "log", "--pretty=format:'|%B|'", "-n", "1", "origin/"+branch)
		stdout, err = cmd.Output()
		if err != nil {
			logger.Debug("Cannot get log " + branch)
			continue
		}
		foo := strings.Split(string(stdout), "|")
		decrypted, _ := DecryptString(strings.TrimSpace(strings.Replace(foo[1], "'", "", -1)), Passphrase)
		logger.Debug("Decrypted message:%s", decrypted)
		if strings.TrimSpace(decrypted) == "deleted" {
			cmd := exec.Command("git", "rebase")
			out2, _ := cmd.StderrPipe()
			cmd.Start()
			out2b, _ := ioutil.ReadAll(out2)
			cmd.Wait()
			logger.Debug("git rebase : " + string(out2b))
		} else {
			logger.Debug("pulling...")
			cmd = exec.Command("git", "pull")
			out2, _ := cmd.StdoutPipe()
			cmd.Start()
			out2b, _ := ioutil.ReadAll(out2)
			cmd.Wait()
			logger.Debug("git pull : " + string(out2b))
			if strings.Contains(string(out2b), "Merge conflict") {
				files, _ := ioutil.ReadDir("./")
				for _, f := range files {
					if strings.Contains(f.Name(), ".cache") {
						continue
					}
					bText, _ := ioutil.ReadFile(f.Name())
					logger.Debug("TEXT: [%s]", string(bText))
					logger.Debug("Merging branch %s", branch)
					fmt.Printf("Merging branch %s\n", DecryptOTP(branch))
					merged := MergeEncrypted(string(bText), Passphrase)
					if len(merged) > 0 {
						ioutil.WriteFile(f.Name(), []byte(merged), 0644)
					}
					EncryptFile(f.Name(), Passphrase)
				}
				logger.Debug("committing (just in case)...")
				cmd = exec.Command("git", "commit", "-am", "'merged'")
				_, err = cmd.Output()
			}
		}

		cmd = exec.Command("git", "checkout", "master")
		stdout, err = cmd.Output()
		if err != nil {
			logger.Debug("Cannot checkout master")
			continue
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

	return gotErr
}

func NewDocument(gitfolder string, documentname string, fulltext string, message string, datestring string, branchNameOverride string) (string, error) {
	defer timeTrack(time.Now(), "NewDocument")
	var err error
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(gitfolder)

	newBranch := branchNameOverride
	if len(branchNameOverride) == 0 {
		newBranch = GenerateEntryName()
	}
	newBranch = newBranch

	// Encrypt everything
	documentname = EncryptOTP(documentname)
	newBranch = EncryptOTP(newBranch)
	message = EncryptString(message, Passphrase)

	cmd := exec.Command("git", "reset", "--hard", "HEAD")
	cmd.Output()

	cmd = exec.Command("git", "checkout", "--orphan", newBranch)
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
	err = EncryptFile(documentname, Passphrase)
	if err != nil {
		return newBranch, err
	}

	cmd = exec.Command("git", "add", documentname)
	_, err = cmd.Output()
	if err != nil {
		return newBranch, errors.New("Cannot add file " + documentname)
	}

	cmd = exec.Command("git", "commit", "--date", datestring, "-m", "'"+message+"'", documentname)
	out2, _ := cmd.StderrPipe()
	cmd.Start()
	out2b, _ := ioutil.ReadAll(out2)
	cmd.Wait()
	if strings.Contains(string(out2b), "error") || strings.Contains(string(out2b), "***") {
		fmt.Println(string(out2b))
		return newBranch, errors.New("Cannot commit " + documentname + " error: " + err.Error())
	}

	logger.Debug("Updated document %s in branch %s", documentname, newBranch)
	_, errExistence := GetTextOfOne(gitfolder, "master", "."+documentname)
	if errExistence != nil {
		WriteToMaster(gitfolder, "."+documentname, "new file")
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
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	if !exists(folder) {
		logger.Debug("Cloning %s into directory at %s", remote, folder)

		cmd := exec.Command("git", "clone", remote, folder)
		_, err = cmd.Output()
		if err != nil {
			fmt.Println("Cloning failed, will not continue")
			os.Exit(-1)
			return errors.New("Cloning command failed: '" + strings.Join([]string{"git", "clone", remote, folder}, " ") + "'")
		}
	}

	Fetch(folder)

	return nil
}
