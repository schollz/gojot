package gitsdees

import "github.com/jcelliott/lumber"

var logger *lumber.ConsoleLogger

func init() {
	logger = lumber.NewConsoleLogger(lumber.DEBUG)
	logger.Level(0)
}

func DebugMode() {
	logger.Level(2)
}
