package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kardianos/osext"
	jot "github.com/schollz/jot/src"
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
	defer jot.CleanUp()

	// Handle Ctl+C for cleanUp
	// from http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		jot.CleanUp()
		os.Exit(1)
	}()

	// App information
	setBuild()
	app := cli.NewApp()
	app.Name = "jot"
	app.Version = Version + " " + Build + " " + BuildTime + " " + OS
	app.Usage = `jot is for distributed editing of encrypted stuff

	 https://github.com/schollz/jot

FOLDERS:
	'` + jot.CachePath + `' stores all encrypted files and repositories
	'` + jot.ConfigPath + `' stores all configuration files

EXAMPLE USAGE:
   jot new.txt # create new / edit a document, 'new.txt'
   jot Entry123 # edit a entry, 'Entry123'
   jot --summary # list a summary of all entries
   jot --search "dogs cats" # find entries that mention 'dogs' or 'cats'`

	app.Action = func(c *cli.Context) error {
		// Set the log level
		if Debug {
			jot.DebugMode()
		}

		CheckIfGitIsInstalled()

		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			jot.InputDocument = workingFile
		}

		// Check if its Windows
		if runtime.GOOS == "windows" {
			jot.Extension = ".exe"
		} else {
			jot.Extension = ""
		}
		jot.Version = Version

		// Check new Version
		programPath, _ := osext.Executable()
		jot.CheckNewVersion(programPath, Version, LastCommit, OS)
		jot.ProgramPath, _ = osext.ExecutableFolder()

		// Process some flags
		if ResetConfig {
			jot.SetupConfig()
		} else if Clean {
			jot.CleanAll()
		} else {
			jot.Run()
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
			Usage:       "Deletes all jot files",
			Destination: &Clean,
		},
		cli.StringFlag{
			Name:        "search",
			Usage:       "Search for `word`",
			Destination: &jot.Search,
		},
		cli.BoolFlag{
			Name:        "importold",
			Usage:       "Import `document` (JRNL-format)",
			Destination: &jot.ImportOldFlag,
		},
		cli.BoolFlag{
			Name:        "import",
			Usage:       "Import `document`",
			Destination: &jot.ImportFlag,
		},
		cli.BoolFlag{
			Name:        "export",
			Usage:       "Export `document`",
			Destination: &jot.Export,
		},
		cli.BoolFlag{
			Name:        "config",
			Usage:       "Configure",
			Destination: &ResetConfig,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Edit all of the document",
			Destination: &jot.All,
		},
		cli.BoolFlag{
			Name:        "delete",
			Usage:       "Delete `X`, where X is a document or entry",
			Destination: &jot.DeleteFlag,
		},
		cli.BoolFlag{
			Name:        "summary",
			Usage:       "Gets summary",
			Destination: &jot.Summarize,
		},
		cli.BoolFlag{
			Name:        "stats",
			Usage:       "Print stats",
			Destination: &jot.ShowStats,
		},
	}
	app.Run(os.Args)
}

func CheckIfGitIsInstalled() {
	cmd := exec.Command("git", "--version")
	stdout, err := cmd.Output()
	versionNums := strings.Split(strings.Split(strings.TrimSpace(string(stdout)), " ")[2], ".")
	major, _ := strconv.Atoi(versionNums[0])
	minor, _ := strconv.Atoi(versionNums[1])
	if major < 2 || (major == 2 && minor < 5) {
		fmt.Printf("\n%s detected.\n\nPlease install git version 2.5+ before proceeding. To install, go to \n\n    https://git-scm.com/downloads \n\nand install the version for your operating system.\nPress enter to continue... \n", strings.TrimSpace(string(stdout)))
		var input string
		fmt.Scanln(&input)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("\ngit is not detected.\n\nPlease install git version 2.5+ before proceeding. To install, go to \n\n    https://git-scm.com/downloads \n\nand install the version for your operating system.\nPress enter to continue... ")
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
		err := os.Chdir(path.Join(os.Getenv("GOPATH"), "src", "github.com", "schollz", "jot"))
		if err != nil {
			return
		}
		cmd := exec.Command("git", "log", "-1", "--pretty=format:'%h||%ad'")
		stdout, err := cmd.Output()
		if err != nil {
			return
		}
		items := strings.Split(string(stdout),"||")
		LastCommit = strings.Replace(items[1], "'", "", -1)
		Build = strings.Replace(items[0], "'", "", -1)
		BuildTime = LastCommit
	} else {
		Build = Build[0:7]
	}
}
