package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func test() {

}

func printFileList() {
	fmt.Println("Available documents:\n")
	for i, f := range listFiles() {
		files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, f))
		if len(files) > 0 {
			fmt.Printf("[%d] %s (%d entries)\n", i, f, len(files))
		}
	}
	fmt.Print("\n")
}

func listFiles() []string {
	files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath))
	fileNames := []string{}
	for _, f := range files {
		fileNameSplit := strings.Split(f.Name(), "/")
		fileName := fileNameSplit[len(fileNameSplit)-1]
		if fileName == "config.json" || fileName == "temp" || strings.Contains(fileName, ".cache") {
			continue
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}

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
		logger.Debug("No cache.")
		wholeText := ""
		for _, file := range allFiles {
			wholeText += decrypt(file) + "\n"
			cache.Files = append(cache.Files, file)
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

		// If file doesn't exist, add it
		cache.Files = []string{}
		for _, file := range allFiles {
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
	return strings.TrimSpace(fullEntry), cache.Entries

}

func promptPassword() {
	// Get password for working file
	passwordAccepted := false
	for passwordAccepted == false {
		fmt.Printf("Enter password for editing '%s': ", ConfigArgs.WorkingFile)
		bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password := strings.TrimSpace(string(bytePassword))
		RuntimeArgs.Passphrase = password
		if exists(path.Join(RuntimeArgs.FullPath, ConfigArgs.WorkingFile+".pass")) {
			// Check old password
			fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.FullPath, ConfigArgs.WorkingFile+".pass"))
			err := CheckPasswordHash(string(fileContents), password)
			if err == nil {
				passwordAccepted = true
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		} else {
			// Generate new passwrod
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 := strings.TrimSpace(string(bytePassword2))
			if password == password2 {
				// Write password to file
				passwordAccepted = true
				passwordHashed, _ := HashPassword(password)
				err := ioutil.WriteFile(path.Join(RuntimeArgs.FullPath, ConfigArgs.WorkingFile+".pass"), passwordHashed, 0644)
				if err != nil {
					log.Fatal("Could not write to file.")
				}
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		}
	}
	fmt.Println("")
}

func run() {
	logger.Debug("Available files: %s", strings.Join(listFiles(), ", "))

	// Check if VIM exists
	_, err := exec.Command("vim", "--version").Output()
	if err != nil {
		fmt.Println(`You need to download vim. If your using Unix:

	apt-get install vim

If you're using Windows:

	wget ftp://ftp.vim.org/pub/vim/pc/vim74w32.zip
	unzip vim74w32.zip
	mv vim/vim74/vim.exe ./
`)
		return
	}

	if !RuntimeArgs.DontSync {
		if HasInternetAccess() {
			syncDown()
		} else {
			logger.Info("Unable to pull, no internet access.")
		}
	}

	promptPassword()

	if RuntimeArgs.Summarize {
		_, entries := getFullEntry()
		totalEntries := len(entries)
		numberToShow := totalEntries
		if len(RuntimeArgs.NumberToShow) > 0 {
			numberToShow, _ = strconv.Atoi(RuntimeArgs.NumberToShow)
		}
		for i, entry := range entries {
			if i > totalEntries-numberToShow {
				lines := strings.Split(entry, "\n")
				fmt.Println(lines[0])
			}
		}
		return
	}

	fullEntry := ""
	if len(RuntimeArgs.TextSearch) == 0 && RuntimeArgs.EditWhole {
		fullEntry, _ = getFullEntry()
		if len(fullEntry) > 0 {
			fullEntry += "\n\n"
		}
	} else if len(RuntimeArgs.TextSearch) > 0 {
		searchTerms := strings.Split(RuntimeArgs.TextSearch, " ")
		for i := range searchTerms {
			searchTerms[i] = " " + searchTerms[i]
		}
		logger.Debug("Search terms: %v", searchTerms)
		_, entries := getFullEntry()
		for _, entry := range entries {
			shouldAdd := true
			for _, term := range searchTerms {
				if !strings.Contains(strings.ToLower(entry), strings.ToLower(term)) {
					shouldAdd = false
					break
				}
			}
			if shouldAdd {
				fullEntry += entry + "\n\n"
			}
		}
	}
	totalWords := len(strings.Split(fullEntry, " "))

	t := time.Now()
	fullEntry += string(t.Format("2006-01-02 15:04:05")) + "  "
	err = ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "temp"), []byte(fullEntry), 0644)
	if err != nil {
		log.Fatal(err)
	}

	newEntry := editEntry()
	entries, _ := parseEntries(newEntry)
	totalNewWords := 0
	for _, entry := range entries {
		if writeEntry(entry, false) {
			totalNewWords = totalNewWords + len(strings.Split(entry, " ")) - 2
		}
	}
	if totalWords > 1 && totalNewWords > 0 {
		logger.Info("+%d words. %d total.", totalNewWords, totalWords)
	} else if totalNewWords > 0 {
		logger.Info("+%d words.", totalNewWords)
	}

	if !RuntimeArgs.DontSync {
		if HasInternetAccess() {
			syncUp()
		} else {
			logger.Info("Unable to push, no internet access.")
		}
	}
	return

}
