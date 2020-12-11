package main

import (
	"flag"
	"fmt"
	"log"
)

var verbosityLevel int = 0

func main() {
    // Ensure we exit with an error code and log message
    // when needed after deferred cleanups have run.
    // Credit: https://medium.com/@matryer/golang-advent-calendar-day-three-fatally-exiting-a-command-line-tool-with-grace-874befeb64a4
    var err error
    defer func() {
        if err != nil {
            log.Fatalln(err)
        }
    }()

    //options := getArgs()
}

func getArgs() map[string]string {
    options := make(map[string]string)
    var tplPath string
    var appPath string
    var answers string

    flag.StringVar(&answers, "answers", "", "Path to an answer file.")

    options["tplPath"] = flag.Arg(1)
    options["appPath"] = flag.Arg(2)

    verboseF(1,"downloading \"%v\"", tplPath)
    verboseF(1, "will make \"%v\"",  appPath)

    if answers != "" {
        options["answers"] = flag.Arg(2)
        verboseF(1, "answer are located here: \"%v\"", answers)
    }

    return options
}

func verboseF(lvl int, message string, a ...interface{}) {
    if verbosityLevel > lvl {
        fmt.Printf(message, a...)
    }
}
