package main

import (
	"path/filepath"
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
		want2     string
	}{
		{"cloneRepo1", fixturesDir + "/repo-01.git", testTmp + "/repo-01", "refs/heads/main", false, "b7e42844c597d2beaf774eddfdcb653a2a4b0050", testTmp + "/repo-01"},
		{"clone404", fixturesDir + "/does-not-exit.git", testTmp, "", true, "", ""},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {

			gotPath, gotHash, err := gitClone(tc.repo, tc.outPath, tc.branch)

			if tc.shouldErr == true && err == nil {
				t.Error("did not get expected err")
			}

			if tc.shouldErr == false && err != nil {
				t.Errorf("got an unexpected err: %s", err)
			}

			if gotHash != tc.want {
				t.Errorf("got %v, want %v", gotHash, tc.want)
			}

			if gotPath != tc.want2 {
				t.Errorf("got %v, want %v", gotPath, tc.want2)
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
		{"bareRepo", "/repo-01.git", "repo-01.git"},
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
