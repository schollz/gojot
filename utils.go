package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/jcelliott/lumber"
)

var logger *lumber.ConsoleLogger

func init() {
	logger = lumber.NewConsoleLogger(lumber.TRACE)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Debug("%s took %s.", name, elapsed)
}

// getPassword gets masked password
// from http://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
func getPassword() string {
	if len(RuntimeArgs.Passphrase) > 0 {
		return string(RuntimeArgs.Passphrase)
	}
	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// HasInternetAccess determines whether or not the internet is accessible
func HasInternetAccess() bool {
	logger.Debug("Checking internet connection...")
	_, connectionError := http.Get("http://www.google.com/")
	internetAccess := true
	if connectionError != nil {
		internetAccess = false
		logger.Debug("...Unavailable")
	} else {
		logger.Debug("...OK")
	}
	return internetAccess
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// from http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
// from http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
