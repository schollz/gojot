package sdees

import (
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func ListFiles(gitfolder string) []string {
	defer timeTrack(time.Now(), "Listing files")
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(gitfolder)
	if err != nil {
		logger.Error("ListFiles Cannot chdir into " + gitfolder)
	}

	cmd := exec.Command("git", "ls-tree", "--name-only", "master")
	stdout, err := cmd.Output()
	if err != nil {
		logger.Error("Problem doing ls-tree")
	}
	documents := []string{}
	for _, document := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if document[0] == '.' && document != ".key" && document != ".new" && HashIDToString(document) != ".deleted" {
			if len(strings.TrimLeft(strings.TrimRight(document, "\x00"), "\x00")) < 2 {
				continue
			}
			logger.Debug("Found document: %s", document)
			documents = append(documents, strings.Replace(document, ".", "", -1))
		}
	}

	sort.Strings(documents)
	return documents
}

func ListFileOfOne(gitfolder string, branch string) (string, bool) {
	defer timeTrack(time.Now(), "Listing files")
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(gitfolder)
	if err != nil {
		logger.Error("Cannot chdir into " + gitfolder)
	}

	cmd := exec.Command("git", "ls-tree", "--name-only", branch)
	stdout, err := cmd.Output()
	if err != nil {
		logger.Error("Problem doing ls-tree")
	}
	docMap := make(map[string]bool)
	for _, document := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if document[0] == '.' {
			continue
		}
		docMap[strings.Replace(document, ".gpg", "", -1)] = strings.Contains(document, ".gpg")
	}

	documents := make([]string, len(docMap))
	i := 0
	for k, _ := range docMap {
		documents[i] = k
		i++
	}
	sort.Strings(documents)
	encrypted := make([]bool, len(documents))
	for i, val := range documents {
		encrypted[i] = docMap[val]
	}
	return documents[0], encrypted[0]
}
