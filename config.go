package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"io/ioutil"
	"os"
)

type Config struct {
	answers               tplVars
	AllowedUrls           []string // URLs allowed to download from.
	answersPath           string   //
	appPath               string
	cacheDir              string
	tplPath               string
	tmpl                  string
	ExcludeFileExtensions []string // Files to skip when sending to the go parsing engine.
	IncludeFileExtensions []string // Files to include when sending to the go parsing engine.
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

	defer func() {
		err = f.Close()
	}()

	_, err = f.WriteString(DEFAULT_CFG)

	verboseF(1, "made a new config file %q exist", file)

	return
}

// settings runtime options are a mix of config and command line arguments.
func settings(filename string, cfg *Config) (err error) {
	content, er := ioutil.ReadFile(filename)

	if os.IsNotExist(er) {
		err = fmt.Errorf("could not %s", er.Error())
		return
	}

	er = json.Unmarshal(content, cfg)
	if er != nil {
		err = fmt.Errorf("could not decode %q, error: %s", filename, er.Error())
		return
	}

	return
}

// loadAnswers Load key/value pairs from a JSON file to fill in placeholders when processing Go templates.
func loadAnswers(filename string) (answers tplVars, err error) {
	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &answers)
	if err != nil {
		err = fmt.Errorf("could not decode JSON file %q, because of error: %s", filename, err.Error())
		return
	}

	return
}
