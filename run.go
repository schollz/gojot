package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func test() {

}

func listFiles() []string {
	files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath))
	fileNames := []string{}
	for _, f := range files {
		fileNameSplit := strings.Split(f.Name(), "/")
		fileName := fileNameSplit[len(fileNameSplit)-1]
		if fileName == "config.json" || fileName == "temp" {
			continue
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}

func getFullEntry() string {
	fullEntry := ""
	wholeText := decryptAll()
	allEntries := parseEntries(wholeText)
	for _, entry := range allEntries {
		fullEntry += entry + "\n\n"
	}
	return strings.TrimSpace(fullEntry)
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
		os.Exit(-1)
	}

	if !RuntimeArgs.EditLocally && HasInternetAccess() {
		syncDown()
	}

	promptPassword()

	fullEntry := ""
	if RuntimeArgs.EditWhole {
		fullEntry = getFullEntry()
		if len(fullEntry) > 0 {
			fullEntry += "\n\n"
		}
	}

	t := time.Now()
	fullEntry += string(t.Format("2006-01-02 15:04:05")) + "  "
	err = ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "temp"), []byte(fullEntry), 0644)
	if err != nil {
		log.Fatal(err)
	}

	newEntry := editEntry()
	entries := parseEntries(newEntry)
	for _, entry := range entries {
		writeEntry(entry, false)
	}

	if !RuntimeArgs.EditLocally && HasInternetAccess() {
		syncUp()
	}

}
