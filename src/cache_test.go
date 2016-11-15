package jot

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
	_, _, err := UpdateCache(CACHE_TEST_GITFOLDER, EncryptOTP("test.txt"), true)
	if !exists(path.Join(CACHE_TEST_GITFOLDER, EncryptOTP("test.txt")+".cache")) || err != nil {
		t.Errorf("Error creating cache: %s, %v", path.Join(CACHE_TEST_GITFOLDER, EncryptOTP("test.txt")+".cache"), err)
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

	_, updatedBranches, _ := UpdateCache(gitfolder, EncryptOTP("test2.txt"), false)
	if len(updatedBranches) < 2 {
		t.Errorf("Error updating branches, got %v", updatedBranches)
	}
}

func TestLoadCache(t *testing.T) {
	log.Println("Testing LoadCache...")
	UpdateCache(CACHE_TEST_GITFOLDER, EncryptOTP("other.txt"), true)
	cache, err := LoadCache(CACHE_TEST_GITFOLDER, EncryptOTP("other.txt"))
	if err != nil {
		t.Errorf("Got error: " + err.Error())
	}
	if len(cache.Branch) == 0 {
		t.Errorf("Error loading cache, got: %v", cache.Branch)
	}
}
