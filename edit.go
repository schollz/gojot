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
	"time"
)

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

	t := time.Now()
	dateString := string(t.Format(time.RFC3339)) + "   "
	err = ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "temp"), []byte(dateString), 0644)
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

func writeEntry(fileContents string, forceWrite bool) string {
	if len(fileContents) < 32 && !forceWite {
		return ""
	}
	// Hash contents to get filename
	h := sha1.New()
	h.Write(fileContents)
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	fileName := string(sha1_hash) + ".gpg"

	encryptedText := encryptString(string(fileContents), getPassword())
	err := ioutil.WriteFile(path.Join(RuntimeArgs.FullPath, fileName), []byte(encryptedText), 0644)
	return fileName
}
