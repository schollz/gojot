package sdees

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// Hash generates a hash of data using HMAC-SHA-512/256. The tag is intended to
// be a natural-language string describing the purpose of the hash, such as
// "hash file for lookup key" or "master secret to client secret".  It serves
// as an HMAC "key" and ensures that different purposes will have different
// hash output. This function is NOT suitable for hashing passwords.
func HMACHash(data string, salt string) string {
	h := hmac.New(sha512.New512_256, []byte(salt))
	h.Write([]byte(data))
	return string(h.Sum(nil))
}

func ShortEncrypt(text string) string {
	passphrase := Passphrase
	key := HMACHash(passphrase, "passphrase")
	fmt.Println(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(text)
	iv := []byte(HMACHash(key, "iv")[:16]) // 16 bytes
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return hex.EncodeToString(ciphertext)
}

func ShortDecrypt(text string) string {
	passphrase := Passphrase
	key := HMACHash(passphrase, "passphrase")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext, _ := hex.DecodeString(text)
	iv := []byte(HMACHash(key, "iv")[:16]) // 16 bytes
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}
