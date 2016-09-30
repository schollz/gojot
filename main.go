package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	home "github.com/mitchellh/go-homedir"
)

var (
	Version, BuildTime, Build, OS, Program string
)

func main() {
	// Set the paths
	Program = "sdees"
	homeDir, _ := home.Dir()
	if !exists(path.Join(homeDir, ".cache")) {
		err := os.MkdirAll(path.Join(homeDir, ".cache"), 0711)
		if err != nil {
			log.Fatal("Could not create cache path: " + path.Join(homeDir, ".cache"))
		}
	}
	programPath := path.Join(homeDir, ".cache", "sdees-binary")
	if !exists(programPath) {
		err := os.MkdirAll(programPath, 0711)
		if err != nil {
			log.Fatal("Could not create program path: " + programPath)
		}
	}

	if !exists(path.Join(programPath, Program)) {
		data, err := Asset("bin/" + Program)
		if err == nil {
			err = ioutil.WriteFile(path.Join(programPath, Program), data, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Could not extract program: '" + Program + "'")
		}
	}
	cmd := exec.Command(path.Join(programPath, Program))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
