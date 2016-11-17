package gojot

import "testing"

func TestGetInfo(t *testing.T) {
	branchNames, _ := ListBranches("./gittest")
	entries, err := GetInfo("./gittest", branchNames)
	if err != nil {
		t.Errorf("Got error for GetInfo " + err.Error())
	}
	foundOne := false
	for i, entry := range entries {
		if i == 0 {
			m, _ := DecryptString(entry.Message, Passphrase)
			if m != "Hi" {
				t.Errorf("Problem decoding message")
			}
		}
		if entry.Document == EncryptOTP("test.txt") {
			foundOne = true
			break
		}
	}
	if !foundOne {
		t.Errorf("Could not get info!")
	}

}
