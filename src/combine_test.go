package gitsdees

import (
	"strings"
	"testing"
)

func TestCombine(t *testing.T) {
	var cache Cache
	cache.Branch = make(map[string]Entry)
	cache.Branch["1"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "one"}
	cache.Branch["2"] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: "two"}
	cache.Branch["3"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "three"}
	var branchHashes = make(map[string]string)
	branchHashes["1"] = "uKn3Fdu"
	branchHashes["2"] = "NdbTNGe"
	branchHashes["3"] = "+XxdKZQ"
	combined, _ := CombineEntries(cache)
	if strings.Contains(combined[0], "one") != true {
		t.Errorf("Got [%s] instead of the correct", strings.Join(combined, "\n\n"))
	}
}
