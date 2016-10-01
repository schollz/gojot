package sdees

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

func Run() {

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

	// List available documents to choose from
	availableFiles, encrypted := ListFiles(RemoteFolder)
	if len(InputDocument) == 0 {
		var editDocument string
		fmt.Printf("\nCurrently available documents: ")
		logger.Debug("Last documents was %s", CurrentDocument)
		for i, file := range availableFiles {
			fmt.Printf("\n- %s ", file)
			if encrypted[i] {
				fmt.Print("[encrypted] ")
			}
			if file == CurrentDocument {
				fmt.Print("(default) ")
			}
		}
		if len(CurrentDocument) == 0 {
			CurrentDocument = "notes.txt"
		}
		fmt.Printf("\n\nWhich document (press enter for '%s', or type name): ", CurrentDocument)
		fmt.Scanln(&editDocument)
		if len(editDocument) == 0 && len(CurrentDocument) > 0 {
			// Pass
		} else if len(editDocument) == 0 && len(availableFiles) > 0 {
			CurrentDocument = availableFiles[0]
		} else if len(CurrentDocument) == 0 && len(editDocument) == 0 && len(availableFiles) == 0 {
			CurrentDocument = "notes.txt"
		} else if len(editDocument) > 0 {
			CurrentDocument = editDocument
		}
	} else {
		CurrentDocument = InputDocument
	}
	// Save choice of current document
	SaveConfiguration(Editor, Remote, CurrentDocument)

	// Check if encryption is needed
	isNew := true
	Encrypt = false
	for i, file := range availableFiles {
		if CurrentDocument == file {
			isNew = false
			Encrypt = encrypted[i]
			break
		}
	}
	if isNew {
		// Prompt whether encryption is wanted for new files
		var yesencryption string
		fmt.Print("\nDo you want to add encryption (default: y)? (y/n) ")
		fmt.Scanln(&yesencryption)
		if yesencryption == "n" {
			Encrypt = false
		} else {
			Encrypt = true
		}
	} else if !All && !Summarize && !Export {
		// Prompt for whether to load whole document
		var yesnoall string
		fmt.Print("\nLoad all entries (press enter for 'n')? (y/n) ")
		fmt.Scanln(&yesnoall)
		if yesnoall == "y" {
			All = true
		}
	}

	// Prompt for passphrase if encrypted
	if Encrypt {
		Passphrase = PromptPassword(RemoteFolder, CurrentDocument)
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
	var branchHashes map[string]string
	if All || Export || Summarize || len(Search) > 0 {
		texts, branchHashes = CombineEntries(cache)
		// Conduct the search
		if len(Search) > 0 {
			searchWords := GetWordsFromText(Search)
			textFoo := []string{}
			for i := range texts {
				for _, searchWord := range searchWords {
					if strings.Contains(texts[i], searchWord) {
						textFoo = append(textFoo, texts[i])
						break
					}
				}
			}
			texts = textFoo
		}
	}

	// Case-switch for what to do with fulltext
	if Export {
		fmt.Println("Exporting to " + CurrentDocument)
		ioutil.WriteFile(CurrentDocument, []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
		return
	} else if Summarize {
		fmt.Println("\nSummary:")
		fmt.Println(SummarizeEntries(texts))
		return
	} else {
		texts = append(texts, HeadMatter(GetCurrentDate(), MakeAlliteration()))
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
