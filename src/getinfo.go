package gojot

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

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

	branchesToGet := make(map[string]bool)
	for _, branch := range branchNames {
		branchesToGet[branch] = true
	}

	entries := []Entry{}

	cmd := exec.Command("git", "log", "--name-only", "--pretty=format:'-==-%d -=-%H -=-%ad -=-%B -=-'", "--all")
	stdout, err := cmd.Output()
	if err != nil {
		logger.Error(`Couldn't run git log --name-only --pretty=format:'-==-%%d-=-%%H-=-%%ad-=-%%B-=-' --all`)
		return entries, errors.New("Problem running git log")
	}
	branchStrings := strings.Split(strings.Replace(string(stdout), "'", "", -1), "-==-")
	logger.Debug("Got info for %d branches, now will sift for %d branches", len(branchStrings), len(branchesToGet))
	for _, branchString := range branchStrings {
		items := strings.Split(branchString, "-=-")
		if len(items) < 5 {
			continue
		}
		var result Entry
		result.Branch = ""
		if strings.Contains(items[0], ",") {
			foo := strings.Split(items[0], ",")
			result.Branch = strings.TrimSpace(strings.Replace(foo[len(foo)-1], ")", "", -1))
		} else {
			result.Branch = strings.TrimSpace(strings.Replace(strings.Replace(items[0], "(", "", -1), ")", "", -1))
		}
		result.Branch = strings.TrimSpace(strings.Replace(result.Branch, "origin/", "", -1))
		if len(result.Branch) == 0 {
			continue
		}
		result.Hash = items[1]
		parsedDate, _ := ParseDate(items[2])
		result.Date = FormatDate(parsedDate)
		result.Message = strings.TrimSpace(items[3])
		result.Document = strings.TrimSpace(items[4])
		if _, ok := branchesToGet[result.Branch]; ok {
			entries = append(entries, result)
		}
	}

	return entries, nil
}
