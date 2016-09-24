package main

import (
	"strings"
	"testing"
)

func TestCombine(t *testing.T) {
	var cache Cache
	cache.Branch = make(map[string]Entry)
	cache.Branch["asdfasdf"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "one"}
	cache.Branch["asdfa"] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: "two"}
	cache.Branch["fajskdf"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "three"}
	combined := CombineEntries(cache)
	if strings.Contains(combined[0], "one") != true {
		t.Errorf("Got [%s] instead of the correct", strings.Join(combined, "\n\n"))
	}
}
