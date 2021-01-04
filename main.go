package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kohirens/go-gitter/stdlib"
	"github.com/kohirens/go-gitter/template"
)

const (
	tplPathKey = "templatePath"
	appPathKey = "applicationPath"
	answersKey = "applicationPath"
	PS         = string(os.PathSeparator)
)

var verbosityLevel int = 0
var errMsgs = [...]string{
	"please specify a path (or URL) to a template",
	"enter a local path to output the app",
	"the following error occurred trying to get the app data directory: %q",
}

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
	appDataDir, er := stdlib.AppDataDir()
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

	client := http.Client{}
	err = template.Download(options[tplPathKey].(string), options[appPathKey].(string), &client)
}

var answers string

func init() {
	flag.StringVar(&answers, "answers", "", "Path to an answer file.")
	flag.IntVar(&verbosityLevel, "verbose", 0, "extra detail processing info.")
}

func getArgs() (map[string]interface{}, error) {
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

func verboseF(lvl int, message string, a ...interface{}) {
	if verbosityLevel >= lvl {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
