package main

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestSubCmdConfigExitCode(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		args     []string
		expected string
	}{
		{"noArgs", 1, []string{"config"}, "usage: config"},
		{"keyDoesNotExist", 1, []string{"config", "set", "key", "value"}, "no \"key\" setting found"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

			out, sce := cmd.CombinedOutput()

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			// using bytes.NewBuffer(out).String() allows the byte array to be
			// converted as-is and keeps chars such as newlines and quotes from
			// being converted to '\n' or '\"', etc.
			if !strings.Contains(bytes.NewBuffer(out).String(), test.expected) {
				t.Errorf("std error from program did not contain %q", test.expected)
			}

			if testing.Verbose() && sce != nil {
				fmt.Print("\nBEGIN sub-command\n")
				fmt.Printf("stdout:\n%s\n", out)
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}
		})
	}
}

func xTestSubCmdConfigBadExit(tester *testing.T) {
	oldArgs := os.Args
	var tests = []struct {
		name     string
		wantCode int
		args     []string
		expected string
	}{
		{"noArgs", 1, []string{oldArgs[0], "config"}, "usage: config"},
		{"keyDoesNotExist", 1, []string{oldArgs[0], "config", "set", "key", "value"}, "no config setting \"keyDoesNotExist\" found"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			defer func() {
				os.Args = oldArgs
			}()

			// replace the test flaga and arts with the test fixture flags and args.
			os.Args = test.args

			testConf := &Config{}
			testConf.define()
			e := flagMain(testConf)
			if e != nil {
				t.Error(e)
			}
		})
	}
}
func TestSubCmdConfigSuccess(tester *testing.T) {
	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")
		_ = os.Setenv("LOCALAPPDATA", testTmp)
		defer func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}()
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", testTmp)
		defer func() {
			_ = os.Setenv("HOME", oldHome)
		}()
	}

	dir, _ := stdlib.AppDataDir()
	// Debug
	fmt.Printf("app data dir: %q", dir)

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"setCache", 0, []string{"config", "set", "cacheDir", "/tmp"}, ""},
		//{"help", 0, []string{"config", "-help"}, ""},
		//{"getCache", 0, []string{"config", "get", "cacheDir"}, "/tmp"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

			out, sce := cmd.CombinedOutput()

			// Debug
			if sce != nil {
				fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n", string(out))
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			if !strings.Contains(string(out), test.contains) {
				t.Errorf("did not contain %q", test.contains)
			}
		})
	}
}
