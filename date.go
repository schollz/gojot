package main

import "time"

func GetCurrentDate() string {
	// git date format: Thu, 07 Apr 2005 22:13:13 +0200
	formattedTime := time.Now().Format("Mon, 02 Jan 2006 15:04:05 +0200")
	return formattedTime
}
