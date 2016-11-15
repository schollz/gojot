package jot

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func IsItDocumentOrEntry(doc string) (bool, string, string) {
	availableFiles := ListFiles(RemoteFolder)
	for _, file := range availableFiles {
		logger.Debug(doc, file)
		if EncryptOTP(doc) == EncryptOTP(file) {
			return true, EncryptOTP(doc), ""
		}
	}
	branchList, _ := ListBranches(RemoteFolder)
	for _, branch := range branchList {
		if EncryptOTP(branch) == EncryptOTP(doc) {
			for _, doc := range ListFilesOfOne(RemoteFolder, branch) {
				logger.Debug("You've entered a entry '%s' which is in document '%s'", branch, doc)
				return true, EncryptOTP(doc), EncryptOTP(branch)
			}
		}
	}
	return false, "", ""
}

func GoDelete() {
	// Check if user is deleting a entry or a document
	if len(InputDocument) == 0 {
		fmt.Printf("Which document or entry would you like to delete? ")
		fmt.Scanln(&InputDocument)
	}
	gotOne, document, entry := IsItDocumentOrEntry(InputDocument)
	if !gotOne {
		fmt.Printf("%s is not a document or entry, did you type it correctly?\n", DecryptOTP(InputDocument))
		return
	}

	// Get the cache
	cache, _, err := UpdateCache(RemoteFolder, document, false)
	if err != nil {
		logger.Error("Error updating cache: %s", err.Error())
		return
	}
	if len(entry) == 0 {
		GoDeleteDocument(document, cache)
	} else {
		GoDeleteEntry(document, entry, cache)
	}
	CurrentDocument = document
	DeleteCache()
}

func GoDeleteEntry(document string, entry string, cache Cache) {
	var yesno string
	fmt.Printf("Are you sure you want to delete the entry '%s' in document '%s'? (y/n) ", DecryptOTP(entry), DecryptOTP(document))
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		deleteSuccess := false
		for _, branch := range cache.Branch {
			if branch.Branch == entry {
				err := DeleteBranch(entry)
				deleteSuccess = true
				if err == nil {
					fmt.Printf("Deleted entry %s\n", DecryptOTP(entry))
				} else {
					fmt.Printf("Error deleting %s, does it exist?\n", DecryptOTP(entry))
				}
			}
		}
		if !deleteSuccess {
			fmt.Printf("Error deleting %s, it does not exist\n", DecryptOTP(entry))
		}
	} else {
		fmt.Printf("Did not delete %s\n", DecryptOTP(entry))
	}
}

func GoDeleteDocument(document string, cache Cache) error {
	var yesno string
	fmt.Printf("Are you sure you want to delete the document %s? (y/n) ", DecryptOTP(document))
	fmt.Scanln(&yesno)
	if string(yesno) == "y" {
		for _, branch := range cache.Branch {
			err := Delete(RemoteFolder, branch.Branch)
			if err != nil {
				logger.Debug(err.Error())
			}
			if err == nil {
				fmt.Printf("Deleted entry %s\n", DecryptOTP(branch.Branch))
			} else {
				fmt.Printf("Error deleting %s\n", DecryptOTP(branch.Branch))
			}
		}
	} else {
		fmt.Printf("Did not delete %s\n", DecryptOTP(document))
	}

	logger.Debug("Deleting master index file: .%s", document)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(RemoteFolder)

	// Make sure we aren't on that branch
	cmd := exec.Command("git", "checkout", "master")
	_, err := cmd.Output()
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
