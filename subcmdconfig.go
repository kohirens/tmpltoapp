package main

import (
	"fmt"
	"os"
)

// userOptions Options the user can set
type userOptions struct {
	ExcludeFileExtensions []string
	IncludeFileExtensions []string
	cacheDir              string
}

var usrOpts = &userOptions{
	ExcludeFileExtensions: []string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"},
	IncludeFileExtensions: []string{},
}

// runs during flag parsing
func (cfg *Config) subCmdConfigMain(osArgs []string) error {
	if e := cfg.subCmdConfig.flagSet.Parse(osArgs); e != nil {
		return fmt.Errorf("error pasing sub command config flags: %v", e.Error())
	}

	if len(osArgs) < 2 {
		subCmdConfigUsage()
		return fmt.Errorf("invalid number of arguments passed to config sub-command, please try config -h for usage")
	}

	cfg.subCmd = cmdConfig

	cfg.subCmdConfig.method = osArgs[0]
	cfg.subCmdConfig.key = osArgs[1]
	cfg.subCmdConfig.value = osArgs[2]

	if cfg.help {
		Usage(cfg)
		os.Exit(0)
	}

	return nil
}

func runSubCmdConfig(cfg *Config) error {
	switch cfg.subCmdConfig.method {
	case "set":
		if e := setSetting(cfg.subCmdConfig.key, cfg.subCmdConfig.value); e != nil {
			return e
		}
		break
	case "get":
		fmt.Printf("%v", setting(cfg.subCmdConfig.key))
	}

	return saveConfigFile(cfg.path, usrOpts)
}

func subCmdConfigUsage() {
	fmt.Println("usage: config set|get <args>")
	fmt.Println("example: config set \"cache\" \"./path/to/a/directory\"")
}

func setSetting(key, val string) error {
	switch key {
	case "cacheDir":
		usrOpts.cacheDir = val
		break
	default:
		return fmt.Errorf("no setting found %v", key)
	}
	return nil
}

// setting get the value of a user setting.
func setting(key string) error {
	switch key {
	case "cacheDir":
		fmt.Printf("%v", usrOpts.cacheDir)
		break
	default:
		return fmt.Errorf("no setting %v found", key)
	}
	return nil
}
