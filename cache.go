package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	home "github.com/mitchellh/go-homedir"
)

func CleanFolderName(gitfolder string) string {
	return RemoteFolder
}

func init() {
	// Set the paths
	homeDir, _ := home.Dir()

	if !exists(path.Join(homeDir, ".cache")) {
		err := os.MkdirAll(path.Join(homeDir, ".cache"), 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	CachePath = path.Join(homeDir, ".cache", "gitsdees")
	if !exists(CachePath) {
		err := os.MkdirAll(CachePath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	TempPath = path.Join(homeDir, ".cache", "gitsdees", "temp")
	if !exists(TempPath) {
		err := os.MkdirAll(TempPath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func UpdateCache(gitfolder string, forceUpdate bool) (map[string]Entry, []string) {
	defer timeTrack(time.Now(), "Updating cache")
	cache := make(map[string]Entry)
	cacheFile := RemoteFolder + ".cache"
	logger.Debug("Using cacheFile: %s", cacheFile)
	branchNames, _ := ListBranches(gitfolder)
	entriesToUpdate := []Entry{} // which branches to update in cache
	entries, _ := GetInfo(gitfolder, branchNames)

	if !exists(cacheFile) || forceUpdate {
		// Generate new cache
		logger.Debug("Generating new cache")
		entriesToUpdate = entries
	} else {
		// Load current cache
		logger.Debug("Loading and updating cache")
		cache = LoadCache(gitfolder)
		for _, info := range entries {
			if _, ok := cache[info.Branch]; !ok {
				entriesToUpdate = append(entriesToUpdate, info)
				continue
			}
			if info.Hash != cache[info.Branch].Hash {
				entriesToUpdate = append(entriesToUpdate, info)
			}
		}
	}

	entriesToUpdate, _ = GetText(gitfolder, entriesToUpdate)
	updatedBranches := make([]string, len(entriesToUpdate))
	for i, info := range entriesToUpdate {
		logger.Debug("Updating branch %s", info.Branch)
		cache[info.Branch] = info
		updatedBranches[i] = info.Branch
	}
	go WriteCache(gitfolder, cache)
	return cache, updatedBranches
}

func WriteCache(gitfolder string, cache map[string]Entry) {
	cacheFile := RemoteFolder + ".cache"
	b, err := json.Marshal(cache)
	if err != nil {
		logger.Error("Error: " + err.Error())
	}
	err = ioutil.WriteFile(cacheFile, b, 0644)
	if err != nil {
		logger.Error("Error: " + err.Error())
	}
	logger.Debug("Wrote config file: %s", cacheFile)
}

func LoadCache(gitfolder string) map[string]Entry {
	defer timeTrack(time.Now(), "Loading cache")
	cacheFile := RemoteFolder + ".cache"
	b, _ := ioutil.ReadFile(cacheFile)
	var cache map[string]Entry
	err := json.Unmarshal(b, &cache)
	if err != nil {
		logger.Error("Error: " + err.Error())
	}
	return cache
}
