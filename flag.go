// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
	Package main.flags implements command-line flag parsing.

	Usage

	```shell
	tmpltoapp [options] "<dir/url>" "<outputDir>" "<answers>"
	```
*/
package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"os"
)

// define All application flags.
func (cfg *Config) define() {
	flag.StringVar(&cfg.answersPath, "answers", "", usageMsgs["answers"])
	flag.StringVar(&cfg.answersPath, "a", "", usageMsgs["answers"]+" (shorthand)")
	flag.StringVar(&cfg.appPath, "appPath", "", usageMsgs["appPath"])
	flag.StringVar(&cfg.appPath, "p", "", usageMsgs["appPath"]+" (shorthand)")
	flag.StringVar(&cfg.branch, "branch", "main", usageMsgs["branch"])
	flag.BoolVar(&cfg.help, "help", false, usageMsgs["help"])
	flag.BoolVar(&cfg.help, "h", false, usageMsgs["help"]+" (shorthand)")
	flag.StringVar(&cfg.tmplPath, "tmplPath", "", usageMsgs["tmplPath"])
	flag.StringVar(&cfg.tmplPath, "t", "", usageMsgs["tmplPath"]+" (shorthand)")
	flag.StringVar(&cfg.tmplType, "tmplType", "zip", usageMsgs["tmplType"])
	flag.IntVar(&verbosityLevel, "verbosity", 0, usageMsgs["verbosity"])
	flag.BoolVar(&cfg.version, "version", false, usageMsgs["version"])
	flag.BoolVar(&cfg.version, "v", false, usageMsgs["version"]+" (shorthand)")
}

// flagMain Process and validate all CLI flags.
func flagMain(config *Config) error {
	flag.Parse()

	if config.help {
		flag.PrintDefaults()
		fmt.Printf(usageMsgs["subCommands"])
		os.Exit(0)
	}

	if config.version {
		fmt.Printf("%v, %v\n", config.CurrentVersion, config.CommitHash)
		os.Exit(0)
	}

	pArgs := flag.Args()[0:]
	errf("number of arguments passed in: %d", len(pArgs))
	errf("arguments passed in: %v", pArgs)

	numArgs := len(pArgs)
	if numArgs >= 1 {
		config.tmplPath = pArgs[0]
	}
	if numArgs >= 2 {
		config.appPath = pArgs[1]
	}
	if numArgs >= 3 {
		config.answersPath = pArgs[3]
	}

	err2 := config.validate()
	if err2 != nil {
		return err2
	}

	return nil
}

// validate parses command line flags into program options.
func (cfg *Config) validate() error {
	if cfg.tmplPath == "" {
		return fmt.Errorf(errors.tmplPath)
	}

	if cfg.appPath == "" {
		return fmt.Errorf(errors.localOutPath)
	}

	if stdlib.DirExist(cfg.appPath) {
		return fmt.Errorf("appPath already exits %q", cfg.appPath)
	}

	if cfg.answersPath != "" && !stdlib.PathExist(cfg.answersPath) {
		return fmt.Errorf(errors.answerFile404, cfg.answersPath)
	}

	if !regExpTmplType.MatchString(cfg.tmplType) {
		return fmt.Errorf(errors.badTmplType, cfg.tmplType)
	}

	return nil
}
