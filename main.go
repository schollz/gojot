package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/maxwellhealth/go-gpg"
	"github.com/pkg/sftp"
)

// .vimrc
//
// func! WordProcessorModeCLI()
//     setlocal formatoptions=t1
//     setlocal textwidth=80
//     map j gj
//     map k gk
//     set formatprg=par
//     setlocal wrap
//     setlocal linebreak
//     setlocal noexpandtab
//     normal G$
// endfu
// com! WPCLI call WordProcessorModeCLI()
var privateKey []byte
var passphrase []byte
var publicKey []byte
var userPass string
var userName string
var serverName string

func init() {
	passphrase = []byte("")
	privateKey = []byte(``)
	publicKey = []byte(``)

}

func main() {
	// cmdArgs := []string{"-c", "WPCLI", "+startinsert", "new.txt"}
	// cmd := exec.Command("vim", cmdArgs...)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// err := cmd.Run()
	// fmt.Println(cmdArgs)
	// fmt.Println(err)

	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPass),
		},
	}
	connection, err := ssh.Dial("tcp", serverName+":22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	sftp, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	err = sftp.Mkdir("test")
	if err != nil {
		fmt.Println("Directory exists?")
	}

	// walk a directory
	files := []string{}
	w := sftp.Walk("/home/" + userName + "/test")
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
		files = append(files, w.Path())
	}
	fmt.Println(files)

	// leave your mark
	f, err := sftp.Create("test/hi.txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Hello world! Again")); err != nil {
		log.Fatal(err)
	}

	// check it's there
	fi, err := sftp.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)

	fp, err := sftp.Open("test/hi.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, fp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
	s := string(buf.Bytes())
	fmt.Println(s)
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	fmt.Println(s, sha1_hash)

	encrypt()
	start := time.Now()

	for i := 0; i < 10; i++ {
		decrypt()
	}
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func encrypt() {
	toEncrypt, err := os.OpenFile("test.txt", os.O_RDONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	destination, err := os.OpenFile("test.txt.gpg", os.O_WRONLY, 0660)
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
