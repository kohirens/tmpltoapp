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
		got, err := settings(FIXTURES_DIR + "/config-01.json")
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
		got, err := settings(FIXTURES_DIR + "/config-01.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}

		if got.allowedUrls[0] != want {
			t.Errorf("got %v, want [%v]", got, want)
		}
	})
}

func TestInitConfigFile(t *testing.T) {
	var tests = []struct {
		name, file string
		want       error
	}{
		{"NotExist", TEST_TMP + PS + "config-fix-01.json", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// exec code.
			got := initConfigFile(tt.file)

			if tt.want != got {
				t.Errorf("got %v; want %v", got, tt.want)
			}

			_, err := os.Stat(tt.file)

			if err != nil {
				t.Errorf("got %v; want %v", got, tt.want)
			}
		})
	}
}
