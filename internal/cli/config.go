package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	AnswersJson    *AnswersJson // data use for template processing
	AnswersPath    string       // flag to get the path to a file containing values to variables to be parsed.
	OutPath        string       // flag to set the location of the processed template output.
	DataDir        string       // Directory to store app data.
	DefaultVal     string       // Flag to set a default placeholder value when a placeholder is empty.
	TmplPath       string       // flag to set the URL or local template path to a template.
	Tmpl           string       // Path to template, this will be the cached path.
	TmplJson       *TmplJson    // Data about the template such as placeholders, their descriptions, version, etc.
	Branch         string       // flag to set the desired branch of the template to .
	SubCmd         string       // sub-command to execute
	TmplLocation   string       // Indicates local or remote location to downloaded
	TmplType       string       // Flag to indicate the type of package for a template, such as a zip to Extract or a repository to Download.
	CurrentVersion string       // Current semantic version of the application.
	CommitHash     string       // Git commit has of the current version.
	Help           bool         // flag to show the usage for all flags.
	Path           string       // Path to configuration file.
	Version        bool         // flag to show the current version
	UsrOpts        *UserOptions // options that can configured by the user.
	SubCmdConfig   struct {
		FlagSet *flag.FlagSet
		Key     string // config setting
		Method  string // Method to call
		Value   string // value to update config setting
	}
	SubCmdManifest struct {
		FlagSet *flag.FlagSet
		Path    string // path to generate a manifest for.
	}
}

// Setup All application configuration.
func (cfg *Config) Setup(appName, ps string, dirMode os.FileMode) error {
	osDataDir, err1 := stdlib.AppDataDir()
	log.Dbugf("app data dir = %q\n", osDataDir)
	if err1 != nil {
		return err1
	}

	// Make a hidden directory in userspace to store data.
	cfg.DataDir = osDataDir + ps + "." + appName
	if e := os.MkdirAll(cfg.DataDir, dirMode); e != nil {
		return e
	}

	cfg.UsrOpts.CacheDir = cfg.DataDir + ps + "cache"
	if e := os.MkdirAll(cfg.UsrOpts.CacheDir, dirMode); e != nil {
		return fmt.Errorf(Errors.CouldNotMakeCacheDir, e.Error())
	}

	cfg.Path = cfg.DataDir + ps + "config.json"
	// Make a configuration file when there is none.
	if e := cfg.initFile(); e != nil {
		return e
	}

	if e := cfg.LoadUserSettings(cfg.Path); e != nil {
		return e
	}

	// Determine if the template is on the local file system or a remote server.
	cfg.TmplLocation = cfg.getTmplLocation()

	if cfg.TmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		cfg.Tmpl = filepath.Clean(cfg.TmplPath)
	}

	return nil
}

// Initialize a configuration file.
func (cfg *Config) initFile() error {
	if stdlib.PathExist(cfg.Path) {
		log.Infof(Messages.ConfigFileExist, cfg.Path)
		return nil
	}

	f, err1 := os.Create(cfg.Path)
	if err1 != nil {
		return fmt.Errorf(Errors.CouldNotSaveConf, err1.Error())
	}

	data, err2 := json.Marshal(cfg.UsrOpts)
	if err2 != nil {
		return fmt.Errorf(Errors.CouldNotEncodeConfig, err2.Error())
	}

	b, err3 := f.Write(data)
	if err3 != nil {
		return fmt.Errorf(Errors.CouldNotWriteFile, cfg.Path, err3.Error())
	}

	if e := f.Close(); e != nil {
		return fmt.Errorf(Errors.CouldNotCloseFile, cfg.Path, e.Error())
	}

	log.Infof(Messages.MadeNewConfig, b, cfg.Path)

	return nil
}

// save configuration file.
func (cfg *Config) saveUserSettings(mode os.FileMode) error {
	data, err1 := json.Marshal(cfg.UsrOpts)

	log.Dbugf(Messages.SaveData, data)

	if err1 != nil {
		return fmt.Errorf(Errors.CouldNotEncodeConfig, err1.Error())
	}

	if e := os.WriteFile(cfg.Path, data, mode); e != nil {
		return e
	}

	return nil
}

// getTmplLocation Determine if the template is on the local file system or a remote server.
func (cfg *Config) getTmplLocation() string {
	tmplPath := cfg.TmplPath
	regExpAbsolutePath := regexp.MustCompile(`^/([a-zA-Z._\-][a-zA-Z/._\-].*)?`)
	regExpRelativePath := regexp.MustCompile(`^(\.\.|\.|~)(/[a-zA-Z/._\-].*)?`)
	regExpWinDrive := regexp.MustCompile(`^[a-zA-Z]:\\[a-zA-Z/._\\-].*$`)
	pathType := "remote"

	if regExpAbsolutePath.MatchString(tmplPath) ||
		regExpRelativePath.MatchString(tmplPath) ||
		regExpWinDrive.MatchString(tmplPath) {
		pathType = "local"
	}

	return pathType
}

// LoadUserSettings from a file, replacing the default built-in settings.
func (cfg *Config) LoadUserSettings(filename string) error {
	log.Infof(Messages.ReadConfig, filename)
	content, er := ioutil.ReadFile(filename)

	if os.IsNotExist(er) {
		return fmt.Errorf(Errors.CouldNot, er.Error())
	}

	if e := json.Unmarshal(content, &cfg.UsrOpts); e != nil {
		return fmt.Errorf(Errors.CouldNotDecode, filename, er.Error())
	}

	return nil
}

// LoadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func LoadAnswers(filename string) (*AnswersJson, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf(Errors.CannotReadAnswerFile, filename, err.Error())

	}

	var aj *AnswersJson
	if e := json.Unmarshal(content, &aj); e != nil {
		return nil, fmt.Errorf(Errors.CannotDecodeAnswerFile, filename, e.Error())
	}

	return aj, nil
}

// UserOptions Options the user can set
type UserOptions struct {
	ExcludeFileExtensions *[]string
	CacheDir              string
}

func UpdateUserSettings(cfg *Config, mode os.FileMode) error {
	switch cfg.SubCmdConfig.Method {
	case "set":
		if e := cfg.set(cfg.SubCmdConfig.Key, cfg.SubCmdConfig.Value); e != nil {
			return e
		}
		break
	case "get":
		v, e := cfg.get(cfg.SubCmdConfig.Key)
		if e != nil {
			return e
		}
		fmt.Printf("%v", v)
	}

	return cfg.saveUserSettings(mode)
}

// set the value of a user setting
func (cfg *Config) set(key, val string) error {
	switch key {
	case "CacheDir":
		log.Dbugf("setting CacheDir = %q", val)
		cfg.UsrOpts.CacheDir = val
		break
	case "ExcludeFileExtensions":
		log.Dbugf("adding exclusions %q to config", val)
		tmp := strings.Split(val, ",")
		cfg.UsrOpts.ExcludeFileExtensions = &tmp
		break
	default:
		return fmt.Errorf("no %q setting found", key)
	}
	return nil
}

// get the value of a user setting.
func (cfg *Config) get(key string) (interface{}, error) {
	var val interface{}

	switch key {
	case "CacheDir":
		val = cfg.UsrOpts.CacheDir
		break
	case "ExcludeFileExtensions":
		v2 := fmt.Sprintf("%v", val)
		ok, _ := regexp.Match("^[a-zA-Z0-9-.]+(?:,[a-zA-Z0-9-.]+)*", []byte(v2))
		if !ok {
			return nil, fmt.Errorf(Errors.BadExcludeFileExt, val)
		}
		val = strings.Join(*cfg.UsrOpts.ExcludeFileExtensions, ",")
		break
	default:
		return "", fmt.Errorf("no setting %v found", key)
	}

	return val, nil
}

// Validate parses command line flags into program options.
func (cfg *Config) Validate() error {
	if cfg.TmplPath == "" {
		return fmt.Errorf(Errors.TmplPath)
	}

	if cfg.OutPath == "" {
		return fmt.Errorf(Errors.LocalOutPath)
	}

	if cfg.TmplPath == cfg.OutPath {
		return fmt.Errorf(Errors.OutPathCollision, cfg.TmplPath, cfg.OutPath)
	}

	if stdlib.DirExist(cfg.OutPath) {
		return fmt.Errorf(Messages.OutPathExist, cfg.OutPath)
	}

	if cfg.AnswersPath != "" && !stdlib.PathExist(cfg.AnswersPath) {
		return fmt.Errorf(Errors.AnswerFile404, cfg.AnswersPath)
	}

	regExpTmplType := regexp.MustCompile("^(zip|git|dir)$")

	if !regExpTmplType.MatchString(cfg.TmplType) {
		return fmt.Errorf(Errors.BadTmplType, cfg.TmplType)
	}

	return nil
}
