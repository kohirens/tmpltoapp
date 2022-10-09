package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// userOptions Options the user can set
type userOptions struct {
	ExcludeFileExtensions *[]string
	CacheDir              string
}

// parseConfigCmd parse the config sub-command flags/options/args but do not execute the command itself
func (cfg *Config) parseConfigCmd(osArgs []string) error {
	if e := cfg.subCmdConfig.flagSet.Parse(osArgs); e != nil {
		return fmt.Errorf(errors.parsingConfigArgs, e.Error())
	}

	cfg.subCmd = cmdConfig

	if cfg.help {
		return nil
	}

	if len(osArgs) < 2 {
		subCmdConfigUsage(cfg)
		return fmt.Errorf(errors.invalidNoArgs)
	}

	cfg.subCmdConfig.method = osArgs[0]
	cfg.subCmdConfig.key = osArgs[1]

	if len(osArgs) > 2 {
		cfg.subCmdConfig.value = osArgs[2]
	}

	dbugf("cfg.subCmdConfig.method = %v\n", cfg.subCmdConfig.method)
	dbugf("cfg.subCmdConfig.key = %v\n", cfg.subCmdConfig.key)
	dbugf("cfg.subCmdConfig.value = %v\n", cfg.subCmdConfig.value)

	return nil
}

// parseManifestCmd parse the config sub-command flags/options/args but do not execute the command itself
func (cfg *Config) parseManifestCmd(osArgs []string) error {
	cfg.subCmd = cmdManifest
	if e := cfg.subCmdManifest.flagSet.Parse(osArgs); e != nil {
		fmt.Printf("here")
		os.Exit(0)
		return fmt.Errorf(errors.parsingConfigArgs, e.Error())
	}

	if cfg.help {
		return Usage(cfg)
	}

	if len(osArgs) < 2 {
		Usage(cfg)
		return fmt.Errorf(errors.invalidNoSubCmdArgs, cmdManifest, 1)
	}

	cfg.subCmdManifest.path = osArgs[0]

	dbugf("cfg.subCmdManifest.path = %v\n", cfg.subCmdManifest.path)

	return nil
}

func updateUserSettings(cfg *Config, mode os.FileMode) error {
	switch cfg.subCmdConfig.method {
	case "set":
		if e := cfg.set(cfg.subCmdConfig.key, cfg.subCmdConfig.value); e != nil {
			return e
		}
		break
	case "get":
		v, e := cfg.get(cfg.subCmdConfig.key)
		if e != nil {
			return e
		}
		fmt.Printf("%v", v)
	}

	return cfg.saveUserSettings(mode)
}

// subCmdConfigUsage print config command usage
func subCmdConfigUsage(cfg *Config) {
	fmt.Printf("usage: config set|get <args>\n\n")
	fmt.Println("examples:")
	fmt.Printf("\tconfig set \"CacheDir\" \"./path/to/a/directory\"\n")
	fmt.Printf("\tconfig get \"CacheDir\"\n\n")
	fmt.Printf("Settings: \n")
	fmt.Printf("\tCacheDir - Path to store template downloaded\n")
	fmt.Printf("\tExcludeFileExtensions - Files ending with these extensions will be excluded from parsing and copied as-is\n\n")
	fmt.Printf("Options: \n")
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
func (cfg *Config) set(key, val string) error {
	switch key {
	case "CacheDir":
		cfg.usrOpts.CacheDir = val
		break
	case "ExcludeFileExtensions":
		tmp := strings.Split(val, ",")
		cfg.usrOpts.ExcludeFileExtensions = &tmp
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
		val = cfg.usrOpts.CacheDir
		break
	case "ExcludeFileExtensions":
		v2 := fmt.Sprintf("%v", val)
		ok, _ := regexp.Match("^[a-zA-Z0-9-.]+(?:,[a-zA-Z0-9-.]+)*", []byte(v2))
		if !ok {
			return nil, fmt.Errorf(errors.badExcludeFileExt, val)
		}
		val = strings.Join(*cfg.usrOpts.ExcludeFileExtensions, ",")
		break
	default:
		return "", fmt.Errorf("no setting %v found", key)
	}

	return val, nil
}
