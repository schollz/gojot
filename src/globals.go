package sdees

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	HomePath, CachePath, ConfigPath, TempPath                     string
	CurrentDocument, Editor, Remote, InputDocument                string
	All, Export, Summarize, ImportFlag, ImportOldFlag, DeleteFlag bool
	Search                                                        string
	RemoteFolder, CacheFile                                       string
	Extension                                                     string
	Passphrase, Cryptkey                                          string
	Debug, ResetConfig, DeleteDocument                            bool
)
