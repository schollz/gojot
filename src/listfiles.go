package sdees

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func ListFiles(gitfolder string) ([]string, []bool) {
	defer timeTrack(time.Now(), "Listing files")
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(gitfolder)
	if err != nil {
		logger.Error("Cannot chdir into " + gitfolder)
	}

	cmd := exec.Command("git", "ls-tree", "--name-only", "master")
	stdout, err := cmd.Output()
	if err != nil {
		logger.Error("Problem doing ls-tree")
	}
	documents := []string{}
	encrypted := []bool{}
	for _, document := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if document[0] == '.' {
			continue
		}
		encrypted = append(encrypted, strings.Contains(document, ".gpg"))
		documents = append(documents, strings.Replace(document, ".gpg", "", -1))
	}
	return documents, encrypted
}
