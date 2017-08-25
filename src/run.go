package main

import "fmt"

func run() (err error) {
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
	err = gj.SetRepo(repoString)
	if err != nil {
		return
	}
	err = gj.LoadConfig(identity)
	if err != nil {
		return
	}
	err = gj.LoadRepo()

	err = gj.NewEntry(true)
	if err != nil {
		return
	}

	err = gj.Save()
	if err != nil {
		return
	}
	return

}

func main() {
	fmt.Println(run())
}
