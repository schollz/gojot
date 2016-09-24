package main

import (
	"path"
	"strings"
)

func Run() {

	// Load configuration
	LoadConfiguration()

	// Check if cloning needs to occur
	RemoteFolder = path.Join(CachePath, HashString(Remote))
	if !exists(RemoteFolder) {
		logger.Debug("Remote folder does not exist: %s", RemoteFolder)
		Clone(RemoteFolder, Remote)
	} else {
		Fetch(RemoteFolder)
	}
	// cache, _ := UpdateCache(RemoteFolder, false)

	logger.Debug("Getting ready to edit %s", CurrentDocument)
	fulltext := WriteEntry()
	NewDocument(RemoteFolder, CurrentDocument, fulltext, GetMessage(fulltext), GetCurrentDate(), "")
	Push(RemoteFolder)
}

func GetMessage(m string) string {
	ms := strings.Split(m, " ")
	if len(ms) < 18 {
		return strings.Join(ms, " ")
	} else {
		return strings.Join(ms[:18], " ")
	}
}
