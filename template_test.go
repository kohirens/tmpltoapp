package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kohirens/stdlib"
)

type HttpMock struct {
	Resp *http.Response
	Err  error
}

func (h HttpMock) Get(url string) (*http.Response, error) {
	return h.Resp, h.Err
}

func (h HttpMock) Head(url string) (*http.Response, error) {
	return h.Resp, h.Err
}

func TestDownload(tester *testing.T) {
	defer quiet()()

	var err error
	fixtures := HttpMock{
		&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("200 OK")),
			StatusCode: 200,
		},
		err,
	}

	tester.Run("canDownload", func(test *testing.T) {
		got, err := download("/fake_dl", testTmp, &fixtures)
		if err != nil {
			test.Errorf("got %q, want nil", err.Error())
		}
		_, err = os.Stat(got)

		if os.IsNotExist(err) {
			test.Errorf("got %q, want nil", got)
		}
	})
}

func ExampleDownload() {
	client := http.Client{}
	_, err := download(
		"https://github.com/kohirens/tmpltoapp-test-tpl/archive/main.zip",
		testTmp,
		&client,
	)

	if err != nil {
		return
	}
}

func TestExtract(test *testing.T) {
	test.Run("canExtractDownload", func(t *testing.T) {
		wd, _ := os.Getwd()
		fixture := wd + PS + fixturesDir + PS + "001.zip"
		want := testTmp + PS + "001"
		_, err := extract(fixture)

		if err != nil {
			t.Errorf("could not extract %s, error: %v", want, err.Error())
		}
	})
}

func ExampleExtract() {
	_, err := extract(
		testTmp + PS + "001.zip",
	)

	if err != nil {
		return
	}
}

func TestParseDir2(tester *testing.T) {
	defer quiet()()

	fixtures := []struct {
		dstDir,
		name,
		srcDir string
		want func(error) bool
		vars tplVars
	}{
		{
			testTmp + PS + "template-04-out",
			"dir1IsEmpty",
			fixturesDir + "/template-04",
			func(e error) bool {
				return !stdlib.PathExist(testTmp + PS + "template-04-out" + PS + "dir1" + PS + ".empty")
			},
			tplVars{},
		},
	}

	fileChkr, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := parseDir(fxtr.srcDir, fxtr.dstDir, fxtr.vars, fileChkr, []string{})
			isAllGood := fxtr.want(err)

			if !isAllGood {
				test.Error("all is not good")
			}
		})
	}
}

func TestParse(tester *testing.T) {
	type wanted func(error) bool

	fixtures := []struct {
		name, tplFile, appDir string
		vars                  tplVars
		gotWhatIWant          wanted
		failMsg               string
	}{
		{
			"emptyInput",
			"",
			"",
			tplVars{"var1": "1234"},
			func(err error) bool { return err != nil },
			"failed with no input.",
		},
		{
			"validInput",
			fixturesDir + "/template-02/file-01.tpl",
			testTmp + "/appDirParse-01",
			tplVars{"var1": "1234"},
			func(err error) bool {
				f, _ := ioutil.ReadFile(testTmp + "/appDirParse-01/file-01.tpl")
				s := string(f)
				return s == "testings 1234"
			},
			"failed with valid input",
		},
	}

	err := os.MkdirAll(testTmp+"/appDirParse-01", os.FileMode(DIR_MODE))
	if err != nil {
		tester.Errorf("Could not copy dir, err: %s", err.Error())
	}

	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := parse(fxtr.tplFile, fxtr.appDir, fxtr.vars)

			if !fxtr.gotWhatIWant(err) {
				test.Error(fxtr.failMsg)
			}
		})
	}
}

func TestGetTmplLocation(runner *testing.T) {
	fixtures := []struct {
		name, path, want string
	}{
		{"relative", "./", "local"},
		{"relative2", ".", "local"},
		{"relativeUp", "..", "local"},
		{"absolute", "/home/myuser", "local"},
		{"windows", "C:\\Temp", "local"},
		{"http", "http://example.com/repo1", "remote"},
		{"https", "https://example.com/repo1", "remote"},
		{"git", "git://example.com/repo1", "remote"},
		{"file", "file://example.com/repo1", "remote"},
		{"hiddenRelative", ".m/example.com/repo1", "local"},
		{"tildeRelative", "~/repo1.git", "local"},
	}

	for _, tc := range fixtures {
		runner.Run(tc.name, func(test *testing.T) {
			got := getTmplLocation(tc.path)

			if got != tc.want {
				test.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseDir(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(fixturesDir + "/parse-dir-01")
	tmpDir, _ := filepath.Abs(testTmp)

	fixtures := []struct {
		name, tmplPath, outPath string
		tplVars                 tplVars
		fileToCheck, want       string
	}{
		{
			"parse-dir-01", fixturePath1, tmpDir + "/parse-dir-01",
			tplVars{"APP_NAME": "SolarPolar"},
			tmpDir + "/parse-dir-01/dir1/README.md", "SolarPolar\n",
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		fec, _ := stdlib.NewFileExtChecker(nil, &[]string{"md", "yml"})

		err := parseDir(fxtr.tmplPath, fxtr.outPath, fxtr.tplVars, fec, nil)

		if err != nil {
			test.Errorf("got an error %q", err.Error())
			return
		}

		got, err := ioutil.ReadFile(tmpDir + "/parse-dir-01/dir1/README.md")

		if err != nil {
			test.Errorf("got an error %q", err.Error())
		}

		//TODO: Verify the files in the parsed dir was processed
		if string(got) != fxtr.want {
			test.Errorf("got %q, but want %q", string(got), fxtr.want)
		}
	})
}

func TestReadTemplateJson(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(fixturesDir + "/template-03")

	fixtures := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			"canBeFound",
			&Config{
				tplPath: fixturePath1,
			},
			false,
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		got, err := readTemplateJson(fxtr.config.tplPath + PS + TMPL_MANIFEST)

		if fxtr.shouldErr && err == nil {
			test.Errorf("expected an error, but got nil")
		}

		if got.Version != "0.1.0" {
			test.Error("could not get version from template.json")
			return
		}
	})
}

func TestQuestionsInput(tester *testing.T) {
	// Use a temp file to simulate input on the command line.
	tmpFile, err := ioutil.TempFile(testTmp, "qi-01")
	if err != nil {
		tester.Errorf("failed to make temp file %v", err.Error())
	}

	defer os.Remove(tmpFile.Name())

	if _, e := tmpFile.Write([]byte("1\n")); e != nil {
		tester.Errorf("failed to write content to temp file %v, error: %v", tmpFile.Name(), e.Error())
	}

	resetTmpFile := func() {
		if _, e := tmpFile.Seek(0, 0); e != nil {
			tester.Errorf("failed to reset to beginning of temp file %v, error: %v", tmpFile.Name(), e.Error())
		}
	}

	fixtures := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			"missingAnAnswer",
			&Config{
				answers: tplVars{"var1": "", "var2": ""},
				Questions: questions{
					Version:   "0.1.0",
					Variables: tplVars{"var1": "var1", "var2": "var2", "var3": "var3"},
					Excludes:  nil,
				},
			},
			"var3",
		},
		{
			"noMissingAnswers",
			&Config{
				answers: tplVars{"var1": "1", "var2": "2", "var3": "3"},
				Questions: questions{
					Version:   "0.1.0",
					Variables: tplVars{"var1": "var1", "var2": "var2", "var3": "var3"},
					Excludes:  nil,
				},
			},
			"var3",
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		resetTmpFile()
		err := getInput(&fxtr.config.Questions, &fxtr.config.answers, tmpFile)

		if err != nil {
			test.Errorf("got an error %q", err.Error())
		}

		if fxtr.config.answers[fxtr.want] != "1" {
			test.Errorf("failed to answer missing question %v using file as input", fxtr.want)
		}
	})
}
