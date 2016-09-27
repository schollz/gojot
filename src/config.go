package gitsdees

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

type Config struct {
	Remote, Editor, CurrentDocument string
}

func SetupConfig() {
	var yesno string
	var configParamaters Config
	fmt.Print("sdees has capability to use a remote git repository to keep documents in sync.\nWould you like to set this up? (y/n) ")
	fmt.Scanln(&yesno)
	if strings.TrimSpace(strings.ToLower(yesno)) == "y" {
		fmt.Print("Enter remote (e.g.: git@github.com:USER/REPO.git): ")
		fmt.Scanln(&configParamaters.Remote)
		// logger.Debug("configParamaters.Remote: %s", configParamaters.Remote)
		if len(configParamaters.Remote) == 0 {
			configParamaters.Remote = "local"
		}
	} else {
		configParamaters.Remote = "local"
	}

	fmt.Printf("Which editor do you want to use: vim (default), nano, or emacs? ")
	fmt.Scanln(&yesno)
	if strings.TrimSpace(strings.ToLower(yesno)) == "nano" {
		configParamaters.Editor = "nano"
	} else if strings.TrimSpace(strings.ToLower(yesno)) == "emacs" {
		configParamaters.Editor = "emacs"
	} else {
		configParamaters.Editor = "vim"
	}
	configParamaters.CurrentDocument = ""

	b, err := json.Marshal(configParamaters)
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
	} else {
		c.CurrentDocument = CurrentDocument
	}
	Editor = c.Editor
	Remote = c.Remote
	RemoteFolder = path.Join(CachePath, HashString(Remote))
	SaveConfiguration(Editor, Remote, CurrentDocument)
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
