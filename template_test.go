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
		"https://github.com/kohirens/go-gitter-test-tpl/archive/main.zip",
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
		fixture := wd + "/" + fixturesDir + "/001.zip"
		want := testTmp + "/sample_main"
		_, err := extract(fixture, want)

		if err != nil {
			t.Errorf("could not extract %s, error: %v", want, err.Error())
		}
	})
}

func ExampleExtract() {
	_, err := extract(
		testTmp+"/001.zip",
		testTmp+"/sample",
	)

	if err != nil {
		return
	}
}

func TestCopyFiles(test *testing.T) {
	err := copyDir(fixturesDir+"/template-01/tt.tpl", testTmp+"/tt.app")
	if err == nil {
		test.Errorf("copyDir did not err")
	}
}

func TestCopyDirSuccess(tester *testing.T) {
	fixtures := []struct {
		dstDir, name, srcDir string
		want                 error
		IsVerified           func(string) bool
	}{
		{testTmp + "/template-01-out", "success", fixturesDir + "/template-01", nil, func(p string) bool { return !stdlib.PathExist(p) }},
	}

	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := copyDir(fxtr.srcDir, fxtr.dstDir)
			isAllGood := fxtr.IsVerified(fxtr.dstDir)

			if err != fxtr.want {
				test.Errorf("Could not copy dir, err: %s", err.Error())
			}

			if isAllGood {
				test.Errorf("all is not good: %v", isAllGood)
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

func TestGetPathType(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(fixturesDir + "/template-01")

	fixtures := []struct {
		name, tmplPath, want string
	}{
		{"localAbsolutePath", fixturePath1, "local"},
		{"localRelativePath", fixturesDir + "/template-01", "local"},
		{"httpPath", "http://example.com", "http"},
		{"httpSecurePath", "https://example.com", "http"},
	}

	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			got := getPathType(fxtr.tmplPath)

			if got != fxtr.want {
				test.Errorf("got %q, want %q", got, fxtr.want)
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

		err := parseDir(fxtr.tmplPath, fxtr.outPath, fxtr.tplVars, fec)

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
