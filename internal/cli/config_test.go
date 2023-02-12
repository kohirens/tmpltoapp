package cli

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetTmplLocation(runner *testing.T) {
	fixtures := []struct {
		name, want string
		cfg        *Config
	}{
		{"relative", "local", &Config{TmplPath: "./"}},
		{"relative2", "local", &Config{TmplPath: "."}},
		{"relativeUp", "local", &Config{TmplPath: ".."}},
		{"absolute", "local", &Config{TmplPath: "/home/myuser"}},
		{"windows", "local", &Config{TmplPath: "C:\\Temp"}},
		{"http", "remote", &Config{TmplPath: "http://example.com/repo1"}},
		{"https", "remote", &Config{TmplPath: "https://example.com/repo1"}},
		{"git", "remote", &Config{TmplPath: "git://example.com/repo1"}},
		{"file", "remote", &Config{TmplPath: "file://example.com/repo1"}},
		{"hiddenRelative", "local", &Config{TmplPath: ".m/example.com/repo1"}},
		{"tildeRelative", "local", &Config{TmplPath: "~/repo1.git"}},
	}

	for _, tc := range fixtures {
		runner.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.getTmplLocation()

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetSettings(t *testing.T) {

	t.Run("configNotFound", func(t *testing.T) {
		cfgFixture := &Config{}
		gotErr := cfgFixture.LoadUserSettings("does-not-exist")

		if !strings.Contains(gotErr.Error(), "could not open") {
			t.Errorf("got %q, want %q", gotErr, "could not open")
		}
	})

	t.Run("canReadConfig", func(t *testing.T) {
		cfgFixture := &Config{}
		err := cfgFixture.LoadUserSettings(FixtureDir + PS + "config-01.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}
	})
}

func TestInitConfigFile(t *testing.T) {
	var testCases = []struct {
		name string
		want error
		cfg  *Config
	}{
		{"NotExist", nil, &Config{Path: TmpDir + PS + "config-fix-01.json"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.cfg.initFile()

			if tc.want != gotErr {
				t.Errorf("got %t; want %t", gotErr, tc.want)
			}

			_, err := os.Stat(tc.cfg.Path)

			if err != nil {
				t.Errorf("got %v; want %v", gotErr, tc.want)
			}
		})
	}
}

func TestLoadAnswers(tester *testing.T) {
	var fixtures = []struct {
		name, file, want string
	}{
		{"goodJson", FixtureDir + PS + "answers-01.json", "value1"},
		{"badJson", FixtureDir + PS + "answers-02.json", ""},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(t *testing.T) {
		got, err := LoadAnswers(fxtr.file)

		if err == nil && got.Placeholders["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got.Placeholders["var1"], fxtr.want)
		}
	})

	fxtr = fixtures[1]
	tester.Run(fxtr.name, func(t *testing.T) {
		_, err := LoadAnswers(fxtr.file)

		if err == nil {
			t.Error("did not get an error")
		}
	})
}

//func xTestSubCmdConfigBadExit(tester *testing.T) {
//	oldArgs := os.Args
//	var tests = []struct {
//		name     string
//		wantCode int
//		args     []string
//		expected string
//	}{
//		{"noArgs", 1, []string{oldArgs[0], CmdConfig}, "usage: config"},
//		{"keyDoesNotExist", 1, []string{oldArgs[0], CmdConfig, "set", "key", "value"}, "no config setting \"keyDoesNotExist\" found"},
//	}
//
//	for _, tc := range tests {
//		tester.Run(tc.name, func(t *testing.T) {
//			defer func() {
//				os.Args = oldArgs
//			}()
//
//			// replace the test flaga and arts with the test fixture flags and args.
//			os.Args = tc.args
//
//			testConf := &Config{}
//			defineFlags(testConf)
//			e := parseFlags(testConf)
//			if e != nil {
//				t.Error(e)
//			}
//		})
//	}
//}

func xTestLoadUserSettings(tester *testing.T) {
	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")
		_ = os.Setenv("LOCALAPPDATA", TmpDir)
		defer func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}()
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", TmpDir)
		defer func() {
			_ = os.Setenv("HOME", oldHome)
		}()
	}

	var tests = []struct {
		name     string
		filename string
		want     *Config
	}{
		{
			"goodFile",
			FixtureDir + PS + "good-config-01.json",
			&Config{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
		{
			"badFile",
			FixtureDir + PS + "bad-config-01.json",
			&Config{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			gotCfg := &Config{}
			err := gotCfg.LoadUserSettings(tc.filename)

			if err != nil { // test bad values

			} else { // test good values

			}
		})
	}
}
