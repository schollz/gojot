package sdees

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

func StringToHashID(s string) string {
	key := []byte(Cryptkey)
	plaintext := []byte(s)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return hex.EncodeToString(ciphertext)
}

func HashIDToString(e string) string {
	if len(e) == 0 {
		return ""
	}
	e = strings.Replace(e, ".", "", -1)
	e = strings.Replace(e, ".cache", "", -1)
	key := []byte(Cryptkey)
	ciphertext, err := hex.DecodeString(e)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext2 := make([]byte, len(ciphertext)-aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])
	return string(plaintext2)
}
