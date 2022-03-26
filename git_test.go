package main

import (
	"testing"
)

// Clone a repo
func TestGitClone(tester *testing.T) {
	var testCases = []struct {
		name      string
		repo      string
		outPath   string
		branch    string
		shouldErr bool
		want      string
	}{
		{"cloneRepo1", fixturesDir + "/repo-01.git", testTmp + "/repo-01", "main", false, "b7e42844c597d2beaf774eddfdcb653a2a4b0050"},
		{"clone404", fixturesDir + "/does-not-exit.git", testTmp + "/does-not-exit", "", true, ""},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {

			got, err := gitClone(tc.repo, tc.outPath, tc.branch)

			if tc.shouldErr == true && err == nil {
				t.Error("did not get expected err")
			}

			if tc.shouldErr == false && err != nil {
				t.Errorf("got an unexpected err: %s", err)
			}

			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
