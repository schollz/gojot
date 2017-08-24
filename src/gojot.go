package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	uuid "github.com/satori/go.uuid"
	"github.com/schollz/gogit"
	"github.com/schollz/gogpg"
	"github.com/sirupsen/logrus"
)

type gojot struct {
	root   string
	repo   *gogit.GitRepo
	gpg    *gogpg.GPGStore
	logger *logrus.Logger
	log    *logrus.Entry
	config Config
}

type Config struct {
	Salt     string
	Identity string
}

var cacheFolder string

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		return
	}
	cacheFolder = path.Join(homedir, ".cache", "gojot2")
}

func New(repo string, debug ...bool) (gj *gojot, err error) {
	gj = new(gojot)
	gj.logger = logrus.New()
	gj.log = gj.logger.WithFields(logrus.Fields{
		"source": "gojot",
	})

	// check debug
	toDebug := false
	if len(debug) > 0 {
		toDebug = debug[0]
	}

	// setup GPG
	gj.log.Info("Setting up GPG")
	gj.gpg, err = gogpg.New(toDebug)
	if err != nil {
		return
	}

	// setup Git
	gj.log.Info("Setting up Git")
	gj.root = path.Join(cacheFolder, gogit.ParseRepoFolder(repo))
	gj.repo, err = gogit.New(repo, gj.root)
	if err != nil {
		return
	}
	gj.repo.Debug(toDebug)
	gj.Debug(toDebug)
	return
}

func (gj *gojot) NewConfig() (err error) {
	config := Config{
		Salt:     uuid.NewV4().String(),
		Identity: gj.gpg.Identity(),
	}
	configB, err := json.Marshal(config)
	if err != nil {
		return
	}
	enc, err := gj.gpg.Encrypt(configB)
	if err != nil {
		return
	}
	err = gj.repo.AddData(enc, path.Join(gj.root, "config.asc"))
	return
}

func (gj *gojot) LoadConfig() (err error) {
	if !exists(path.Join(gj.root, "config.asc")) {
		return errors.New("Need to make config file")
	}
	data, err := ioutil.ReadFile(path.Join(gj.root, "config.asc"))
	if err != nil {
		return
	}
	dec, err := gj.gpg.Decrypt(data)
	if err != nil {
		return
	}
	gj.log.Debugf("config: %s", dec)
	fmt.Println(string(dec))
	return json.Unmarshal(dec, &gj.config)
}

func (gj *gojot) Debug(on bool) {
	gj.gpg.Debug(on)
	gj.repo.Debug(on)
	if on {
		gj.logger.SetLevel(logrus.DebugLevel)
	} else {
		gj.logger.SetLevel(logrus.WarnLevel)
	}
}

func (gj *gojot) Init() (err error) {
	err = gj.repo.Update()
	if err != nil {
		return
	}

	// Check config file
	if !exists(path.Join(gj.root, "config.asc")) {
		return errors.New("Need to make config file")
	}
	return
}

func ListDirs() (repos map[string]string, err error) {
	repos = make(map[string]string)
	files, err := ioutil.ReadDir(cacheFolder)
	if err != nil {
		return
	}
	for _, f := range files {
		p := path.Join(cacheFolder, f.Name())
		fi, err2 := os.Stat(p)
		if err2 != nil {
			err = err2
			return
		}
		if fi.IsDir() {
			repo, err := gogit.GetRemoteOriginURL(p)
			if err != nil {
				continue
			}
			repos[repo] = p
		}
	}
	return
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
