package main

import (
	"os"
	"strings"
	"testing"
)

const (
	fixturesDir = "testdata"
	testTmp     = "tmp"
)


func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.RemoveAll(testTmp)
	// Set up a temporary dir for generate files
	os.Mkdir(testTmp, DIR_MODE) // set up a temporary dir for generate files
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
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
			gotErr := parseArgs(tt.config[0], tt.config[1:], &cfg)

			if !strings.Contains(gotErr.Error(), tt.want) {
				t.Errorf("got %q, want %q", gotErr, tt.want)
			}
		})
	}

	t.Run("allGood", func(t *testing.T) {
		cfg := Config{}
		want := fixturesDir + "/ans-1.yml"
		// set args for test.
		cfgFixture := []string{"go-gitter", "-answers=" + want, fixturesDir + "/tpl-1", "appPath4"}
		// exec code.
		err := parseArgs(cfgFixture[0], cfgFixture[1:], &cfg)
		if err != nil {
			t.Errorf("got unexpected error: %v", err.Error())
		}

		if cfg.answersPath != want {
			t.Errorf("got %q, want %q", cfg.answersPath, want)
		}
	})
}
