package stdlib

import (
	"testing"
)

const (
	FIXTURES_DIR = "testdata"
	TEST_TMP     = "testtmp"
)

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

func TestIsTextFile(t *testing.T) {

	cases := []struct {
		name, path string
		want       bool
	}{
		{"IsATextFile", FIXTURES_DIR + "/text-file-01.txt", true},
		{"notATextFile", FIXTURES_DIR + "/text-file-02.jpg", false},
		{"notATextFile", FIXTURES_DIR + "/text-file-03.gif", false},
		{"notATextFile", FIXTURES_DIR + "/text-file-04.png", false},
		{"notATextFile", FIXTURES_DIR + "/text-file-05.json", true},
		{"notATextFile", FIXTURES_DIR + "/text-file-06.md", true},
		{"notATextFile", FIXTURES_DIR + "/text-file-07.xml", true},
	}

	for _, sbj := range cases {

		got := IsTextFile(sbj.path)

		if got != sbj.want {
			t.Errorf("got %v, want %v, for %v", got, sbj.want, sbj.path)
		}
	}
}
