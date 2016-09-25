package main

import (
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
		Fetch(RemoteFolder)
	}

	if Encrypt {
		Passphrase = PromptPassword(RemoteFolder, CurrentDocument)
	}
	cache, _, err := UpdateCache(RemoteFolder, CurrentDocument, false)
	if err != nil {
		logger.Error("Error updating cache: %s", err.Error())
		return
	}

	logger.Debug("Getting ready to edit %s", CurrentDocument)
	texts := []string{}
	var branchHashes map[string]string
	if All {
		texts, branchHashes = CombineEntries(cache)
	}
	texts = append(texts, HeadMatter(GetCurrentDate(), ""))
	ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	fulltext := WriteEntry()
	ProcessEntries(fulltext, branchHashes)
	Push(RemoteFolder)
}
