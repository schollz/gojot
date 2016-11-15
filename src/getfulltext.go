package jot

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getTextWorker(id int, jobs <-chan Entry, results chan<- Entry) {
	for job := range jobs {
		result := new(Entry)
		result.Branch = job.Branch
		result.Document = job.Document
		result.Date = job.Date
		result.Hash = job.Hash
		result.Message = job.Message

		cmd := exec.Command("git", "show", result.Branch+":"+result.Document)
		stdout, err := cmd.Output()
		if err != nil {
			logger.Error("git show %s:%s did not work", result.Branch, result.Document)
		}
		result.Text = strings.TrimSpace(string(stdout))
		result.Text, _ = DecryptString(result.Text, Passphrase)

		results <- *result
	}
}

func getTextInParallel(inputs []Entry) []Entry {
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan Entry, len(inputs))
	results := make(chan Entry, len(inputs))
	//This starts up 50 workers, initially blocked because there are no jobs yet.
	for w := 0; w < 50; w++ {
		go getTextWorker(w, jobs, results)
	}
	//Here we send len(branchNames) jobs and then close that channel to indicate that’s all the work we have.
	for _, name := range inputs {
		jobs <- name
	}
	close(jobs)
	//Finally we collect all the results of the work.
	entries := make([]Entry, len(inputs))
	for a := 0; a < len(inputs); a++ {
		entries[a] = <-results
	}
	return entries
}

func GetText(folder string, entries []Entry) ([]Entry, error) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Getting Text %s", id, folder)
	defer timeTrack(time.Now(), "["+id+"]Getting Text "+folder)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []Entry{}, errors.New("Cannot chdir into " + folder)
	}

	return getTextInParallel(entries), nil
}

func GetTextOfOne(gitfolder string, branch string, document string) (string, error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(gitfolder)
	if err != nil {
		return "", errors.New("Cannot chdir into " + gitfolder)
	}
	cmd := exec.Command("git", "show", branch+":"+document)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}
