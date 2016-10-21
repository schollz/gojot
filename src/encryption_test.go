package sdees

import (
	"io/ioutil"
	"testing"
)

func TestEncryptStringCorrectPassword(t *testing.T) {
	encrypted := EncryptString("encrypt me", "test")
	decrypted, err := DecryptString(encrypted, "test")

	if err != nil || decrypted != "encrypt me" {
		t.Errorf("Something went wrong. Encryption: \n%s\n Decryption:\n%s\n", encrypted, decrypted)
	}
}

func TestEncryptStringWrongPassword(t *testing.T) {
	encrypted := EncryptString("encrypt me", "test")
	decrypted, err := DecryptString(encrypted, "wrong password")

	if err == nil {
		t.Errorf("Something went wrong. Encryption: \n%s\n Decryption:\n%s\n", encrypted, decrypted)
	}
}

func TestDecryptFileCorrectPassword(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("Some random text"), 0644)
	EncryptFile("test.txt", "abcd")
	err := DecryptFile("test.txt", "abcd")
	if err != nil {
		t.Errorf("Error decrypting file: %s", err.Error())
	}
	b, err := ioutil.ReadFile("test.txt")
	if string(b) != "Some random text" {
		t.Errorf("Expected 'Some random text', and instead got '%s'", string(b))
	}
}

func TestDecryptFileWrongPassword(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("Some random text"), 0644)
	EncryptFile("test.txt", "abcd")
	err := DecryptFile("test.txt", "asdfasdf")
	if err == nil {
		t.Errorf("Error decrypting file: %s", err.Error())
	}
	b, err := ioutil.ReadFile("test.txt")
	if string(b) == "Some random text" {
		t.Errorf("Expected NOT 'Some random text', and instead got '%s'", string(b))
	}
}

// go test -run=Decrypt -bench=.
func BenchmarkDecrypt(b *testing.B) {
	encrypted := EncryptString("encrypt me", "test")
	for n := 0; n < b.N; n++ {
		DecryptString(encrypted, "test")
	}
}
