package main

import (
	"fmt"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/test"
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
	}{
		{
			"cloneRepo1",
			"repo-01.git",
			TmpDir + cli.PS + "repo-01-refs-heads-main",
			"refs/heads/main",
			false,
			"b7e42844c597d2beaf774eddfdcb653a2a4b0050",
		},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {
			repoPath := test.SetupARepository(tc.repo, TmpDir, FixtureDir, cli.PS)

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

			if gotPath != tc.outPath {
				t.Errorf("got %v, want %v", gotPath, tc.outPath)
			}
		})
	}
}

// Clone a repo
func TestGitCannotClone(tester *testing.T) {
	var testCases = []struct {
		name      string
		repo      string
		outPath   string
		branch    string
		shouldErr bool
		wantHash  string
	}{
		{"clone404", "does-not-exit.git", "", "", true, ""},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {
			repoPath := test.SetupARepositoryOld(tc.repo)

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

			if gotPath != tc.outPath {
				t.Errorf("got %v, want %v", gotPath, tc.outPath)
			}
		})
	}
}

func TestGetRepoDir(tester *testing.T) {
	var testCases = []struct {
		name    string
		repo    string
		refName string
		want    string
	}{
		{"local", "/repo-01", "", "repo-01"},
		{"url", "https://example.com/repo-01.git", "", "repo-01"},
		{"bareRepo", "/repo-01.git", "", "repo-01"},
		{"bareRepo", "/repo-02.git", "refs/heads/main", "repo-02-refs-heads-main"},
		{"bareRepo", "/repo-02.git", "refs/tags/0.1.0", "repo-02-refs-tags-0.1.0"},
		{"bareRepo", "/repo-02.git", "refs/remotes/origin/HEAD", "repo-02-refs-remotes-origin-HEAD"},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {

			got := getRepoDir(tc.repo, tc.refName)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetRepoDir2(tester *testing.T) {
	absTestTmp, _ := filepath.Abs(TmpDir)
	var testCases = []struct {
		name     string
		bundle   string
		branch   string
		want     string
		wantHash string
	}{
		{"localFullBranchRefArg", "repo-02", "refs/remotes/origin/third-commit", absTestTmp + cli.PS + "repo-02", "bfeb0a45c027420e4df286dc089965599e350bf9"},
	}

	for _, tc := range testCases {
		repoPath := test.SetupARepository(tc.bundle, TmpDir, FixtureDir, cli.PS)

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
		repoPath := test.SetupARepository(tc.bundle, TmpDir, FixtureDir, cli.PS)

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
		repoPath := test.SetupARepositoryOld(tc.bundle)

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
		repoPath := test.SetupARepository(tc.bundle, TmpDir, FixtureDir, cli.PS)

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
