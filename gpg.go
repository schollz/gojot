package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// HashPassword generates a bcrypt hash of the password using work factor 14.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte("alskdjcoimecalks3234kj"+password), 12)
}

// CheckPassword securely compares a bcrypt hashed password with its possible
// plaintext equivalent.  Returns nil on success, or an error on failure.
func CheckPasswordHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte("alskdjcoimecalks3234kj"+password))
}

// decryptString returns the decrypted string using a passphrase and
// GPG symmetric encryption
func decryptString(decryptionString string, encryptionPassphraseString string) (string, error) {
	encryptionPassphrase := []byte(encryptionPassphraseString)
	decbuf := bytes.NewBuffer([]byte(decryptionString))
	result, err := armor.Decode(decbuf)
	if err != nil {
		return "", err
	}

	alreadyPrompted := false
	md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if alreadyPrompted {
			os.Remove(path.Join(RuntimeArgs.FullPath, ConfigArgs.WorkingFile+".pass"))
			return nil, errors.New("Could not decrypt using passphrase")
		} else {
			alreadyPrompted = true
		}
		return encryptionPassphrase, nil
	}, nil)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// decryptString returns the encrypted string using a passphrase and
// GPG symmetric encryption
func encryptString(encryptionText string, encryptionPassphraseString string) string {
	encryptionPassphrase := []byte(encryptionPassphraseString)
	encbuf := bytes.NewBuffer(nil)
	w, err := armor.Encode(encbuf, "PGP SIGNATURE", nil)
	if err != nil {
		log.Fatal(err)
	}

	plaintext, err := openpgp.SymmetricallyEncrypt(w, encryptionPassphrase, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte(encryptionText)
	_, err = plaintext.Write(message)

	plaintext.Close()
	w.Close()
	return encbuf.String()
}

// decrypt returns the decrypted contents of a GPG symmetric encrypted file
func decrypt(file string) string {
	fileContents, _ := ioutil.ReadFile(file)
	decrypted, _ := decryptString(string(fileContents), RuntimeArgs.Passphrase)
	return decrypted
}
