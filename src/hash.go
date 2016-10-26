package sdees

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	b64 "encoding/base64"
	"hash/fnv"
)

// HashString generates a 6-character random string from integer hash of string
func HashString(s string) string {
	seed := integerHash(s)
	return RandStringBytesMaskImprSrc(5, seed)
}

// integerHash generates a integer hash
func integerHash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

// GetMD5Hash from http://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return b64.StdEncoding.EncodeToString([]byte(hasher.Sum(nil)))[0:7]
}

// Hash generates a hash of data using HMAC-SHA-512/256. The tag is intended to
// be a natural-language string describing the purpose of the hash, such as
// "hash file for lookup key" or "master secret to client secret".  It serves
// as an HMAC "key" and ensures that different purposes will have different
// hash output. This function is NOT suitable for hashing passwords.
func HashWithSalt(s string, tag string) string {
	data := []byte(s)
	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)
	return string(h.Sum(nil))
}
