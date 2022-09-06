package main

import (
	"os"
	"strings"
	"testing"
)

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
		err := cfgFixture.loadUserSettings(fixturesDir + PS + "config-01.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}
	})
}

func TestInitConfigFile(t *testing.T) {
	var tests = []struct {
		name string
		want error
		cfg  *Config
	}{
		{"NotExist", nil, &Config{path: testTmp + PS + "config-fix-01.json"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.cfg.initFile()

			if tt.want != gotErr {
				t.Errorf("got %t; want %t", gotErr, tt.want)
			}

			_, err := os.Stat(tt.cfg.path)

			if err != nil {
				t.Errorf("got %v; want %v", gotErr, tt.want)
			}
		})
	}
}

func TestLoadAnswers(test *testing.T) {
	var fixtures = []struct {
		name, file, want string
	}{
		{"goodJson", fixturesDir + PS + "answers-01.json", "value1"},
		{"badJson", fixturesDir + PS + "answers-02.json", ""},
	}

	fxtr := fixtures[0]
	test.Run(fxtr.name, func(t *testing.T) {
		got, err := loadAnswers(fxtr.file)

		if err == nil && got.Placeholders["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got.Placeholders["var1"], fxtr.want)
		}
	})

	fxtr = fixtures[1]
	test.Run(fxtr.name, func(t *testing.T) {
		_, err := loadAnswers(fxtr.file)

		if err == nil {
			t.Error("did not get an error")
		}
	})
}
