package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
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

func UpdateCache(gitfolder string, document string, forceUpdate bool) (Cache, []string, error) {
	id := RandStringBytesMaskImprSrc(4, time.Now().UnixNano())
	logger.Debug("[%s]Updating cache for document %s in %s", id, document, gitfolder)
	defer timeTrack(time.Now(), "["+id+"]Updating cache")
	var cache Cache
	var err error
	err = nil

	cacheFile := path.Join(RemoteFolder, document+".cache")

	// FIrst colelct branches to get info from
	branchNames, _ := ListBranches(gitfolder)
	var branchesToGetInfo []string
	logger.Debug("Using cacheFile: %s", cacheFile)
	cacheTest, err2 := LoadCache(gitfolder, document)
	if err2 != nil || forceUpdate {
		logger.Debug("Generating new cache")
		branchesToGetInfo = branchNames
		cache.Ignore = make(map[string]bool)
		cache.Branch = make(map[string]Entry)
	} else {
		cache = cacheTest
		for _, branch := range branchNames {
			ignore, ok := cache.Ignore[branch]
			if !ok || !ignore {
				branchesToGetInfo = append(branchesToGetInfo, branch)
			}
		}
	}

	// From those branches, determine which entries need fulltext updating
	if Encrypt {
		document += ".gpg"
	}
	entriesToUpdate := []Entry{} // which branches to update in cache
	entries, _ := GetInfo(gitfolder, branchesToGetInfo)
	for _, entry := range entries {
		cache.Ignore[entry.Branch] = entry.Document != document
		ignore, ok := cache.Ignore[entry.Branch]
		if !ok {
			entriesToUpdate = append(entriesToUpdate, entry)
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
	if Encrypt {
		document = strings.Split(document, ".gpg")[0]
	}
	WriteCache(gitfolder, document, cache)

	return cache, updatedBranches, err
}

func WriteCache(gitfolder string, document string, cache Cache) {
	cacheFile := path.Join(RemoteFolder, document+".cache")
	b, err := json.Marshal(cache)
	if err != nil {
		logger.Error("Error marshaling " + cacheFile + ": " + err.Error())
	}
	err = ioutil.WriteFile(cacheFile, b, 0644)
	if err != nil {
		logger.Error("Error writing " + cacheFile + ": " + err.Error())
	}
	if Encrypt {
		EncryptFile(cacheFile, Passphrase)
	}
	logger.Debug("Wrote cache file: %s", cacheFile)
}

func LoadCache(gitfolder string, document string) (Cache, error) {
	var cache Cache
	cacheFile := path.Join(RemoteFolder, document+".cache")
	if Encrypt {
		err := DecryptFile(cacheFile, Passphrase)
		if err != nil {
			logger.Error("Error decrypting %s", cacheFile)
			return cache, err
		}
	}
	defer timeTrack(time.Now(), "Loading cache")
	b, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		logger.Error("Error loading " + cacheFile + ": " + err.Error())
		return cache, err
	}
	err = json.Unmarshal(b, &cache)
	if err != nil {
		logger.Error("Error umarshling " + cacheFile + ": " + err.Error())
		return cache, err
	}
	return cache, err
}
