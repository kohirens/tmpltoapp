package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var verbosityLevel int = 0
var errMsgs = [...]string{
	"please specify a path (or URL) to a template",
	"enter a local path to output the app",
}

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
}

func getArgs() (map[string]string, error) {
	var tplPath string
	var appPath string
	var answers string
	var err error
	options := make(map[string]string)

	verboseF(1, "no. arguments passed in: %d\n", len(os.Args))
	verboseF(1, "arguments passed in: %v\n", os.Args)

	if len(os.Args) < 2 {
		err = fmt.Errorf("display help")
	}

	flag.StringVar(&answers, "answers", "", "Path to an answer file.")
	flag.IntVar(&verbosityLevel, "verbosity", 0, "extra detail processing info.")

	options["tplPath"] = flag.Arg(0)
	options["appPath"] = flag.Arg(1)

	verboseF(1, "downloading \"%v\"\n", tplPath)
	verboseF(1, "will make \"%v\"\n", appPath)

	if options["tplPath"] == "" {
		err = fmt.Errorf(errMsgs[0])
		return options, err
	}

	if options["appPath"] == "" {
		err = fmt.Errorf(errMsgs[1])
		return options, err
	}

	if answers != "" {
		verboseF(1, "answer are located here: \"%v\"", answers)
	}

	return options, err
}

func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
	}
}
