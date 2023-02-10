// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"os"
	"text/template"
)

// define All application flags.
func defineFlags(cfg *cli.Config) {
	// Note: These are defined in alphabetical order.
	flag.StringVar(&cfg.AnswersPath, "answer-path", "", usageMsgs["answer-path"])
	flag.StringVar(&cfg.Branch, "branch", "main", usageMsgs["branch"])
	flag.StringVar(&cfg.DefaultVal, "default-val", " ", usageMsgs["default-val"])
	flag.BoolVar(&cfg.Help, "help", false, usageMsgs["help"])
	flag.BoolVar(&cfg.Help, "h", false, usageMsgs["help"]+" (shorthand)")
	flag.StringVar(&cfg.OutPath, "out-path", "", usageMsgs["out-path"])
	flag.StringVar(&cfg.TmplPath, "tmpl-path", "", usageMsgs["tmpl-path"])
	flag.StringVar(&cfg.TmplType, "tmpl-type", "git", usageMsgs["tmpl-type"])
	flag.IntVar(&verbosityLevel, "verbosity", 0, usageMsgs["verbosity"])
	flag.BoolVar(&cfg.Version, "version", false, usageMsgs["version"])
	cfg.SubCmdConfig.FlagSet = flag.NewFlagSet(cli.CmdConfig, flag.ExitOnError)
	cfg.SubCmdConfig.FlagSet.BoolVar(&cfg.Help, "help", false, usageMsgs["help"])
	cfg.SubCmdManifest.FlagSet = flag.NewFlagSet(cli.CmdManifest, flag.ExitOnError)
	cfg.SubCmdManifest.FlagSet.BoolVar(&cfg.Help, "help", false, usageMsgs["help"])
	cfg.SubCmdManifest.FlagSet.Usage = func() {
		Usage(cfg)
	}
}

// Parse Process and validate all CLI flags.
func parseFlags(cfg *cli.Config) error {
	// Remember that flag parsing stops just before the first argument that does not have a "-" and is also NOT the
	// value of a flag or comes after the terminator "--".
	// It was planed to allow for flags/arguments in any order, but it may be less confusing to only support flag first
	// and then arguments; it may also require less code to debug and document for not very much gain.
	flag.Parse()

	infof(cli.Messages.VerboseLevelInfo, verbosityLevel)

	pArgs := flag.Args()
	dbugf(cli.Messages.NumNonFlagArgs, len(pArgs))
	dbugf(cli.Messages.ActualArgs, pArgs)
	dbugf(cli.Messages.NumParsedFlags, flag.NFlag())
	if verbosityLevel == verboseLvlDbug {
		fmt.Println(cli.Messages.PrintAllFlags)
		flag.Visit(func(f *flag.Flag) {
			fmt.Printf(cli.Messages.PrintFlag, f.Name, f.Value, f.DefValue)
		})
	}

	// process sub-commands
	if len(pArgs) > 0 {
		switch pArgs[0] {
		case cli.CmdConfig:
			return parseSubCmd(cfg, pArgs[1:])
		case cli.CmdManifest:
			return parseManifestCmd(cfg, pArgs[1:])
		}
	}

	// throw an error when a flag comes after any arguments.
	for i := 0; i < len(pArgs); i++ {
		v := pArgs[i]
		if v[0] == '-' {
			return fmt.Errorf(cli.Errors.FlagOrderErr, v)
		}
	}

	if cfg.Help {
		return nil
	}

	if cfg.Version {
		logf(cli.Messages.CurrentVersion, cfg.CurrentVersion, cfg.CommitHash)
		os.Exit(0)
	}

	numArgs := len(pArgs)
	if numArgs >= 1 {
		cfg.TmplPath = pArgs[0]
	}
	if numArgs >= 2 {
		cfg.OutPath = pArgs[1]
	}
	if numArgs >= 3 {
		cfg.AnswersPath = pArgs[3]
	}

	if e := cfg.Validate(); e != nil {
		return e
	}

	return nil
}

// parseSubCmd Parse the sub-command flags/options/args.
func parseSubCmd(cfg *cli.Config, osArgs []string) error {
	if e := cfg.SubCmdConfig.FlagSet.Parse(osArgs); e != nil {
		return fmt.Errorf(cli.Errors.ParsingConfigArgs, e.Error())
	}

	cfg.SubCmd = cli.CmdConfig

	if cfg.Help {
		return nil
	}

	if len(osArgs) < 2 {
		subCmdConfigUsage(cfg)
		return fmt.Errorf(cli.Errors.InvalidNoArgs)
	}

	cfg.SubCmdConfig.Method = osArgs[0]
	cfg.SubCmdConfig.Key = osArgs[1]

	if len(osArgs) > 2 {
		cfg.SubCmdConfig.Value = osArgs[2]
	}

	log.Dbugf("cfg.SubCmdConfig.method = %v\n", cfg.SubCmdConfig.Method)
	log.Dbugf("cfg.SubCmdConfig.key = %v\n", cfg.SubCmdConfig.Key)
	log.Dbugf("cfg.SubCmdConfig.value = %v\n", cfg.SubCmdConfig.Value)

	return nil
}

// ParseManifestCmd Parse the config sub-command flags/options/args but do not execute the command itself
func parseManifestCmd(cfg *cli.Config, osArgs []string) error {
	cfg.SubCmd = cli.CmdManifest
	if e := cfg.SubCmdManifest.FlagSet.Parse(osArgs); e != nil {
		fmt.Printf("here")
		os.Exit(0)
		return fmt.Errorf(cli.Errors.ParsingConfigArgs, e.Error())
	}

	if cfg.Help {
		return Usage(cfg)
	}

	if len(osArgs) < 2 {
		Usage(cfg)
		return fmt.Errorf(cli.Errors.InvalidNoSubCmdArgs, cli.CmdManifest, 1)
	}

	cfg.SubCmdManifest.Path = osArgs[0]

	log.Dbugf("cfg.SubCmdManifest.path = %v\n", cfg.SubCmdManifest.Path)

	return nil
}

// subCmdConfigUsage print config command usage
func subCmdConfigUsage(cfg *cli.Config) {
	fmt.Printf("usage: config set|get <args>\n\n")
	fmt.Println("examples:")
	fmt.Printf("\tconfig set \"CacheDir\" \"./path/to/a/directory\"\n")
	fmt.Printf("\tconfig get \"CacheDir\"\n\n")
	fmt.Printf("Settings: \n")
	fmt.Printf("\tCacheDir - Path to store template downloaded\n")
	fmt.Printf("\tExcludeFileExtensions - Files ending with these extensions will be excluded from parsing and copied as-is\n\n")
	fmt.Printf("Options: \n")
	// print options usage
	cfg.SubCmdConfig.FlagSet.VisitAll(func(f *flag.Flag) {
		um, ok := usageMsgs[f.Name]
		if ok {
			fmt.Printf("  -%-11s %v\n\n", f.Name, um)
			f.Value.String()
		}
	})
}

// Usage Print app usage documentation.
func Usage(cfg *cli.Config) error {
	tmpl := template.New("usage")

	switch cfg.SubCmd {
	case cli.CmdConfig:
		subCmdConfigUsage(cfg)
		return nil
	case cli.CmdManifest:
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

func UsageTmpl(cfg *cli.Config, tmpl *template.Template) error {
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
