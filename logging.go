package main

import "fmt"

var verbosityLevel = 0

// Show additional logging based on the verbosity level. Prints a newline after every message.
func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
