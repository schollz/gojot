package gojot

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	Version, HomePath, CachePath, ConfigPath, TempPath, ProgramPath string
	CurrentDocument, Editor, Remote, InputDocument                  string
	All, Export, Summarize, ImportFlag, ImportOldFlag, DeleteFlag   bool
	Search                                                          string
	RemoteFolder, CacheFile                                         string
	Extension                                                       string
	Passphrase, Cryptkey, HashSalt                                  string
	Debug, ResetConfig, DeleteDocument, ShowStats                   bool
)
