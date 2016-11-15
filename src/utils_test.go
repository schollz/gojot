package jot

import "testing"

func TestStrExtract(t *testing.T) {
	text1 := `!!!!!some text<<<<<<`
	extracted := StrExtract(text1, "!!!!!", "<<<<<<", 1)

	if extracted != "some text" {
		t.Errorf("Incorrect extracted: %s", extracted)
	}
}

func TestGetRandomMD5Hash(t *testing.T) {
	if GetRandomMD5Hash() == GetRandomMD5Hash() {
		t.Errorf("MD5 hashes not random!")
	}
	if len(GetRandomMD5Hash()) < 8 {
		t.Errorf("MD5 hashes are longer than expected")
	}
}
