package stdlib

import "testing"

const FIXTURES_DIR = "testdata"

func TestPathexists(t *testing.T) {

	cases := []struct {
		name, path string
		want       bool
	}{
		{"existIsTrue", FIXTURES_DIR + "/file-exist-01.md", true},
		{"existIsTrue", FIXTURES_DIR + "/file-does-not-exist-01.md", false},
	}

	for _, sbj := range cases {
		got := PathExist(sbj.path)

		if got != sbj.want {
			t.Errorf("got %v, want %v", got, sbj.want)
		}
	}
}
