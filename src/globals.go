package sdees

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	CachePath, ConfigPath, TempPath                   string
	CurrentDocument, Editor, Remote, InputDocument    string
	All, Export, Summarize, ImportFlag, ImportOldFlag bool
	DeleteEntry, Search                               string
	RemoteFolder, CacheFile                           string
	Extension                                         string
	Passphrase, Cryptkey                              string
	Debug, ResetConfig, DeleteDocument                bool
)
