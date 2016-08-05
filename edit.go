package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	entries := parseEntries(string(fileContents))
	for _, entry := range entries {
		writeEntry(entry, true)
	}
}

func parseEntries(text string) []string {
	entry := ""
	entries := make(map[int]string)
	var gt int
	for _, line := range strings.Split(text, "\n") {
		entryWords := strings.Split(strings.TrimSpace(line), " ")
		if len(entryWords) > 1 {
			t1, e1 := time.Parse("2006-01-02 15:04:05", entryWords[0]+" "+entryWords[1])
			t2, e2 := time.Parse("2006-01-02 15:04", entryWords[0]+" "+entryWords[1])
			if e1 == nil || e2 == nil {
				if len(entry) > 0 {
					if _, ok := entries[gt]; ok {
						logger.Warn("Duplicate entry for %s", entryWords[0]+" "+entryWords[1])
					}
					entries[gt] = strings.TrimSpace(entry)
				}
				entry = ""
				if e1 == nil {
					gt = int(t1.Unix())
				} else {
					gt = int(t2.Unix())
				}
			}
		}
		entry += strings.TrimRight(line, " ") + "\n"
	}
	if len(entry) > 0 {
		entries[gt] = strings.TrimSpace(entry)
	}

	// Sort the entries in order
	var keys []int
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	entriesInOrder := []string{}
	for _, k := range keys {
		entriesInOrder = append(entriesInOrder, entries[k])
	}
	logger.Debug("Parsed %d entries.", len(entriesInOrder))
	return entriesInOrder
}

func cleanUp() error {
	logger.Debug("Cleaning...")
	dir := RuntimeArgs.TempPath
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func editEntry() string {
	logger.Debug("Editing file")
	err := ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "vimrc"), []byte(`func! WordProcessorModeCLI()
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
com! WPCLI call WordProcessorModeCLI()`), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmdArgs := []string{"-u", path.Join(RuntimeArgs.TempPath, "vimrc"), "-c", "WPCLI", "+startinsert", path.Join(RuntimeArgs.TempPath, "temp")}
	cmd := exec.Command("vim", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.TempPath, "temp"))
	cleanUp()
	return string(fileContents)
}

func writeEntry(fileContents string, forceWrite bool) {
	logger.Debug("Entry contains %d bytes.", len(fileContents))
	if len(fileContents) < 22 && !forceWrite {
		logger.Info("No data appended.")
		return
	}
	// Hash contents to get filename
	h := sha1.New()
	h.Write([]byte(fileContents))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	fileName := string(sha1_hash) + ".gpg"
	if exists(path.Join(RuntimeArgs.FullPath, fileName)) {
		return
	}

	encryptedText := encryptString(string(fileContents), RuntimeArgs.Passphrase)
	err := ioutil.WriteFile(path.Join(RuntimeArgs.FullPath, fileName), []byte(encryptedText), 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Wrote %s.", fileName)
}
