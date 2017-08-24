package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func Write(document string, entry ...string) (err error) {
	entryName := "RANDOM DEFAULT"
	if len(entry) > 0 {
		entryName = entry[0]
	}
	d := NewDocument(document, entryName)
	dString, err := d.String()
	if err != nil {
		return
	}
	err = ioutil.WriteFile("tempfff", []byte(dString), 0644)
	if err != nil {
		return
	}
	cmd := exec.Command("vim", "tempfff")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return
}

func main() {
	Write("some new", "test")
}
