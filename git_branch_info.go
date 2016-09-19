package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getBranchesInfoWorker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		cmd := exec.Command("git", "log", "--pretty=format:'%H-=-%ad-=-%s'", j)
		stdout, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		results <- strings.TrimSpace(string(stdout))
	}
}

func getBranchesInfoInParallel() []string {
	start := time.Now()
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan string, 100)
	results := make(chan string, 100)
	//This starts up 50 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 50; w++ {
		go getBranchesInfoWorker(w, jobs, results)
	}
	//Here we send 9 jobs and then close that channel to indicate that’s all the work we have.
	for j := 0; j < 100; j++ {
		jobs <- strconv.Itoa(j)
	}
	close(jobs)
	//Finally we collect all the results of the work.
	resultStrings := make([]string, 100)
	for a := 0; a < 100; a++ {
		resultStrings[a] = <-results
	}
	elapsed := time.Since(start)
	log.Printf("testWorkersBranchInfo took %s", elapsed/101)
	return resultStrings
}

func GetBranchesInfo(folder string) ([]string, error) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(folder)
	if err != nil {
		return []string{}, errors.New("Cannot chdir into " + folder)
	}

	branches := getBranchesInfoInParallel()
	return branches, nil
}
