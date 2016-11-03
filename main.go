package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/kardianos/osext"
	sdees "github.com/schollz/sdees/src"
	"github.com/urfave/cli"
)

var (
	Version, BuildTime, Build, OS, LastCommit string
	Debug                                     bool
	DontEncrypt, Clean                        bool
	ResetConfig                               bool
	ImportOldFile, ImportFile                 bool
)

func main() {
	// Delete temp files upon exit
	defer sdees.CleanUp()

	// Handle Ctl+C for cleanUp
	// from http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		sdees.CleanUp()
		os.Exit(1)
	}()

	// App information
	setBuild()
	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime + " " + OS
	app.Usage = `SDEES Does Editing, Encryption, and Synchronization

	 https://github.com/schollz/sdees

FOLDERS:
	'` + sdees.CachePath + `' stores all encrypted files and repositories
	'` + sdees.ConfigPath + `' stores all configuration files

EXAMPLE USAGE:
   sdees new.txt # edit a new document, new.txt
	 sdees UglyDuckling # edit a entry in a document
   sdees --summary -n 5 # list a summary of last five entries
   sdees --search "dogs cats" # find entries that mention 'dogs' or 'cats'`

	app.Action = func(c *cli.Context) error {
		// Set the log level
		if Debug {
			sdees.DebugMode()
		}

		CheckIfGitIsInstalled()

		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			sdees.InputDocument = workingFile
		}

		// Check if its Windows
		if runtime.GOOS == "windows" {
			sdees.Extension = ".exe"
		} else {
			sdees.Extension = ""
		}
		sdees.Version = Version

		// Check new Version
		programPath, _ := osext.Executable()
		sdees.CheckNewVersion(programPath, Version, LastCommit, OS)

		// Process some flags
		if ResetConfig {
			sdees.SetupConfig()
		} else if Clean {
			sdees.CleanAll()
		} else {
			sdees.Run()
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Turn on debug mode",
			Destination: &Debug,
		},
		cli.BoolFlag{
			Name:        "clean",
			Usage:       "Deletes all sdees files",
			Destination: &Clean,
		},
		cli.StringFlag{
			Name:        "search",
			Usage:       "Search for `word`",
			Destination: &sdees.Search,
		},
		cli.BoolFlag{
			Name:        "importold",
			Usage:       "Import `document` (JRNL-format)",
			Destination: &sdees.ImportOldFlag,
		},
		cli.BoolFlag{
			Name:        "import",
			Usage:       "Import `document`",
			Destination: &sdees.ImportFlag,
		},
		cli.BoolFlag{
			Name:        "export",
			Usage:       "Export `document`",
			Destination: &sdees.Export,
		},
		cli.BoolFlag{
			Name:        "config",
			Usage:       "Configure",
			Destination: &ResetConfig,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Edit all of the document",
			Destination: &sdees.All,
		},
		cli.BoolFlag{
			Name:        "delete",
			Usage:       "Delete `X`, where X is a document or entry",
			Destination: &sdees.DeleteFlag,
		},
		cli.BoolFlag{
			Name:        "summary",
			Usage:       "Gets summary",
			Destination: &sdees.Summarize,
		},
	}
	app.Run(os.Args)
}

func CheckIfGitIsInstalled() {
	cmd := exec.Command("git", "--version")
	stdout, err := cmd.Output()
	fmt.Println(string(stdout))
	os.Exit(1)
	if err != nil {
		fmt.Println("\ngit is not detected.\n\nPlease install git before proceeding. To install, go to \n\n    https://git-scm.com/downloads \n\nand install the version for your operating system.\nPress enter to continue... ")
		var input string
		fmt.Scanln(&input)
		os.Exit(1)
	}
}

func setBuild() {
	if len(Build) == 0 {
		cwd, _ := os.Getwd()
		defer os.Chdir(cwd)
		Build = "dev"
		Version = "dev"
		BuildTime = time.Now().String()
		err := os.Chdir(path.Join(os.Getenv("GOPATH"), "src", "github.com", "schollz", "sdees"))
		if err != nil {
			return
		}
		cmd := exec.Command("git", "log", "-1", "--pretty=format:'%ad'")
		stdout, err := cmd.Output()
		if err != nil {
			return
		}
		LastCommit = strings.Replace(string(stdout), "'", "", -1)
		BuildTime = LastCommit
	} else {
		Build = Build[0:7]
	}
}
