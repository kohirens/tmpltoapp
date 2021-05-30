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

const (
	FIXTURES_DIR = "testdata"
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

func TestDownload(t *testing.T) {
	var err error
	c := HttpMock{
		&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("200 OK")),
			StatusCode: 200,
		},
		err,
	}

	t.Run("canDownload", func(t *testing.T) {
		got, err := download("/fake_dl", TEST_TMP, &c)
		if err != nil {
			t.Errorf("got %q, want nil", err.Error())
		}
		_, err = os.Stat(got)

		if os.IsNotExist(err) {
			t.Errorf("got %q, want nil", got)
		}
	})
}

func ExampleDownload() {
	client := http.Client{}
	_, err := download(
		"https://github.com/kohirens/go-gitter-test-tpl/archive/main.zip",
		TEST_TMP,
		&client,
	)

	if err != nil {
		return
	}
}

func TestExtract(t *testing.T) {
	t.Run("canExtractDownload", func(t *testing.T) {
		wd, _ := os.Getwd()
		fixture := wd + "/" + FIXTURES_DIR + "/001.zip"
		want := TEST_TMP + "/sample_main"
		err := extract(fixture, want)

		if err != nil {
			t.Errorf("could not extract %s, error: %v", want, err.Error())
		}
	})
}

func ExampleExtract() {
	err := extract(
		TEST_TMP+"/001.zip",
		TEST_TMP+"/sample",
	)

	if err != nil {
		return
	}
}

func TestCopyFiles(t *testing.T) {
	err := copyDir(FIXTURES_DIR+"/template-01/tt.tpl", TEST_TMP+"/tt.app")
	if err == nil {
		t.Errorf("copyDir did nit err")
	}
}

func TestCopyDirSuccess(test *testing.T) {
	cases := []struct {
		dstDir, name, srcDir string
		want                 error
		IsVerified           func(string) bool
	}{
		{TEST_TMP + "/template-01-out", "success", FIXTURES_DIR + "/template-01", nil, func(p string) bool { return !stdlib.PathExist(p) }},
	}

	for _, sbj := range cases {
		test.Run(sbj.name, func(t *testing.T) {
			err := copyDir(sbj.srcDir, sbj.dstDir)
			isAllGood := sbj.IsVerified(sbj.dstDir)

			if err != sbj.want {
				t.Errorf("Could not copy dir, err: %s", err.Error())
			}

			if isAllGood {
				t.Errorf("all is not good: %v", isAllGood)
			}
		})
	}
}

func TestParse(test *testing.T) {
	type wanted func(error) bool

	cases := []struct {
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
			FIXTURES_DIR + "/template-02/file-01.tpl",
			TEST_TMP + "/appDirParse-01",
			tplVars{"var1": "1234"},
			func(err error) bool {
				f, _ := ioutil.ReadFile(TEST_TMP + "/appDirParse-01/file-01.tpl")
				s := string(f)
				return s == "testings 1234"
			},
			"failed with valid input",
		},
	}

	err := os.MkdirAll(TEST_TMP+"/appDirParse-01", os.FileMode(DIR_MODE))
	if err != nil {
		test.Errorf("Could not copy dir, err: %s", err.Error())
	}

	for _, sbj := range cases {
		test.Run(sbj.name, func(t *testing.T) {
			err := parse(sbj.tplFile, sbj.appDir, sbj.vars)

			if !sbj.gotWhatIWant(err) {
				t.Error(sbj.failMsg)
			}
		})
	}
}

func TestGetPathType(test *testing.T) {
	fixturePath1, _ := filepath.Abs(FIXTURES_DIR + "/template-01")

	cases := []struct {
		name, tmplPath, want string
	}{
		{"localAbsolutePath", fixturePath1, "local"},
		{"localRelativePath", FIXTURES_DIR + "/template-01", "local"},
		{"httpPath", "http://example.com", "http"},
		{"httpSecurePath", "https://example.com", "http"},
	}

	for _, sbj := range cases {
		test.Run(sbj.name, func(t *testing.T) {
			got := getPathType(sbj.tmplPath)

			if got != sbj.want {
				t.Errorf("got %q, want %q", got, sbj.want)
			}
		})
	}
}

func TestParseDir(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(FIXTURES_DIR + "/parse-dir-01")
	tmpDir, _ := filepath.Abs(TEST_TMP)

	fixtures := []struct {
		name, tmplPath, outPath string
		tplVars tplVars
		fileToCheck, want string
	}{
		{
			"parse-dir-01", fixturePath1, tmpDir + "/parse-dir-01",
			tplVars{"APP_NAME": "SolarPolar"},
			tmpDir + "/parse-dir-01/dir1/README.md", "SolarPolar\n",
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		err := parseDir(fxtr.tmplPath, fxtr.outPath, fxtr.tplVars)

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
