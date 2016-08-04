package main

import (
	"bytes"
	"log"
	"os"
	"path"

	"github.com/maxwellhealth/go-gpg"
)

func encrypt() {
	w, err := os.Create(path.Join(RuntimeArgs.TempPath, "temp.gpg"))
	if err != nil {
		panic(err)
	}
	w.Close()

	toEncrypt, err := os.OpenFile(path.Join(RuntimeArgs.TempPath, "temp"), os.O_RDONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	destination, err := os.OpenFile(path.Join(RuntimeArgs.TempPath, "temp.gpg"), os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	err = gpg.Encode(publicKey, toEncrypt, destination)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Encrypted file!")
}

func decrypt() string {
	toDecrypt, err := os.OpenFile("test.txt.gpg", os.O_RDONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	destination := new(bytes.Buffer)
	err = gpg.Decode(privateKey, passphrase, toDecrypt, destination)
	if err != nil {
		log.Fatal(err)
	}
	return destination.String()
}
