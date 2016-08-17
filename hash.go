package main

import (
	"hash/fnv"
	"math/rand"

	"github.com/speps/go-hashids"
)

func encodeNumber(num int) string {
	hd := hashids.NewData()
	hd.Salt = "sdeeseeds"
	hd.MinLength = 6
	h := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{num})
	return string(e)
}

func decodeNumber(s string) int {
	hd := hashids.NewData()
	hd.Salt = "sdeeseeds"
	hd.MinLength = 6
	h := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(s)
	return int(d[0])
}

func hashString(s string) string {
	seed := integerHash(s)
	return RandStringBytesMaskImprSrc(6, seed)
}

func integerHash(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int, seed int64) string {
	src := rand.NewSource(seed)
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
