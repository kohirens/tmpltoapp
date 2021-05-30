package main

import (
	"os"
	"strings"
	"testing"
)

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
		got, err := settings(fixturesDir + "/config-01.json")
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
		got, err := settings(fixturesDir + "/config-01.json")
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
		{"NotExist", testTmp + PS + "config-fix-01.json", nil},
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

		if err == nil && got["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got["var1"], fxtr.want)
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
