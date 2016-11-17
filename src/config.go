package gojot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Remote, Editor, CurrentDocument string
}

func SetupConfig() {
	var configParameters Config

	var yesno string
	err := errors.New("Incorrect git repo")
	for {
		fmt.Print("Enter git repo (e.g.: git@github.com:USER/REPO.git): ")
		fmt.Scanln(&yesno)
		cwd, _ := os.Getwd()
		os.Chdir(CachePath)
		if !exists(EncodeBase58(yesno)) {
			fmt.Println("Cloning " + yesno + " ...")
			cmd := exec.Command("git", "clone", yesno, EncodeBase58(yesno))
			out2, _ := cmd.StderrPipe()
			cmd.Start()
			out2b, _ := ioutil.ReadAll(out2)
			cmd.Wait()
			os.Chdir(cwd)
			if strings.Contains(string(out2b), "fatal: ") {
				fmt.Println(strings.TrimSpace(string(out2b)))
				fmt.Println("Could not clone, please re-enter")
			} else {
				// Remove, this cloning only was to check that it was a valid thing
				os.RemoveAll(path.Join(CachePath, EncodeBase58(yesno)))
				break
			}
		} else {
			fmt.Println("Already exists, doing nothing.")
			break
		}
	}
	configParameters.Remote = yesno

	// Loop until user chooses an available program
	for {
		fmt.Printf("Which editor do you want to use: micro (default), vim,  nano, or emacs? ")
		fmt.Scanln(&yesno)
		if strings.TrimSpace(strings.ToLower(yesno)) == "nano" {
			configParameters.Editor = "nano"
		} else if strings.TrimSpace(strings.ToLower(yesno)) == "emacs" {
			configParameters.Editor = "emacs"
		} else if strings.TrimSpace(strings.ToLower(yesno)) == "vim" {
			configParameters.Editor = "vim"
		} else {
			configParameters.Editor = "micro"
		}
		if Version != "dev" && configParameters.Editor == "micro" {
			break
		}
		// check if it actually exists
		if runtime.GOOS == "windows" {
			Extension = ".exe"
		}
		cmd := exec.Command(path.Join(ProgramPath, configParameters.Editor+Extension), "--version")
		_, err2 := cmd.Output()
		if err2 == nil {
			break
		}
		cmd = exec.Command(path.Join(".", configParameters.Editor+Extension), "--version")
		_, err2 = cmd.Output()
		if err2 == nil {
			break
		}
		_, err2 = Asset("bin/" + configParameters.Editor + Extension)
		if err2 == nil {
			break
		}
		fmt.Printf("\n%s not found, are you sure its installed?\n\n", configParameters.Editor+Extension)
	}

	configParameters.CurrentDocument = ""

	b, err := json.Marshal(configParameters)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(ConfigPath, "config.json"), b, 0644)
}

func LoadConfiguration() {
	defer timeTrack(time.Now(), "Loaded and saved configuration")
	var c Config
	data, err := ioutil.ReadFile(path.Join(ConfigPath, "config.json"))
	if err != nil {
		logger.Error("Could not load config.json")
		return
	}
	json.Unmarshal(data, &c)
	if len(CurrentDocument) == 0 {
		CurrentDocument = c.CurrentDocument
	}
	Editor = c.Editor
	Remote = c.Remote
	RemoteFolder = path.Join(CachePath, EncodeBase58(Remote))
	if len(Remote) == 0 {
		SetupConfig()
	}
}

func SaveConfiguration(editor string, remote string, currentdoc string) {
	defer timeTrack(time.Now(), "Saved configuration")
	c := Config{Editor: editor, Remote: remote, CurrentDocument: currentdoc}
	b, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(ConfigPath, "config.json"), b, 0644)
}
