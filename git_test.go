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
	log.Println("Creating branches for testing...")
	createBranches()
	if _, err := os.Stat("./gittest"); os.IsNotExist(err) {
		createBranches()
	}
	exitVal := m.Run()
	log.Println("This gets run AFTER any tests get run!")

	os.Exit(exitVal)
}

func TestOne(t *testing.T) {
	log.Println("TestOne running")
}

func TestTwo(t *testing.T) {
	log.Println("TestTwo running")
}

func TestListBranches(t *testing.T) {
	log.Println("Testing list branches...")
	branches, err := ListBranches("./gittest")
	if _, err := os.Stat("./gittest"); os.IsNotExist(err) {
		createBranches()
	}
	fmt.Println(len(branches), err)
	// Output: 100 <nil>
}

func TestGetInfo(t *testing.T) {
	log.Println("Testing Get info...")
	branchNames, _ := ListBranches("./gittest")
	entries, _ := GetInfo("./gittest", branchNames)
	for _, entry := range entries {
		if entry.Branch == "12" {
			fmt.Println(entry.Fulltext)
		}
	}
	// Output: hello, world branch #12
}

func createBranches() {
	os.RemoveAll("./gittest")
	os.Mkdir("gittest", 0644)
	os.Chdir("gittest")
	cmd := exec.Command("git", "init")
	_, err := cmd.Output()
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
