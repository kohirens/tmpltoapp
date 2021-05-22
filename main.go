package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kohirens/stdlib"
)

const (
	PS       = string(os.PathSeparator)
	CONF     = "config.json"
	DIR_MODE = 0774
)

var (
	errMsgs = [...]string{
		"please specify a path (or URL) to a template",
		"enter a local path to output the app",
		"the following error occurred trying to get the app data directory: %q",
		"path/URL to template is not in the allow-list",
		"template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
	}
)

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Fatalln(err)
		}
	}()

	configFile := "config.json"
	appDataDir, err := stdlib.HomeDir()
	if err == nil {
		configFile = appDataDir + PS + "config.json"
	}

	err = initConfigFile(configFile)
	if err != nil {
		return
	}

	options, err := settings(configFile)
	if err != nil {
		return
	}

	err = parseArgs(os.Args[0], os.Args[1:], &options)
	if err != nil {
		return
	}

	verboseF(1, "config location %q", configFile)

	isUrl, isAllowed := urlIsAllowed(options.tplPath, options.allowedUrls)
	if isUrl && !isAllowed {
		err = fmt.Errorf(errMsgs[3])
		return
	}
	verboseF(1, "isUrl %v", isUrl)

	options.cacheDir = appDataDir + PS + "cache"
	err = os.MkdirAll(options.cacheDir, DIR_MODE)
	if err != nil {
		err = fmt.Errorf("could not make cache directory, error: %s", err.Error())
		return
	}

	if isUrl {
		client := http.Client{}
		zipFile, iErr := Download(options.tplPath, options.cacheDir, &client)
        if iErr != nil {
            err = iErr
            return
        }

        iErr = Extract(zipFile, os.TempDir())
        if iErr != nil {
            err = iErr
            return
        }
	}
	// TODO: local copy.

	// Parse
}

// Process any program flags fed into the program.
func parseArgs(progName string, pArgs []string, options *Config) (err error) {

	pFlags := flag.NewFlagSet(progName, flag.ExitOnError)

	pFlags.StringVar(&options.answersPath, "answers", "", "Path to an answer file.")
	pFlags.IntVar(&verbosityLevel, "verbose", 0, "extra detail processing info.")

	pFlags.Parse(pArgs)

	verboseF(2, "running program %q", progName)
	verboseF(1, "verbose level: %v", verbosityLevel)
	verboseF(1, "number of arguments passed in: %d", len(os.Args))
	verboseF(1, "arguments passed in: %v", os.Args)

	options.tplPath = pFlags.Arg(0)
	options.appPath = pFlags.Arg(1)
	options.verbosityLevel = verbosityLevel

	if options.tplPath == "" {
		err = fmt.Errorf(errMsgs[0])
		return
	}

	if options.appPath == "" {
		err = fmt.Errorf(errMsgs[1])
		return
	}

	if options.answersPath != "" {
		verboseF(1, "will use answers in the file %q", options.answersPath)
	}

	return
}

// Check to see if a URL is in the allowed list to download template from.
func urlIsAllowed(loc string, urls []string) (isUrl, isAllowed bool) {
	isUrl = strings.HasPrefix(loc, "https://")
	isAllowed = false

	if isUrl {
		for _, url := range urls {
			if strings.HasPrefix(loc, url) {
				isAllowed = true
				break
			}
		}
	}

	return
}
