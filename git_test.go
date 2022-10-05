package main

import (
	"fmt"
	"path/filepath"
	"strings"
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
		wantHash  string
		wantOut   string
	}{
		{"cloneRepo1", "repo-01.git", testTmp, "refs/heads/main", false, "b7e42844c597d2beaf774eddfdcb653a2a4b0050", testTmp + PS + "repo-01"},
		{"clone404", fixturesDir + "/does-not-exit.git", testTmp, "", true, "", ""},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {
			repoPath := setupARepository(tc.repo)

			gotPath, gotHash, err := gitClone(repoPath, tc.outPath, tc.branch)

			if tc.shouldErr == true && err == nil {
				t.Error("did not get expected err")
			}

			if tc.shouldErr == false && err != nil {
				t.Errorf("got an unexpected err: %s", err)
			}

			if gotHash != tc.wantHash {
				t.Errorf("got %v, want %v", gotHash, tc.wantHash)
			}

			if gotPath != tc.wantOut {
				t.Errorf("got %v, want %v", gotPath, tc.wantOut)
			}
		})
	}
}

func TestGetRepoDir(tester *testing.T) {
	var testCases = []struct {
		name string
		repo string
		want string
	}{
		{"local", "/repo-01", "repo-01"},
		{"url", "https://example.com/repo-01.git", "repo-01"},
		{"bareRepo", "/repo-01.git", "repo-01"},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {

			got := getRepoDir(tc.repo)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetRepoDir2(tester *testing.T) {
	absTestTmp, _ := filepath.Abs(testTmp)
	var testCases = []struct {
		name     string
		bundle   string
		branch   string
		want     string
		wantHash string
	}{
		{"local", "repo-02", "refs/remotes/origin/third-commit", absTestTmp + PS + "remotes" + PS + "repo-02", "bfeb0a45c027420e4df286dc089965599e350bf9"},
	}

	for _, tc := range testCases {
		repoPath := setupARepository(tc.bundle)

		tester.Run(tc.name, func(t *testing.T) {
			gotRepo, gotHash, gotErr := gitCheckout(repoPath, tc.branch)

			if gotErr != nil {
				t.Errorf("unexpected error in test %q", gotErr.Error())
			}

			if gotRepo != tc.want {
				t.Errorf("got %v, want %v", gotRepo, tc.want)
			}

			if gotHash != tc.wantHash {
				t.Errorf("got %v, want %v", gotHash, tc.wantHash)
			}
		})
	}
}

func TestGetLatestTag(tester *testing.T) {
	var testCases = []struct {
		name   string
		bundle string
		want   string
	}{
		{"found", "repo-03", "0.1.0"},
	}

	for i, tc := range testCases {
		repoPath := setupARepository(tc.bundle)

		tester.Run(fmt.Sprintf("%v.%v", i+1, tc.name), func(t *testing.T) {
			got, gotErr := getLatestTag(repoPath)

			if gotErr != nil {
				t.Errorf("unexpected error in test %q", gotErr.Error())
			}

			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetLatestTagError(tester *testing.T) {
	var testCases = []struct {
		name   string
		bundle string
		want   string
	}{
		{"doesNotExist", "repo-dne", ""},
	}

	for i, tc := range testCases {
		repoPath := setupARepository(tc.bundle)

		tester.Run(fmt.Sprintf("%v.%v", i+1, tc.name), func(t *testing.T) {
			got, gotErr := getLatestTag(repoPath)

			if gotErr == nil {
				t.Errorf("unexpected error in test %q", gotErr.Error())
			}

			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetRemoteTags(tester *testing.T) {
	var testCases = []struct {
		name      string
		bundle    string
		want      string
		shouldErr bool
	}{
		{"hasTags", "repo-04", "1.0.0,0.2.0,0.1.1,0.1.0", false},
		{"noTags", "repo-05", "", true},
	}

	for i, tc := range testCases {
		repoPath := setupARepository(tc.bundle)

		tester.Run(fmt.Sprintf("%v.%v", i+1, tc.name), func(t *testing.T) {
			got, gotErr := getRemoteTags(repoPath)

			if !tc.shouldErr && gotErr != nil {
				t.Errorf("unexpected error in test %q", gotErr.Error())
			}

			t1 := strings.Join(got, ",")
			if got != nil && t1 != tc.want {
				t.Errorf("got %v, want %v", t1, tc.want)
			}
		})
	}
}
