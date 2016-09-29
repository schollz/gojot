package sdees

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// PromptPassword prompts for password and tests against the file in input,
// use "" for no file, in which a new password will be generated
func PromptPassword(gitfolder string, document string) string {
	password1 := "1"
	textToTest, err := GetTextOfOne(gitfolder, "master", document+".gpg")
	if err != nil {
		fmt.Printf("Getting new password for %s\n", document)
		password2 := "2"
		for password1 != password2 {
			fmt.Printf("Enter new password for %s: ", document)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 = strings.TrimSpace(string(bytePassword2))
			if password1 != password2 {
				fmt.Println("\nPasswords do not match.")
			}
		}
		Passphrase = password1
	} else {
		logger.Debug("Testing with master:%s.gpg", document)
		passwordAccepted := false
		for passwordAccepted == false {
			fmt.Printf("\nEnter password to open %s: ", document)
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			_, err := DecryptString(textToTest, password1)
			if err == nil {
				passwordAccepted = true
			} else {
				fmt.Println("\nPasswords do not match.")
				logger.Debug("Got error: %s", err.Error())
			}
		}
	}
	fmt.Println("")
	return password1
}
