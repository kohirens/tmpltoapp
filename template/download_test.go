package template

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
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
