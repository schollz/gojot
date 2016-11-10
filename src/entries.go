package sdees

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func Import(filename string) error {
	fmt.Printf("Importing %s\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Error reading file: %s", err.Error())
		return err
	}
	blobs := ProcessEntries(string(data))
	if len(blobs) < 1 {
		fmt.Println("No entries found, perhaps you should run with --importold")
	}
	for i := range blobs {
		fmt.Println("Importing entry " + DecryptOTP(blobs[i].Branch))
		_, err = NewDocument(RemoteFolder, filename, blobs[i].Text, GetMessage(blobs[i].Text), blobs[i].Date, blobs[i].Branch)
		if err != nil {
			logger.Error("Error creating new document: %s", err.Error())
		}
	}
	err = Push(RemoteFolder)
	if err == nil {
		fmt.Println("Pushed changes")
	} else {
		fmt.Println("No internet, not pushing")
	}
	return nil
}

type BlobEntry struct {
	Date, Branch, Hash, Text string
}

func UpdateEntryFromText(fulltext string, branchHashes map[string]string) []string {
	var branchesUpdated []string
	blobs := ProcessEntries(fulltext)
	for _, blob := range blobs {
		if _, ok := branchHashes[blob.Branch]; !ok {
			logger.Debug("Branch not present updating entry for branch %s ", blob.Branch)
			if len(blob.Text) < 10 {
				logger.Debug("No new data, not committing entry for branch %s\n", blob.Branch)
				fmt.Printf("No new data, not committing entry for branch %s\n", DecryptOTP(blob.Branch))
				continue
			}
			_, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, blob.Branch)
			branchesUpdated = append(branchesUpdated, blob.Branch)
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Debug("Created entry %s (+%d words)\n", blob.Branch, len(strings.Split(blob.Text, " ")))
				fmt.Printf("Created entry %s (+%d words)\n", DecryptOTP(blob.Branch), len(strings.Split(blob.Text, " ")))
			}
		} else if blob.Hash != branchHashes[blob.Branch] {
			logger.Debug("Current hash (%s) != Previous hash (%s), updating entry for %s ", blob.Hash, branchHashes[blob.Branch], blob.Branch)
			_, err := NewDocument(RemoteFolder, CurrentDocument, blob.Text, GetMessage(blob.Text), blob.Date, blob.Branch)
			branchesUpdated = append(branchesUpdated, blob.Branch)
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Debug("Updated entry for %s\n", blob.Branch)
				fmt.Printf("Updated entry for %s\n", DecryptOTP(blob.Branch))
			}
		}
	}

	return branchesUpdated
}

func ProcessEntries(fulltext string) []BlobEntry {
	var blobs []BlobEntry
	var currentBlob BlobEntry
	for _, line := range strings.Split(fulltext, "\n") {
		if strings.Count(line, " -==- ") == 1 && len(strings.Split(line, " -==- ")) == 2 {
			if len(currentBlob.Date) > 0 {
				currentBlob.Text = strings.TrimSpace(currentBlob.Text)
				currentBlob.Hash = GetMD5Hash(currentBlob.Text)
				blobs = append(blobs, currentBlob)
				currentBlob.Text = ""
			}
			items := strings.Split(line, " -==- ")
			currentBlob.Date = strings.TrimSpace(items[0])
			currentBlob.Branch = EncryptOTP(strings.TrimSpace(items[1]))
		} else {
			currentBlob.Text = currentBlob.Text + line + "\n"
		}
	}
	if len(currentBlob.Date) > 0 {
		currentBlob.Text = strings.TrimSpace(currentBlob.Text)
		currentBlob.Hash = GetMD5Hash(currentBlob.Text)
		blobs = append(blobs, currentBlob)
	}

	return blobs
}

func HeadMatter(date string, branch string) string {
	if len(branch) == 0 {
		branch = "NEW"
	}
	return date + " -==- " + DecryptOTP(branch) + "\n\n"
}

func GetMessage(m string) string {
	ms := strings.Split(m, " ")
	if len(ms) < 18 {
		return strings.Join(ms, " ")
	} else {
		return strings.Join(ms[:18], " ")
	}
}
