package gitsdees

import "fmt"

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

	err := DeleteCache()
	if err != nil {
		logger.Debug(err.Error())
		return err
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
