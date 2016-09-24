package main

import (
	"crypto/md5"
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
	return b64.StdEncoding.EncodeToString([]byte(hasher.Sum(nil)))[0:8]
}
