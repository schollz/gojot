package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

func importFile(filename string) {
	promptPassword()
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("%v", err)
		os.Exit(-1)
	}
	entries, _ := parseEntries(string(fileContents))
	for _, entry := range entries {
		writeEntry(entry, true)
	}
	logger.Info("Imported '%s' to %s.", filename, ConfigArgs.WorkingFile)
}

func exportFile(filename string) {
	promptPassword()
	fullText, _ := getFullEntry()
	err := ioutil.WriteFile(filename, []byte(fullText), 0644)
	if err != nil {
		logger.Error("%v", err)
		os.Exit(-1)
	}
	logger.Info("Exported '%s' to %s.", ConfigArgs.WorkingFile, filename)
}

func parseEntries(text string) ([]string, []int) {
	defer timeTrack(time.Now(), "Parsing entries")
	entry := ""
	entries := make(map[int]string)
	var gt int
	for _, line := range strings.Split(text, "\n") {
		entryWords := strings.Split(strings.TrimSpace(line), " ")
		if len(entryWords) > 1 {
			isDate, new_gt := parseDate(entryWords[0] + " " + entryWords[1])
			if isDate {
				if len(entry) > 0 {
					if _, ok := entries[gt]; ok {
						logger.Debug("Duplicate entry for %s", entryWords[0]+" "+entryWords[1])
					} else {
						entries[gt] = strings.TrimSpace(entry)
					}
				}
				entry = ""
				gt = new_gt
			}
		}
		entry += strings.TrimRight(line, " ") + "\n"
	}
	if len(entry) > 0 {
		entries[gt] = strings.TrimSpace(entry)
	}
	if len(entries) == 1 {
		return append([]string{}, entries[gt]), append([]int{}, gt)
	}

	entriesInOrder, gtsInOrder := sortEntries(entries)
	return entriesInOrder, gtsInOrder
}

func sortEntries(entries map[int]string) ([]string, []int) {
	// Sort the entries in order
	var keys []int
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	entriesInOrder := []string{}
	gtsInOrder := []int{}
	for _, k := range keys {
		entriesInOrder = append(entriesInOrder, entries[k])
		gtsInOrder = append(gtsInOrder, k)
	}
	logger.Debug("Sorted %d entries.", len(entriesInOrder))
	return entriesInOrder, gtsInOrder
}

func editEntry() string {
	logger.Debug("Editing file")
	vimrc := `func! WordProcessorModeCLI()
		setlocal formatoptions=t1
		setlocal textwidth=80
		map j gj
		map k gk
		set formatprg=par
		setlocal wrap
		setlocal linebreak
		setlocal noexpandtab
		normal G$
endfu
com! WPCLI call WordProcessorModeCLI()`
	// Append to .vimrc file
	if exists(path.Join(RuntimeArgs.HomePath, ".vimrc")) {
		// Check if .vimrc file contains code
		logger.Debug("Found .vimrc.")
		fileContents, err := ioutil.ReadFile(path.Join(RuntimeArgs.HomePath, ".vimrc"))
		if err != nil {
			log.Fatal(err)
		}
		if !strings.Contains(string(fileContents), "com! WPCLI call WordProcessorModeCLI") {
			// Append to fileContents
			logger.Debug("WPCLI not found in .vimrc, adding it...")
			newvimrc := string(fileContents) + "\n" + vimrc
			err := ioutil.WriteFile(path.Join(RuntimeArgs.HomePath, ".vimrc"), []byte(newvimrc), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			logger.Debug("WPCLI found in .vimrc.")
		}
	} else {
		logger.Debug("Can not find .vimrc, creating new .vimrc...")
		err := ioutil.WriteFile(path.Join(RuntimeArgs.HomePath, ".vimrc"), []byte(vimrc), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmdArgs := []string{"-c", "WPCLI", "+startinsert", path.Join(RuntimeArgs.TempPath, "temp")}
	if len(RuntimeArgs.TextSearch) > 0 {
		searchTerms := strings.Split(RuntimeArgs.TextSearch, " ")
		cmdArgs = append([]string{"-c", "2match Keyword /\\c\\v(" + strings.Join(searchTerms, "|") + ")/"}, cmdArgs...)
	}
	cmd := exec.Command("vim", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.TempPath, "temp"))
	return string(fileContents)
}

func writeEntry(fileContents string, forceWrite bool) bool {
	// logger.Debug("Entry contains %d bytes.", len(fileContents))
	if len(fileContents) < 22 && !forceWrite {
		logger.Info("No data appended.")
		return false
	}

	// Hash date to get fileName
	dateString := ""
	for _, line := range strings.Split(fileContents, "\n") {
		s := strings.Split(line, " ")
		dateString = s[0] + " " + s[1]
		break
	}
	_, dateVal := parseDate(dateString)
	fileNameFrontMatter := encodeNumber(dateVal) + "." + hashString(fileContents)
	for _, file := range RuntimeArgs.CurrentFileList {
		if strings.Contains(file, fileNameFrontMatter) {
			// File already exists
			return false
		}
	}

	_, dateVal = parseDate(time.Now().Format("2006-01-02 15:04:05"))
	fileName := fileNameFrontMatter + "." + encodeNumber(dateVal) + ".gpg"
	encryptedText := encryptString(string(fileContents), RuntimeArgs.Passphrase)
	err := ioutil.WriteFile(path.Join(RuntimeArgs.FullPath, fileName), []byte(encryptedText), 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Wrote %s.", fileName)
	return true
}
