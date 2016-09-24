package main

import (
	"fmt"
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
	cache, _ := UpdateCache(RemoteFolder, false)

	logger.Debug("Getting ready to edit %s", CurrentDocument)
	texts := []string{}
	if All {
		texts = CombineEntries(cache)
	}
	texts = append(texts, HeadMatter(GetCurrentDate(), "NEW", RandStringBytesMaskImprSrc(6, time.Now().UnixNano())))
	ioutil.WriteFile(path.Join(TempPath, "temp"), []byte(strings.Join(texts, "\n\n")+"\n"), 0644)
	fulltext := WriteEntry()
	ProcessFullText(fulltext)
	// NewDocument(RemoteFolder, CurrentDocument, fulltext, GetMessage(fulltext), GetCurrentDate(), "")
	// Push(RemoteFolder)
}

func ProcessFullText(fulltext string) {
	type Blob struct {
		Date, Branch, Hash, Text string
	}

	var blobs []Blob
	var currentBlob Blob
	for _, line := range strings.Split(fulltext, "\n") {
		if strings.Count(line, " -==- ") == 2 && len(strings.Split(line, " -==- ")) == 3 {
			if len(currentBlob.Hash) > 0 {
				currentBlob.Text = strings.TrimSpace(currentBlob.Text)
				blobs = append(blobs, currentBlob)
				currentBlob.Text = ""
			}
			items := strings.Split(line, " -==- ")
			currentBlob.Date = strings.TrimSpace(items[0])
			currentBlob.Branch = strings.TrimSpace(items[1])
			currentBlob.Hash = strings.TrimSpace(items[2])
		} else {
			currentBlob.Text = currentBlob.Text + line + "\n"
		}
	}
	if len(currentBlob.Hash) > 0 {
		currentBlob.Text = strings.TrimSpace(currentBlob.Text)
		blobs = append(blobs, currentBlob)
	}
	for _, blob := range blobs {
		fmt.Println(blob.Hash, GetMD5Hash(blob.Text), blob.Hash == GetMD5Hash(blob.Text))
		if blob.Branch == "NEW" {
			if len(blob.Text) < 10 {
				continue
			}
			logger.Debug("Writing new entry for " + blob.Branch)
			_, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, "")
			if err != nil {
				logger.Error(err.Error())
			}
		} else if blob.Hash != GetMD5Hash(blob.Text) {
			logger.Debug("Updating entry for " + blob.Branch)
			_, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, blob.Branch)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}

}

func HeadMatter(date string, branch string, text string) string {
	return date + " -==- " + branch + " -==- " + GetMD5Hash(text) + "\n\n"
}

func GetMessage(m string) string {
	ms := strings.Split(m, " ")
	if len(ms) < 18 {
		return strings.Join(ms, " ")
	} else {
		return strings.Join(ms[:18], " ")
	}
}
