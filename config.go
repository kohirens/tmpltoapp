package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"io/ioutil"
	"os"
)

type Config struct {
	answers tplVars
	allowedUrls    []string // URLs allowed to download from.
	answersPath    string //
	appPath        string
	cacheDir       string
	tplPath        string
	verbosityLevel int
}

// Load configuration file.
func initConfigFile(file string) (err error) {

	if stdlib.PathExist(file) {
		verboseF(1, "config file %q exist", file)
		return
	}

	f, er := os.Create(file)

	if er != nil {
		err = er
		return
	}

	defer func () {
		err = f.Close()
	}()

	_, err = f.WriteString(DEFAULT_CFG)

	verboseF(1, "made a new config file %q exist", file)

	return
}

// settings runtime options are a mix of config and command line arguments.
func settings(filename string) (cfg Config, err error) {
	var data map[string]interface{}

	content, er := ioutil.ReadFile(filename)

	if os.IsNotExist(er) {
		err = fmt.Errorf("could not %s", er.Error())
		return
	}

	er = json.Unmarshal(content, &data)
	if er != nil {
		err = fmt.Errorf("could not decode %q, error: %s", filename, er.Error())
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

// loadAnswers Loan answers from a JSON file to be used to fill in placeholders when processing Go templates.
func loadAnswers(filename string) (answers tplVars, err error) {
	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &answers)
	if err != nil {
		return
	}

	return
}
