package main

import (
	"os"
	"strings"
	"testing"
)

const TEST_TMP = "tmp"

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.RemoveAll(TEST_TMP)
	// Set up a temporary dir for generate files
	os.Mkdir(TEST_TMP, 0774) // set up a temporary dir for generate files
	// Run all tests
	exitcode := m.Run()
	// Clean up
	os.Exit(exitcode)
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
		want := FIXTURES_DIR + "/ans-1.yml"
		// set args for test.
		cfgFixture := []string{"go-gitter", "-answers=" + want, FIXTURES_DIR + "/tpl-1", "appPath4"}
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
