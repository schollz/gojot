package main

import (
	"fmt"
	"os"

	home "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var privateKey []byte
var passphrase []byte
var publicKey []byte
var userPass string
var userName string
var serverName string

var Version string
var BuildTime string
var Build string

var RuntimeArgs struct {
	HomeDir     string
	ImportFile  string
	ExportFile  string
	WorkingFile string
	Debug       bool
	EditWhole   bool
	EditLocally bool
	ListFiles   bool
	UpdateSdees bool
}

func init() {
	passphrase = []byte("")
	privateKey = []byte(``)
	publicKey = []byte(``)

}

func main() {
	fmt.Println(Version, Build, BuildTime)
	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = "sync, decrypt, edit, encrypt, and sync"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Usage:       "Work on `FILE`",
			Destination: &RuntimeArgs.WorkingFile,
		},
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
	RuntimeArgs.HomeDir, _ = home.Dir()
	fmt.Println(RuntimeArgs)
	if RuntimeArgs.Debug == false {
		logger.Level(2)
	} else {
		logger.Level(0)
	}
}
