package main

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

func TestEncryptFile(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("Some random text"), 0644)
	err := EncryptFile("test.txt", "abcd")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
	if exists("test.txt") {
		t.Errorf("test.txt should have been shredded!")
	}
	if !exists("test.txt.gpg") {
		t.Errorf("test.txt.gpg should have been created!")
	}
}

func TestDecryptFileCorrectPassword(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("Some random text"), 0644)
	EncryptFile("test.txt", "abcd")
	text, err := DecryptFile("test.txt.gpg", "abcd")
	if err != nil {
		t.Errorf("Error decrypting file: %s", err.Error())
	}
	if text != "Some random text" {
		t.Errorf("Expected 'Some random text', and instead got '%s'", text)
	}
}

func TestDecryptFileWrongPassword(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("Some random text"), 0644)
	EncryptFile("test.txt", "abcd")
	text, err := DecryptFile("test.txt.gpg", "asdfasdf")
	if err == nil {
		t.Errorf("Error decrypting file: %s", err.Error())
	}
	if text == "Some random text" {
		t.Errorf("Expected NOT 'Some random text', and instead got '%s'", text)
	}
}
