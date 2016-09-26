package gitsdees

import (
	"strings"
	"testing"
)

func TestProcessFiles(t *testing.T) {
	var cache Cache
	cache.Branch = make(map[string]Entry)
	cache.Branch["1"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "one"}
	cache.Branch["2"] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: "two"}
	cache.Branch["3"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "three"}
	_, hashes := CombineEntries(cache)
	cache.Branch["1"] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "oneEDIT7"}
	cache.Branch["3"] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "threeEDIT7"}
	combined, _ := CombineEntries(cache)
	testEntry := strings.Join(combined, "\n\n")
	branchesUpdated := ProcessEntries(testEntry, hashes)
	if branchesUpdated[0] != "1" && branchesUpdated[1] != "3" {
		t.Errorf("Error processing files")
	}
}
