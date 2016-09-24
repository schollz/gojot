package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// PromptPassword prompts for password and tests against the file in input,
// use "" for no file, in which a new password will be generated
func PromptPassword(fileToTest string) string {
	password1 := "1"
	if !exists(fileToTest + ".gpg") {
		password2 := "2"
		for password1 != password2 {
			fmt.Printf("Enter password: ")
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 = strings.TrimSpace(string(bytePassword2))
			if password1 != password2 {
				fmt.Println("\nPasswords do not match.")
			}
		}
	} else {
		logger.Debug("Testing with %s", fileToTest)
		passwordAccepted := false
		for passwordAccepted == false {
			fmt.Printf("Enter password: ")
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			err := DecryptFile(fileToTest, password1)
			if err == nil {
				passwordAccepted = true
				EncryptFile(fileToTest, password1)
			} else {
				fmt.Println("\nPasswords do not match.")
			}
		}
	}
	fmt.Println("")
	return password1
}
