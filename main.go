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

func init() {
	passphrase = []byte("")
	privateKey = []byte(``)
	publicKey = []byte(``)

}

func main() {
	var importfile, exportfile, changeToFile string
	var editwhole, localedit, listFiles, updateSdees bool
	fmt.Println(Version, Build, BuildTime)
	fmt.Println(home.Dir())
	app := cli.NewApp()
	app.Name = "sdees"
	app.Version = Version + " " + Build + " " + BuildTime
	app.Usage = "sync, decrypt, edit, encrypt, and sync"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Usage:       "Work on `FILE`",
			Destination: &changeToFile,
		},
		cli.StringFlag{
			Name:        "import",
			Usage:       "Import text from `FILE`",
			Destination: &importfile,
		},
		cli.StringFlag{
			Name:        "export",
			Usage:       "Export text from `FILE`",
			Destination: &exportfile,
		},
		cli.BoolFlag{
			Name:        "edit, e",
			Usage:       "Edit whole document",
			Destination: &editwhole,
		},
		cli.BoolFlag{
			Name:        "local, l",
			Usage:       "Work locally",
			Destination: &localedit,
		},
		cli.BoolFlag{
			Name:        "update, u",
			Usage:       "Update sdees",
			Destination: &updateSdees,
		},
		cli.BoolFlag{
			Name:        "list, ls",
			Usage:       "List available files",
			Destination: &listFiles,
		},
	}
	app.Run(os.Args)
}
