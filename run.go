package main

import (
	"io/ioutil"
	"path"
	"strings"
	"time"
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
	if All {
		texts = CombineEntries(cache)
	}
	texts = append(texts, HeadMatter(GetCurrentDate(), " NEW ", RandStringBytesMaskImprSrc(5, time.Now().UnixNano())))
	ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	fulltext := WriteEntry()
	ProcessEntries(fulltext)
	Push(RemoteFolder)
}
