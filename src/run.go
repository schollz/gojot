package sdees

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func Run() {
	// Some variables to be set later
	filterBranch := ""
	Passphrase = "alskdfjalskdjfalskjdflajsdfljasd"

	// Load the configuration
	LoadConfiguration()

	// Check if cloning needs to occur
	logger.Debug("Current remote: %s", Remote)
	logger.Debug("Current remote folder: %s", RemoteFolder)
	if !exists(RemoteFolder) {
		logger.Debug("Remote folder does not exist: %s", RemoteFolder)
		err := Clone(RemoteFolder, Remote)
		if err != nil {
			logger.Warn("Problems cloning remote '%s': %s", Remote, err.Error())
		}
	} else {
		logger.Debug("Remote folder does exist: %s", RemoteFolder)
		errFetch := Fetch(RemoteFolder)
		if errFetch != nil {
			fmt.Println("Unable to fetch latest:")
			fmt.Println(errFetch.Error())
		}
	}

	// If deleting, Delete
	if DeleteFlag {
		GoDelete()
		return
	}

	// If importing, import
	if ImportFlag || ImportOldFlag {
		if len(InputDocument) == 0 {
			fmt.Println("Must enter name of file to import")
			return
		}
		var err1 error
		if ImportFlag {
			err1 = Import(InputDocument)
		} else if ImportOldFlag {
			err1 = ImportOld(InputDocument)
		}
		if err1 != nil {
			logger.Error(err1.Error())
		}
		return
	}

	// List available documents to choose from
	availableFiles := ListFiles(RemoteFolder)
	if ShowStats {
		DisplayStats(availableFiles)
		return
	}
	if len(InputDocument) == 0 {
		var editDocument string
		logger.Debug("Last documents was %s", (CurrentDocument))
		data := [][]string{}
		for fileNum, file := range availableFiles {
			quickCache, errCache := LoadCache(RemoteFolder, EncryptOTP(file))
			entryString := "N/A"
			if errCache == nil {
				entryString = Comma(int64(len(quickCache.Branch)))
			}
			data = append(data, []string{strconv.Itoa(fileNum + 1), file, entryString})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Document", "# Entries"})
		for _, v := range data {
			table.Append(v)
		}
		fmt.Printf("\n")
		table.Render()

		if len(CurrentDocument) == 0 {
			CurrentDocument = "notes.txt"
		}
		fmt.Printf("\n\nWhich document (press enter for '%s', or type name): ", DecryptOTP(CurrentDocument))
		fmt.Scanln(&editDocument)
		if len(editDocument) == 0 && len(CurrentDocument) > 0 {
			// Pass
		} else if len(editDocument) == 0 && len(availableFiles) > 0 {
			CurrentDocument = availableFiles[0]
		} else if len(CurrentDocument) == 0 && len(editDocument) == 0 && len(availableFiles) == 0 {
			CurrentDocument = ("notes.txt")
		} else if len(editDocument) > 0 {
			fileNum, err := strconv.Atoi(strings.TrimSpace(editDocument))
			if err == nil {
				CurrentDocument = availableFiles[fileNum-1]
			} else {
				CurrentDocument = strings.TrimSpace(editDocument)
			}
		}
	} else {
		InputDocument = EncryptOTP(InputDocument)
		branchList, _ := ListBranches(RemoteFolder)
		for _, branch := range branchList {
			if branch == InputDocument {
				for _, doc := range ListFilesOfOne(RemoteFolder, branch) {
					logger.Debug("You've entered a branch %s which is in document %s", DecryptOTP(branch), DecryptOTP(doc))
					InputDocument = doc
					filterBranch = branch
				}
			}
		}
		CurrentDocument = InputDocument
	}
	CurrentDocument = EncryptOTP(CurrentDocument)
	logger.Debug("Current document: %s", DecryptOTP(CurrentDocument))
	// Save choice of current document
	SaveConfiguration(Editor, Remote, CurrentDocument)

	// Check if encryption is needed
	isNew := true
	for _, file := range availableFiles {
		if CurrentDocument == EncryptOTP(file) {
			isNew = false
			break
		}
	}

	if !isNew && !All && !Summarize && !Export && len(filterBranch) == 0 && len(Search) == 0 {
		// Prompt for whether to load whole document
		var yesnoall string
		fmt.Print("\nLoad all entries (press enter for 'n')? (y/n) ")
		fmt.Scanln(&yesnoall)
		if yesnoall == "y" {
			All = true
		}
	}

	// Load fulltext
	texts := []string{}
	textsBranch := []string{}
	var branchHashes map[string]string
	if All || Export || Summarize || len(Search) > 0 || len(filterBranch) > 0 {
		// Update the cache
		cache, _, err := UpdateCache(RemoteFolder, CurrentDocument, false)
		if err != nil {
			logger.Error("Error updating cache: %s", err.Error())
			return
		}

		texts, textsBranch, branchHashes = CombineEntries(cache)
		// Conduct the search
		if len(Search) > 0 {
			searchWords := GetWordsFromText(Search)
			textFoo := []string{}
			for i := range texts {
				for _, searchWord := range searchWords {
					if strings.Contains(strings.ToLower(texts[i]), strings.ToLower(searchWord)) {
						textFoo = append(textFoo, texts[i])
						break
					}
				}
			}
			texts = textFoo
		}
		if len(filterBranch) > 0 {
			for i, branch := range textsBranch {
				if branch == filterBranch {
					logger.Debug("Filtering out everything but branch %s", filterBranch)
					texts = []string{texts[i]}
					textsBranch = []string{textsBranch[i]}
				}
			}
		}
	}

	// Case-switch for what to do with fulltext
	if Export {
		fmt.Println("Exporting to " + DecryptOTP(CurrentDocument))
		ioutil.WriteFile(DecryptOTP(CurrentDocument), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
		return
	} else if Summarize {
		fmt.Println("\nSummary:")
		fmt.Println(SummarizeEntries(texts, textsBranch))
		return
	} else {
		if len(filterBranch) == 0 {
			texts = append(texts, HeadMatter(GetCurrentDate(), GenerateEntryName()))
		} else {
			logger.Debug("Loaded entry '%s' on document '%s'\n", filterBranch, CurrentDocument)
			fmt.Printf("Loaded entry '%s' on document '%s'\n", DecryptOTP(filterBranch), DecryptOTP(CurrentDocument))
		}
		ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	}
	fulltext := WriteEntry()
	UpdateEntryFromText(fulltext, branchHashes)

	// Push new changes
	measureTime := time.Now()
	fmt.Print("Pushing changes")
	err := Push(RemoteFolder)
	if err == nil {
		fmt.Print("...done")
	} else {
		fmt.Print("...no internet, not pushing")
	}
	fmt.Printf(" (%s)\n", time.Since(measureTime).String())
}
