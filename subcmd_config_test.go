package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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
		{"noArgs", 1, []string{cmdConfig}, "usage: config"},
		{"keyDoesNotExist", 1, []string{cmdConfig, "set", "key", "value"}, "no \"key\" setting found"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

			out, sce := cmd.CombinedOutput()

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if testing.Verbose() && sce != nil {
				fmt.Print("\nBEGIN sub-command\n")
				fmt.Printf("stdout:\n%s\n", out)
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			// using bytes.NewBuffer(out).String() allows the byte array to be
			// converted as-is and keeps chars such as newlines and quotes from
			// being converted to '\n' or '\"', etc.
			if !strings.Contains(bytes.NewBuffer(out).String(), test.expected) {
				t.Errorf("std error from program did not contain %q", test.expected)
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
		{"noArgs", 1, []string{oldArgs[0], cmdConfig}, "usage: config"},
		{"keyDoesNotExist", 1, []string{oldArgs[0], cmdConfig, "set", "key", "value"}, "no config setting \"keyDoesNotExist\" found"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			defer func() {
				os.Args = oldArgs
			}()

			// replace the test flaga and arts with the test fixture flags and args.
			os.Args = test.args

			testConf := &Config{}
			testConf.defineFlags()
			e := testConf.parseFlags()
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

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"help", 0, []string{cmdConfig, "-help"}, "usage: config"},
		{"getCache", 0, []string{cmdConfig, "get", "CacheDir"}, "tmp"},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

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

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			if !strings.Contains(bytes.NewBuffer(out).String(), test.contains) {
				//fmt.Printf("stdout = %s\n", out)
				t.Errorf("did not contain %q", test.contains)
			}
		})
	}
}

func TestSetUserOptions(tester *testing.T) {
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

	var tests = []struct {
		name     string
		wantCode int
		args     []string
	}{
		{"setCache", 0, []string{cmdConfig, "set", "CacheDir", "setCache"}},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {

			cmd := getTestBinCmd(test.args)

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

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			file := testTmp + PS + "tmpltoapp" + PS + "config.json"
			content, _ := ioutil.ReadFile(file)
			usrOpts := &userOptions{}
			_ = json.Unmarshal(content, usrOpts)

			rUsrOpts := reflect.ValueOf(usrOpts)
			f := reflect.Indirect(rUsrOpts).FieldByName(test.args[2])

			if f.String() != test.args[3] {
				t.Errorf("did not contai")
			}
		})
	}
}
