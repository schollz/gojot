package sdees

import (
	"strings"
	"testing"
)

func TestProcessFiles(t *testing.T) {
	var cache Cache
	cache.Branch = make(map[string]Entry)
	cache.Branch[ShortEncrypt("1")] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "one"}
	cache.Branch[ShortEncrypt("2")] = Entry{Date: "Fri, 08 Apr 2005 22:13:13 +0200", Text: "two"}
	cache.Branch[ShortEncrypt("3")] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "three"}
	_, _, hashes := CombineEntries(cache)
	cache.Branch[ShortEncrypt("1")] = Entry{Date: "Thu, 07 Apr 2005 22:13:13 +0200", Text: "oneEDIT7"}
	cache.Branch[ShortEncrypt("3")] = Entry{Date: "Sat, 09 Apr 2005 22:13:13 +0200", Text: "threeEDIT7"}
	combined, _, _ := CombineEntries(cache)
	testEntry := strings.Join(combined, "\n\n")
	branchesUpdated := UpdateEntryFromText(testEntry, hashes)
	if branchesUpdated[0] != ShortEncrypt("1") && branchesUpdated[1] != ShortEncrypt("3") && len(branchesUpdated) == 2 {
		t.Errorf("Error processing files")
	}
}
