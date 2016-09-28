package gitsdees

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	CachePath, ConfigPath, TempPath                          string
	CurrentDocument, Editor, Remote, InputDocument           string
	All, Export, Summarize                                   bool
	DeleteEntry, Search                                      string
	RemoteFolder, CacheFile                                  string
	Extension                                                string
	Passphrase                                               string
	Debug, Encrypt, DontEncrypt, ResetConfig, DeleteDocument bool
)
