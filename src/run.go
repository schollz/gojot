package gojot

import (
	"fmt"

	"github.com/fatih/color"
)

func Run() (err error) {
	// TODO: Unbundle vim

	gj, err := New(true)
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
		highlight := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Loading settings for '%s' \nin repo '%s'\n\nTo load new repo, use -new\n\n", highlight(identity), highlight(repoString))
	} else {
		repoString = ""
		identity = ""
	}

	err = gj.SetRepo(repoString)
	if err != nil {
		return
	}
	// TODO: Check to see if it still works if the internet is unconnected

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
			if err.Error() == "Quitting time" {
				err = nil
				break
			} else {
				return
			}
		}
	}

	// fmt.Print("Pushing...")
	// err = gj.Push()
	// if err == nil {
	// 	fmt.Println("...done.")
	// } else {
	// 	fmt.Println("...failed.")
	// }

	return
}
