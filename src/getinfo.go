package sdees

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getInfoWorker(id int, jobs <-chan string, results chan<- Entry) {
	for branch := range jobs {
		result := new(Entry)
		result.Branch = branch

		// cmd = exec.Command("git", "show", branch+":"+result.Document)
		// stdout, err = cmd.Output()
		// if err != nil {
		// 	logger.Error("git show %s:%s did not work", branch, result.Document)
		// }
		//strings.TrimSpace(string(stdout))

		cmd := exec.Command("git", "log", "--name-only", "--pretty=format:'%H-=-%ad-=-%s'", branch)
		stdout, err := cmd.Output()
		if err != nil {
			logger.Error(`Couldn't run git log --pretty=format:'%%H-=-%%ad-=-%%s'` + branch)
		}
		lines := strings.Split(strings.Replace(string(stdout), "'", "", -1), "\n")
		items := strings.Split(lines[0], "-=-")
		result.Document = strings.TrimSpace(lines[1])
		result.Hash = items[0]
		result.Date = items[1]
		result.Message = strings.TrimSpace(items[2])

		results <- *result
	}
}

func getInfoInParallel(branchNames []string) []Entry {
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan string, len(branchNames))
	results := make(chan Entry, len(branchNames))
	//This starts up 50 workers, initially blocked because there are no jobs yet.
	for w := 0; w < 100; w++ {
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
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Getting info %s", id, folder)
	defer timeTrack(time.Now(), "["+id+"]Getting info "+folder)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []Entry{}, errors.New("Cannot chdir into " + folder)
	}

	entries := getInfoInParallel(branchNames)
	return entries, nil
}
