package main

import "fmt"

const (
	verboseLvlInfo = 1
	verboseLvlWarn = 2
	verboseLvlErr  = 3
	verboseLvlDbug = 4
)

var verbosityLevel = 0

// Show additional logging based on the verbosity level. Prints a newline after every message.
func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}

func infof(message string, vars ...interface{}) {
	verboseF(verboseLvlInfo, message, vars...)
}

func warnf(message string, vars ...interface{}) {
	verboseF(verboseLvlWarn, message, vars...)
}

func errf(message string, vars ...interface{}) {
	verboseF(verboseLvlErr, message, vars...)
}

func dbugf(message string, vars ...interface{}) {
	verboseF(verboseLvlDbug, message, vars...)
}
