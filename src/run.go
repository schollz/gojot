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
			fmt.Print("...no internet, not fetching")
		}
	}
	fmt.Printf(" (%s)\n", time.Since(measureTime).String())

	// Get files
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
	SaveConfiguration(Editor, Remote, CurrentDocument)

	if !All && !Summarize && !Export {
		var yesnoall string
		fmt.Print("\nLoad all entries (press enter for 'n')? (y/n) ")
		fmt.Scanln(&yesnoall)
		if yesnoall == "y" {
			All = true
		}
	}

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
		var yesencryption string
		fmt.Print("\nDo you want to add encryption (default: y)? (y/n) ")
		fmt.Scanln(&yesencryption)
		if yesencryption == "n" {
			Encrypt = false
		} else {
			Encrypt = true
		}
	}
	if Encrypt {
		Passphrase = PromptPassword(RemoteFolder, CurrentDocument)
	}
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
