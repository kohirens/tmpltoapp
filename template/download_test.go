package template

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

const TEST_TMP = "go_gitter_test_tmp"

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
	os.Mkdir(TEST_TMP, 0774) // set up a temporary dir for generate files

	// Create whatever test files are needed.

	// Run all tests and clean up.
	exitcode := m.Run()
	os.RemoveAll(TEST_TMP)
	os.Exit(exitcode)
}

func TestDownload(t *testing.T) {
	var err error
	c := HttpMock{
		&http.Response{
			Body: ioutil.NopCloser(strings.NewReader("200 OK")),
		},
		err,
	}

	t.Run("canDownload", func(t *testing.T) {
		got := Download("fake_path", &c)

		if got != nil {
			t.Errorf("got %q, want nil", got)
		}
	})
}

func ExampleDownload() {
	client := http.Client{}
	err := Download(
		"https://github.com/kohirens/go-gitter-test-tpl/archive/main.zip",
		&client,
	)

	if err != nil {
		return
	}
}
