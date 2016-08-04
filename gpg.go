package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/maxwellhealth/go-gpg"
)

func encrypt() string {
	w, err := os.Create(path.Join(RuntimeArgs.TempPath, "temp.gpg"))
	if err != nil {
		cleanUp()
		panic(err)
	}
	w.Close()

	fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.TempPath, "temp"))
	if len(fileContents) < 32 {
		return ""
	}
	h := sha1.New()
	h.Write(fileContents)
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	fileName := string(sha1_hash) + ".gpg"

	w, err = os.Create(path.Join(RuntimeArgs.FullPath, fileName))
	if err != nil {
		cleanUp()
		panic(err)
	}
	w.Close()

	toEncrypt, err := os.OpenFile(path.Join(RuntimeArgs.TempPath, "temp"), os.O_RDONLY, 0660)
	if err != nil {
		cleanUp()
		log.Fatal(err)
	}
	destination, err := os.OpenFile(path.Join(RuntimeArgs.FullPath, fileName), os.O_WRONLY, 0660)
	if err != nil {
		cleanUp()
		log.Fatal(err)
	}
	err = gpg.Encode(publicKey, toEncrypt, destination)
	if err != nil {
		cleanUp()
		log.Fatal(err)
	}
	log.Println("Encrypted file!")
	return fileName
}

func decrypt() string {
	toDecrypt, err := os.OpenFile("test.txt.gpg", os.O_RDONLY, 0660)
	if err != nil {
		cleanUp()
		log.Fatal(err)
	}
	destination := new(bytes.Buffer)
	err = gpg.Decode(privateKey, passphrase, toDecrypt, destination)
	if err != nil {
		cleanUp()
		log.Fatal(err)
	}
	return destination.String()
}
