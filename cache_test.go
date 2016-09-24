package main

import (
	"log"
	"os"
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

func TestUpdateCache(t *testing.T) {
	log.Println("Testing Update...")
	gitfolder := "testOld"
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

	_, updatedBranches, _ := UpdateCache(gitfolder, "test2.txt", false)
	if len(updatedBranches) != 2 {
		t.Errorf("Error updating branches, got %v", updatedBranches)
	}
}

func TestLoadCache(t *testing.T) {
	log.Println("Testing LoadCache...")
	gitfolder := "testOld"
	cache := LoadCache(gitfolder, "test.txt")
	if _, ok := cache.Branch["test3"]; !ok {
		t.Errorf("Error loading cache, got: %v", cache.Branch["test3"])
	}
}
