package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetTmplLocation(runner *testing.T) {
	fixtures := []struct {
		name, want string
		cfg        *Config
	}{
		{"relative", "local", &Config{TmplPath: "./"}},
		{"relative2", "local", &Config{TmplPath: "."}},
		{"relativeUp", "local", &Config{TmplPath: ".."}},
		{"absolute", "local", &Config{TmplPath: "/home/myuser"}},
		{"windows", "local", &Config{TmplPath: "C:\\Temp"}},
		{"http", "remote", &Config{TmplPath: "http://example.com/repo1"}},
		{"https", "remote", &Config{TmplPath: "https://example.com/repo1"}},
		{"git", "remote", &Config{TmplPath: "git://example.com/repo1"}},
		{"file", "remote", &Config{TmplPath: "file://example.com/repo1"}},
		{"hiddenRelative", "local", &Config{TmplPath: ".m/example.com/repo1"}},
		{"tildeRelative", "local", &Config{TmplPath: "~/repo1.git"}},
	}

	for _, tc := range fixtures {
		runner.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.getTmplLocation()

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetSettings(t *testing.T) {

	t.Run("configNotFound", func(t *testing.T) {
		cfgFixture := &Config{}
		gotErr := cfgFixture.loadUserSettings("does-not-exist")

		if !strings.Contains(gotErr.Error(), "could not open") {
			t.Errorf("got %q, want %q", gotErr, "could not open")
		}
	})

	t.Run("canReadConfig", func(t *testing.T) {
		cfgFixture := &Config{}
		err := cfgFixture.loadUserSettings(test.FixturesDir + PS + "config-01.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}
	})
}

func TestInitConfigFile(t *testing.T) {
	var testCases = []struct {
		name string
		want error
		cfg  *Config
	}{
		{"NotExist", nil, &Config{Path: test.TmpDir + PS + "config-fix-01.json"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.cfg.initFile()

			if tc.want != gotErr {
				t.Errorf("got %t; want %t", gotErr, tc.want)
			}

			_, err := os.Stat(tc.cfg.Path)

			if err != nil {
				t.Errorf("got %v; want %v", gotErr, tc.want)
			}
		})
	}
}

func TestLoadAnswers(tester *testing.T) {
	var fixtures = []struct {
		name, file, want string
	}{
		{"goodJson", test.FixturesDir + PS + "answers-01.json", "value1"},
		{"badJson", test.FixturesDir + PS + "answers-02.json", ""},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(t *testing.T) {
		got, err := LoadAnswers(fxtr.file)

		if err == nil && got.Placeholders["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got.Placeholders["var1"], fxtr.want)
		}
	})

	fxtr = fixtures[1]
	tester.Run(fxtr.name, func(t *testing.T) {
		_, err := LoadAnswers(fxtr.file)

		if err == nil {
			t.Error("did not get an error")
		}
	})
}

func TestSubCmdConfigExitCode(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		args     []string
		expected string
	}{
		{"noArgs", 1, []string{CmdConfig}, "usage: config"},
		{"keyDoesNotExist", 1, []string{CmdConfig, "set", "key", "value"}, "no \"key\" setting found"},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := test.GetTestBinCmd(tc.args)

			out, sce := cmd.CombinedOutput()

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if testing.Verbose() && sce != nil {
				fmt.Print("\nBEGIN sub-command\n")
				fmt.Printf("stdout:\n%s\n", out)
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			// using bytes.NewBuffer(out).String() allows the byte array to be
			// converted as-is and keeps chars such as newlines and quotes from
			// being converted to '\n' or '\"', etc.
			if !strings.Contains(bytes.NewBuffer(out).String(), tc.expected) {
				t.Errorf("std error from program did not contain %q", tc.expected)
			}
		})
	}
}

//func xTestSubCmdConfigBadExit(tester *testing.T) {
//	oldArgs := os.Args
//	var tests = []struct {
//		name     string
//		wantCode int
//		args     []string
//		expected string
//	}{
//		{"noArgs", 1, []string{oldArgs[0], CmdConfig}, "usage: config"},
//		{"keyDoesNotExist", 1, []string{oldArgs[0], CmdConfig, "set", "key", "value"}, "no config setting \"keyDoesNotExist\" found"},
//	}
//
//	for _, tc := range tests {
//		tester.Run(tc.name, func(t *testing.T) {
//			defer func() {
//				os.Args = oldArgs
//			}()
//
//			// replace the test flaga and arts with the test fixture flags and args.
//			os.Args = tc.args
//
//			testConf := &Config{}
//			defineFlags(testConf)
//			e := parseFlags(testConf)
//			if e != nil {
//				t.Error(e)
//			}
//		})
//	}
//}

func TestSubCmdConfigSuccess(tester *testing.T) {
	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")
		_ = os.Setenv("LOCALAPPDATA", test.TmpDir)
		defer func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}()
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", test.TmpDir)
		defer func() {
			_ = os.Setenv("HOME", oldHome)
		}()
	}

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"help", 0, []string{CmdConfig, "-help"}, "usage: config"},
		{"getCache", 0, []string{CmdConfig, "get", "CacheDir"}, "tmp"},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := test.GetTestBinCmd(tc.args)

			out, sce := cmd.CombinedOutput()

			// Debug
			if sce != nil {
				fmt.Print("\nBEGIN sub-command\n")
				fmt.Printf("stdout:\n%s\n", out)
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			if !strings.Contains(bytes.NewBuffer(out).String(), tc.contains) {
				t.Errorf("did not contain %q", tc.contains)
			}
		})
	}
}

func TestSetUserOptions(tester *testing.T) {
	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")
		_ = os.Setenv("LOCALAPPDATA", test.TmpDir)
		defer func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}()
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", test.TmpDir)
		defer func() {
			_ = os.Setenv("HOME", oldHome)
		}()
	}

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		want     string
	}{
		{"setCache", 0, []string{CmdConfig, "set", "CacheDir", "setCache"}, "setCache"},
		{"setExcludeFileExtensions", 0, []string{CmdConfig, "set", "ExcludeFileExtensions", "md,txt"}, `"ExcludeFileExtensions":["md","txt"]`},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := test.GetTestBinCmd(tc.args)

			gotOut, sce := cmd.CombinedOutput()

			// Debug
			if sce != nil {
				fmt.Print("\nBEGIN sub-command\n")
				fmt.Printf("stdout:\n%s\n", gotOut)
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}

			gotExit := cmd.ProcessState.ExitCode()

			if gotExit != tc.wantCode {
				t.Errorf("got %q, want %q", gotExit, tc.wantCode)
			}

			file := test.TmpDir + PS + "tmpltoapp" + PS + "config.json"

			gotCfg := &Config{}
			_ = gotCfg.loadUserSettings(file)
			ec, _ := json.Marshal(gotCfg.UsrOpts)

			if !strings.Contains(string(ec), tc.want) {
				t.Errorf("the config %s did not contain %v set to %v", ec, tc.args[2], tc.want)
			}
		})
	}
}

func xTestLoadUserSettings(tester *testing.T) {
	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")
		_ = os.Setenv("LOCALAPPDATA", test.TmpDir)
		defer func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}()
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", test.TmpDir)
		defer func() {
			_ = os.Setenv("HOME", oldHome)
		}()
	}

	var tests = []struct {
		name     string
		filename string
		want     *Config
	}{
		{
			"goodFile",
			test.FixturesDir + PS + "good-config-01.json",
			&Config{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
		{
			"badFile",
			test.FixturesDir + PS + "bad-config-01.json",
			&Config{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			gotCfg := &Config{}
			err := gotCfg.loadUserSettings(tc.filename)

			if err != nil { // test bad values

			} else { // test good values

			}
		})
	}
}
