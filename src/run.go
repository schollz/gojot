package sdees

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

func Run() {
	// Some variables to be set later
	filterBranch := ""

	// Check if cloning needs to occur
	logger.Debug("Current remote: %s", Remote)
	measureTime := time.Now()
	fmt.Print("Fetching latest")
	if !exists(RemoteFolder) {
		logger.Debug("Remote folder does not exist: %s", RemoteFolder)
		err := Clone(RemoteFolder, Remote)
		if err != nil {
			logger.Warn("Problems cloning remote '%s': %s", Remote, err.Error())
		}
	} else {
		errFetch := Fetch(RemoteFolder)
		if errFetch == nil {
			fmt.Print("...done")
		} else {
			fmt.Println("...unable to fetch:")
			fmt.Println(errFetch.Error())
		}
	}
	fmt.Printf(" (%s)\n", time.Since(measureTime).String())

	// Prompt for passphrase
	Passphrase = PromptPassword(RemoteFolder)

	// List available documents to choose from
	availableFiles := ListFiles(RemoteFolder)
	if len(InputDocument) == 0 {
		var editDocument string
		fmt.Printf("\nCurrently available documents: ")
		logger.Debug("Last documents was %s", HashIDToString(CurrentDocument))
		for _, file := range availableFiles {
			fmt.Printf("\n- %s ", HashIDToString(file))
			if file == CurrentDocument {
				fmt.Print("(default) ")
			}
		}
		if len(CurrentDocument) == 0 {
			CurrentDocument = StringToHashID("notes.txt")
		}
		fmt.Printf("\n\nWhich document (press enter for '%s', or type name): ", HashIDToString(CurrentDocument))
		fmt.Scanln(&editDocument)
		if len(editDocument) == 0 && len(CurrentDocument) > 0 {
			// Pass
		} else if len(editDocument) == 0 && len(availableFiles) > 0 {
			CurrentDocument = availableFiles[0]
		} else if len(CurrentDocument) == 0 && len(editDocument) == 0 && len(availableFiles) == 0 {
			CurrentDocument = StringToHashID("notes.txt")
		} else if len(editDocument) > 0 {
			CurrentDocument = StringToHashID(editDocument)
		}
	} else {
		branchList, _ := ListBranches(RemoteFolder)
		for _, branch := range branchList {
			if branch == StringToHashID(InputDocument) {
				doc, _ := ListFileOfOne(RemoteFolder, branch)
				logger.Debug("You've entered a branch %s which is in document %s", branch, doc)
				InputDocument = doc
				filterBranch = branch
			}
		}
		CurrentDocument = InputDocument
	}
	logger.Debug("Current document: %s", CurrentDocument)
	// Save choice of current document
	SaveConfiguration(Editor, Remote, CurrentDocument)

	// Check if encryption is needed
	isNew := true
	for _, file := range availableFiles {
		if CurrentDocument == file {
			isNew = false
			break
		}
	}
	if !isNew && !All && !Summarize && !Export && !DeleteDocument && len(DeleteEntry) == 0 && len(filterBranch) == 0 && len(Search) == 0 {
		// Prompt for whether to load whole document
		var yesnoall string
		fmt.Print("\nLoad all entries (press enter for 'n')? (y/n) ")
		fmt.Scanln(&yesnoall)
		if yesnoall == "y" {
			All = true
		}
	}

	// Update the cache using the passphrase if needed
	cache, _, err := UpdateCache(RemoteFolder, CurrentDocument, false)
	if err != nil {
		logger.Error("Error updating cache: %s", err.Error())
		return
	}

	// Do deletions
	if DeleteDocument {
		GoDeleteDocument(cache)
		return
	} else if len(DeleteEntry) > 0 {
		GoDeleteEntry(cache)
		return
	}

	// Load fulltext
	texts := []string{}
	textsBranch := []string{}
	var branchHashes map[string]string
	if All || Export || Summarize || len(Search) > 0 || len(filterBranch) > 0 {
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
		fmt.Println("Exporting to " + CurrentDocument)
		ioutil.WriteFile(CurrentDocument, []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
		return
	} else if Summarize {
		fmt.Println("\nSummary:")
		fmt.Println(SummarizeEntries(texts, textsBranch))
		return
	} else {
		if len(filterBranch) == 0 {
			texts = append(texts, HeadMatter(GetCurrentDate(), MakeAlliteration()))
		} else {
			fmt.Printf("Loaded entry '%s' on document '%s'\n", filterBranch, CurrentDocument)
		}
		ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	}
	fulltext := WriteEntry()
	UpdateEntryFromText(fulltext, branchHashes)

	// Push new changes
	measureTime = time.Now()
	fmt.Print("Pushing changes")
	err = Push(RemoteFolder)
	if err == nil {
		fmt.Print("...done")
	} else {
		fmt.Print("...no internet, not pushing")
	}
	fmt.Printf(" (%s)\n", time.Since(measureTime).String())
}
