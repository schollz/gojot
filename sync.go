package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	logger.Debug("Using %s", file)
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func syncDown() {
	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: ConfigArgs.ServerUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(path.Join(RuntimeArgs.HomeDir, ".ssh", "id_rsa")),
		},
	}
	logger.Debug("Connecting to %s...", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort)
	connection, err := ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	sftp, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	err = sftp.Mkdir("/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir)
	if err != nil {
		// has directory
	}
	err = sftp.Mkdir("/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + ConfigArgs.WorkingFile)
	if err != nil {
		// has directory
	}

	// walk a directory
	RuntimeArgs.ServerFileSet = make(map[string]bool)
	files := []string{}
	dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + ConfigArgs.WorkingFile
	logger.Debug("Walking %s", dirToWalk)
	w := sftp.Walk(dirToWalk)
	first := true
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		if first {
			first = !first
			continue
		}
		files = append(files, w.Path())
	}

	for _, file := range files {
		fp, err := sftp.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, fp)
		if err != nil {
			log.Fatal(err)
		}
		fp.Close()

		fileNameSplit := strings.Split(file, "/")
		fileName := fileNameSplit[len(fileNameSplit)-1]
		RuntimeArgs.ServerFileSet[fileName] = true
		if !exists(path.Join(RuntimeArgs.FullPath, fileName)) {
			logger.Info("Syncing %s.", fileName)
			err = ioutil.WriteFile(path.Join(RuntimeArgs.FullPath, fileName), buf.Bytes(), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			logger.Debug("Skipping %s", fileName)
		}

	}

	logger.Info("Download complete.")
}

func syncUp() {
	filesToSync := []string{}
	files, _ := ioutil.ReadDir(path.Join(RuntimeArgs.FullPath))
	for _, f := range files {
		if _, ok := RuntimeArgs.ServerFileSet[f.Name()]; !ok {
			filesToSync = append(filesToSync, f.Name())
		}
	}

	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: ConfigArgs.ServerUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(path.Join(RuntimeArgs.HomeDir, ".ssh", "id_rsa")),
		},
	}
	logger.Debug("Connecting to %s...", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort)
	connection, err := ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	sftp, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + ConfigArgs.WorkingFile

	for _, file := range filesToSync {
		logger.Info("Writing %s.", file)
		f, err := sftp.Create(path.Join(dirToWalk, file))
		if err != nil {
			log.Fatal(err)
		}
		fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.FullPath, file))
		if _, err = f.Write(fileContents); err != nil {
			log.Fatal(err)
		}
	}
	logger.Info("Upload complete.")
}

func sshtest() {
	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: ConfigArgs.ServerUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(path.Join(RuntimeArgs.HomeDir, ".ssh", "id_rsa")),
		},
	}
	logger.Debug("Connecting to %s...", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort)
	connection, err := ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
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
	dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + ConfigArgs.WorkingFile
	logger.Debug("Walking %s", dirToWalk)
	w := sftp.Walk(dirToWalk)
	first := true
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		if first {
			first = !first
			continue
		}
		logger.Debug(w.Path())
		files = append(files, w.Path())
	}
	logger.Debug("%v", files)

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
