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
	"strings"
	"syscall"

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
	Passphrase       string
	ServerPassphrase string
	ImportFile       string
	ExportFile       string
	HomePath         string // home path, usually "~/"
	SSHKey           string // path to key, usually "~/.ssh/id_rsa"
	WorkingPath      string // main path, usually "~/.sdees/"
	FullPath         string // path with working file, usuallly "~/.sdees/notes.txt/"
	TempPath         string // usually "~/.sdees/temp/"
	SdeesDir         string // name of sdees dir, like ".sdees"
	NumberToShow     string
	TextSearch       string
	ServerFileSet    map[string]bool
	Debug            bool
	EditWhole        bool
	EditLocally      bool
	ListFiles        bool
	UpdateSdees      bool
	Summarize        bool
}

var ConfigArgs struct {
	WorkingFile string
	ServerHost  string
	ServerPort  string
	ServerUser  string
	SdeesDir    string
}

func main() {

	// Handle Ctl+C from http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanUp()
		os.Exit(1)
	}()

	defer cleanUp()
	RuntimeArgs.SdeesDir = ".sdeesgo"
	if len(Build) == 0 {
		Build = "devdevdevdevdev"
	} else {
		Build = Build[0:7]
	}

	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = "sync, decrypt, edit, encrypt, and sync"
	app.Action = func(c *cli.Context) error {
		// Set the log level
		fmt.Printf("sdees version %s (%s)\n", Version, Build)
		if RuntimeArgs.Debug == false {
			logger.Level(2)
		} else {
			logger.Level(0)
		}
		// Set the paths
		homeDir, _ := home.Dir()
		RuntimeArgs.HomePath = homeDir
		RuntimeArgs.WorkingPath = path.Join(homeDir, RuntimeArgs.SdeesDir)
		RuntimeArgs.SSHKey = path.Join(homeDir, ".ssh", "id_rsa")

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

		if workingFile == "pull" {
			if HasInternetAccess() {
				syncDown()
			} else {
				logger.Info("No internet.")
			}
			return nil
		} else if workingFile == "push" {
			if HasInternetAccess() {
				syncUp()
			} else {
				logger.Info("No internet.")
			}
			return nil
		} else {
			logger.Debug("Working file: %s", ConfigArgs.WorkingFile)
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

		// Run Importing/Exporting
		if len(RuntimeArgs.ImportFile) > 0 {
			importFile(RuntimeArgs.ImportFile)
			return nil
		}
		if len(RuntimeArgs.ExportFile) > 0 {
			exportFile(RuntimeArgs.ExportFile)
			return nil
		}

		// Updating
		if RuntimeArgs.UpdateSdees {
			update()
			return nil
		}

		if RuntimeArgs.ListFiles {
			printFileList()
			return nil
		}

		// run main app (run.go)
		run()
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
			Name:        "summary",
			Usage:       "Summarize",
			Destination: &RuntimeArgs.Summarize,
		},
		cli.StringFlag{
			Name:        "number, n",
			Usage:       "Limit number shown",
			Destination: &RuntimeArgs.NumberToShow,
		},
		cli.StringFlag{
			Name:        "search, s",
			Usage:       "Search for text",
			Destination: &RuntimeArgs.TextSearch,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Turn on debug mode",
			Destination: &RuntimeArgs.Debug,
		},
		cli.BoolFlag{
			Name:        "update, u",
			Usage:       "Update sdees (requires Linux, Go1.6+)",
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

	b, err := json.Marshal(ConfigArgs)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, "config.json"), b, 0644)

}

func update() {
	out, err := exec.Command("sdees", "--version").Output()
	if err == nil {
		fmt.Println("Current version:")
		fmt.Println(string(out))
	}
	fullCommand := strings.Split("git clone https://github.com/schollz/sdees.git tempsdees", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Chdir("./tempsdees")

	fullCommand = strings.Split("make install", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Chdir("../")
	fullCommand = strings.Split("rm -rf ./tempsdees", " ")
	if err := exec.Command(fullCommand[0], fullCommand[1:]...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err = exec.Command("sdees", "--version").Output()
	if err == nil {
		fmt.Println("Updated to version:")
		fmt.Println(string(out))
	}

}
