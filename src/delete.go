package sdees

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func GoDeleteEntry(cache Cache) {
	var yesno string
	fmt.Printf("Are you sure you want to delete the entry %s in document %s? (y/n) ", DeleteEntry, CurrentDocument)
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		deleteSuccess := false
		for _, branch := range cache.Branch {
			if branch.Branch == DeleteEntry {
				err := DeleteBranch(DeleteEntry)
				deleteSuccess = true
				if err == nil {
					fmt.Printf("Deleted entry %s\n", DeleteEntry)
				} else {
					fmt.Printf("Error deleting %s, does it exist?\n", DeleteEntry)
				}
			}
		}
		if !deleteSuccess {
			fmt.Printf("Error deleting %s, it does not exist\n", DeleteEntry)
		}
	} else {
		fmt.Printf("Did not delete %s\n", DeleteEntry)
	}
}

func GoDeleteDocument(cache Cache) error {
	var yesno string
	fmt.Printf("Are you sure you want to delete the document %s? (y/n) ", CurrentDocument)
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		for _, branch := range cache.Branch {
			err := Delete(RemoteFolder, branch.Branch)
			if err != nil {
				logger.Debug(err.Error())
				return err
			}
			if err == nil {
				fmt.Printf("Deleted entry %s\n", branch.Branch)
			} else {
				fmt.Printf("Error deleting %s\n", branch.Branch)
			}
		}
	} else {
		fmt.Printf("Did not delete %s\n", CurrentDocument)
	}

	logger.Debug("Deleting cache")
	err := DeleteCache()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}

	logger.Debug("Deleting master index file: %s", CurrentDocument)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(RemoteFolder)

	// Make sure we aren't on that branch
	cmd := exec.Command("git", "checkout", "master")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem switching to master")
	}

	document := CurrentDocument
	// Remove file from index
	logger.Debug("git rm -f %s", document)
	cmd = exec.Command("git", "rm", "-f", document)
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem git rm -f ")
	}

	logger.Debug("git commit -m %s", document)
	cmd = exec.Command("git", "commit", "-m", "removed '"+document+"'")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem git commit -m ")
	}

	fmt.Print("Deleting on remote")
	err = Push(RemoteFolder)
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	fmt.Println("...done")
	return nil
}
