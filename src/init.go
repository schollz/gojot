package jot

import (
	"os"
	"path"

	home "github.com/mitchellh/go-homedir"
)

func init() {
	setupPaths()
}

func setupPaths() {
	// Set the paths
	homeDir, _ := home.Dir()
	HomePath = homeDir

	if !exists(path.Join(homeDir, ".cache")) {
		err := os.MkdirAll(path.Join(homeDir, ".cache"), 0711)
		if err != nil {
			logger.Error("Error creating %s", path.Join(homeDir, ".cache"))
		}
	}

	CachePath = path.Join(homeDir, ".cache", "jot")
	if !exists(CachePath) {
		err := os.MkdirAll(CachePath, 0711)
		if err != nil {
			logger.Error("Error creating %s", path.Join(homeDir, ".cache", "jot"))
		}
	}

	TempPath = path.Join(homeDir, ".cache", "jot", "temp")
	if !exists(TempPath) {
		err := os.MkdirAll(TempPath, 0711)
		if err != nil {
			logger.Error("Error creating %s", path.Join(homeDir, ".cache", "jot", "temp"))
		}
	}

	if !exists(path.Join(homeDir, ".config")) {
		err := os.MkdirAll(path.Join(homeDir, ".config"), 0711)
		if err != nil {
			logger.Error("Error creating %s", path.Join(homeDir, ".config"))
		}
	}

	ConfigPath = path.Join(homeDir, ".config", "jot")
	if !exists(ConfigPath) {
		err := os.MkdirAll(ConfigPath, 0711)
		if err != nil {
			logger.Error("Error creating %s", path.Join(homeDir, ".config", "jot"))
		}
	}

	if !exists(path.Join(ConfigPath, "config.json")) {
		SetupConfig()
	}
}
