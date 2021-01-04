package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const testTmp = "go_gitter_test_tmp"

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Mkdir(testTmp, 0774) // set up a temporary dir for generate files

	// Create whatever test files are needed.

	// Run all tests and clean up
	exitcode := m.Run()
	os.RemoveAll(testTmp) // remove the directory and its contents.
	os.Exit(exitcode)
}

func xTestMainOutput(t *testing.T) {
	var tests = []struct {
		tpl, appname, ans, want string
	}{
		{"", "", "", "help"},
		// {"non-existent", "", "", "The new application path is required."},
		// {"non-existent", "", "", "path/URL to template does not exist"},
		// {"./fixtures/tpl-1", "", "", "path/URL to template does not exist"},
	}

	for _, tt := range tests {

		testName := fmt.Sprintf("%s,%s,%s", tt.tpl, tt.appname, tt.ans)
		t.Run(testName, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = []string{tt.tpl, tt.appname, tt.ans}
			main()
			//cmd := exec.Command("go-gitter", tt.tpl, tt.appName, tt.ans)
			//bGot, err := cmd.CombinedOutput()
			//got := string(bGot)
			//
			//if err != nil {
			//   t.Errorf("%v", err)
			//}
			//fmt.Printf("%v", got)
			//if got != tt.want {
			//	t.Errorf("got %s, want %s", got, tt.want)
			//}
		})
	}
}

func TestInput(t *testing.T) {
	var tests = []struct {
		name, tpl, appName, ans, want string
	}{
		{"noArgs", "", "", "", errMsgs[0]},
		{"noAppPath", "test-template-path", "", "", errMsgs[1]},
		// {"https://example.com/dummy-template", "appPath3", "", "path/URL to template is not in the allow-list"},
		// {"./fixtures/tpl-1", "", "", "path/URL to template does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// backup args and restore after test run.
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			// set args for test.
			os.Args = []string{"dummyPath", tt.tpl, tt.appName, tt.ans}
			// exec code.
			_, gotErr := getArgs()

			if !strings.Contains(gotErr.Error(), tt.want) {
				t.Errorf("got %q, want %q", gotErr, tt.want)
			}
		})
	}
}

func TestGetSettings(t *testing.T) {
	var tests = []struct {
		name, file, want string
	}{
		{"configNotFound", "does-not-exist", "could not open"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// exec code.
			_, gotErr := settings(tt.file)

			if !strings.Contains(gotErr.Error(), tt.want) {
				t.Errorf("got %q, want %q", gotErr, tt.want)
			}
		})
	}

	t.Run("canReadConfig", func(t *testing.T) {
		// exec code.
		want := "test.com"
		cfg, err := settings("fixtures/config.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}

		got, ok := cfg.Array("urlsAllowed")
		if !ok {
			t.Errorf("got %v, want %v", ok, !ok)
		}

		if got[0] != want {
			t.Errorf("got %v, want [%v]", got, want)
		}
	})
}
