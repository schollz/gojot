package sdees

import (
	"fmt"
	"os"
	"path/filepath"
)

// CleanUp deletes all temporary files and also deletes documents that were
// made accidently (documents with no data)
func CleanUp() error {
	logger.Debug("Cleaning...")
	dir := TempPath
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		logger.Debug("Shredding %s", filepath.Join(dir, name))
		Shred(filepath.Join(dir, name))
	}
	return nil
}

func CleanAll() {
	var yesno string
	fmt.Print("\n\nThis will remove all local files, but not remote. Are you sure? (y/n) ")
	fmt.Scanln(&yesno)
	if yesno == "y" {
		logger.Debug("Removing cache: %s", CachePath)
		os.RemoveAll(CachePath)
		logger.Debug("Removing config: %s", ConfigPath)
		os.RemoveAll(ConfigPath)
		fmt.Println("All local files removed")
	}
}
