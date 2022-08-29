package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	answersJson           *answersJson // data use for template processing
	answersPath           string       // flag to get the path to a file containing values to variables to be parsed.
	outPath               string       // flag to set the location of the processed template output.
	cacheDir              string
	tmplPath              string    // flag to set the URL or local path to a template.
	tmpl                  string    // Path to template, this will be the cached path.
	ExcludeFileExtensions *[]string // Files to skip when sending to the go parsing engine.
	IncludeFileExtensions *[]string // Files to include when sending to the go parsing engine.
	TmplJson              *tmplJson // Question for requesting input for the template.
	branch                string    // flag to set the desired branch to clone.
	tmplLocation          string    // Indicates local or remote location to downloaded
	tmplType              string    // Flag to indicate the type of package for a template, such as a zip to extract or a repository to download.
	CurrentVersion        string
	CommitHash            string
	help                  bool // flag to show the usage for all flags.
	version               bool // flag to show the current version
}

// configMain initialize the application configuration
func configMain(appDataDir string) error {
	appConfig.cacheDir = appDataDir + PS + "cache"
	e1 := os.MkdirAll(appConfig.cacheDir, DirMode)
	if e1 != nil {
		return fmt.Errorf("could not make cache directory, error: %s", e1.Error())
	}

	appConfig.tmplLocation = getTmplLocation(appConfig.tmplPath)

	if appConfig.tmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		appConfig.tmpl = filepath.Clean(appConfig.tmplPath)
	}

	return nil
}

// Load configuration file.
func initConfigFile(file string) (err error) {
	if stdlib.PathExist(file) {
		infof(messages.configFileExist, file)
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

	infof(messages.madeNewConfig, file)

	return
}

// settings runtime options are a mix of config and command line arguments.
func settings(filename string, cfg *Config) (err error) {
	content, er := ioutil.ReadFile(filename)

	if os.IsNotExist(er) {
		err = fmt.Errorf(errors.couldNot, er.Error())
		return
	}

	er = json.Unmarshal(content, cfg)
	if er != nil {
		err = fmt.Errorf(errors.couldNotDecode, filename, er.Error())
		return
	}

	return
}

// loadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func loadAnswers(filename string) (aj *answersJson, err error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		err = fmt.Errorf(errors.cannotReadAnswerFile, filename, err.Error())
		return
	}

	err = json.Unmarshal(content, &aj)
	if err != nil {
		err = fmt.Errorf(errors.cannotDecodeAnswerFile, filename, err.Error())
		return
	}

	return
}
