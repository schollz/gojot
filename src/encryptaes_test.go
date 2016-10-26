package sdees

import (
	"fmt"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	Cryptkey = "asdfasdfasdfasdf"
	Passphrase = "test"
	fmt.Printf("\nEncrypted:[%s]", ("test"))
	fmt.Printf("\nDecrypted:[%s]\n", ("test"))
	if "some kind of string" != DecryptAES(EncryptAES("some kind of string")) {
		t.Errorf("HashID not working")
	}
}

func Benchmark(b *testing.B) {
	s := EncryptAES("some kind of string ")
	for n := 0; n < b.N; n++ {
		DecryptAES(s)
	}
}

func Benchmark(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncryptAES("some kind of string")
	}
}
