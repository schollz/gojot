package main

// Structures
type Entry struct {
	Document, Branch, Date, Hash, Message, Text string
}

// Global parameters
var (
	CachePath, ConfigPath, TempPath          string
	CurrentDocument, Editor, Remote          string
	All                                      bool
	DeleteDocument                           string
	RemoteFolder, CacheFile                  string
	Extension                                string
	Passphrase                               string
	Debug, Encrypt, DontEncrypt, ResetConfig bool
)
