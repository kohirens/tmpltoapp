package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kohirens/go-gitter/stdlib"
	"github.com/kohirens/go-gitter/template"
)

const (
	PS = string(os.PathSeparator)
)

var (
	verbosityLevel int = 0
	errMsgs            = [...]string{
		"please specify a path (or URL) to a template",
		"enter a local path to output the app",
		"the following error occurred trying to get the app data directory: %q",
		"path/URL to template is not in the allow-list",
	}
)

type Config struct {
	allowedUrls    []string
	answersPath    string
	appPath        string
	tplPath        string
	verbosityLevel int
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

	configFile := "settings.json"
	appDataDir, er := stdlib.HomeDir()
	if er == nil {
		configFile = appDataDir + PS + "settings.json"
		return
	}

	verboseF(0, "config location %q", configFile)

	options, err := settings(configFile)
	if err != nil {
		return
	}

	err = parseArgs(os.Args[0], os.Args[1:], &options)
	if err != nil {
		return
	}

	isUrl, isAllowed := urlIsAllowed(options.tplPath, options.allowedUrls)
	if isUrl && !isAllowed {
		err = fmt.Errorf(errMsgs[3])
		return
	}

	if isUrl {
		client := http.Client{}
		err = template.Download(options.tplPath, options.appPath, &client)
	}
	// TODO: local copy.
}

// Process any program flags fed into the program.
func parseArgs(progName string, pArgs []string, options *Config) (err error) {

	verboseF(1, "running program %q", progName)

	pFlags := flag.NewFlagSet(progName, flag.ExitOnError)

	pFlags.StringVar(&options.answersPath, "answers", "", "Path to an answer file.")
	pFlags.IntVar(&options.verbosityLevel, "verbose", 0, "extra detail processing info.")

	verboseF(1, "verbose level: %v", verbosityLevel)
	verboseF(1, "number of arguments passed in: %d\n", len(os.Args))
	verboseF(1, "arguments passed in: %v\n", os.Args)

	pFlags.Parse(pArgs)

	options.tplPath = pFlags.Arg(0)
	options.appPath = pFlags.Arg(1)

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

func settings(filename string) (cfg Config, err error) {
	var data map[string]interface{}

	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}

	val, ok := data["allowedUrls"].([]interface{})
	if ok && len(data) > 0 {
		cfg.allowedUrls = make([]string, len(data))
		for i, v := range val {
			cfg.allowedUrls[i] = v.(string)
		}
	}

	return
}

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

func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
