package gitsdees

import (
	"strings"
	"time"
)

func ListFiles(gitfolder string) []string {
	defer timeTrack(time.Now(), "Listing files")
	branchNames, _ := ListBranches(gitfolder)
	infos, _ := GetInfo(gitfolder, branchNames)
	foundDocuments := make(map[string]bool)
	documents := []string{}
	for _, info := range infos {
		fileName := strings.Replace(info.Document, ".gpg", "", -1)
		if fileName == ".deleted" {
			continue
		}
		if _, ok := foundDocuments[fileName]; !ok {
			foundDocuments[fileName] = true
			documents = append(documents, fileName)
		}
	}
	return documents
}
