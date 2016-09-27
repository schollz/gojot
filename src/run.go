package gitsdees

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

func Run() {

	// Check if cloning needs to occur
	measureTime := time.Now()
	fmt.Print("Fetching latest")
	if !exists(RemoteFolder) {
		logger.Debug("Remote folder does not exist: %s", RemoteFolder)
		Clone(RemoteFolder, Remote)
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
	if len(InputDocument) == 0 {
		var editDocument string
		fmt.Printf("\nCurrently available documents: ")
		logger.Debug("Last documents was %s", CurrentDocument)
		availableFiles := ListFiles(RemoteFolder)
		for _, file := range availableFiles {
			fmt.Printf("\n- %s ", file)
			if file == CurrentDocument {
				fmt.Print("(default) ")
			}
		}
		fmt.Print("\n\nWhat is the name of the document you want to edit (enter for default)? ")
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

	if Encrypt {
		Passphrase = PromptPassword(RemoteFolder, CurrentDocument)
	}
	cache, _, err := UpdateCache(RemoteFolder, CurrentDocument, false)
	if err != nil {
		logger.Error("Error updating cache: %s", err.Error())
		return
	}

	// if Summarize {
	// 	SummarizeCache(cache)
	// }

	texts := []string{}
	var branchHashes map[string]string
	if All || Export || Summarize || len(Search) > 0 {
		texts, branchHashes = CombineEntries(cache)
		if len(Search) > 0 {
			textFoo := []string{}
			for i := range texts {
				if strings.Contains(texts[i], Search) {
					textFoo = append(textFoo, texts[i])
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
