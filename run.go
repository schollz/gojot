// run.go handles the main functionality after the CLI flags are determined

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// Imports a file into a document, flag --import
func importFile(filename string) {
	promptPassword()
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("%v", err)
		os.Exit(-1)
	}
	entries, _ := parseEntries(string(time.Now().Format("2006-01-02 15:04:05")) + "\n" + string(fileContents))
	for _, entry := range entries {
		writeEntry(entry, true)
	}
	fmt.Printf("Imported '%s' to %s.", filename, ConfigArgs.WorkingFile)
}

// Exports a file into a document, flag --export
func exportFile(filename string) {
	promptPassword()
	fullText, _ := getFullEntry()
	err := ioutil.WriteFile(filename, []byte(fullText), 0644)
	if err != nil {
		logger.Error("%v", err)
		os.Exit(-1)
	}
	fmt.Printf("Exported '%s' to %s.", ConfigArgs.WorkingFile, filename)
}

// Prompt for password (cross-compatiable, except cygwin)
func promptPassword() {
	possibleFiles := getEntryList()
	password1 := "1"
	if len(possibleFiles) == 0 {
		password2 := "2"
		for password1 != password2 {
			fmt.Printf("Enter password for editing '%s': ", ConfigArgs.WorkingFile)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 = strings.TrimSpace(string(bytePassword2))
			if password1 != password2 {
				fmt.Println("\nPasswords do not match.")
			}
		}
	} else {
		testFile := possibleFiles[0]
		logger.Debug("Testing with %s", path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile, testFile))
		fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile, testFile))
		passwordAccepted := false
		for passwordAccepted == false {
			fmt.Printf("Enter password for editing '%s': ", ConfigArgs.WorkingFile)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			_, err := decryptString(string(fileContents), password1)
			if err == nil {
				passwordAccepted = true
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		}
	}
	RuntimeArgs.Passphrase = password1
	fmt.Println("")
}

// Run the Syncing, and Editing
func run() {
	logger.Debug("Available files: %s", strings.Join(listFiles(), ", "))

	// // Check if editor exists
	// _, err := exec.Command(ConfigArgs.Editor, "--version").Output()
	// if err != nil {
	// 	if ConfigArgs.Editor == "vim" {
	// 		fmt.Println(`You need to download vim. If your using Unix:
	//
	// 	apt-get install vim
	//
	// If you're using Windows:
	//
	// 	wget ftp://ftp.vim.org/pub/vim/pc/vim74w32.zip
	// 	unzip vim74w32.zip
	// 	mv vim/vim74/vim.exe ./
	// `)
	// 	} else {
	// 		fmt.Printf("You need to download %s or switch editors using `sdees --config`.\n", RuntimeArgs.Editor)
	// 	}
	// 	return
	// }

	// Pull latest copies
	logger.Debug("RuntimeArgs.DontSync: %v", RuntimeArgs.DontSync)
	if !RuntimeArgs.DontSync && !RuntimeArgs.OnlyPush {
		if HasInternetAccess() {
			syncDown()
		} else {
			fmt.Println("Unable to pull, no internet access.")
		}
	}

	// Get password for access to GPG-encryption
	promptPassword()

	// Get current entry if needed
	fullEntry := ""
	if (len(RuntimeArgs.TextSearch) == 0 && RuntimeArgs.EditWhole) || len(RuntimeArgs.NumberToShow) > 0 {
		// Get full entry
		_, allEntries := getFullEntry()
		totalEntries := len(allEntries)
		numberToShow := totalEntries
		if len(RuntimeArgs.NumberToShow) > 0 {
			numberToShow, _ = strconv.Atoi(RuntimeArgs.NumberToShow)
			logger.Debug("Showing latest %d of %d entries.", numberToShow, totalEntries)
		}
		numberToShow = numberToShow + 1
		for i, entry := range allEntries {
			if i > totalEntries-numberToShow {
				fullEntry += entry + "\n\n"
			}
		}
	} else if len(RuntimeArgs.TextSearch) > 0 {
		// Get only entries that match search terms
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
	RuntimeArgs.Lines = len(strings.Split(fullEntry, "\n"))

	if RuntimeArgs.Summarize {
		// If summarizing, use only the first lines
		_, entries := getFullEntry()
		totalEntries := len(entries)
		numberToShow := totalEntries + 10
		if len(RuntimeArgs.NumberToShow) > 0 {
			numberToShow, _ = strconv.Atoi(RuntimeArgs.NumberToShow)
		}
		for i, entry := range entries {
			if i > totalEntries-numberToShow {
				lines := strings.Split(entry, "\n")
				fullEntry += lines[0] + "\n"
			}
		}
	} else {
		// Add the timestamp for the new entry
		t := time.Now()
		fullEntry += string(t.Format("2006-01-02 15:04:05")) + " "
		if ConfigArgs.Editor == "vim" {
			fullEntry += " "
		}
	}
	// Write the data contents to the tempfile
	err := ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "temp"), []byte(fullEntry), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Edit the entry
	newEntry := editEntry()

	if !RuntimeArgs.Summarize {
		// Parse and save the new entry
		entries, _ := parseEntries(newEntry)
		totalNewWords := 0
		for _, entry := range entries {
			if writeEntry(entry, false) {
				totalNewWords = totalNewWords + len(strings.Split(entry, " ")) - 2
			}
		}
		if totalWords > 1 && totalNewWords > 0 {
			fmt.Printf("+%d words. %s total.\n", totalNewWords, Comma(int64(totalWords)))
		} else if totalNewWords > 0 {
			fmt.Printf("+%d words.\n", totalNewWords)
		}
	}

	// Sync it back up
	if !RuntimeArgs.DontSync || RuntimeArgs.OnlyPush {
		if HasInternetAccess() {
			syncUp()
		} else {
			fmt.Println("Unable to push, no internet access.")
		}
	}
	return
}
