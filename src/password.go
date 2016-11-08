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
func PromptPassword(gitfolder string) string {
	var err error
	password1 := "1"
	textToTest, _ := GetTextOfOne(gitfolder, "master", ".key")
	if len(textToTest) == 0 {
		password2 := "2"
		for password1 != password2 {
			fmt.Printf("Enter new password: ")
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
		logger.Debug("It seems key doesn't exist yet, making it")
		Cryptkey = GenerateCryptkey()
		WriteToMaster(gitfolder, ".key", Cryptkey)
	} else {
		logger.Debug("Testing with master:key")
		passwordAccepted := false
		for passwordAccepted == false {
			fmt.Printf("Enter password: ")
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			Cryptkey, err = DecryptString(textToTest, password1)
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
