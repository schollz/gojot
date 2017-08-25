package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/fatih/color"
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
		if runtime.GOOS == "windows" {
			data, err := Asset("src/bundle/vim.exe")
			if err == nil {
				ioutil.WriteFile("vim.exe", data, 0777)
			}
		}

		if version == "" {
			p := path.Join(os.Getenv("GOPATH"), "src", "github.com", "schollz", "gojot")
			_, err := os.Stat(p)
			fmt.Println(p, err)
			if err == nil {
				cwd, err2 := os.Getwd()
				if err2 != nil {
					return err2
				}
				os.Chdir(p)
				cmd := exec.Command("git", "rev-parse", "HEAD")
				stdoutStderr, err2 := cmd.CombinedOutput()
				if err2 != nil {
					return err2
				}
				version = string(stdoutStderr)[:8]
				os.Chdir(cwd)
			} else {
				version = "?"
			}
		}

		color.Set(color.FgYellow, color.Bold)
		fmt.Println(`
  ___   __     __   __  ____ 
 / __) /  \  _(  ) /  \(_  _)
( (_ \(  O )/ \) \(  O ) )(  
 \___/ \__/ \____/ \__/ (__) 
`)

		fmt.Printf("version %s\n\n", version)
		color.Unset()
		return gojot.Run()
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
