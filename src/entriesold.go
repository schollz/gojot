package gitsdees

import "strings"

func ProcessEntriesOld(fulltext string) ([]string, []string) {
	type Blob struct {
		Date, Text string
	}

	var blobs []Blob
	var currentBlob Blob
	currentBlob.Text = ""
	for _, line := range strings.Split(fulltext, "\n") {
		splitLine := strings.Split(line, " ")
		if len(splitLine) >= 2 {
			possibleDate := strings.Join(splitLine[0:2], " ")
			parsedDate, err := ParseDate(possibleDate)
			if err == nil {
				if len(currentBlob.Date) > 0 {
					currentBlob.Text = strings.TrimSpace(currentBlob.Text)
					blobs = append(blobs, currentBlob)
				}
				currentBlob.Date = FormatDate(parsedDate)
				if len(splitLine) > 2 {
					currentBlob.Text = strings.Join(splitLine[2:], " ") + "\n"
				} else {
					currentBlob.Text = ""
				}
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
	}
	return texts, dates
}
