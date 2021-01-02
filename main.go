package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	tplPathKey = "templatePath"
	appPathKey = "applicationPath"
	answersKey = "applicationPath"
)

var verbosityLevel int = 0
var errMsgs = [...]string{
	"please specify a path (or URL) to a template",
	"enter a local path to output the app",
}

//TODO: add verbositry messages array.

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

	_, err = getArgs()
}

var answers string

func init() {
	flag.StringVar(&answers, "answers", "", "Path to an answer file.")
	flag.IntVar(&verbosityLevel, "verbose", 0, "extra detail processing info.")
}

func getArgs() (map[string]string, error) {
	var err error
	options := make(map[string]string)

	verboseF(0, "verbose level: %v", verbosityLevel)
	verboseF(1, "number of arguments passed in: %d\n", len(os.Args))
	verboseF(1, "arguments passed in: %v\n", os.Args)

	flag.Parse()

	options[tplPathKey] = flag.Arg(0)
	options[appPathKey] = flag.Arg(1)

	if options[tplPathKey] == "" {
		err = fmt.Errorf(errMsgs[0])
		return options, err
	}

	if options[appPathKey] == "" {
		err = fmt.Errorf(errMsgs[1])
		return options, err
	}

	if answers != "" {
		verboseF(1, "will use answers in the file %q", answers)
		options[answersKey] = answers
	}

	return options, err
}

func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
