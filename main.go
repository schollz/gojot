// -----------------------------------------------------------------------------
//                           SDEES information
// -----------------------------------------------------------------------------
//
// All configuration files an documents are stored in the folder ~/.config/sdees/.
// Documents are encrypted with symmetric GPG-encryption.
// Each document, X.txt, is stored as a new folder ~/.config/sdees/X.txt/.
// Only files in the document folder are synced remotely.
// Passwords for GPG-encryption are hashed with bcrypt and then stored
// in the document folder ~/.config/sdees/X.txt/X.txt.pass.
//
// Individual entries for a document are stored as GPG encoded files.
// Entries are edited in a temporary file that is always deleted upon exiting,
// which ensures that entries never leave a trace on disk or terminal.
//
// The name of the files for each document contain information about the entry date,
// the file contents, and the modification date. A typical filename is:
//
// yAkbAnL.onLBFi.dew9E6W.gpg
//    ^      ^      ^
//    |------------------- reversible hash-id of the entry date
//           |------------ irreversible hash of the file contents
//                  |----- reversible hash-id of the modification date
//
// Multiple edits of the same entry will result in the same reversible hash-id
// of the entry date. A change in the file contents is determined when
// when the 6-letter irreversible hash of the file contents changes.
// In this cases, the modification date (the third hash) is used to sort
// the entries so that only the newest is displayed when loading the full document.

// -----------------------------------------------------------------------------
//                           SDEES code structure
// -----------------------------------------------------------------------------
//
// main() (main.go):
// - Program entry
// - Processes flags
// - Initiates cleanup on close
// - Run program -> run()
//
// run() (run.go):
// - Handles importing/exporting
// - Pulls latest copy from remote
// - Prompts for password
// - Starts new entry -> editEntry() (entries.go)
// - Pushes latest to remote
//
// editEntry() (entries.go):
// - Edits with vim
// - Encrypts
// - Writes new entry
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	home "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

// App parameters
var Version string
var BuildTime string
var Build string
var Extension string

// Global parameters
var RuntimeArgs struct {
	Passphrase       string
	ServerPassphrase string
	ImportFile       string
	ExportFile       string
	HomePath         string // home path, usually "~/"
	SSHKey           string // path to key, usually "~/.ssh/id_rsa"
	WorkingPath      string // main path, usually "~/.config/sdees/"
	FullPath         string // path with working file, usuallly "~/.config/sdees/notes.txt/"
	TempPath         string // usually "~/.config/sdees/temp/"
	SdeesDir         string // name of sdees dir, like ".config/sdees"
	NumberToShow     string
	TextSearch       string
	DeleteDirectory  string
	Editor           string
	TryPassword      string
	CurrentFileList  []string
	ServerFileSet    map[string]bool
	DontSync         bool
	OnlyPush         bool
	Push             bool
	Pull             bool
	Debug            bool
	EditWhole        bool
	EditLocally      bool
	ListFiles        bool
	UpdateSdees      bool
	Summarize        bool
	ConfigAgain      bool
	Lines            int
}

// Permanent global parameters
var ConfigArgs struct {
	WorkingFile string
	ServerHost  string
	ServerPort  string
	ServerUser  string
	SdeesDir    string
	Editor      string
}

func main() {

	// Delete temp files upon exit
	defer cleanUp()
	// Handle Ctl+C for cleanUp
	// from http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanUp()
		os.Exit(1)
	}()

	// Default directory to store temp files, config, and document
	RuntimeArgs.SdeesDir = path.Join(".config", "sdees")

	// App information
	app := cli.NewApp()
	app.Name = "sdees"
	if len(Build) == 0 {
		Build = "dev"
		out, err := exec.Command("git", []string{"rev-parse", "HEAD"}...).Output()
		if err != nil {
			log.Fatal(err)
		}
		bString := string(out)
		Build = bString[0:7]
		Version = "dev"
		BuildTime = time.Now().String()
	} else {
		Build = Build[0:7]
	}
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = `SDEES Does Editing, Encryption, and Synchronization

	 https://github.com/schollz/sdees

EXAMPLE USAGE:
   sdees new.txt # edit a new document, new.txt
   sdees --summary -n 5 # list a summary of last five entries
   sdees --search "dogs cats" # find entries that mention 'dogs' or 'cats'`

	app.Action = func(c *cli.Context) error {
		// ----------------------
		// Process flags from CLI
		// ----------------------

		// Set the log level
		if RuntimeArgs.Debug == false {
			logger.Level(2)
		} else {
			logger.Level(0)
		}

		// Check if its Windows
		if runtime.GOOS == "windows" {
			Extension = ".exe"
		} else {
			Extension = ""
		}

		// Set the paths
		homeDir, _ := home.Dir()
		RuntimeArgs.HomePath = homeDir
		RuntimeArgs.WorkingPath = path.Join(homeDir, RuntimeArgs.SdeesDir)
		RuntimeArgs.SSHKey = path.Join(homeDir, ".ssh", "id_rsa")

		// Get configuration parameters, or initialize if they don't exist
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

		// Determine the current document to work on - the "working file"
		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			num, isNum := isNumber(workingFile)
			if isNum {
				allFiles := getFileList()
				workingFile = allFiles[num]
			}
			ConfigArgs.WorkingFile = workingFile
		}
		logger.Debug("Working file: %s", ConfigArgs.WorkingFile)

		// Save current config parameters for next time
		b, err := json.Marshal(ConfigArgs)
		if err != nil {
			log.Println(err)
		}
		ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

		// Create the path to the document if it doesn't exist
		RuntimeArgs.FullPath = path.Join(RuntimeArgs.WorkingPath, ConfigArgs.WorkingFile)
		if !exists(RuntimeArgs.FullPath) {
			err := os.MkdirAll(RuntimeArgs.FullPath, 0711)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Create the path to the tepm storage if it doesn't exist
		RuntimeArgs.TempPath = path.Join(RuntimeArgs.WorkingPath, "temp")
		if !exists(RuntimeArgs.TempPath) {
			err := os.MkdirAll(RuntimeArgs.TempPath, 0711)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Run Importing/Exporting
		if len(RuntimeArgs.ImportFile) > 0 {
			importFile(RuntimeArgs.ImportFile)
			// Sync it back up
			if !RuntimeArgs.DontSync || RuntimeArgs.OnlyPush {
				if HasInternetAccess() {
					syncUp()
				} else {
					fmt.Println("Unable to push, no internet access.")
				}
			}
			return nil
		}
		if len(RuntimeArgs.ExportFile) > 0 {
			// Pull latest copies
			logger.Debug("RuntimeArgs.DontSync: %v", RuntimeArgs.DontSync)
			if !RuntimeArgs.DontSync && !RuntimeArgs.OnlyPush {
				if HasInternetAccess() {
					syncDown()
				} else {
					fmt.Println("Unable to pull, no internet access.")
				}
			}
			exportFile(RuntimeArgs.ExportFile)
			return nil
		}

		// Run deletion
		if len(RuntimeArgs.DeleteDirectory) > 0 {
			logger.Debug("Removing %s", path.Join(RuntimeArgs.WorkingPath, RuntimeArgs.DeleteDirectory))
			var yesno string
			fmt.Printf("Are you sure you want to delete '%s'? (y/n): ", RuntimeArgs.DeleteDirectory)
			fmt.Scanln(&yesno)
			if yesno == "y" || yesno == "yes" {
				fmt.Printf("Deleting locally...")
				fmt.Println("done.")
				os.RemoveAll(path.Join(RuntimeArgs.WorkingPath, RuntimeArgs.DeleteDirectory))
				os.RemoveAll(path.Join(RuntimeArgs.WorkingPath, RuntimeArgs.DeleteDirectory+".cache.json"))
				if deleteRemote(RuntimeArgs.DeleteDirectory) {
					fmt.Printf("Deleted %s.\n", RuntimeArgs.DeleteDirectory)
				} else {
					fmt.Printf("Did not delete remote copy of %s.\n", RuntimeArgs.DeleteDirectory)
				}
			} else {
				fmt.Printf("Did not delete %s.\n", RuntimeArgs.DeleteDirectory)
			}
			return nil
		}

		// Run Re-initialization
		if RuntimeArgs.ConfigAgain {
			initialize()
			return nil
		}

		// Run updating
		if RuntimeArgs.UpdateSdees {
			update()
			return nil
		}

		// Run list files
		if RuntimeArgs.ListFiles {
			printFileList()
			return nil
		}

		// Get current file list
		RuntimeArgs.CurrentFileList = getEntryList()

		// Check whether user setup remote syncing
		if ConfigArgs.ServerHost == "do not sync" {
			RuntimeArgs.DontSync = true
		}
		// Run main app (run.go)
		run()

		// Exit
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Edit all, loads whole document",
			Destination: &RuntimeArgs.EditWhole,
		},
		cli.BoolFlag{
			Name:        "list, ls",
			Usage:       "List available files",
			Destination: &RuntimeArgs.ListFiles,
		},
		cli.BoolFlag{
			Name:        "local, l",
			Usage:       "Local editing (no syncing)",
			Destination: &RuntimeArgs.DontSync,
		},
		cli.BoolFlag{
			Name:        "push, p",
			Usage:       "Only push, prevents pulling",
			Destination: &RuntimeArgs.OnlyPush,
		},
		cli.StringFlag{
			Name:        "search, s",
			Usage:       "View only entries that contain `TEXT`",
			Destination: &RuntimeArgs.TextSearch,
		},
		cli.BoolFlag{
			Name:        "summary",
			Usage:       "Summarize",
			Destination: &RuntimeArgs.Summarize,
		},
		cli.StringFlag{
			Name:        "number, n",
			Usage:       "Show up to `N` entries when summarizing",
			Destination: &RuntimeArgs.NumberToShow,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Turn on debug mode",
			Destination: &RuntimeArgs.Debug,
		},
		cli.BoolFlag{
			Name:        "config",
			Usage:       "Edit configuration parameters",
			Destination: &RuntimeArgs.ConfigAgain,
		},
		cli.BoolFlag{
			Name:        "update",
			Usage:       "Update sdees (requires Linux, Go1.6+)",
			Destination: &RuntimeArgs.UpdateSdees,
		},
		cli.StringFlag{
			Name:        "import",
			Usage:       "Generate document from `FILE`",
			Destination: &RuntimeArgs.ImportFile,
		},
		cli.StringFlag{
			Name:        "export",
			Usage:       "Export text from a `DOCUMENT`",
			Destination: &RuntimeArgs.ExportFile,
		},
		cli.StringFlag{
			Name:        "delete",
			Usage:       "Delete a `DOCUMENT`",
			Destination: &RuntimeArgs.DeleteDirectory,
		},
	}
	app.Run(os.Args)
}

// initialize asks the user for the remote user and remote and the editor preference
// and saves these parameters to ~/.config/sdees/config.json
func initialize() {
	var yesno string
	currentUser, _ := user.Current()
	ConfigArgs.ServerHost = ""
	ConfigArgs.ServerUser = ""
	ConfigArgs.ServerPort = ""

	// Make .config directory in home if doesn't exist
	err := os.MkdirAll(path.Join(RuntimeArgs.HomePath, ".config"), 0711)
	if err != nil {
		log.Println("Error creating directory")
		log.Println(err)
		return
	}

	// Make directory
	err = os.MkdirAll(RuntimeArgs.WorkingPath, 0711)
	if err != nil {
		log.Println("Error creating directory")
		log.Println(err)
		return
	}
	fmt.Print("sdees has capability to SSH tunnel to a remote host in order to \nkeep documents synced across devices.\nWould you like to set this up? (y/n) ")
	fmt.Scanln(&yesno)
	if strings.TrimSpace(strings.ToLower(yesno)) == "y" {
		fmt.Print("Enter remote address (default: localhost): ")
		fmt.Scanln(&ConfigArgs.ServerHost)
		logger.Debug("ConfigArgs.ServerHost: [%v]", ConfigArgs.ServerHost)
		if len(ConfigArgs.ServerHost) == 0 {
			ConfigArgs.ServerHost = "localhost"
		}

		fmt.Printf("Enter remote user (default: %s): ", currentUser.Username)
		fmt.Scanln(&ConfigArgs.ServerUser)
		if len(ConfigArgs.ServerUser) == 0 {
			ConfigArgs.ServerUser = currentUser.Username
		}

		fmt.Printf("Enter remote port (default: %s): ", "22")
		fmt.Scanln(&ConfigArgs.ServerPort)
		if len(ConfigArgs.ServerPort) == 0 {
			ConfigArgs.ServerPort = "22"
		}
	} else {
		ConfigArgs.ServerHost = "do not sync"
		ConfigArgs.ServerUser = currentUser.Username
		ConfigArgs.ServerPort = "22"
	}

	fmt.Printf("Which editor do you want to use: vim (default), nano, or emacs? ")
	fmt.Scanln(&yesno)
	if strings.TrimSpace(strings.ToLower(yesno)) == "nano" {
		ConfigArgs.Editor = "nano"
	} else if strings.TrimSpace(strings.ToLower(yesno)) == "emacs" {
		ConfigArgs.Editor = "emacs"
	} else {
		ConfigArgs.Editor = "vim"
	}

	if len(ConfigArgs.WorkingFile) == 0 {
		fmt.Printf("Enter new document name (default: %s): ", "notes.txt")
		fmt.Scanln(&ConfigArgs.WorkingFile)
		if len(ConfigArgs.WorkingFile) == 0 {
			ConfigArgs.WorkingFile = "notes.txt"
		}
	}

	b, err := json.Marshal(ConfigArgs)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

}

// update does `git pull` to collect the latest version of sdees, and does a
// make install to copy the new version into the local directory
func update() {
	logger.Debug("Updating sdees...")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}
	logger.Debug(dir)
	if !HasInternetAccess() {
		fmt.Println("Cannot access internet to update.")
		return
	}
	logger.Debug("Checking current version...")
	out, err := exec.Command("sdees", "--version").Output()
	if err == nil {
		fmt.Println("Current version:")
		fmt.Println(string(out))
	} else {
		logger.Error("Something went wrong: %s", err.Error())
	}

	// Remove git directory if it exists
	os.RemoveAll("tempsdees")

	logger.Debug("Cloning latest...")
	fullCommand := strings.Split("git clone https://github.com/schollz/sdees.git tempsdees", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		logger.Error("Could not clone latest: %s", err.Error())
		os.Exit(-1)
	}

	os.Chdir("./tempsdees")

	fullCommand = strings.Split("make", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		logger.Error("Could not make: %s", err.Error())
		os.Exit(-1)
	}

	fullCommand = strings.Split("make install", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		logger.Error("Could not make install: %s", err.Error())
		os.Exit(-1)
	}

	os.Chdir("../")
	os.RemoveAll("tempsdees")

	out, err = exec.Command("sdees", "--version").Output()
	if err == nil {
		fmt.Println("Updated to version:")
		fmt.Println(string(out))
	}

}
