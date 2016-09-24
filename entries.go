package main

import (
	"fmt"
	"strings"
)

func ProcessEntries(fulltext string) []string {
	var branchesUpdated []string
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
		if blob.Branch == "NEW" {
			if len(blob.Text) < 10 {
				fmt.Println("No new data, not commiting new document.")
				continue
			}
			logger.Debug("Writing new entry for " + blob.Branch)
			newBranch, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, "")
			branchesUpdated = append(branchesUpdated, newBranch)
			if err != nil {
				logger.Error(err.Error())
			}
		} else if blob.Hash != GetMD5Hash(blob.Text) {
			logger.Debug("Updating entry for " + blob.Branch)
			_, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, blob.Branch)
			branchesUpdated = append(branchesUpdated, blob.Branch)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}

	return branchesUpdated
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
