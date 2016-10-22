package sdees

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func IsItDocumentOrEntry(doc string) (bool, string, string) {
	availableFiles := ListFiles(RemoteFolder)
	for _, file := range availableFiles {
		if doc == file {
			return true, doc, ""
		}
	}
	branchList, _ := ListBranches(RemoteFolder)
	for _, branch := range branchList {
		if branch == doc {
			doc, _ := ListFileOfOne(RemoteFolder, branch)
			logger.Debug("You've entered a entry %s which is in document %s", branch, doc)
			return true, doc, branch
		}
	}
	return false, "", ""
}

func GoDelete() {
	// Check if user is deleting a entry or a document
	if len(InputDocument) == 0 {
		fmt.Printf("Which document would you like to delete? ")
		fmt.Scanln(&InputDocument)
	}
	gotOne, doc, entry := IsItDocumentOrEntry(InputDocument)
	fmt.Println(gotOne, doc, entry)
}

func GoDeleteEntry(entry string, cache Cache) {
	var yesno string
	fmt.Printf("Are you sure you want to delete the entry %s in document '%s'? (y/n) ", HashIDToString(entry), HashIDToString(entry))
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		deleteSuccess := false
		for _, branch := range cache.Branch {
			if branch.Branch == entry {
				err := DeleteBranch(entry)
				deleteSuccess = true
				if err == nil {
					fmt.Printf("Deleted entry %s\n", HashIDToString(entry))
				} else {
					fmt.Printf("Error deleting %s, does it exist?\n", HashIDToString(entry))
				}
			}
		}
		if !deleteSuccess {
			fmt.Printf("Error deleting %s, it does not exist\n", HashIDToString(entry))
		}
	} else {
		fmt.Printf("Did not delete %s\n", HashIDToString(entry))
	}
}

func GoDeleteDocument(document string, cache Cache) error {
	var yesno string
	fmt.Printf("Are you sure you want to delete the document %s? (y/n) ", HashIDToString(document))
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		for _, branch := range cache.Branch {
			err := Delete(RemoteFolder, branch.Branch)
			if err != nil {
				logger.Debug(err.Error())
			}
			if err == nil {
				fmt.Printf("Deleted entry %s\n", HashIDToString(branch.Branch))
			} else {
				fmt.Printf("Error deleting %s\n", HashIDToString(branch.Branch))
			}
		}
	} else {
		fmt.Printf("Did not delete %s\n", HashIDToString(document))
	}

	logger.Debug("Deleting cache")
	err := DeleteCache()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}

	logger.Debug("Deleting master index file: .%s", document)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(RemoteFolder)

	// Make sure we aren't on that branch
	cmd := exec.Command("git", "checkout", "master")
	_, err = cmd.Output()
	if err != nil {
		return errors.New("Problem switching to master")
	}

	// Remove file from index
	logger.Debug("git rm -f '.%s'", document)
	cmd = exec.Command("git", "rm", "-f", "."+document)
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
