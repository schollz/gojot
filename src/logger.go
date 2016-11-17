package gojot

import "github.com/jcelliott/lumber"

var logger *lumber.ConsoleLogger

func init() {
	logger = lumber.NewConsoleLogger(lumber.DEBUG)
	logger.Level(2)
}

func DebugMode() {
	logger.Level(0)
}
