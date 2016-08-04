package main

import (
	"fmt"

	home "github.com/mitchellh/go-homedir"
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
	fmt.Println(Version, Build, BuildTime)
	fmt.Println(home.Dir())
}
