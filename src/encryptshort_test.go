package sdees

import (
	"fmt"
	"testing"
)

func TestEncryptOTP(t *testing.T) {
	fmt.Printf("\nEncrypted:[%s]", EncryptOTP("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", DecryptOTP(EncryptOTP("some kind of string")))
	fmt.Printf("\nEncrypted:[%s]", EncryptOTP("some kind of string"))
	fmt.Printf("\nDecrypted:[%s]\n", DecryptOTP(EncryptOTP("some kind of string")))
	if EncryptOTP("some kind of string") != EncryptOTP("some kind of string") {
		t.Errorf("EncryptOTP not the same for same input")
	}
	if "some kind of string" != DecryptOTP(EncryptOTP("some kind of string")) {
		t.Errorf("EncryptOTP not working")
	}
}

func BenchmarkDecryptOTP(b *testing.B) {
	s := EncryptOTP("some kind of string ")
	for n := 0; n < b.N; n++ {
		DecryptOTP(s)
	}
}

func BenchmarkEncryptOTP(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncryptOTP("some kind of string")
	}
}
