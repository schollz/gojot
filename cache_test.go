package main

import (
	"log"
	"os"
	"path"
	"testing"
	"time"
)

var CACHE_TEST_PATH = "./gittest10"

func TestCache(t *testing.T) {
	if _, err := os.Stat(CACHE_TEST_PATH); os.IsNotExist(err) {
		log.Println("Creating branches for testing...")
		createBranches(CACHE_TEST_PATH, 100)
	}
	UpdateCache(CACHE_TEST_PATH, true)
	if !exists(path.Join(CachePath, CleanFolderName(CACHE_TEST_PATH)+".cache")) {
		t.Errorf("Error creating cache")
	}
}

func TestLoadCache(t *testing.T) {
	cache := LoadCache(CACHE_TEST_PATH)
	if cache["6"].Text != "hello, world branch #6" {
		t.Errorf("Error loading cache")
	}
}

func TestUpdateCache(t *testing.T) {
	gitfolder := "testOld"
	os.RemoveAll(gitfolder)
	err := Clone(gitfolder, "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}
	UpdateCache(gitfolder, true)
	newLocalBranch, err := NewDocument(gitfolder, "test2.txt", "hiii!", "some other message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	logger.Debug("Created new local branch: %s", newLocalBranch)

	newLocalBranch2, err := NewDocument(gitfolder, "test.txt", RandStringBytesMaskImprSrc(10, time.Now().UnixNano()), "some other message", "Thu, 07 Apr 2005 22:13:13 +0200", "test3")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	logger.Debug("Updated local branch: %s", newLocalBranch2)

	_, updatedBranches := UpdateCache(gitfolder, false)
	if len(updatedBranches) != 2 {
		t.Errorf("Error updating branches")
	}
}
