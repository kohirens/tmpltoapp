package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	answersJson    *answersJson // data use for template processing
	answersPath    string       // flag to get the path to a file containing values to variables to be parsed.
	outPath        string       // flag to set the location of the processed template output.
	dataDir        string       // Directory to store app data.
	defaultVal     string       // Flag to set a default placeholder value when a placeholder is empty.
	tmplPath       string       // flag to set the URL or local template path to a template.
	tmpl           string       // Path to template, this will be the cached path.
	TmplJson       *tmplJson    // Data about the template such as placeholders, their descriptions, version, etc.
	branch         string       // flag to set the desired branch of the template to .
	subCmd         string       // sub-command to execute
	tmplLocation   string       // Indicates local or remote location to downloaded
	tmplType       string       // Flag to indicate the type of package for a template, such as a zip to extract or a repository to download.
	CurrentVersion string       // Current semantic version of the application.
	CommitHash     string       // Git commit has of the current version.
	help           bool         // flag to show the usage for all flags.
	path           string       // Path to configuration file.
	version        bool         // flag to show the current version
	usrOpts        *userOptions // options that can configured by the user.
	subCmdConfig   struct {
		flagSet *flag.FlagSet
		key     string // config setting
		method  string // method to call
		value   string // value to update config setting
	}
}

// setup All application configuration.
func (cfg *Config) setup(appName, ps string, dirMode os.FileMode) error {
	osDataDir, err1 := stdlib.AppDataDir()
	if err1 != nil {
		return err1
	}

	// Make a directory to store data.
	cfg.dataDir = osDataDir + ps + appName
	if e := os.MkdirAll(cfg.dataDir, dirMode); e != nil {
		return e
	}

	cfg.usrOpts.CacheDir = cfg.dataDir + ps + "cache"
	if e := os.MkdirAll(cfg.usrOpts.CacheDir, dirMode); e != nil {
		return fmt.Errorf(errors.couldNotMakeCacheDir, e.Error())
	}

	cfg.path = cfg.dataDir + ps + "config.json"
	// Make a configuration file when there is none.
	if e := cfg.initFile(); e != nil {
		return e
	}

	if e := cfg.loadUserSettings(cfg.path); e != nil {
		return e
	}

	// Determine if the template is on the local file system or a remote server.
	cfg.tmplLocation = cfg.getTmplLocation()

	if cfg.tmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		cfg.tmpl = filepath.Clean(cfg.tmplPath)
	}

	return nil
}

// Initialize a configuration file.
func (cfg *Config) initFile() error {
	if stdlib.PathExist(cfg.path) {
		infof(messages.configFileExist, cfg.path)
		return nil
	}

	f, err1 := os.Create(cfg.path)
	if err1 != nil {
		return fmt.Errorf(errors.couldNotSaveConf, err1.Error())
	}

	data, err2 := json.Marshal(cfg.usrOpts)
	if err2 != nil {
		return fmt.Errorf(errors.couldNotEncodeConfig, err2.Error())
	}

	b, err3 := f.Write(data)
	if err3 != nil {
		return fmt.Errorf(errors.couldNotWriteFile, cfg.path, err3.Error())
	}

	if e := f.Close(); e != nil {
		return fmt.Errorf(errors.couldNotCloseFile, cfg.path, e.Error())
	}

	infof(messages.madeNewConfig, b, cfg.path)

	return nil
}

// save configuration file.
func (cfg *Config) saveUserSettings(mode os.FileMode) error {
	data, err1 := json.Marshal(cfg.usrOpts)

	dbugf("%s\n", data)

	if err1 != nil {
		return fmt.Errorf(errors.couldNotEncodeConfig, err1.Error())
	}

	if e := os.WriteFile(cfg.path, data, mode); e != nil {
		return e
	}

	return nil
}

// getTmplLocation Determine if the template is on the local file system or a remote server.
func (cfg *Config) getTmplLocation() string {
	tmplPath := cfg.tmplPath
	regExpAbsolutePath := regexp.MustCompile(`^/([a-zA-Z._\-][a-zA-Z/._\-].*)?`)
	pathType := "remote"

	if regExpAbsolutePath.MatchString(tmplPath) || regExpRelativePath.MatchString(tmplPath) || regExpWinDrive.MatchString(tmplPath) {
		pathType = "local"
	}

	return pathType
}

// loadUserSettings from a file, replacing the default built-in settings.
func (cfg *Config) loadUserSettings(filename string) error {
	content, er := ioutil.ReadFile(filename)

	if os.IsNotExist(er) {
		return fmt.Errorf(errors.couldNot, er.Error())
	}

	if e := json.Unmarshal(content, cfg); e != nil {
		return fmt.Errorf(errors.couldNotDecode, filename, er.Error())
	}

	return nil
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
