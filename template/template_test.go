package template

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kohirens/go-gitter/stdlib"
)

const (
	FIXTURES_DIR = "testdata"
	TEST_TMP     = "testtmp"
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

func TestMain(m *testing.M) {
	os.RemoveAll(TEST_TMP)
	// Set up a temporary dir for generate files
	os.Mkdir(TEST_TMP, 0774)
	// Run all tests
	exitcode := m.Run()
	// Clean up
	os.Exit(exitcode)
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
		got := Download("/fake_dl", TEST_TMP, &c)

		_, err := os.Stat(TEST_TMP + "/fake_dl")

		if got != nil || os.IsNotExist(err) {
			t.Errorf("got %q, want nil", got)
		}
	})
}

func ExampleDownload() {
	client := http.Client{}
	err := Download(
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
		err := Extract(fixture, want)

		if err != nil {
			t.Errorf("could not extract %s, error: %v", want, err.Error())
		}
	})
}

func ExampleExtract() {
	err := Extract(
		TEST_TMP+"/001.zip",
		TEST_TMP+"/sample",
	)

	if err != nil {
		return
	}
}

func TestCopyFiles(t *testing.T) {
	err := CopyDir(FIXTURES_DIR+"/template-01/tt.tpl", TEST_TMP+"/tt.app")
	if err == nil {
		t.Errorf("CopyDir did nit err")
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
			err := CopyDir(sbj.srcDir, sbj.dstDir)
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
	cases := []struct {
		name, tplFile, appDir string
		vars                  tplVars
		err                   error
		want                  string
	}{
		{"emptyInput", "", "", tplVars{"var1": "1234"}, fmt.Errorf("no good"), "testing 1234"},
		{"validInput", FIXTURES_DIR + "/template-02/file-01.tpl", TEST_TMP + "/appDirParse-01", tplVars{"var1": "1234"}, nil, "testing 1234"},
	}

	err := os.MkdirAll(TEST_TMP+"/appDirParse-01", os.FileMode(0774))
	if err != nil {
		test.Errorf("Could not copy dir, err: %s", err.Error())
	}

	for _, sbj := range cases {
		test.Run(sbj.name, func(t *testing.T) {
			err := Parse(sbj.tplFile, sbj.appDir, sbj.vars)

			if err != nil && sbj.err != nil {
				t.Errorf("Could not copy dir, err: %s", err.Error())
			}
		})
	}
}
