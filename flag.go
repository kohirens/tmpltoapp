// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"os"
	"text/template"
)

const cmdConfig = "config"
const cmdManifest = "manifest"

// define All application flags.
func (cfg *Config) defineFlags() {
	// Note: These are defined in alphabetical order.
	flag.StringVar(&cfg.answersPath, "answer-path", "", usageMsgs["answer-path"])
	flag.StringVar(&cfg.branch, "branch", "main", usageMsgs["branch"])
	flag.StringVar(&cfg.defaultVal, "default-val", " ", usageMsgs["default-val"])
	flag.BoolVar(&cfg.help, "help", false, usageMsgs["help"])
	flag.BoolVar(&cfg.help, "h", false, usageMsgs["help"]+" (shorthand)")
	flag.StringVar(&cfg.outPath, "out-path", "", usageMsgs["out-path"])
	flag.StringVar(&cfg.tmplPath, "tmpl-path", "", usageMsgs["tmpl-path"])
	flag.StringVar(&cfg.tmplType, "tmpl-type", "git", usageMsgs["tmpl-type"])
	flag.IntVar(&verbosityLevel, "verbosity", 0, usageMsgs["verbosity"])
	flag.BoolVar(&cfg.version, "version", false, usageMsgs["version"])
	cfg.subCmdConfig.flagSet = flag.NewFlagSet(cmdConfig, flag.ExitOnError)
	cfg.subCmdConfig.flagSet.BoolVar(&cfg.help, "help", false, usageMsgs["help"])
	cfg.subCmdManifest.flagSet = flag.NewFlagSet(cmdManifest, flag.ExitOnError)
	cfg.subCmdManifest.flagSet.BoolVar(&cfg.help, "help", false, usageMsgs["help"])
	cfg.subCmdManifest.flagSet.Usage = func() {
		Usage(cfg)
	}
}

// parse Process and validate all CLI flags.
func (cfg *Config) parseFlags() error {
	// Remember that flag parsing stops just before the first argument that does not have a "-" and is also NOT the
	// value of a flag or comes after the terminator "--".
	// It was planed to allow for flags/arguments in any order, but it may be less confusing to only support flag first
	// and then arguments; it may also require less code to debug and document for not very much gain.
	flag.Parse()

	infof(messages.verboseLevelInfo, verbosityLevel)

	pArgs := flag.Args()
	dbugf(messages.numNonFlagArgs, len(pArgs))
	dbugf(messages.actualArgs, pArgs)
	dbugf(messages.numParsedFlags, flag.NFlag())
	if verbosityLevel == verboseLvlDbug {
		fmt.Println(messages.printAllFlags)
		flag.Visit(func(f *flag.Flag) {
			fmt.Printf(messages.printFlag, f.Name, f.Value, f.DefValue)
		})
	}

	// process sub-commands
	if len(pArgs) > 0 {
		switch pArgs[0] {
		case cmdConfig:
			return cfg.parseConfigCmd(pArgs[1:])
		case cmdManifest:
			return cfg.parseManifestCmd(pArgs[1:])
		}
	}

	// throw an error when a flag comes after any arguments.
	for i := 0; i < len(pArgs); i++ {
		v := pArgs[i]
		if v[0] == '-' {
			return fmt.Errorf(errors.flagOrderErr, v)
		}
	}

	if cfg.help {
		return nil
	}

	if cfg.version {
		logf(messages.currentVersion, cfg.CurrentVersion, cfg.CommitHash)
		os.Exit(0)
	}

	numArgs := len(pArgs)
	if numArgs >= 1 {
		cfg.tmplPath = pArgs[0]
	}
	if numArgs >= 2 {
		cfg.outPath = pArgs[1]
	}
	if numArgs >= 3 {
		cfg.answersPath = pArgs[3]
	}

	if e := cfg.validate(); e != nil {
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
		return fmt.Errorf(messages.outPathExist, cfg.outPath)
	}

	if cfg.answersPath != "" && !stdlib.PathExist(cfg.answersPath) {
		return fmt.Errorf(errors.answerFile404, cfg.answersPath)
	}

	if !regExpTmplType.MatchString(cfg.tmplType) {
		return fmt.Errorf(errors.badTmplType, cfg.tmplType)
	}

	return nil
}

// Usage Print app usage documentation.
func Usage(cfg *Config) error {
	tmpl := template.New("usage")

	switch cfg.subCmd {
	case cmdConfig:
		subCmdConfigUsage(cfg)
		return nil
	case cmdManifest:
		template.Must(tmpl.Parse(usageManifest))
		return UsageTmpl(cfg, tmpl)
	}

	uTmplData := map[string]string{
		"appName": AppName,
	}

	_, err := tmpl.Parse(usageTmpl)
	if err != nil {
		return fmt.Errorf("error parsing the usage template: %v", err.Error())
	}

	if e := tmpl.Execute(os.Stdout, uTmplData); e != nil {
		return fmt.Errorf("error executing the usage template %v", e.Error())
	}

	var mE error
	flag.VisitAll(func(f *flag.Flag) {
		um, ok := usageMsgs[f.Name]
		if ok {
			td := map[string]string{
				"option": f.Name, "info": um, "dv": f.Value.String(),
			}

			if e := tmpl.ExecuteTemplate(os.Stdout, "option", td); e != nil {
				mE = e
				return
			}
		}
	})

	if mE != nil {
		return mE
	}

	return nil
}

func UsageTmpl(cfg *Config, tmpl *template.Template) error {
	uTmplData := map[string]string{
		"appName": AppName,
	}

	if e := tmpl.Execute(os.Stdout, uTmplData); e != nil {
		return fmt.Errorf("error executing the usage template %v", e.Error())
	}

	var mE error
	flag.VisitAll(func(f *flag.Flag) {
		um, ok := usageMsgs[f.Name]
		if ok {
			td := map[string]string{
				"option": f.Name, "info": um, "dv": f.Value.String(),
			}

			if e := tmpl.ExecuteTemplate(os.Stdout, "option", td); e != nil {
				mE = e
				return
			}
		}
	})

	if mE != nil {
		return mE
	}

	return nil
}
