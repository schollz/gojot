package main

import (
	"strings"
	"testing"
)

func TestCombine(t *testing.T) {
	cache := make(map[string]Entry)
	cache["asdfasdf"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "one"}
	cache["asdfa"] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: "two"}
	cache["fajskdf"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "three"}
	combined := CombineEntries(cache)
	if strings.Contains(combined[0], "one") != true {
		t.Errorf("Got [%s] instead of the correct", strings.Join(combined, "\n\n"))
	}
}
