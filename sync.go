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

// PublicKeyFiles is used to accessing the remote server
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

// syncDown pulls the latest copies of all the documents from the remote server
func syncDown() {
	fmt.Println("Pulling from remote...")
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
			fmt.Printf("\nSSH key to connect to %s@%s not found, \nperhaps add it with `ssh-copy-id %s@%s`?\n\n", ConfigArgs.ServerUser, ConfigArgs.ServerHost, ConfigArgs.ServerUser, ConfigArgs.ServerHost)
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

	// walk a directory
	RuntimeArgs.ServerFileSet = make(map[string]bool)
	files := []string{}
	serverFolders := make(map[string]bool)
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
		if !strings.Contains(w.Path(), ".gpg") && !strings.Contains(w.Path(), ".pass") {
			serverFolders[w.Path()] = true
		}
	}

	// Check whether local documents are the same as server documents
	for _, localFolder := range listFiles() {
		if _, ok := serverFolders[dirToWalk+"/"+localFolder]; !ok {
			logger.Debug("Server doesn't have %s", localFolder)
			// newFiles := files
			// for i := range files {
			// 	fmt.Println(files[i])
			// 	if strings.Contains(files[i], dirToWalk+"/"+localFolder) {
			// 		fmt.Println(files[i])
			// 		newFiles = append(newFiles[:i], newFiles[i+1:]...)
			// 	}
			// }
		}
	}

	filesToSync := []string{}
	fileNamesToSync := []string{}
	folderNamesToSync := []string{}
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
			folderNamesToSync = append(folderNamesToSync, folderName)
			filesToSync = append(filesToSync, file)
			fileNamesToSync = append(fileNamesToSync, fileName)

		} else {
			// logger.Debug("Skipping %s", fileName)
		}

	}

	for i := range filesToSync {
		folderName := folderNamesToSync[i]
		fileName := fileNamesToSync[i]
		file := filesToSync[i]
		fmt.Printf("%d/%d)\tSyncing %s/%s.\n", i+1, len(filesToSync), folderName, fileName)

		fp, err := sftp.Open(file)
		if err != nil {
			logger.Error("Could not open %s", fp)
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, fp)
		if err != nil {
			logger.Error("Could not write to %s", fp)
			log.Fatal(err)
		}
		fp.Close()

		err = ioutil.WriteFile(path.Join(RuntimeArgs.WorkingPath, folderName, fileName), buf.Bytes(), 0644)
		if err != nil {
			logger.Error("Could copy down %s", fp)
			log.Fatal(err)
		}
	}
	fmt.Println("...done.")
}

// syncUp pushes the latest versions of all the documents to the server
func syncUp() {
	fmt.Println("Pushing to remote...")
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
			fmt.Printf("\nSSH key to connect to %s@%s not found, \nperhaps add it with `ssh-copy-id %s@%s`?\n\n", ConfigArgs.ServerUser, ConfigArgs.ServerHost, ConfigArgs.ServerUser, ConfigArgs.ServerHost)
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
		for i, file := range filesToSync {
			f, err := sftp.Create(path.Join(dirToWalk, file))
			if err != nil {
				logger.Debug("%s does not exist on server, skipping.", dirToWalk)
			} else {
				fmt.Printf("%d/%d)\tSyncing %s/%s.\n", i+1, len(filesToSync), folder, file)
				fileContents, _ := ioutil.ReadFile(path.Join(RuntimeArgs.WorkingPath, folder, file))
				if _, err = f.Write(fileContents); err != nil {
					logger.Error("Could not write to %s", file)
					log.Fatal(err)
				}
			}
		}

	}
	fmt.Println("...done.")

}

// deleteRemote deletes the entire document from the remote server
func deleteRemote(folderToDelete string) bool {
	if !HasInternetAccess() {
		fmt.Println("No internet access.")
		return false
	}
	fmt.Printf("Deleting from remote...")
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
	dirToWalk := "/home/" + ConfigArgs.ServerUser + "/" + RuntimeArgs.SdeesDir + "/" + folderToDelete
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
		logger.Debug(file)
		err := sftp.Remove(file)
		if err != nil {
			logger.Debug("Error removing file %s: %s", file, err.Error())
		}
	}
	err = sftp.Remove(dirToWalk)
	if err != nil {
		logger.Debug("Error removing %s: %s", dirToWalk, err.Error())
	}
	fmt.Println("done.")
	return true
}
