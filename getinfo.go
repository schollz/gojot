package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Entry struct {
	Document, Branch, Date, Hash, Message, Fulltext string
}

func getInfoWorker(id int, jobs <-chan string, results chan<- Entry) {
	for branch := range jobs {
		result := new(Entry)
		result.Branch = branch
		// cmd := exec.Command("git", "ls-tree", "--name-only", branch)
		// stdout, err := cmd.Output()
		// if err != nil {
		// 	logger.Error("git ls-tree --name-only %s did not work", branch)
		// }
		// result.Document = strings.TrimSpace(string(stdout))

		// cmd = exec.Command("git", "show", branch+":"+result.Document)
		// stdout, err = cmd.Output()
		// if err != nil {
		// 	logger.Error("git show %s:%s did not work", branch, result.Document)
		// }
		result.Fulltext = "" //strings.TrimSpace(string(stdout))

		cmd := exec.Command("git", "log", "--name-only", "--pretty=format:'%H-=-%ad-=-%s'", branch)
		stdout, err := cmd.Output()
		if err != nil {
			logger.Error(`Couldn't run git log --pretty=format:'%%H-=-%%ad-=-%%s'` + branch)
		}
		lines := strings.Split(string(stdout), "\n")
		items := strings.Split(lines[0], "-=-")
		result.Document = strings.TrimSpace(lines[1])
		result.Hash = items[0]
		result.Date = items[1]
		result.Message = items[2][1 : len(items[2])-1]

		results <- *result
	}
}

func getInfoInParallel(branchNames []string) []Entry {
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan string, len(branchNames))
	results := make(chan Entry, len(branchNames))
	//This starts up 50 workers, initially blocked because there are no jobs yet.
	for w := 0; w < 50; w++ {
		go getInfoWorker(w, jobs, results)
	}
	//Here we send len(branchNames) jobs and then close that channel to indicate that’s all the work we have.
	for _, name := range branchNames {
		jobs <- name
	}
	close(jobs)
	//Finally we collect all the results of the work.
	entries := make([]Entry, len(branchNames))
	for a := 0; a < len(branchNames); a++ {
		entries[a] = <-results
	}
	return entries
}

func GetInfo(folder string, branchNames []string) ([]Entry, error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []Entry{}, errors.New("Cannot chdir into " + folder)
	}
	start := time.Now()
	entries := getInfoInParallel(branchNames)
	logger.Debug("Got info from %d branches in %s", len(entries), time.Since(start).String())
	return entries, nil
}
