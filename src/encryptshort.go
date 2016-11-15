package jot

import (
	"hash/fnv"
	"math/rand"
	"strings"

	"github.com/codahale/chacha20"
)

func GenerateCryptkey() string {
	key := make([]byte, chacha20.KeySize)
	rand.Read(key)
	return EncodeToString(key)
}

// EncryptOTP runs a XOR encryption on the input string using ChaCha20
// The nonce is generate from a hash so that its reproducible
func EncryptOTP(s string) string {
	if strings.Contains(s, ".otp") || len(s) == 0 {
		return s
	}
	key := DecodeString(Cryptkey)

	// Get integer hash of input, using some of the random bytes in key as salt
	h := fnv.New32a()
	h.Write(append([]byte(s), []byte(key)[:10]...))
	inputToNum := h.Sum32()

	// Use random integer to seed and generate nonce
	rand.Seed(int64(inputToNum))
	nonce := make([]byte, chacha20.NonceSize)
	rand.Read(nonce)

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
