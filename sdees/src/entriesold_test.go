package sdees

import "testing"

func TestProcessOld(t *testing.T) {
	fulltext := `2014-09-26 21:08:48 Some new entry

Some Text1

2015-09-26 Another entry

Another text2

`
	texts, dates := ProcessEntriesOld(fulltext)
	if texts[0] != `Some new entry
Some Text1` {
		t.Errorf(texts[0])
	}
	if dates[0] != "Fri Sep 26 21:08:48 2014 +0000" {
		t.Errorf(dates[0])
	}
	if texts[1] != `Another entry
Another text2` {
		t.Errorf(texts[1])
	}
	if dates[1] != "Sat Sep 26 00:00:00 2015 +0000" {
		t.Errorf(dates[1])
	}
}
