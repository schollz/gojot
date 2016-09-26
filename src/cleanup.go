package gitsdees

import (
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
