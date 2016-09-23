package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	home "github.com/mitchellh/go-homedir"
)

var (
	CachePath string
)

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
}

func CleanFolderName(gitfolder string) string {
	return strings.Replace(strings.Replace(gitfolder, "/", "", -1), ".", "", -1)
}

func UpdateCache(gitfolder string, currentCache map[string]Entry) map[string]Entry {
	defer timeTrack(time.Now(), "Updating cache")
	cache := make(map[string]Entry)
	cacheFile := path.Join(CachePath, CleanFolderName(gitfolder)+".cache")

	branchNames, _ := ListBranches(gitfolder)
	entries, _ := GetInfo(gitfolder, branchNames)

	// New cache
	if !exists(cacheFile) || len(currentCache) == 0 {
		logger.Debug("Generating new cache")
		entries, _ = GetText(gitfolder, entries)
		for _, entry := range entries {
			cache[entry.Branch] = entry
		}
		go WriteCache(gitfolder, cache)
		return cache
	}

	// Load and update cache
	logger.Debug("Loading and updating cache")
	cache = LoadCache(gitfolder)
	branchesToUpdate := []Entry{} // which branches to update in cache
	for _, info := range entries {
		if _, ok := cache[info.Branch]; !ok {
			branchesToUpdate = append(branchesToUpdate, info)
			continue
		}
		if info.Hash != cache[info.Branch].Hash {
			branchesToUpdate = append(branchesToUpdate, info)
		}
	}

	branchesToUpdate, _ = GetText(gitfolder, branchesToUpdate)
	for _, info := range branchesToUpdate {
		logger.Debug("Updating branch %s", info.Branch)
		cache[info.Branch] = info
	}
	go WriteCache(gitfolder, cache)
	return cache
}

func WriteCache(gitfolder string, cache map[string]Entry) {
	cacheFile := path.Join(CachePath, CleanFolderName(gitfolder)+".cache")
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
	cacheFile := path.Join(CachePath, CleanFolderName(gitfolder)+".cache")
	b, _ := ioutil.ReadFile(cacheFile)
	var cache map[string]Entry
	err := json.Unmarshal(b, &cache)
	if err != nil {
		logger.Error("Error: " + err.Error())
	}
	return cache
}
