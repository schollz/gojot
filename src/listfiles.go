package gojot

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
		if len(document) < 1 {
			continue
		}
		if document[0] == '.' {
			document = DecryptOTP(document[1:])
			logger.Debug("Found document: %s", document)
			if document == ".deleted" || document == "jot-key" || document == "jot-hs" {
				continue
			}
			documents = append(documents, document)
		}
	}

	sort.Strings(documents)
	return documents
}

func ListFilesOfOne(gitfolder string, branch string) []string {
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
	documents := []string{}
	for _, document := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if document[0] == '.' {
			continue
		}
		documents = append(documents, document)
	}

	return documents
}
