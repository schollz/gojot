package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ResultText struct {
	Document, Branch, Text string
}

type JobText struct {
	Branch, Document string
}

func getTextWorker(id int, jobs <-chan JobText, results chan<- ResultText) {
	for job := range jobs {
		result := new(ResultText)
		result.Branch = job.Branch
		result.Document = job.Document

		cmd := exec.Command("git", "show", result.Branch+":"+result.Document)
		stdout, err := cmd.Output()
		if err != nil {
			logger.Error("git show %s:%s did not work", result.Branch, result.Document)
		}
		result.Text = strings.TrimSpace(string(stdout))

		results <- *result
	}
}

func getTextInParallel(inputs []JobText) []ResultText {
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan JobText, len(inputs))
	results := make(chan ResultText, len(inputs))
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
	entries := make([]ResultText, len(inputs))
	for a := 0; a < len(inputs); a++ {
		entries[a] = <-results
	}
	return entries
}

func GetText(folder string, branchNames []string, documentNames []string) ([]ResultText, error) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Getting Text %s", id, folder)
	defer timeTrack(time.Now(), "["+id+"]Getting Text "+folder)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []ResultText{}, errors.New("Cannot chdir into " + folder)
	}

	inputs := make([]JobText, len(documentNames))
	for i := range documentNames {
		inputs[i].Branch = branchNames[i]
		inputs[i].Document = documentNames[i]
	}
	entries := getTextInParallel(inputs)
	return entries, nil
}
