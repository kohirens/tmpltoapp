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
	tplPathKey = "templatePath"
	appPathKey = "applicationPath"
	answersKey = "applicationPath"
	PS         = string(os.PathSeparator)
)

var (
	verbosityLevel int    = 0
	answers        string = ""
	errMsgs               = [...]string{
		"please specify a path (or URL) to a template",
		"enter a local path to output the app",
		"the following error occurred trying to get the app data directory: %q",
		"path/URL to template is not in the allow-list",
	}
)

type config map[string]interface{}

func (c config) Array(key string) (v []string, ok bool) {
	x, ok := c[key]
	if ok {
		y := x.([]interface{})
		v = make([]string, len(y))
		for k, z := range y {
			v[k] = z.(string)
		}
	}
	// TODO: always return an initialized array
	return
}

func (c config) String(key string) string {
	return fmt.Sprintf("%v", c[key])
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

	verboseF(0, "go-gitter config location %q", configFile)

	options, err := settings(configFile)
	if err != nil {
		return
	}

	options, err = getArgs()
	if err != nil {
		return
	}

	tplLoc := options[tplPathKey].(string)
	allowedUrls, _ := options.Array("allowedUrls")
	isUrl, isAllowed := urlIsAllowed(tplLoc, allowedUrls)

	if isUrl && !isAllowed {
		err = fmt.Errorf(errMsgs[3])
		return
	}

	if isUrl {
		client := http.Client{}
		err = template.Download(tplLoc, options[appPathKey].(string), &client)
	}
	// TODO: local copy.
}

func init() {
	flag.StringVar(&answers, "answers", "", "Path to an answer file.")
	flag.IntVar(&verbosityLevel, "verbose", 0, "extra detail processing info.")
}

func getArgs() (config, error) {
	var err error
	options := make(map[string]interface{})

	verboseF(1, "verbose level: %v", verbosityLevel)
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

func settings(filename string) (cfg config, err error) {
	var data interface{}

	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}

	// Convert the interface to a map.
	cfg = data.(map[string]interface{})

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
