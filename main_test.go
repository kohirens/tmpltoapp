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
		name, want string
		config     []string
	}{
		{"noArgs", errMsgs[0], []string{"go-gitter", "", ""}},
		{"noAppPath", errMsgs[1], []string{"go-gitter", "templatePath", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// exec code.
			cfg := Config{}
			gotErr := getArgs(tt.config[0], tt.config[1:], &cfg)

			if !strings.Contains(gotErr.Error(), tt.want) {
				t.Errorf("got %q, want %q", gotErr, tt.want)
			}
		})
	}

	t.Run("allGood", func(t *testing.T) {
		cfg := Config{}
		want := "./fixtures/ans-1.yml"
		// set args for test.
		cfgFixture := []string{"go-gitter", "-answers=" + want, "./fixtures/tpl-1", "appPath4"}
		// exec code.
		err := getArgs(cfgFixture[0], cfgFixture[1:], &cfg)
		if err != nil {
			t.Errorf("got unexpected error: %v", err.Error())
		}

		if cfg.answersPath != want {
			t.Errorf("got %q, want %q", cfg.answersPath, want)
		}
	})
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
		got, err := settings("fixtures/config.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}

		if got.allowedUrls[0] != want {
			t.Errorf("got %v, want [%v]", got, want)
		}
	})
}

func TestUrlIsAllowed(t *testing.T) {
	var tests = []struct {
		name, url    string
		want1, want2 bool
	}{
		{"isAllowed", "https://test.com", true, true},
		{"notAllowed", "https://gitit.com", true, false},
		{"notAUrl", "/local/path", false, false},
	}

	fixtures := []string{"https://test.com"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// exec code.
			got1, got2 := urlIsAllowed(tt.url, fixtures)

			if tt.want1 != got1 && tt.want2 != got2 {
				t.Errorf("got [%v,%v]; want [%v, %v]", got1, got2, tt.want1, tt.want2)
			}
		})
	}

	t.Run("canReadConfig", func(t *testing.T) {
		// exec code.
		want := "test.com"
		got, err := settings("fixtures/config.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}

		if got.allowedUrls[0] != want {
			t.Errorf("got %v, want [%v]", got, want)
		}
	})
}
