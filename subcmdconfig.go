package main

import (
	"flag"
	"fmt"
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

// parseConfigCmd parse the config sub-command flags/options/args but do not execute the command itself
func (cfg *Config) parseConfigCmd(osArgs []string) error {
	if e := cfg.subCmdConfig.flagSet.Parse(osArgs); e != nil {
		return fmt.Errorf("error pasing sub command config flags: %v", e.Error())
	}

	if cfg.help {
		Usage(cfg)
		return nil
	}

	if len(osArgs) < 2 {
		subCmdConfigUsage(cfg)
		return fmt.Errorf("invalid number of arguments passed to config sub-command, please try config -h for usage")
	}

	cfg.subCmd = cmdConfig

	cfg.subCmdConfig.method = osArgs[0]
	cfg.subCmdConfig.key = osArgs[1]
	cfg.subCmdConfig.value = osArgs[2]

	return nil
}

func updateUserSettings(cfg *Config) error {
	fmt.Printf("\ncfg.subCmdConfig.method = %v", cfg.subCmdConfig.method)
	switch cfg.subCmdConfig.method {
	case "set":
		if e := usrOpts.set(cfg.subCmdConfig.key, cfg.subCmdConfig.value); e != nil {
			return e
		}
		break
	case "get":
		fmt.Printf("%v", usrOpts.get(cfg.subCmdConfig.key))
	}

	return saveConfigFile(cfg.path, usrOpts)
}

// subCmdConfigUsage print config command usage
func subCmdConfigUsage(cfg *Config) {
	fmt.Println("usage: config set|get <args>")
	fmt.Println("examples:")
	fmt.Printf("\tconfig set \"cacheDir\" \"./path/to/a/directory\"\n")
	fmt.Printf("\tconfig get \"cacheDir\"\n")
	fmt.Printf("Options: \n\n")
	// print options usage
	cfg.subCmdConfig.flagSet.VisitAll(func(f *flag.Flag) {
		um, ok := usageMsgs[f.Name]
		if ok {
			fmt.Printf("  -%-11s %v\n\n", f.Name, um)
			f.Value.String()
		}
	})
}

// set the value of a user setting
func (uo *userOptions) set(key, val string) error {
	switch key {
	case "cacheDir":
		usrOpts.cacheDir = val
		break
	default:
		return fmt.Errorf("no %q setting found", key)
	}
	return nil
}

// get the value of a user setting.
func (uo *userOptions) get(key string) error {
	switch key {
	case "cacheDir":
		fmt.Printf("%v", usrOpts.cacheDir)
		break
	default:
		return fmt.Errorf("no setting %v found", key)
	}
	return nil
}
