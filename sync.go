package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func sshtest() {

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
	if _, err = f.Write([]byte("Hello world! Again")); err != nil {
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
