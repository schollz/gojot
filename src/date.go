package gitsdees

import (
	"strings"
	"time"
)

func GetCurrentDate() string {
	// git date format: Thu, 07 Apr 2005 22:13:13 +0200
	formattedTime := time.Now().Format("Mon Jan 02 15:04:05 2006 -0700")
	return formattedTime
}

func ParseDate(date string) (time.Time, error) {
	date = strings.TrimSpace(date)
	newTime, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		newTime, err = time.Parse("Mon Jan 02 15:04:05 2006 -0700", date)
	}
	if err != nil {
		newTime, err = time.Parse("Mon Jan 2 15:04:05 2006 -0700", date)
	}
	if err != nil {
		newTime, err = time.Parse("Mon, Jan 02 15:04:05 2006 -0700", date)
	}
	if err != nil {
		newTime, err = time.Parse("Mon 02 Jan 2006 15:04:05 -0700", date)
	}
	if err != nil {
		newTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", date)
	}
	if err != nil {
		newTime, err = time.Parse("2006-01-02 15:04:05", date)
	}
	if err != nil {
		newTime, err = time.Parse("2006-01-02 15:04", date)
	}
	if err != nil {
		newTime, err = time.Parse("2006-01-02", date)
	}
	return newTime, err
}

func FormatDate(date time.Time) string {
	return date.Format("Mon Jan 02 15:04:05 2006 -0700")
}
