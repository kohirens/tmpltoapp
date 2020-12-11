package main

import (
	"fmt"
	"testing"
	"os"
)

const testTmp = "go_gitter_test_tmp"

func TestMain(m *testing.M) {
    // call flag.Parse() here if TestMain uses flags
    os.Mkdir(testTmp, 0777) // set up a temporary dir for generate files

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
		tpl, appName, ans, want string
	}{
		{"", "", "", "help"},
		// {"non-existent", "", "", "The new application path is required."},
		// {"htttps://github.com/dummy-template", "", "", "path/URL to template does not exist"},
		// {"./fixtures/tpl-1", "", "", "path/URL to template does not exist"},
	}

	for _, tt := range tests {

		testName := fmt.Sprintf("%s,%s,%s", tt.tpl, tt.appName, tt.ans)
		t.Run(testName, func(t *testing.T) {
		    // backup args and restore after test run.
            oldArgs := os.Args
            defer func() { os.Args = oldArgs }()
            // set args for test.
		    os.Args = []string{"dummyPath", tt.tpl, tt.appName, tt.ans}
		    //
		    options := getArgs()

			if options["tplPath"] != tt.want {
				t.Errorf("got %q, want %q", options["tplPath"], tt.want)
			}
		})
	}
}
