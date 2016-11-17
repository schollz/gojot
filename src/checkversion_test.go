package gojot

import "testing"

func TestCheckGithub(t *testing.T) {
	_, _, versions := checkGithub("1.0.0")
	if versions[0] != 2 {
		t.Errorf("Github says the major version is %d", versions[0])
	}
}
func TestUpdateDevVersion(t *testing.T) {
	err := updateDevVersion("2000-01-01")
	if err != nil {
		t.Errorf("Github error" + err.Error())
	}
}
