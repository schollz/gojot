package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat("./gittest"); os.IsNotExist(err) {
		log.Println("Creating branches for testing...")
		createBranches()
	}
	exitVal := m.Run()
	log.Println("Testing completed.")

	os.Exit(exitVal)
}

func TestListBranches(t *testing.T) {
	log.Println("Testing ListBranches()...")
	branches, err := ListBranches("./gittest")
	if len(branches) != 100 && err != nil {
		t.Errorf("Expected 100 branches, got %d, and error %s", len(branches), err.Error())
	}
}

func TestGetInfo(t *testing.T) {
	log.Println("Testing GetInfo()...")
	branchNames, _ := ListBranches("./gittest")
	entries, _ := GetInfo("./gittest", branchNames)
	for _, entry := range entries {
		if entry.Branch == "12" {
			if entry.Document != "test.txt" {
				t.Errorf("Expected %s, got %s", "test.txt", entry.Document)
			}
			if entry.Message != "added test.txt" {
				t.Errorf("Expected %s, got %s", "added test.txt", entry.Message)
			}
			break
		}
	}
}

func TestClone(t *testing.T) {
	log.Println("Testing CloneRepo()...")
	err := os.RemoveAll("test")
	if err != nil {
		t.Errorf("Got error while removing directory: " + err.Error())
	}
	err = Clone("test", "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	branches, _ := ListBranches("test")
	if len(branches) > 2 && err != nil {
		t.Errorf("Something unexpected " + err.Error())
	}

}

func TestNewDocument(t *testing.T) {
	log.Println("Testing NewDocument()...")
	_, err := NewDocument("test", "test2.txt", "hi", "some message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
}

func TestPush(t *testing.T) {
	log.Println("Testing Push()...")
	err := Push("test")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}
}

func TestDelete(t *testing.T) {
	log.Println("Testing Delete()...")

	os.RemoveAll("testDelete1")
	err := Clone("testDelete1", "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}
	branches, _ := ListBranches("testDelete1")

	err = Delete("testDelete1", branches[1])
	if err != nil {
		t.Errorf("Got error deleting: " + err.Error())
	}

	err = Push("testDelete1")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}

	os.RemoveAll("testDelete")
	err = Clone("testDelete", "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	info, _ := GetInfo("testDelete", []string{branches[1]})
	if info[0].Message != "deleted" {
		t.Errorf("Error while deleting, got %v", info[0])
	}
	os.RemoveAll("testDelete")
	os.RemoveAll("testDelete1")
}

func TestGetLatest(t *testing.T) {
	log.Println("Testing GetLatest()...")

	os.RemoveAll("testOld")
	err := Clone("testOld", "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	os.RemoveAll("testNew")
	err = Clone("testNew", "https://github.com/schollz/test.git")
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	branch, err := NewDocument("testNew", "test2.txt", "hi", "some message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	logger.Debug("Created new branch %s", branch)

	err = Push("testNew")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}
	newBranches, _, err := GetLatest("testOld")
	if err != nil {
		t.Errorf("Got error GetLatest: " + err.Error())
	}
	logger.Debug("Fetched new branches: %v", newBranches)

	if newBranches[0] != branch {
		t.Errorf("Expected seeing %s but got %v instead", branch, newBranches)
	}

	// Test deletion
	err = Delete("testNew", branch)
	if err != nil {
		t.Errorf("Got error deleting: " + err.Error())
	}
	err = Push("testNew")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}
	logger.Debug("Deleted new branch %s", branch)

	_, _, err = GetLatest("testOld")
	if err != nil {
		t.Errorf("Got error GetLatest: " + err.Error())
	}

	info, _ := GetInfo("testOld", []string{branch})
	if info[0].Message != "deleted" {
		t.Errorf("Error while deleting %s, got %v", branch, info[0])
	}

	// os.RemoveAll("testNew")
	// os.RemoveAll("testOld")

}

func createBranches() {
	os.RemoveAll("./gittest")
	os.Mkdir("gittest", 0755)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir("./gittest")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("git", "init")
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	d1 := []byte("hello, world")
	err = ioutil.WriteFile("test.txt", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("git", "add", "test.txt")
	_, err = cmd.Output()

	cmd = exec.Command("git", "commit", "-am", "'added test.txt'")
	_, err = cmd.Output()

	start := time.Now()
	for i := 0; i < 100; i++ {
		cmd := exec.Command("git", "checkout", "--orphan", strconv.Itoa(i))
		_, err := cmd.Output()

		d1 = []byte("hello, world branch #" + strconv.Itoa(i))
		err = ioutil.WriteFile("test.txt", d1, 0644)
		if err != nil {
			fmt.Println("Can't checkout")
			log.Fatal(err)
		}
		cmd = exec.Command("git", "add", "test.txt")
		_, err = cmd.Output()

		cmd = exec.Command("git", "commit", "-am", "'added test.txt'")
		_, err = cmd.Output()
	}

	elapsed := time.Since(start)
	log.Printf("createBranches took %s", elapsed)
}
