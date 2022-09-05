package main

import (
	"encoding/json"
	"flag"
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
	cacheDir              string       // Cache for downloaded templates.
	dataDir               string       // Directory to store app data.
	defaultVal            string       // Flag to set a default placeholder value when a placeholder is empty.
	tmplPath              string       // flag to set the URL or local path to a template.
	tmpl                  string       // Path to template, this will be the cached path.
	ExcludeFileExtensions *[]string    // Files to skip when sending to the go parsing engine.
	IncludeFileExtensions *[]string    // Files to include when sending to the go parsing engine.
	TmplJson              *tmplJson    // Question for requesting input for the template.
	branch                string       // flag to set the desired branch to clone.
	subCmd                string       // sub-command to execute
	tmplLocation          string       // Indicates local or remote location to downloaded
	tmplType              string       // Flag to indicate the type of package for a template, such as a zip to extract or a repository to download.
	CurrentVersion        string       // Current semantic version of the application.
	CommitHash            string       // Git commit has of the current version.
	help                  bool         // flag to show the usage for all flags.
	path                  string       // Path to configuration file.
	version               bool         // flag to show the current version
	subCmdConfig          struct {
		flagSet *flag.FlagSet
		key     string // config setting
		method  string // method to call
		value   string // value to update config setting
	}
}

// setup All application configuration.
func (cfg *Config) setup() error {
	osDataDir, err1 := stdlib.AppDataDir()
	if err1 != nil {
		return err1
	}

	// Make a directory to store data.
	cfg.dataDir = osDataDir + PS + AppName
	if e := os.MkdirAll(cfg.dataDir, DirMode); e != nil {
		return e
	}

	// Make a configuration file when there is none.
	cfg.path = cfg.dataDir + PS + "config.json"
	if e := initConfigFile(cfg.path); e != nil {
		return e
	}

	if e := settings(cfg.path, cfg); e != nil {
		return e
	}

	cfg.cacheDir = cfg.dataDir + PS + "cache"
	if e1 := os.MkdirAll(cfg.cacheDir, DirMode); e1 != nil {
		return fmt.Errorf("could not make cache directory, error: %s", e1.Error())
	}

	cfg.tmplLocation = getTmplLocation(cfg.tmplPath)

	if cfg.tmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		cfg.tmpl = filepath.Clean(cfg.tmplPath)
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

// save configuration file.
func saveConfigFile(file string, defCfg *userOptions) error {
	data, err := json.Marshal(defCfg)
	if err != nil {
		return fmt.Errorf("could not convert coniguration to JSON string: %v", err)
	}

	if e := os.WriteFile(file, data, DirMode); e != nil {
		return fmt.Errorf("could not save configuration: %v", e.Error())
	}

	return nil
}

// TODO: Rename to loadUserSettings as a method of Config
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
