package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	Version, BuildTime, Build       string
	CachePath, ConfigPath, TempPath string
	CurrentDocument, Editor, Remote string
	All                             bool
	DeleteDocument                  string
	RemoteFolder, CacheFile         string
	Extension                       string
	Debug                           bool
)

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
		if Debug == false {
			logger.Level(2)
		} else {
			logger.Level(0)
		}

		workingFile := c.Args().Get(0)
		if len(workingFile) > 0 {
			CurrentDocument = workingFile
		}

		// Check if its Windows
		if runtime.GOOS == "windows" {
			Extension = ".exe"
		} else {
			Extension = ""
		}

		// Load configuration
		LoadConfiguration()

		// Process some flags
		if len(DeleteDocument) > 0 {
			err := Delete(RemoteFolder, DeleteDocument)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			err = Push(RemoteFolder)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			return nil
		}

		Run()
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Turn on debug mode",
			Destination: &Debug,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "Edit all of the document",
			Destination: &All,
		},
		cli.StringFlag{
			Name:        "delete",
			Usage:       "Delete `document`",
			Destination: &DeleteDocument,
		},
	}
	app.Run(os.Args)
}
