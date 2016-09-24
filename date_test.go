package main

import "testing"

func TestParseDate(t *testing.T) {
	date, _ := ParseDate("Thu, 07 Apr 2005 22:13:13 +0200")
	if date.String() != "2005-04-07 22:13:13 +0200 +0200" {
		t.Errorf("Expected %s and got %s", "2005-04-02 22:13:13 +0000 UTC", date.String())
	}
}

func TestFormatDate(t *testing.T) {
	testDate := "Thu Apr 07 22:13:13 2005 +0200"
	date, err := ParseDate(testDate)
	if err != nil {
		t.Errorf(err.Error())
	}
	if FormatDate(date) != testDate {
		t.Errorf("Expected %s and got %s", testDate, FormatDate(date))
	}

	testDate = "Fri Sep 23 21:50:57 2016 -0400"
	date, err = ParseDate(testDate)
	if err != nil {
		t.Errorf(err.Error())
	}
	if FormatDate(date) != testDate {
		t.Errorf("Expected %s and got %s", testDate, FormatDate(date))
	}

}
