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
		t.Error("Expected 100 branches, got %d, and error %s", len(branches), err.Error())
	}
}

func TestGetInfo(t *testing.T) {
	log.Println("Testing GetInfo()...")
	branchNames, _ := ListBranches("./gittest")
	entries, _ := GetInfo("./gittest", branchNames)
	for _, entry := range entries {
		if entry.Branch == "12" {
			if entry.Fulltext != "hello, world branch #12" {
				t.Error("Expected %s, got %s", "hello, world branch #12", entry.Fulltext)
			}
			break
		}
	}
}

func TestClone(t *testing.T) {
	log.Println("Testing CloneRepo()...")
	os.RemoveAll("/tmp/test")
	err := Clone("/tmp/test", "https://github.com:schollz/test.git")
	if err != nil {
		t.Error("Got error while cloning: " + err.Error())
	}

	branches, _ := ListBranches("/tmp/test")
	if len(branches) > 2 && err != nil {
		t.Error("Something unexpected " + err.Error())
	}
}

func TestNewDocument(t *testing.T) {
	log.Println("Testing NewDocument()...")
	_, err := NewDocument("/tmp/test", "test2.txt", "hi", "some message", "Thu, 07 Apr 2005 22:13:13 +0200")
	if err != nil {
		t.Error("Got error while making new document: " + err.Error())
	}
}

func TestPush(t *testing.T) {
	log.Println("Testing Push()...")
	err := Push("/tmp/test")
	if err != nil {
		t.Error("Got error pushing: " + err.Error())
	}
}

func TestGetLatest(t *testing.T) {
	log.Println("Testing GetLatest()...")
	os.RemoveAll("/tmp/testOld")
	err := Clone("/tmp/testOld", "https://github.com:schollz/test.git")
	if err != nil {
		t.Error("Got error while cloning: " + err.Error())
	}

	os.RemoveAll("/tmp/testNew")
	err = Clone("/tmp/testNew", "https://github.com:schollz/test.git")
	if err != nil {
		t.Error("Got error while cloning: " + err.Error())
	}

	branch, err := NewDocument("/tmp/testNew", "test2.txt", "hi", "some message", "Thu, 07 Apr 2005 22:13:13 +0200")
	if err != nil {
		t.Error("Got error while making new document: " + err.Error())
	}
	logger.Debug("Created new branch %s", branch)

	err = Push("/tmp/testNew")
	if err != nil {
		t.Error("Got error pushing: " + err.Error())
	}

	time.Sleep(3 * time.Second)

	_, err = GetLatest("/tmp/testOld")
	if err != nil {
		t.Error("Got error GetLatest: " + err.Error())
	}

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
