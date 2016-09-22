package main

import (
	"log"
	"os"
	"path"
)

func setup() {
	// Set the paths
	homeDir, _ := home.Dir()

	if !exists(path.Join(homeDir, ".cache")) {
		err := os.MkdirAll(path.Join(homeDir, ".cache"), 0711)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !exists(path.Join(homeDir, ".config")) {
		err := os.MkdirAll(path.Join(homeDir, ".config"), 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	RuntimeArgs.CachePath = path.Join(homeDir, ".cache", "gitsdees")
	if !exists(RuntimeArgs.CachePath) {
		err := os.MkdirAll(RuntimeArgs.CachePath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}
	RuntimeArgs.ConfigPath = path.Join(homeDir, ".config", "gitsdees")
	if !exists(RuntimeArgs.ConfigPath) {
		err := os.MkdirAll(RuntimeArgs.ConfigPath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !exists(path.Join(RuntimeArgs.ConfigPath, "synced.json")) {
		config()
	}
}
