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
	AllowedUrls           []string // TODO: Remove this obsolete option URLs you are allowed to download from.
	answersPath           string   // path to a file containing values to variables to be parsed.
	appPath               string   // Location of the processed template output.
	cacheDir              string
	tmplPath              string    // URL or local path to a template.
	tmpl                  string    // Path to template, this will be the cached path.
	ExcludeFileExtensions []string  // Files to skip when sending to the go parsing engine.
	IncludeFileExtensions []string  // Files to include when sending to the go parsing engine.
	Questions             questions // Question for requesting input for the template.
	branch                string    // Desired branch to clone.
	tmplLocation          string    // Indicates local or remote location to downloaded
	tmplType              string    // Indicates a zip to extract or a repository to download.
	CurrentVersion        string
	CommitHash            string
	help                  bool
	version               bool
}

// Load configuration file.
func initConfigFile(file string) (err error) {
	if stdlib.PathExist(file) {
		infof("config file %q exist", file)
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

	infof("made a new config file %q exist", file)

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

// loadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func loadAnswers(filename string) (answers tplVars, err error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		err = fmt.Errorf(errors.cannotReadAnswerFile, filename, err.Error())
		return
	}

	err = json.Unmarshal(content, &answers)
	if err != nil {
		err = fmt.Errorf(errors.cannotDecodeAnswerFile, filename, err.Error())
		return
	}

	return
}
