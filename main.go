package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	sdees "github.com/schollz/gitsdees/src"
	"github.com/urfave/cli"
)

var (
	Version, BuildTime, Build   string
	Debug                       bool
	DontEncrypt, Clean          bool
	DeleteDocument, DeleteEntry string
	ResetConfig                 bool
	ImportOldFile, ImportFile   string
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
		// Set the log level
		if Debug {
			sdees.DebugMode()
		}

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

		sdees.Encrypt = !DontEncrypt

		// Load configuration
		sdees.LoadConfiguration()

		// Process some flags
		if len(DeleteEntry) > 0 {
			fmt.Printf("Deleting entry %s\n", DeleteEntry)
			return sdees.DeleteAndPush(DeleteEntry)
		} else if ResetConfig {
			sdees.SetupConfig()
		} else if len(ImportOldFile) > 0 {
			fmt.Printf("Importing %s using deprecated import file\n", ImportOldFile)
			sdees.CurrentDocument = ImportOldFile
			sdees.ImportOld(ImportOldFile)
		} else if len(ImportFile) > 0 {
			fmt.Printf("Importing %s\n", ImportFile)
			sdees.CurrentDocument = ImportFile
			sdees.Import(ImportFile)
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
			Usage:       "Deletes all gitsdees files",
			Destination: &Clean,
		},
		cli.StringFlag{
			Name:        "search",
			Usage:       "Search for `word`",
			Destination: &sdees.Search,
		},
		cli.StringFlag{
			Name:        "importold",
			Usage:       "Import `document` (JRNL-format)",
			Destination: &ImportOldFile,
		},
		cli.StringFlag{
			Name:        "import",
			Usage:       "Import `document`",
			Destination: &ImportFile,
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
		cli.StringFlag{
			Name:        "delete",
			Usage:       "Delete `entry`",
			Destination: &DeleteEntry,
		},
		cli.BoolFlag{
			Name:        "plaintext",
			Usage:       "Turns off encryption",
			Destination: &DontEncrypt,
		},
		cli.BoolFlag{
			Name:        "summary",
			Usage:       "Gets summary",
			Destination: &sdees.Summarize,
		},
	}
	app.Run(os.Args)
}
