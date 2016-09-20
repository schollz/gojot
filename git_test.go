package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat("./gittest"); os.IsNotExist(err) {
		createBranches()
	}
	os.Exit(m.Run())
}

func TestListBranches(t *testing.T) {
	branches, err := ListBranches("./gittest")
	fmt.Println(len(branches), err)
	// Output: 100 <nil>
}

func TestGetInfo(t *testing.T) {
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
	runCommand("git init")
	d1 := []byte("hello, world")
	err := ioutil.WriteFile("test.txt", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("git add test.txt")
	runCommand("git commit -am 'added test.txt'")

	start := time.Now()
	for i := 0; i < 10; i++ {
		runCommand("git checkout --orphan " + strconv.Itoa(i))
		d1 = []byte("hello, world branch #" + strconv.Itoa(i))
		err = ioutil.WriteFile("test.txt", d1, 0644)
		if err != nil {
			log.Fatal(err)
		}
		runCommand("git add test.txt")
		runCommand("git commit -am 'added test.txt'")
		// 3 commands = 115 ms / command
	}

	elapsed := time.Since(start)
	log.Printf("createBranches took %s", elapsed)
}
