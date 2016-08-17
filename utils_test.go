package main

import "fmt"

func ExampleParseDate() {
	fmt.Println(parseDate("2016-04-12 15:04"))
	// Output: true 1460473440
}

func ExampleParseDate2() {
	fmt.Println(parseDate("2016-0lkjkljl4-12 15:04"))
	// Output: false -1
}

func ExampleIsNumber() {
	num, isNum := isNumber("1")
	fmt.Println(num, isNum)
	// Output: 1 true
}

func ExampleIsNotNumber() {
	num, isNum := isNumber("journal")
	fmt.Println(num, isNum)
	// Output: -1 false
}
