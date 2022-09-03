package main

import (
	"fmt"
	"os"
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
		{"keyDoesNotExist", 1, []string{"config", "set", "key", "value"}, "no setting \"key\" found"},
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

			if !strings.Contains(string(out), test.expected) {
				t.Errorf("error output from program did not contain %q", test.expected)
			}
			if sce != nil {
				fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n", string(out))
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
	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"setCache", 0, []string{"config", "set", "cache", "/tmp"}, ""},
		{"help", 0, []string{"config", "-help"}, ""},
		{"getCache", 0, []string{"config", "get", "cache"}, ""},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

			out, sce := cmd.CombinedOutput()

			// Debug
			fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n\n", string(out))
			fmt.Printf("stderr:\n%v\n", sce.Error())
			fmt.Print("\nEND sub-command\n\n")

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
