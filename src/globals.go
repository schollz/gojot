package gitsdees

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	CachePath, ConfigPath, TempPath                string
	CurrentDocument, Editor, Remote, InputDocument string
	All, Export                                    bool
	DeleteDocument                                 string
	RemoteFolder, CacheFile                        string
	Extension                                      string
	Passphrase                                     string
	Debug, Encrypt, DontEncrypt, ResetConfig       bool
)
