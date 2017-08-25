package main

import (
	"fmt"
	"os"
	"time"

	gojot "github.com/schollz/gojot/src"
	"github.com/urfave/cli"
)

var version string

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Compiled = time.Now()
	app.Name = "gojot"
	app.Usage = ""
	app.UsageText = `	

		`
	// TODO: Add flags to control repos/identities
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "debug mode",
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		return gojot.Run()
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
