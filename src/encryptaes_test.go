package sdees

import (
	"fmt"
	"testing"
)

func TestShortEncrypt(t *testing.T) {

	Passphrase = "test"
	fmt.Printf("\nEncrypted:[%s]", ShortEncrypt("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", ShortDecrypt(ShortEncrypt("some kind of string")))
	fmt.Printf("\nEncrypted:[%s]", ShortEncrypt("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", ShortDecrypt(ShortEncrypt("some kind of string")))
	if "some kind of string" != ShortDecrypt(ShortEncrypt("some kind of string")) {
		t.Errorf("HashID not working")
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
