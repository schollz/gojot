package main

import "fmt"

func run() (err error) {
	gj, err := New(true)
	if err != nil {
		return
	}
	err = gj.SetRepo()
	if err != nil {
		return
	}
	err = gj.LoadConfig()
	if err != nil {
		return
	}
	return

}

func main() {
	fmt.Println(run())

}
