package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

func CleanFolderName(gitfolder string) string {
	return strings.Replace(strings.Replace(RemoteFolder, ".", "", -1), "/", "", -1)
}

type Cache struct {
	Branch map[string]Entry
	Ignore map[string]bool
}

func UpdateCache(gitfolder string, document string, forceUpdate bool) (Cache, []string) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Updating cache for document %s in %s", id, document, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Updating cache")
	var cache Cache

	// FIrst colelct branches to get info from
	branchNames, _ := ListBranches(gitfolder)
	var branchesToGetInfo []string
	if !exists(CacheFile) || forceUpdate {
		logger.Debug("Generating new cache")
		branchesToGetInfo = branchNames
		cache.Ignore = make(map[string]bool)
		cache.Branch = make(map[string]Entry)
	} else {
		logger.Debug("Using CacheFile: %s", CacheFile)
		cache = LoadCache(gitfolder)
		for _, branch := range branchNames {
			ignore, ok := cache.Ignore[branch]
			if !ok || !ignore {
				branchesToGetInfo = append(branchesToGetInfo, branch)
			}
		}
	}

	// From those branches, determine which entries need fulltext updating
	entriesToUpdate := []Entry{} // which branches to update in cache
	entries, _ := GetInfo(gitfolder, branchesToGetInfo)
	for _, entry := range entries {
		cache.Ignore[entry.Branch] = entry.Document != document
		ignore, ok := cache.Ignore[entry.Branch]
		if !ok {
			continue
		}
		if !ignore && entry.Hash != cache.Branch[entry.Branch].Hash {
			entriesToUpdate = append(entriesToUpdate, entry)
		}
	}

	// Update the fulltext of entries
	entriesToUpdate, _ = GetText(gitfolder, entriesToUpdate)
	updatedBranches := make([]string, len(entriesToUpdate))
	if len(entriesToUpdate) > 10 {
		logger.Debug("Updating many entries")
	}
	for i, entry := range entriesToUpdate {
		if len(entriesToUpdate) <= 10 {
			logger.Debug("Updating branch %s", entry.Branch)
		}
		cache.Branch[entry.Branch] = entry
		updatedBranches[i] = entry.Branch
	}

	// Save
	WriteCache(gitfolder, cache)

	return cache, updatedBranches
}

func WriteCache(gitfolder string, cache Cache) {
	b, err := json.Marshal(cache)
	if err != nil {
		logger.Error("Error marshaling " + CacheFile + ": " + err.Error())
	}
	err = ioutil.WriteFile(CacheFile, b, 0644)
	if err != nil {
		logger.Error("Error writing " + CacheFile + ": " + err.Error())
	}
	logger.Debug("Wrote cache file: %s", CacheFile)
}

func LoadCache(gitfolder string) Cache {
	defer timeTrack(time.Now(), "Loading cache")
	b, _ := ioutil.ReadFile(CacheFile)
	var cache Cache
	err := json.Unmarshal(b, &cache)
	if err != nil {
		logger.Error("Error loading " + CacheFile + ": " + err.Error())
	}
	return cache
}
