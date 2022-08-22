// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"os"
)

// define All application flags.
func (cfg *Config) define() {
	// TODO: add flag to set a default values for -skip-un-answered and use a -default-value questions.
	flag.StringVar(&cfg.answersPath, "answers", "", usageMsgs["answers"])
	flag.StringVar(&cfg.outPath, "out-path", "", usageMsgs["out-path"])
	flag.StringVar(&cfg.branch, "branch", "main", usageMsgs["branch"])
	flag.BoolVar(&cfg.help, "help", false, usageMsgs["help"])
	flag.BoolVar(&cfg.help, "h", false, usageMsgs["help"]+" (shorthand)")
	flag.StringVar(&cfg.tmplPath, "tmpl-path", "", usageMsgs["tmplPath"])
	flag.StringVar(&cfg.tmplType, "tmpl-type", "git", usageMsgs["tmplType"])
	flag.IntVar(&verbosityLevel, "verbosity", 0, usageMsgs["verbosity"])
	flag.BoolVar(&cfg.version, "version", false, usageMsgs["version"])
}

// flagMain Process and validate all CLI flags.
func flagMain(config *Config) error {
	// TODO: allow flags in any order
	// Remember that Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--".
	flag.Parse()

	// TODO: Show order of all input here
	pArgs := flag.Args()
	dbugf("number of arguments passed in: %d", len(pArgs))
	dbugf("arguments passed in: %v", pArgs)

	for i := 0; i < len(pArgs); i++ {
		v := pArgs[i]
		if v[0] == '-' {
			return fmt.Errorf(errors.flagOrderErr, v)
		}
	}

	if config.help {
		// TODO: Replace with custom printDefaults function
		flag.PrintDefaults()
		fmt.Printf(usageMsgs["subCommands"])
		os.Exit(0)
	}

	if config.version {
		fmt.Printf("%v, %v\n", config.CurrentVersion, config.CommitHash)
		os.Exit(0)
	}

	numArgs := len(pArgs)
	if numArgs >= 1 {
		config.tmplPath = pArgs[0]
	}
	if numArgs >= 2 {
		config.outPath = pArgs[1]
	}
	if numArgs >= 3 {
		config.answersPath = pArgs[3]
	}

	if e := config.validate(); e != nil {
		return e
	}

	return nil
}

// validate parses command line flags into program options.
func (cfg *Config) validate() error {
	if cfg.tmplPath == "" {
		return fmt.Errorf(errors.tmplPath)
	}

	if cfg.outPath == "" {
		return fmt.Errorf(errors.localOutPath)
	}

	if stdlib.DirExist(cfg.outPath) {
		return fmt.Errorf("outPath already exits %q", cfg.outPath)
	}

	if cfg.answersPath != "" && !stdlib.PathExist(cfg.answersPath) {
		return fmt.Errorf(errors.answerFile404, cfg.answersPath)
	}

	if !regExpTmplType.MatchString(cfg.tmplType) {
		return fmt.Errorf(errors.badTmplType, cfg.tmplType)
	}

	return nil
}
