package jot

import "testing"

func TestIsItDocumentOrEntry(t *testing.T) {
	gotOne, document, entry := IsItDocumentOrEntry("test.txt")
	if gotOne != true || len(document) == 0 || len(entry) != 0 {
		t.Errorf("Incorrectly detected document")
	}
	gotOne, document, entry = IsItDocumentOrEntry("1")
	if gotOne != true || len(document) == 0 || len(entry) == 0 {
		t.Errorf("Incorrectly detected entry")
	}
	gotOne, document, entry = IsItDocumentOrEntry("1kljlkjkl")
	if gotOne != false || len(document) != 0 || len(entry) != 0 {
		t.Errorf("Incorrectly detected entry")
	}
}
