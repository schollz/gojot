package main

import "testing"

func TestProcessFiles(t *testing.T) {
	testEntry := `Sat Sep 24 15:12:50 2016 -0400 -==- AwIZ5 -==- HknvYLLB

Another entry EDITED

Sat Sep 24 15:13:22 2016 -0400 -==- b3VrS -==- LZHOGD1A

New new entry

Sat Sep 24 15:13:31 2016 -0400 -==- 3qgd9 -==- STWCnX8k

Another new entry

Sat Sep 24 15:13:52 2016 -0400 -==-  NEW  -==- rsNLD0eL

jlkjlkjl sadflkja sdflkajsdf alkjs`
	branchesUpdated := ProcessEntries(testEntry)
	if branchesUpdated[0] != "AwIZ5" && branchesUpdated[1] != "qQx0X" {
		t.Errorf("Error processing files")
	}
}
