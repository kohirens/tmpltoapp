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

//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName appConfig

import (
	"flag"
	"fmt"
	"os"
)

// define All application flags.
func (cf *Config) define() {
	flag.StringVar(&cf.answersPath, "answers", "", usageMsgs["answers"])
	flag.StringVar(&cf.answersPath, "a", "", usageMsgs["answers"]+" (shorthand)")
	flag.StringVar(&cf.appPath, "appPath", "", usageMsgs["appPath"])
	flag.StringVar(&cf.appPath, "p", "", usageMsgs["appPath"]+" (shorthand)")
	flag.StringVar(&cf.branch, "branch", "", usageMsgs["branch"])
	flag.BoolVar(&cf.help, "help", false, usageMsgs["help"])
	flag.BoolVar(&cf.help, "h", false, usageMsgs["help"]+" (shorthand)")
	flag.StringVar(&cf.tplPath, "tmplPath", "", usageMsgs["tmplPath"])
	flag.StringVar(&cf.tplPath, "t", "", usageMsgs["tmplPath"]+" (shorthand)")
	flag.StringVar(&cf.tmplType, "tmplType", "zip", usageMsgs["tmplType"])
	flag.IntVar(&verbosityLevel, "verbosity", 0, usageMsgs["verbosity"])
	flag.BoolVar(&cf.version, "version", false, usageMsgs["version"])
	flag.BoolVar(&cf.version, "v", false, usageMsgs["version"]+" (shorthand)")
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
		config.tplPath = pArgs[0]
	}
	if numArgs > 2 {
		config.appPath = pArgs[1]
	}
	if numArgs > 3 {
		config.answersPath = pArgs[3]
	}

	err2 := config.validate()
	if err2 != nil {
		return err2
	}

	return nil
}
