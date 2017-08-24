package main

import "fmt"

func run() (err error) {
	gj, err := New(true)
	if err != nil {
		return
	}
	err = gj.SetRepo("https://github.com/schollz/demo2.git")
	if err != nil {
		return
	}
	err = gj.LoadConfig("Testy McTestFace")
	if err != nil {
		return
	}
	err = gj.NewEntry()
	if err != nil {
		return
	}
	return

}

func main() {
	fmt.Println(run())
}
