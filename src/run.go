package gitsdees

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func Run() {

	// Check if cloning needs to occur
	if !exists(RemoteFolder) {
		logger.Debug("Remote folder does not exist: %s", RemoteFolder)
		Clone(RemoteFolder, Remote)
	} else {
		errFetch := Fetch(RemoteFolder)
		if errFetch == nil {
			fmt.Println("Fetched latest")
		} else {
			fmt.Println("No internet, not fetching")
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

	texts := []string{}
	var branchHashes map[string]string
	if All || Export {
		texts, branchHashes = CombineEntries(cache)
	}
	texts = append(texts, HeadMatter(GetCurrentDate(), MakeAlliteration()))
	if Export {
		fmt.Println("Exporting to " + CurrentDocument)
		ioutil.WriteFile(CurrentDocument, []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
		return
	} else {
		ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	}
	fulltext := WriteEntry()
	UpdateEntryFromText(fulltext, branchHashes)
	err = Push(RemoteFolder)
	if err == nil {
		fmt.Println("Pushed changes")
	} else {
		fmt.Println("No internet, not pushing")
	}
}
