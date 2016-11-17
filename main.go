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
	gojot "github.com/schollz/gojot/src"
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
	defer gojot.CleanUp()

	// Handle Ctl+C for cleanUp
	// from http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		gojot.CleanUp()
		os.Exit(1)
	}()

	// App information
	setBuild()
	app := cli.NewApp()
	app.Name = "gojot"
	app.Version = Version + " " + Build + " " + BuildTime + " " + OS
	app.Usage = `gojot is for distributed editing of encrypted stuff

	 https://github.com/schollz/gojot

FOLDERS:
	'` + gojot.CachePath + `' stores all encrypted files and repositories
	'` + gojot.ConfigPath + `' stores all configuration files

EXAMPLE USAGE:
   gojot new.txt # create new / edit a document, 'new.txt'
   gojot Entry123 # edit a entry, 'Entry123'
   gojot --summary # list a summary of all entries
   gojot --search "dogs cats" # find entries that mention 'dogs' or 'cats'`

	app.Action = func(c *cli.Context) error {
		// Set the log level
		if Debug {
			gojot.DebugMode()
		}

		CheckIfGitIsInstalled()

		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			gojot.InputDocument = workingFile
		}

		// Check if its Windows
		if runtime.GOOS == "windows" {
			gojot.Extension = ".exe"
		} else {
			gojot.Extension = ""
		}
		gojot.Version = Version

		// Check new Version
		programPath, _ := osext.Executable()
		gojot.CheckNewVersion(programPath, Version, LastCommit, OS)
		gojot.ProgramPath, _ = osext.ExecutableFolder()

		// Process some flags
		if ResetConfig {
			gojot.SetupConfig()
		} else if Clean {
			gojot.CleanAll()
		} else {
			gojot.Run()
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
			Usage:       "Deletes all gojot files",
			Destination: &Clean,
		},
		cli.StringFlag{
			Name:        "search",
			Usage:       "Search for `word`",
			Destination: &gojot.Search,
		},
		cli.BoolFlag{
			Name:        "importold",
			Usage:       "Import `document` (JRNL-format)",
			Destination: &gojot.ImportOldFlag,
		},
		cli.BoolFlag{
			Name:        "import",
			Usage:       "Import `document`",
			Destination: &gojot.ImportFlag,
		},
		cli.BoolFlag{
			Name:        "export",
			Usage:       "Export `document`",
			Destination: &gojot.Export,
		},
		cli.BoolFlag{
			Name:        "config",
			Usage:       "Configure",
			Destination: &ResetConfig,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Edit all of the document",
			Destination: &gojot.All,
		},
		cli.BoolFlag{
			Name:        "delete",
			Usage:       "Delete `X`, where X is a document or entry",
			Destination: &gojot.DeleteFlag,
		},
		cli.BoolFlag{
			Name:        "summary",
			Usage:       "Gets summary",
			Destination: &gojot.Summarize,
		},
		cli.BoolFlag{
			Name:        "stats",
			Usage:       "Print stats",
			Destination: &gojot.ShowStats,
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
		err := os.Chdir(path.Join(os.Getenv("GOPATH"), "src", "github.com", "schollz", "gojot"))
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
