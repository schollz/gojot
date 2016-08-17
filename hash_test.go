package main

import "fmt"

func ExampleTestHash() {
	fmt.Println(encodeNumber(int(1230)))
	// Output: 9gdED8
}

func ExampleTestHash2() {
	fmt.Println(decodeNumber(encodeNumber(int(1230))))
	// Output: 1230
}
