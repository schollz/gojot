package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// readAllFiles returns a list of all the files in the sdees path
func readAllFiles() []string {
	files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.FullPath))
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, path.Join(RuntimeArgs.FullPath, f.Name()))
	}
	return fileNames
}

// getEntryList returns a list of the GPG encrypted entries
// for the current working document
func getEntryList() []string {
	files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile))
	fileNames := []string{}
	for _, f := range files {
		fileNameSplit := strings.Split(f.Name(), "/")
		fileName := fileNameSplit[len(fileNameSplit)-1]
		if strings.Contains(fileName, ".gpg") {
			fileNames = append(fileNames, fileName)
		}
	}
	return fileNames
}

// getFileList returns a map of all files in the sdees home directory
// excluding cache, temp, and config files
func getFileList() map[int]string {
	filesIndexed := make(map[int]string)
	for i, f := range listFiles() {
		files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, f))
		if len(files) > 0 {
			filesIndexed[i] = f
		}
	}
	return filesIndexed
}

// printFileList prints out the available files
// available via --list
func printFileList() {
	fmt.Println("Available documents (access using `sdees NUM`):\n")
	for i, f := range listFiles() {
		files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, f))
		if len(files) > 0 {
			fmt.Printf("[%d] %s (%d entries)\n", i, f, len(files))
		}
	}
	fmt.Print("\n")
}

// listFiles returns a list of all the files in the sdees home directory
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

// cleanUp deletes all temporary files and also deletes documents that were
// made accidently (documents with no data)
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

	fileList := listFiles()
	for i, f := range fileList {
		files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, f))
		if len(files) < 2 {
			for _, file := range files {
				logger.Debug("Remove %s.", path.Join(RuntimeArgs.WorkingPath, f, file.Name()))
				err := os.Remove(path.Join(RuntimeArgs.WorkingPath, f, file.Name()))
				if err != nil {
					log.Fatal(err)
				}
			}
			logger.Debug("Remove %s.", path.Join(RuntimeArgs.WorkingPath, f))
			err := os.Remove(path.Join(RuntimeArgs.WorkingPath, f))
			if err != nil {
				log.Fatal(err)
			}
			if ConfigArgs.WorkingFile == f {
				if len(fileList) < 2 {
					ConfigArgs.WorkingFile = "notes.txt"
				} else {
					if i != 0 {
						ConfigArgs.WorkingFile = fileList[0]
					} else {
						ConfigArgs.WorkingFile = fileList[1]
					}
				}
				// Save current config parameters
				b, err := json.Marshal(ConfigArgs)
				if err != nil {
					log.Println(err)
				}
				ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

			}
		}
	}

	return nil
}

// parseDate parses the two possible date formats
func parseDate(s string) (bool, int) {
	t, e := time.Parse("2006-01-02 15:04:05", s)
	if e == nil {
		return true, int(t.Unix())
	}
	t, e = time.Parse("2006-01-02 15:04", s)
	if e == nil {
		return true, int(t.Unix())
	}
	return false, int(-1)
}
