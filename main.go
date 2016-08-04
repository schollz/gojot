package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	home "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var privateKey []byte
var passphrase []byte
var publicKey []byte

var Version string
var BuildTime string
var Build string

var RuntimeArgs struct {
	HomeDir       string
	ImportFile    string
	ExportFile    string
	WorkingFile   string
	WorkingPath   string
	FullPath      string
	TempPath      string
	SdeesDir      string
	ServerFileSet map[string]bool
	Debug         bool
	EditWhole     bool
	EditLocally   bool
	ListFiles     bool
	UpdateSdees   bool
}

var ConfigArgs struct {
	WorkingFile string
	ServerHost  string
	ServerPort  string
	ServerUser  string
	SdeesDir    string
}

func main() {
	defer cleanUp()
	RuntimeArgs.SdeesDir = ".sdeesgo"
	fmt.Println(Version, Build, BuildTime)
	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = "sync, decrypt, edit, encrypt, and sync"
	app.Action = func(c *cli.Context) error {
		RuntimeArgs.WorkingFile = c.Args().Get(0)
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "import",
			Usage:       "Import text from `FILE`",
			Destination: &RuntimeArgs.ImportFile,
		},
		cli.StringFlag{
			Name:        "export",
			Usage:       "Export text from `FILE`",
			Destination: &RuntimeArgs.ExportFile,
		},
		cli.BoolFlag{
			Name:        "edit, e",
			Usage:       "Edit whole document",
			Destination: &RuntimeArgs.EditWhole,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Turn on debug mode",
			Destination: &RuntimeArgs.Debug,
		},
		cli.BoolFlag{
			Name:        "local, l",
			Usage:       "Work locally",
			Destination: &RuntimeArgs.EditLocally,
		},
		cli.BoolFlag{
			Name:        "update, u",
			Usage:       "Update sdees",
			Destination: &RuntimeArgs.UpdateSdees,
		},
		cli.BoolFlag{
			Name:        "list, ls",
			Usage:       "List available files",
			Destination: &RuntimeArgs.ListFiles,
		},
	}
	app.Run(os.Args)

	// Set the log level
	if RuntimeArgs.Debug == false {
		logger.Level(2)
	} else {
		logger.Level(0)
	}

	// Set the paths
	RuntimeArgs.HomeDir, _ = home.Dir()
	RuntimeArgs.WorkingPath = path.Join(RuntimeArgs.HomeDir, RuntimeArgs.SdeesDir)

	// Run Importing/Exporting
	if len(RuntimeArgs.ImportFile) > 0 {
		importFile()
		os.Exit(1)
	}
	if len(RuntimeArgs.ExportFile) > 0 {
		exportFile()
		os.Exit(1)
	}

	// Determine if intialization is needed
	if !exists(RuntimeArgs.WorkingPath) {
		initialize()
	}
	if !exists(path.Join(RuntimeArgs.WorkingPath, "config.json")) {
		initialize()
	} else {
		// Load prevoius parameters
		jsonBlob, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, "config.json"))
		err := json.Unmarshal(jsonBlob, &ConfigArgs)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Reset the working file if it is declared
	if len(RuntimeArgs.WorkingFile) > 0 {
		ConfigArgs.WorkingFile = RuntimeArgs.WorkingFile
	}

	// Save current config parameters
	b, err := json.Marshal(ConfigArgs)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

	RuntimeArgs.FullPath = path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile)
	if !exists(RuntimeArgs.FullPath) {
		err := os.MkdirAll(RuntimeArgs.FullPath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	RuntimeArgs.TempPath = path.Join(RuntimeArgs.WorkingPath, "temp")
	if !exists(RuntimeArgs.TempPath) {
		err := os.MkdirAll(RuntimeArgs.TempPath, 0711)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Set public and private key
	publicKey, err = ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, "public.key"))
	if err != nil {
		fmt.Println(`You need to generate and export GPG keys.
gpg --gen-key
gpg --export -a "Your Name" > ~/.sdeesgo/public.key
gpg --export-secret-keys -a "Your Name" > ~/.sdeesgo/private.key`)
		os.Exit(-1)
	}
	privateKey, err = ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, "private.key"))
	if err != nil {
		fmt.Println(`You need to export GPG keys.

gpg --export-secret-keys -a "Your Name" > ~/.sdeesgo/private.key`)
		os.Exit(-1)
	}

	logger.Debug("ConfigArgs: %+v", ConfigArgs)
	logger.Debug("RuntimeArgs: %+v", RuntimeArgs)
	logger.Info("Working on %s", ConfigArgs.WorkingFile)
	logger.Debug("Full path: %s", RuntimeArgs.FullPath)

	readAllFiles()
	// if HasInternetAccess() && !RuntimeArgs.EditLocally {
	// 	syncDown()
	// }
	// editfile()
	// syncUp()
}

func initialize() {
	// Make directory
	err := os.MkdirAll(RuntimeArgs.WorkingPath, 0711)
	if err != nil {
		log.Println("Error creating directory")
		log.Println(err)
		return
	}
	fmt.Print("Enter server address (default: localhost): ")
	fmt.Scanln(&ConfigArgs.ServerHost)
	if len(ConfigArgs.ServerHost) == 0 {
		ConfigArgs.ServerHost = "localhost"
	}

	currentUser, _ := user.Current()
	fmt.Printf("Enter server user (default: %s): ", currentUser.Username)
	fmt.Scanln(&ConfigArgs.ServerUser)
	if len(ConfigArgs.ServerUser) == 0 {
		ConfigArgs.ServerUser = currentUser.Username
	}

	fmt.Printf("Enter server port (default: %s): ", "22")
	fmt.Scanln(&ConfigArgs.ServerPort)
	if len(ConfigArgs.ServerPort) == 0 {
		ConfigArgs.ServerPort = "22"
	}

	fmt.Printf("Enter new file (default: %s): ", "notes.txt")
	fmt.Scanln(&ConfigArgs.WorkingFile)
	if len(ConfigArgs.WorkingFile) == 0 {
		ConfigArgs.WorkingFile = "notes.txt"
	}

	fmt.Println("Make sure to put your keys into the directory " + RuntimeArgs.SdeesDir)
	b, err := json.Marshal(ConfigArgs)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

}

func importFile() {

}

func exportFile() {

}
