// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
	Package main.flags implements command-line flag parsing.

	Usage

	```shell
	tmpltoapp [options] "<dir/url>" "<outputDir>" "<answers>"
	```

	## Description

	Use a template to initialize a new project. A template can be a local directory
	or a zip from a URL. Zip files will be downloaded and extracted to a local
	directory.

	## Options

	**--tmplPath**, **-t** URL to a zip or a local path to a directory

	**--answers**, **-a** Path to an answer file.

	**--help**, **-h** Output this documentation.

	**--verbosity** Control the level of information/feedback the program will
	output to the user.

	**--version**, **-v** Output version information.

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
	cliFlag{"tmplPath", "t", "URL to a zip or a local path to a directory.", "string"},
	cliFlag{"appPath", "p", "Path to output the new project.", "string"},
	cliFlag{"answers", "a", "Path to an answer file.", "string"},
	cliFlag{"verbosity", "", "extra detail processing infof.", "int"},
	cliFlag{"help", "h", "print usage information.", "bool"},
	cliFlag{"version", "v", "print build version information and exit 0.", "bool"},
	cliFlag{"tmplType", "zip", "Can be of git|zip|local.", "string"},
	cliFlag{"branch", "", "branch to clone.", "string"},
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

// flagMain This is code refactored from main to help keep main readable.
func flagMain() error {
	err1 := flagStore.Flags.Parse(os.Args[1:])
	if err1 != nil {
		return err1
	}

	help, _ := flagStore.GetBool("help")
	if help {
		flagStore.Flags.SetOutput(os.Stdout)
		flagStore.Flags.Usage()
		os.Exit(0)
	}

	version, _ := flagStore.GetBool("version")
	if version {
		flagStore.Flags.SetOutput(os.Stdout)
		fmt.Printf("%v, %v\n", appConfig.CurrentVersion, appConfig.CommitHash)
		os.Exit(0)
	}

	err2 := extractParsedFlags(flagStore, os.Args, appConfig)
	if err2 != nil {
		return err2
	}

	return nil
}
