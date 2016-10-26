package sdees

import (
	"fmt"
	"testing"
	"time"
)

func TestShortEncrypt(t *testing.T) {
	Cryptkey = RandStringBytesMaskImprSrc(500000, time.Now().UnixNano())
	Passphrase = "test"
	fmt.Printf("\nEncrypted:[%s]", ShortEncrypt("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", ShortDecrypt(ShortEncrypt("some kind of string")))
	fmt.Printf("\nEncrypted:[%s]", ShortEncrypt("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", ShortDecrypt(ShortEncrypt("some kind of string")))
	if ShortEncrypt("some kind of string") != ShortEncrypt("some kind of string") {
		t.Errorf("ShortEncrypt not the same for same input")
	}
	if "some kind of string" != ShortDecrypt(ShortEncrypt("some kind of string")) {
		t.Errorf("ShortEncrypt not working")
	}
}

func BenchmarkShortDecrypt(b *testing.B) {
	s := ShortEncrypt("some kind of string ")
	for n := 0; n < b.N; n++ {
		ShortDecrypt(s)
	}
}

func BenchmarkShortEncrypt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ShortEncrypt("some kind of string")
	}
}
