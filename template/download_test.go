package template

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

type HttpMock struct {
	Resp *http.Response
	Err  error
}

func (h HttpMock) Get(url string) (*http.Response, error) {
	return h.Resp, h.Err
}

func TestDownload(t *testing.T) {

	var err error
	c := HttpMock{
		&http.Response{
			Body: ioutil.NopCloser(strings.NewReader("200 OK")),
		},
		err,
	}

	got := Download("fake_path", "fixture_path", &c)

	if got != nil {
		t.Errorf("got %q, want nil", got)
	}
}

func xTestInput(t *testing.T) {
	var tests = []struct {
		name, tpl, appName, ans, want string
	}{
		{"notInAllowList", "https://example.com/dummy-template", "appPath3", "", "path/URL to template is not in the allow-list"},
		// {"./fixtures/tpl-1", "", "", "path/URL to template does not exist"},
	}
	c := HttpMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set args for test.
			os.Args = []string{"dummyPath", tt.tpl, tt.appName, tt.ans}
			// exec code.
			gotErr := Download(tt.tpl, tt.appName, &c)

			if !strings.Contains(gotErr.Error(), tt.want) {
				t.Errorf("got %q, want %q", gotErr, tt.want)
			}
		})
	}
}
