package main

import (
	"log"
	"os"
	"path"
	"testing"
	"time"
)

var CACHE_TEST_GITFOLDER = "./gittest10"

func TestCreateCache(t *testing.T) {
	log.Println("Testing CreateCache...")
	UpdateCache(CACHE_TEST_GITFOLDER, "test.txt", true)
	if !exists(CacheFile) {
		t.Errorf("Error creating cache: %s", CacheFile)
	}
}

func TestLoadCache(t *testing.T) {
	log.Println("Testing LoadCache...")
	cache := LoadCache(CACHE_TEST_GITFOLDER)
	if cache.Branch["6"].Text != "hello, world branch #6" {
		t.Errorf("Error loading cache")
	}
}

func TestUpdateCache(t *testing.T) {
	log.Println("Testing Update...")
	gitfolder := "testOld"
	CacheFile = path.Join(CachePath, "test.txt"+".cache")
	os.RemoveAll(gitfolder)
	err := Clone(gitfolder, GITHUB_TEST_REPO)
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}
	UpdateCache(gitfolder, "test.txt", true)
	newLocalBranch, err := NewDocument(gitfolder, "test2.txt", "hiii!", "some other message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	logger.Debug("Created new local branch: %s", newLocalBranch)

	newLocalBranch2, err := NewDocument(gitfolder, "test2.txt", RandStringBytesMaskImprSrc(10, time.Now().UnixNano()), "some other message", "Thu, 07 Apr 2005 22:13:13 +0200", "test3")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	logger.Debug("Updated local branch: %s", newLocalBranch2)

	_, updatedBranches := UpdateCache(gitfolder, "test.txt", false)
	if len(updatedBranches) != 2 {
		t.Errorf("Error updating branches, got %v", updatedBranches)
	}
}
