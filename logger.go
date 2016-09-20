package main

import "github.com/jcelliott/lumber"

var logger *lumber.ConsoleLogger

func init() {
	logger = lumber.NewConsoleLogger(lumber.DEBUG)
}
