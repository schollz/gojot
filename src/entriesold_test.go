package jot

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
	if dates[0] != "Fri Sep 26 21:08:48 2014" {
		t.Errorf(dates[0])
	}
	if texts[1] != `Another entry
Another text2` {
		t.Errorf(texts[1])
	}
	if dates[1] != "Sat Sep 26 00:00:00 2015" {
		t.Errorf(dates[1])
	}
}

func TestProcessOld2(t *testing.T) {
	fulltext := `
# alsdkfjasldkfjalskjdflaksjdlfjaslkdjflas
2014-09-26 21:08:48 Some new entry

Some Text1

# alsdkfjasldkfjalskjdflaksjdlfjaslkdjflas
2015-09-26 Another entry

Another text2

`
	texts, dates := ProcessEntriesOld(fulltext)
	if texts[0] != `Some new entry
Some Text1` {
		t.Errorf(texts[0])
	}
	if dates[0] != "Fri Sep 26 21:08:48 2014" {
		t.Errorf(dates[0])
	}
	if texts[1] != `Another entry
Another text2` {
		t.Errorf(texts[1])
	}
	if dates[1] != "Sat Sep 26 00:00:00 2015" {
		t.Errorf(dates[1])
	}
}

// # af8dbd0d43fb55458f11aad586ea2abf
// 2013-05-02 15:30 My first DayOne entry
// This is it, I probably won't edit for another yer.
//
// # 2391048fe24111e1983ed49a20be6f9e
// 2014-05-03 03:22 Second @entry
// Wow, looks like its been a whole year!
//
// Thu May 02 15:30:00 2013 -==- FretfulUnusedAirline
//
// My first jot entry
// This is it, I probably won't edit for another yer.
//
// Sat May 03 03:22:00 2014 -==- ViolentZestyDonald
//
// Second @entry
// Wow, looks like its been a whole year!
