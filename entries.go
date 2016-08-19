package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Gets the whole document, using the latest version of each entry
// Defaults to cache if available
func getFullEntry() (string, []string) {
	defer timeTrack(time.Now(), "Got full entry")
	type CachedDoc struct {
		Files      []string
		Entries    []string
		Timestamps []int
	}

	fullEntry := ""
	cache := CachedDoc{[]string{}, []string{}, []int{}}
	allEntries := []string{}
	gts := []int{}
	allFiles := readAllFiles()
	// if cache does not exist
	if !exists(path.Join(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".cache.json"))) {
		entryModifiedDates := make(map[string]int)
		entryStrings := make(map[string]string)
		logger.Debug("No cache.")
		for _, file := range allFiles {
			if strings.Contains(file, ".pass") {
				continue
			}
			foo := strings.Split(file, "/")
			fileName := foo[len(foo)-1]
			info := strings.Split(fileName, ".")
			modifiedTimestamp := decodeNumber(info[2])
			if val, ok := entryModifiedDates[info[0]]; ok {
				if modifiedTimestamp > val {
					entryModifiedDates[info[0]] = modifiedTimestamp
					entryStrings[info[0]] = decrypt(file) + "\n"
				}
			} else {
				entryModifiedDates[info[0]] = modifiedTimestamp
				entryStrings[info[0]] = decrypt(file) + "\n"
			}
			cache.Files = append(cache.Files, file)
		}
		wholeText := ""
		for key := range entryStrings {
			wholeText += entryStrings[key]
		}
		allEntries, gts = parseEntries(wholeText)
	} else {
		logger.Debug("Using cache.")
		fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".cache.json"))
		decryptedFileContents, err := decryptString(string(fileContents), RuntimeArgs.Passphrase)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshal JSON
		err = json.Unmarshal([]byte(decryptedFileContents), &cache)
		if err != nil {
			log.Fatal(err)
		}

		// Make set of files
		hasFile := make(map[string]bool)
		for _, file := range cache.Files {
			hasFile[file] = true
		}

		// If file doesn't exist in cache, add it
		// and then determine all individual entries
		cache.Files = []string{}
		for _, file := range allFiles {
			if strings.Contains(file, ".pass") {
				continue
			}
			cache.Files = append(cache.Files, file)
			if _, ok := hasFile[file]; !ok {
				logger.Debug("New entry %s.", file)
				newEntry, _ := ioutil.ReadFile(file)
				newEntryDecoded, err := decryptString(string(newEntry), RuntimeArgs.Passphrase)
				if err != nil {
					log.Fatal(err)
				}
				text, gt := parseEntries(newEntryDecoded)
				cache.Entries = append(cache.Entries, text...)
				cache.Timestamps = append(cache.Timestamps, gt...)
			}
		}
		entries := make(map[int]string)
		for i, entry := range cache.Entries {
			entries[cache.Timestamps[i]] = entry
		}
		allEntries, gts = sortEntries(entries)
	}

	// Cache the entries for next time
	cache.Entries = []string{}
	cache.Timestamps = []int{}
	for i, entry := range allEntries {
		fullEntry += entry + "\n\n"
		cache.Entries = append(cache.Entries, entry)
		cache.Timestamps = append(cache.Timestamps, gts[i])
	}
	cacheJson, _ := json.Marshal(cache)
	encryptedCacheJson := encryptString(string(cacheJson), RuntimeArgs.Passphrase)
	err := ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".cache.json"), []byte(encryptedCacheJson), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Return the full Entry and the individual entries
	return strings.TrimSpace(fullEntry), cache.Entries

}

// parseEntries is used to parse the full text for any entry and return all those
// entries and their corresponding unix epoch datetimes
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

// sortEntries takes a map of entries and returns a list of entries and a list
// of their dates in ascending order
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

	var cmdArgs []string

	if ConfigArgs.Editor == "vim" {
		// Setup vim
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

		cmdArgs = []string{"-c", "WPCLI", "+startinsert", path.Join(RuntimeArgs.TempPath, "temp")}
		if len(RuntimeArgs.TextSearch) > 0 {
			searchTerms := strings.Split(RuntimeArgs.TextSearch, " ")
			cmdArgs = append([]string{"-c", "2match Keyword /\\c\\v(" + strings.Join(searchTerms, "|") + ")/"}, cmdArgs...)
		}

	} else if ConfigArgs.Editor == "nano" {
		lines := strconv.Itoa(RuntimeArgs.Lines)
		cmdArgs = []string{"+" + lines + ",1000000", "--tempfile", path.Join(RuntimeArgs.TempPath, "temp")}
	} else if ConfigArgs.Editor == "emacs" {
		lines := strconv.Itoa(RuntimeArgs.Lines)
		cmdArgs = []string{"+" + lines + ":1000000", path.Join(RuntimeArgs.TempPath, "temp")}
	}

	// Run the editor
	cmd := exec.Command(ConfigArgs.Editor, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.TempPath, "temp"))
	return string(fileContents)
}

// writeEntry takes the contents of a file and writes the file in the
// specified format (see main.go, top)
// writing will be skipped if there is not much data, but it can be forced with forceWrite
func writeEntry(fileContents string, forceWrite bool) bool {
	// logger.Debug("Entry contains %d bytes.", len(fileContents))
	if len(fileContents) < 22 && !forceWrite {
		fmt.Println("No data appended.")
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
	fmt.Printf("Wrote %s.\n", fileName)
	return true
}
