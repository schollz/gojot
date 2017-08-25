package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	homedir "github.com/mitchellh/go-homedir"
	uuid "github.com/satori/go.uuid"
	"github.com/schollz/gogit"
	"github.com/schollz/gogpg"
	"github.com/sirupsen/logrus"
)

type gojot struct {
	debug  bool
	root   string
	docs   Documents
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
	if !exists(cacheFolder) {
		os.MkdirAll(cacheFolder, 0775)
	}
}

func New(debug ...bool) (gj *gojot, err error) {
	gj = new(gojot)
	gj.logger = logrus.New()
	gj.log = gj.logger.WithFields(logrus.Fields{
		"source": "gojot",
	})

	// check debug
	gj.debug = false
	if len(debug) > 0 {
		gj.debug = debug[0]
	}

	// setup GPG
	gj.log.Info("Setting up GPG")
	gj.gpg, err = gogpg.New(gj.debug)
	if err != nil {
		return
	}

	gj.Debug(gj.debug)
	return
}

func (gj *gojot) Debug(on bool) {
	if gj.gpg != nil {
		gj.gpg.Debug(on)
	}
	if gj.repo != nil {
		gj.repo.Debug(on)
	}
	if on {
		gj.logger.SetLevel(logrus.DebugLevel)
	} else {
		gj.logger.SetLevel(logrus.WarnLevel)
	}
}

// ParseDocuments collects the documents and uses the user salt
// to find the hash. This hash is used to determine whether it is new or not.
func (gj *gojot) ParseDocuments(text string) (docs Documents, err error) {
	docs, err = ParseScroll(text)
	if err != nil {
		return
	}
	for i := 0; i < docs.Len(); i++ {
		hasher := md5.New()
		hasher.Write([]byte(docs[i].Text))
		hasher.Write([]byte(gj.config.Salt))
		docs[i].hash = hex.EncodeToString(hasher.Sum(nil))
		docHashID, err2 := Encode(docs[i].Front.Document, gj.config.Salt)
		if err2 != nil {
			err = err2
			return
		}
		docs[i].file = path.Join(docHashID, docs[i].hash+".asc")
	}
	return
}

func ListAvailableRepos() (repos map[string]string, err error) {
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

func (gj *gojot) SetRepo(repo ...string) (err error) {
	// setup Git
	gj.log.Info("Setting up Git")
	repoString := ""
	if len(repo) > 0 {
		repoString = repo[0]
	} else {
		fmt.Println("Please select a repo (press tab for available):")
		availableRepos, err2 := ListAvailableRepos()
		if err2 != nil {
			return err2
		}
		completer = readline.NewPrefixCompleter()
		for repo := range availableRepos {
			completer.SetChildren(
				[]readline.PrefixCompleterInterface{
					readline.PcItem(repo),
				})
		}

		l, err2 := readline.NewEx(&readline.Config{
			AutoComplete:        completer,
			Prompt:              "\033[31m»\033[0m ",
			InterruptPrompt:     "^C",
			EOFPrompt:           "exit",
			FuncFilterInputRune: filterInput,
		})
		if err2 != nil {
			return err2
		}
		defer l.Close()
		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			repoString = strings.TrimSpace(line)
			if strings.Contains(repoString, ".git") {
				break
			}
			println("'" + repoString + "' is not a valid repo.")
		}
	}
	gj.root = path.Join(cacheFolder, gogit.ParseRepoFolder(repoString))
	gj.repo, err = gogit.New(repoString, gj.root)
	if err != nil {
		return
	}
	gj.repo.Debug(gj.debug)
	err = gj.repo.Update()
	if err != nil {
		return
	}

	return
}

func (gj *gojot) VerifyIdentity(overrideIdentityPassword ...string) (err error) {
	// For testing purposes, you can override it
	if len(overrideIdentityPassword) == 2 {
		return gj.gpg.Init(overrideIdentityPassword[0], overrideIdentityPassword[1])
	}

	// Determine available identities
	availableKeys, err := gj.gpg.ListPrivateKeys()
	if err != nil {
		return
	}

	identity := ""
	// For caching purposes, the identity can be saved to override
	if len(overrideIdentityPassword) == 1 {
		identity = overrideIdentityPassword[0]
		if stringInSlice(overrideIdentityPassword[0], availableKeys) {
			identity = overrideIdentityPassword[0]
		}
	}

	// Setup prompter
	completer = readline.NewPrefixCompleter()
	completer.SetChildren(
		[]readline.PrefixCompleterInterface{
			readline.PcItemDynamic(listThings(availableKeys)),
		})
	l, err2 := readline.NewEx(&readline.Config{
		AutoComplete:        completer,
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		FuncFilterInputRune: filterInput,
	})
	if err2 != nil {
		return err2
	}
	defer l.Close()
	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	if identity == "" {
		// User chooses
		gj.log.Info("Prompting for identity")

		fmt.Println("Please select a GPG identity (tab for available options):")

		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			line = strings.TrimSpace(line)
			if !stringInSlice(line, availableKeys) {
				println("'" + line + "' is not a valid identity.")
			} else {
				identity = line
				break
			}
		}
	}

	fmt.Printf("Please enter password for '%s'\n", identity)
	for {
		pswd, err2 := l.ReadPasswordWithConfig(setPasswordCfg)
		if err2 != nil {
			return err2
		}
		err2 = gj.gpg.Init(identity, string(pswd))
		if err2 != nil {
			println("Password do not match.")
		} else {
			break
		}
	}
	return
}

func (gj *gojot) NewConfig(overrideIdentityPassword ...string) (err error) {
	gj.log.Info("Generating new config")
	if err != nil {
		return err
	}
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

func (gj *gojot) LoadRepo() (err error) {
	filelist := []string{}
	filepath.Walk(gj.root, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file.  ignore.
		}
		matched, err := filepath.Match("*.asc", fi.Name())
		if err != nil {
			fmt.Println(err) // malformed pattern
			return err       // this is fatal.
		}
		if matched {
			_, file := filepath.Split(fp)
			if len(file) == 36 {
				// 36 = 32 character hash + .asc
				// this ensures only actual files go in
				filelist = append(filelist, fp)
			}
		}
		return nil
	})
	data, err := gj.gpg.BulkDecrypt(filelist, true)
	if err != nil {
		return err
	}

	gj.docs = make(Documents, 0, len(data))
	for filename := range data {
		parsedDocs, err2 := gj.ParseDocuments(data[filename])
		if err2 != nil {
			err = err2
			return
		}
		gj.docs = append(gj.docs, parsedDocs[0])
	}
	sort.Sort(gj.docs)

	// TODO: See if this works
	return
}

func (gj *gojot) LoadConfig(overrideIdentityPassword ...string) (err error) {
	err = gj.VerifyIdentity(overrideIdentityPassword...)
	if !exists(path.Join(gj.root, "config.asc")) {
		gj.log.Info("config.asc not found")
		err2 := gj.NewConfig(overrideIdentityPassword...)
		if err2 != nil {
			return err2
		}
	}
	gj.log.Info("Loading config")
	data, err := ioutil.ReadFile(path.Join(gj.root, "config.asc"))
	if err != nil {
		return
	}
	dec, err := gj.gpg.Decrypt(data)
	if err != nil {
		return
	}
	gj.log.Debugf("config: %s", dec)
	return json.Unmarshal(dec, &gj.config)
}

// func (gj *gojot) Open() (err error) {
// 	docs, err := gj.gpg.BulkDecrypt()
// }
func (gj *gojot) NewEntry(showAll bool) (err error) {
	fulltext, err := gj.Write(showAll)
	if err != nil {
		return
	}
	docs, err := gj.ParseDocuments(fulltext)
	if err != nil {
		return
	}
	err = gj.SaveDocuments(docs)
	if err != nil {
		return
	}
	fmt.Println(docs)
	return
}

func (gj *gojot) SaveDocuments(docs Documents) (err error) {
	for i := 0; i < docs.Len(); i++ {
		if !exists(path.Join(gj.root, path.Dir(docs[i].file))) {
			gj.log.Debugf("Creating '%s'", path.Join(gj.root, path.Dir(docs[i].file)))
			err2 := os.MkdirAll(path.Join(gj.root, path.Dir(docs[i].file)), 0775)
			if err2 != nil {
				err = err2
				return
			}
		}
		docs[i].Text = strings.TrimSpace(docs[i].Text)
		if len(docs[i].Text) == 0 {
			continue
		}
		if !exists(path.Join(gj.root, docs[i].file)) {
			gj.log.Debugf("Saving %s", path.Join(gj.root, docs[i].file))
			docString, err2 := docs[i].String()
			if err2 != nil {
				err = err2
				return
			}
			enc, err2 := gj.gpg.Encrypt([]byte(docString))
			if err2 != nil {
				err = err2
				return
			}
			err2 = gj.repo.AddData(enc, path.Join(gj.root, docs[i].file))
			if err2 != nil {
				err = err2
				return
			}
		}
	}
	return
}

func (gj *gojot) Write(showAll bool, documentEntry ...string) (writtenTextString string, err error) {
	var document, entry string
	if len(documentEntry) == 2 {
		document = documentEntry[0]
		entry = documentEntry[1]
	} else if len(documentEntry) == 1 {
		document = documentEntry[1]
	}

	if document == "" {
		document, err = gj.promptForDocument()
		if err != nil {
			return
		}
	}

	var docsString string
	if showAll {
		docsString, err = gj.docs.String(document)
		if err != nil {
			return
		}
	} else {
		docsString = ""
	}

	if entry == "" {
		entry, err = gj.promptForEntry()
		if err != nil {
			return
		}
	}

	d := NewDocument(document, entry)
	dString, err := d.String()
	if err != nil {
		return
	}

	tmpfile, err := ioutil.TempFile("", "write")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up\
	err = ioutil.WriteFile(tmpfile.Name(), []byte(strings.TrimSpace(docsString+"\n\n"+dString)+"\n\n\n"), 0644)
	if err != nil {
		return
	}

	vimrc, err := ioutil.TempFile("", "vimrc")
	if err != nil {
		return
	}
	defer os.Remove(vimrc.Name()) // clean up
	if showAll {
		err = ioutil.WriteFile(vimrc.Name(), []byte(VIMRC), 0644)
	} else {
		err = ioutil.WriteFile(vimrc.Name(), []byte(VIMRC2), 0644)
	}
	if err != nil {
		return
	}

	gj.log.Infof("Running '%s'", strings.Join([]string{"vim", "-u", vimrc.Name(), "-c", "WPCLI", "+", "+startinsert", tmpfile.Name()}, " "))
	cmd := exec.Command("vim.exe", "-u", vimrc.Name(), "-c", "WPCLI", "+", "+startinsert", tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	writtenTextByte, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return
	}
	writtenTextString = string(writtenTextByte)
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (gj *gojot) promptForDocument() (document string, err error) {
	completer = readline.NewPrefixCompleter()
	completer.SetChildren(
		[]readline.PrefixCompleterInterface{
			readline.PcItemDynamic(listThings([]string{"notes"})),
		})
	l, err2 := readline.NewEx(&readline.Config{
		AutoComplete:        completer,
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		FuncFilterInputRune: filterInput,
	})
	if err2 != nil {
		err = err2
		return
	}
	defer l.Close()
	fmt.Println("Please enter a document name:")
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		document = strings.TrimSpace(line)
		break
	}
	return
}

func (gj *gojot) promptForEntry() (entry string, err error) {
	// Setup prompter
	completer = readline.NewPrefixCompleter()
	completer.SetChildren(
		[]readline.PrefixCompleterInterface{
			readline.PcItemDynamic(listThings([]string{"entry1"})),
		})
	l, err2 := readline.NewEx(&readline.Config{
		AutoComplete:        completer,
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		FuncFilterInputRune: filterInput,
	})
	if err2 != nil {
		err = err2
		return
	}
	defer l.Close()
	fmt.Println("Please enter a entry name:")
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		entry = strings.TrimSpace(line)
		break
	}
	return
}
