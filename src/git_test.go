package sdees

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"
	"time"
)

// BEFORE RUNNING TESTS:
// Make sure to set this to a repo you own, that you don't care about
var GITHUB_TEST_REPO = "git@github.com:schollz/test.git"

func TestMain(m *testing.M) {
	DebugMode()
	Passphrase = "test"
	Cryptkey = "test"
	if _, err := os.Stat("./gittest"); os.IsNotExist(err) {
		log.Println("Creating branches for testing...")
		createBranches("./gittest", 100)
	}
	if _, err := os.Stat("./gittest10"); os.IsNotExist(err) {
		log.Println("Creating branches for testing...")
		createBranches("./gittest10", 10)
	}
	setupPaths()
	RemoteFolder = "./gittest10"
	CurrentDocument = "test.txt"
	CacheFile = path.Join(CachePath, CurrentDocument+".cache")

	exitVal := m.Run()

	// Delete all the folders that get Created
	// log.Println("Deleting testing folders")
	// os.RemoveAll("gittest")
	// os.RemoveAll("gittest10")
	// os.RemoveAll("test")
	// os.RemoveAll("testNew")
	// os.RemoveAll("testOld")

	log.Println("Testing complete")
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
	foundOne := false
	for _, entry := range entries {
		if entry.Document == "test.txt" {
			foundOne = true
			break
		}
	}
	if !foundOne {
		t.Errorf("Could not get info!")
	}

}

func TestClone(t *testing.T) {
	log.Println("Testing CloneRepo()...")
	err := os.RemoveAll("test")
	if err != nil {
		t.Errorf("Got error while removing directory: " + err.Error())
	}
	err = Clone("test", GITHUB_TEST_REPO)
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

func TestGetText(t *testing.T) {
	log.Println("Testing GetText()...")
	branchNames, _ := ListBranches("gittest")
	entries, _ := GetInfo("gittest", branchNames)

	entries, _ = GetText("gittest", entries)
	for _, entry := range entries {
		if entry.Branch == "12" {
			if entry.Text != "hello, world branch #12" {
				t.Errorf("Got different text: %s", entry.Text)
			}
		}
	}

}

func TestDelete(t *testing.T) {
	log.Println("Testing Delete()...")

	os.RemoveAll("testDelete1")
	err := Clone("testDelete1", GITHUB_TEST_REPO)
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
	err = Clone("testDelete", GITHUB_TEST_REPO)
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

func TestGetLatestWithLocalEdits(t *testing.T) {
	log.Println("Testing TestGetLatestWithLocalEdits()...")

	os.RemoveAll("testOld")
	err := Clone("testOld", GITHUB_TEST_REPO)
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	newLocalBranch, err := NewDocument("testOld", "test2.txt", "hiii!", "some other message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}

	os.RemoveAll("testNew")
	err = Clone("testNew", GITHUB_TEST_REPO)
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	// Make some new edit and push it
	_, err = NewDocument("testNew", "test2.txt", "hi", "some message", "Thu, 07 Apr 2005 22:13:13 +0200", "")
	if err != nil {
		t.Errorf("Got error while making new document: " + err.Error())
	}
	err = Push("testNew")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}

	// Get latest, should not override the local branch
	_, _, err = GetLatest("testOld")
	if err != nil {
		t.Errorf("Got error GetLatest: " + err.Error())
	}
	// Push the local brance
	err = Push("testOld")
	if err != nil {
		t.Errorf("Got error pushing: " + err.Error())
	}

	// Get new branch
	newBranches, _, err := GetLatest("testNew")
	if err != nil {
		t.Errorf("Got error GetLatest: " + err.Error())
	}
	if newBranches[0] != newLocalBranch {
		t.Errorf("Did the local branch %s get overidden?", newLocalBranch)
	}
}

func TestGetLatest(t *testing.T) {
	log.Println("Testing GetLatest()...")

	os.RemoveAll("testOld")
	err := Clone("testOld", GITHUB_TEST_REPO)
	if err != nil {
		t.Errorf("Got error while cloning: " + err.Error())
	}

	os.RemoveAll("testNew")
	err = Clone("testNew", GITHUB_TEST_REPO)
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

	os.RemoveAll("testNew")
	os.RemoveAll("testOld")

}

func createBranches(gitfolder string, numBranches int) {
	os.RemoveAll(gitfolder)
	os.Mkdir(gitfolder, 0755)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err := os.Chdir(gitfolder)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("git", "init")
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(18)
	start := time.Now()
	for i := 0; i < numBranches; i++ {
		// cmd := exec.Command("git", "checkout", "--orphan", strconv.Itoa(i))
		// _, err := cmd.Output()
		//
		// d1 = []byte("hello, world branch #" + strconv.Itoa(i))
		// fileName = "test.txt"
		// if rand.Float32() < 0.1 {
		// 	fileName = "other.txt"
		// }
		// err = ioutil.WriteFile(fileName, d1, 0644)
		// if err != nil {
		// 	fmt.Println("Can't checkout")
		// 	log.Fatal(err)
		// }
		// cmd = exec.Command("git", "add", fileName)
		// _, err = cmd.Output()
		//
		// cmd = exec.Command("git", "commit", "-am", "'added "+fileName+"'")
		// _, err = cmd.Output()
		fileName := "test.txt"
		if rand.Float32() < 0.1 {
			fileName = "other.txt"
		}
		NewDocument(gitfolder, fileName, "hello, world branch #"+strconv.Itoa(i), "Hi", GetCurrentDate(), strconv.Itoa(i)) // TODO: WHY DOESN"T THIS WORK??
	}

	elapsed := time.Since(start)
	log.Printf("createBranches took %s", elapsed)
}
