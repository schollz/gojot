package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

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
	Passphrase    string
	ImportFile    string
	ExportFile    string
	SSHKey        string // path to key, usually "~/.ssh/id_rsa"
	WorkingPath   string // main path, usually "~/.sdees/"
	FullPath      string // path with working file, usuallly "~/.sdees/notes.txt/"
	TempPath      string // usually "~/.sdees/temp/"
	SdeesDir      string // name of sdees dir, like ".sdees"
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

func start() {
	passwordAccepted := false
	for passwordAccepted == false {
		fmt.Printf("\nEnter password for editing '%s': ", ConfigArgs.WorkingFile)
		bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		password := strings.TrimSpace(string(bytePassword))
		if exists(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".pass")) {
			// Check old password
			fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".pass"))
			err := CheckPasswordHash(string(fileContents), password)
			if err == nil {
				passwordAccepted = true
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		} else {
			// Generate new passwrod
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 := strings.TrimSpace(string(bytePassword2))
			if password == password2 {
				// Write password to file
				passwordAccepted = true
				passwordHashed, _ := HashPassword(password)
				err := ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile+".pass"), passwordHashed, 0644)
				if err != nil {
					log.Fatal("Could not write to file.")
				}
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		}
	}
	fmt.Println("")
	RuntimeArgs.Passphrase = password
}

func main() {
	RuntimeArgs.SdeesDir = ".sdeesgo"
	fmt.Println(Version, Build, BuildTime)
	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = "sync, decrypt, edit, encrypt, and sync"
	app.Action = func(c *cli.Context) error {
		// Set the log level
		if RuntimeArgs.Debug == false {
			logger.Level(2)
		} else {
			logger.Level(0)
		}
		// Set the paths
		homeDir, _ := home.Dir()
		RuntimeArgs.WorkingPath = path.Join(homeDir, RuntimeArgs.SdeesDir)
		RuntimeArgs.SSHKey = path.Join(homeDir, ".ssh", "id_rsa")

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

		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			ConfigArgs.WorkingFile = workingFile
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

		start()
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
