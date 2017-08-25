package main

import (
	"fmt"

	"github.com/fatih/color"
)

var highlight = color.New(color.FgYellow).SprintFunc()

func run() (err error) {
	color.Set(color.FgYellow, color.Bold)
	fmt.Println(`
  ___   __     __   __  ____ 
 / __) /  \  _(  ) /  \(_  _)
( (_ \(  O )/ \) \(  O ) )(  
 \___/ \__/ \____/ \__/ (__) 
`)
	color.Unset()
	gj, err := New(false)
	if err != nil {
		return
	}
	err = gj.Load()
	if err != nil {
		return
	}
	repoString := gj.RepoString
	identity := gj.IdentityString
	if len(repoString) > 0 && len(identity) > 0 {
		fmt.Printf("Loading settings for '%s' \nin repo '%s'\n\nTo load new repo, use -new\n\n", highlight(identity), highlight(repoString))
	} else {
		repoString = ""
		identity = ""
	}

	err = gj.SetRepo(repoString)
	if err != nil {
		return
	}

	err = gj.LoadConfig(identity)
	if err != nil {
		return
	}

	fmt.Println("\nLoading entries:")
	err = gj.LoadRepo()
	if err != nil {
		return
	}

	// Save as last used settings
	err = gj.Save()
	if err != nil {
		return
	}

	// Allow to run around in console forever
	for {
		fmt.Println("")
		err = gj.NewEntry(true)
		if err != nil {
			return
		}
	}

	return
}

func main() {
	fmt.Println(run())
}
