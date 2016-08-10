package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
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
	logger.Info("Pulling from remote...")
	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: ConfigArgs.ServerUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(RuntimeArgs.SSHKey),
		},
	}
	logger.Debug("Connecting to %s...", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort)
	connection, err := ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
	if err != nil {
		if len(RuntimeArgs.ServerPassphrase) == 0 {
			fmt.Printf("Enter password for connecting to '%s': ", ConfigArgs.ServerHost)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Printf("\n")
			RuntimeArgs.ServerPassphrase = strings.TrimSpace(string(bytePassword))
		}
		sshConfig = &ssh.ClientConfig{
			User: ConfigArgs.ServerUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(RuntimeArgs.ServerPassphrase),
			},
		}
		connection, err = ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
		if err != nil {
			log.Fatal(err)
		}
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
	dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir
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
		fileNameSplit := strings.Split(file, "/")
		fileName := fileNameSplit[len(fileNameSplit)-1]
		if !strings.Contains(fileName, ".gpg") && !strings.Contains(fileName, ".pass") {
			continue
		}
		folderName := ""
		for i, s := range fileNameSplit {
			if s == RuntimeArgs.SdeesDir {
				folderName = fileNameSplit[i+1]
				break
			}
		}
		RuntimeArgs.ServerFileSet[fileName] = true

		if !exists(path.Join(RuntimeArgs.WorkingPath, folderName, fileName)) {
			if !exists(path.Join(RuntimeArgs.WorkingPath, folderName)) {
				logger.Debug("Creating directory %s", folderName)
				err := os.MkdirAll(path.Join(RuntimeArgs.WorkingPath, folderName), 0711)
				if err != nil {
					log.Fatal(err)
				}
			}
			logger.Info("Syncing %s/%s.", folderName, fileName)

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

			err = ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, folderName, fileName), buf.Bytes(), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// logger.Debug("Skipping %s", fileName)
		}

	}

	logger.Info("...complete.")
}

func syncUp() {
	logger.Info("Pushing to remote...")
	// open an SFTP session over an existing ssh connection.
	sshConfig := &ssh.ClientConfig{
		User: ConfigArgs.ServerUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(RuntimeArgs.SSHKey),
		},
	}
	logger.Debug("Connecting to %s...", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort)
	connection, err := ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
	if err != nil {
		if len(RuntimeArgs.ServerPassphrase) == 0 {
			fmt.Printf("Enter password for connecting to '%s': ", ConfigArgs.ServerHost)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Printf("\n")
			RuntimeArgs.ServerPassphrase = strings.TrimSpace(string(bytePassword))
		}
		sshConfig = &ssh.ClientConfig{
			User: ConfigArgs.ServerUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(RuntimeArgs.ServerPassphrase),
			},
		}
		connection, err = ssh.Dial("tcp", ConfigArgs.ServerHost+":"+ConfigArgs.ServerPort, sshConfig)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer connection.Close()

	sftp, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	for _, folder := range listFiles() {

		// Collect names of files on Server
		RuntimeArgs.ServerFileSet = make(map[string]bool)
		files := []string{}
		dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + folder
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
		filesToSync := []string{}
		for _, file := range files {
			fileNameSplit := strings.Split(file, "/")
			fileName := fileNameSplit[len(fileNameSplit)-1]
			if !strings.Contains(fileName, ".gpg") && !strings.Contains(fileName, ".pass") {
				continue
			}
			RuntimeArgs.ServerFileSet[fileName] = true
		}

		// Collect local files and check if they are on server
		localFiles, _ := ioutil.ReadDir(path.Join(RuntimeArgs.WorkingPath, folder))
		for _, f := range localFiles {
			if _, ok := RuntimeArgs.ServerFileSet[f.Name()]; !ok {
				filesToSync = append(filesToSync, f.Name())
			} else {
				// logger.Debug("Skipping %s.", f.Name())
			}
		}

		// Sync any local files to server
		for _, file := range filesToSync {
			logger.Info("Syncing %s/%s.", folder, file)
			f, err := sftp.Create(path.Join(dirToWalk, file))
			if err != nil {
				log.Fatal(err)
			}
			fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, folder, file))
			if _, err = f.Write(fileContents); err != nil {
				log.Fatal(err)
			}
		}

	}
	logger.Info("...complete.")

}
