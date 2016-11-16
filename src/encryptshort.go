package jot

import (
	"crypto/hmac"
	cr "crypto/rand"
	"crypto/sha256"
	"strings"

	"github.com/codahale/chacha20"
)

func GenerateCryptkey() string {
	key := make([]byte, chacha20.KeySize)
	cr.Read(key)
	return EncodeToString(key)
}

func GenerateHashSalt() string {
	salt := make([]byte, 20)
	cr.Read(salt)
	return EncodeToString(salt)
}

// EncryptOTP runs a XOR encryption on the input string using ChaCha20
// The nonce is generate from a hash so that its reproducible
func EncryptOTP(s string) string {
	if strings.Contains(s, ".otp") || len(s) == 0 {
		return s
	}
	key := DecodeString(Cryptkey)

	// Get hash of input, using some of the random bytes in key as salt
	h := hmac.New(sha256.New, []byte(HashSalt))
	h.Write([]byte(s))
	nonce := h.Sum(nil)[0:chacha20.NonceSize]

	c, err := chacha20.New(key, nonce)
	if err != nil {
		panic(err)
	}
	src := []byte(s)
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	return EncodeToString(nonce) + "." + EncodeToString(dst) + ".otp"
}

// DecryptOTP runs a XOR encryption on the input string using ChaCha20
func DecryptOTP(input string) string {
	if !strings.Contains(input, ".otp") {
		return input
	}
	items := strings.Split(input, ".")

	key := DecodeString(Cryptkey)

	nonce := DecodeString(items[0])

	c, err := chacha20.New(key, nonce)
	if err != nil {
		panic(err)
	}

	src := DecodeString(items[1])
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	return string(dst)
}
