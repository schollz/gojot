package sdees

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func ImportOld(filename string) error {
	if Encrypt {
		Passphrase = PromptPassword(RemoteFolder)
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Error reading file: %s", err.Error())
		return err
	}
	texts, dates := ProcessEntriesOld(string(data))
	for i := range texts {
		_, err = NewDocument(RemoteFolder, CurrentDocument, texts[i], GetMessage(texts[i]), dates[i], "")
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

func ProcessEntriesOld(fulltext string) ([]string, []string) {
	type Blob struct {
		Date, Text string
	}

	var blobs []Blob
	var currentBlob Blob
	currentBlob.Text = ""
	for _, line := range strings.Split(fulltext, "\n") {
		splitLine := strings.Split(line, " ")
		if len(splitLine) >= 1 {
			var possibleDate string
			var newText string
			var parsedDate time.Time
			err := errors.New("No date")
			if len(splitLine) > 1 {
				possibleDate = strings.Join(splitLine[0:2], " ")
				parsedDate, err = ParseDate(possibleDate)
				newText = strings.Join(splitLine[2:], " ") + "\n"
			}
			if err != nil {
				possibleDate = splitLine[0]
				parsedDate, err = ParseDate(possibleDate)
				newText = strings.Join(splitLine[1:], " ") + "\n"
			}
			if err == nil {
				if len(currentBlob.Date) > 0 {
					currentBlob.Text = strings.TrimSpace(currentBlob.Text)
					blobs = append(blobs, currentBlob)
				}
				currentBlob.Date = FormatDate(parsedDate)
				currentBlob.Text = newText
			} else {
				currentBlob.Text += line
			}
		}
	}
	if len(currentBlob.Date) > 0 {
		currentBlob.Text = strings.TrimSpace(currentBlob.Text)
		blobs = append(blobs, currentBlob)
	}

	texts := make([]string, len(blobs))
	dates := make([]string, len(blobs))
	for i, blob := range blobs {
		texts[i] = blob.Text
		dates[i] = blob.Date
		fmt.Printf("Got entry on date %s\n", blob.Date)
	}
	return texts, dates
}
