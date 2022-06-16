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

type cliFlag struct {
	name        string
	short       string
	description string
	valueType   string
}

type cliFlags []cliFlag

var appFlags = cliFlags{
	cliFlag{"tmplPath", "t", usageMsgs["tmplPath"], "string"},
	cliFlag{"appPath", "p", usageMsgs["appPath"], "string"},
	cliFlag{"answers", "a", usageMsgs["answers"], "string"},
	cliFlag{"verbosity", "", usageMsgs["verbosity"], "int"},
	cliFlag{"help", "h", usageMsgs["help"], "bool"},
	cliFlag{"version", "v", usageMsgs["version"], "bool"},
	cliFlag{"tmplType", "zip", usageMsgs["tmplType"], "string"},
	cliFlag{"branch", "", usageMsgs["branch"], "string"},
}

type flagStorage struct {
	Flags *flag.FlagSet
	ints  map[string]*int
	bools map[string]*bool
	strs  map[string]*string
}

// GetInt Get a flag parsed as an integer.
func (fs *flagStorage) GetInt(key string) (val int, err error) {
	v, ok := fs.ints[key]

	if !ok {
		err = fmt.Errorf("there is no defined int flag %q", key)
	}

	val = *v

	return
}

// GetBool Get a boolean flag.
func (fs *flagStorage) GetBool(key string) (val bool, err error) {
	v, ok := fs.bools[key]

	if !ok {
		err = fmt.Errorf("there is no defined int flag %q", key)
		return
	}

	val = *v

	return
}

// GetString Get a flag parsed as a string.
func (fs *flagStorage) GetString(key string) (val string, err error) {
	v, ok := fs.strs[key]
	if !ok {
		err = fmt.Errorf("there is no defined int flag %q", key)
		return
	}

	val = *v

	return
}

// Process any program flags fed into the program and return an unparsed flag-set.
func defineFlags(programName string, handling flag.ErrorHandling) (flagStore *flagStorage, err error) {
	flags := flag.NewFlagSet(programName, handling)
	ints := map[string]*int{}
	bools := map[string]*bool{}
	strs := map[string]*string{}

	for _, f := range appFlags {
		switch f.valueType {
		default:
			strs[f.name] = flags.String(f.name, "", f.description)
			if len(f.short) == 1 {
				flags.StringVar(strs[f.name], f.short, *strs[f.name], f.description)
			}
		case "bool":
			bools[f.name] = flags.Bool(f.name, false, f.description)
			if len(f.short) == 1 {
				flags.BoolVar(bools[f.name], f.short, *bools[f.name], f.description)
			}
		case "int":
			ints[f.name] = flags.Int(f.name, 0, f.description)
			if len(f.short) == 1 {
				flags.IntVar(ints[f.name], f.short, *ints[f.name], f.description)
			}
		}
	}

	flagStore = &flagStorage{
		Flags: flags,
		ints:  ints,
		bools: bools,
		strs:  strs,
	}

	return
}

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

	err1 := flagStore.Flags.Parse(os.Args[1:])
	if err1 != nil {
		return err1
	}

	if config.help {
		flag.PrintDefaults()
		fmt.Printf(usageMsgs["subCommands"])
		os.Exit(0)
	}

	if config.version {
		fmt.Printf("%v, %v\n", config.CurrentVersion, config.CommitHash)
		os.Exit(0)
	}

	err2 := extractParsedFlags(flag.Args(), config)
	if err2 != nil {
		return err2
	}

	return nil
}
