package manifest

import (
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/tmplpress/internal/press"
	"os"
	"reflect"
	"testing"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
)

func TestGenerateATemplateJson(runner *testing.T) {
	testCases := []struct {
		name string
		repo string
		want map[string]string
	}{
		{
			"onlyDataEvaluations",
			"repo-06",
			map[string]string{
				"appTitle": "",
				"name":     "",
				"age":      "",
			},
		},
		{
			"overwrite-bad-with-default",
			"repo-14",
			nil,
		},
	}

	for _, tc := range testCases {
		runner.Run(tc.name, func(t *testing.T) {
			repoPath := git.CloneFromBundle(tc.repo, tmpDir, fixtureDir, ps)

			got, err := generateATemplateManifest(repoPath)
			if err != nil {
				t.Errorf("want nil, got: %q", err.Error())
			}

			if !fsio.Exist(got) {
				t.Errorf("no template.json found in %v", repoPath)
			}

			b, _ := os.ReadFile(got)
			tm, _ := press.NewTmplManifest(b)
			if !reflect.DeepEqual(tm.Placeholders, tc.want) {
				t.Errorf("got %v, want %v", tm.Placeholders, tc.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		cmd     string
		wantErr bool
		want    map[string]string
	}{
		{"case-1", "repo-07", "generate", false, map[string]string{"Placeholder1": ""}},
	}

	for _, tt := range tests {
		repoPath := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)

		Init()

		t.Run(tt.name, func(t *testing.T) {
			err := Run([]string{tt.cmd, repoPath})
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, _ := os.ReadFile("tmp/repo-07/template.json")
			tm, _ := press.NewTmplManifest(b)

			if !reflect.DeepEqual(tm.Placeholders, tt.want) {
				t.Errorf("got %v, want %v", tm.Placeholders, tt.want)
			}
		})
	}
}

func TestRunValidate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		cmd      string
		wantErr  bool
	}{
		{"replace-dir-missing", fixtureDir + ps + "template-05/template.json", "validate", true},
		{"case-2", fixtureDir + ps + "template-06/template.json", "validate", false},
		{"placeholder-not-found", fixtureDir + ps + "template-2.2.0-02.json", "validate", true},
		{"empty-regexp", fixtureDir + ps + "template-2.2.0-03.json", "validate", true},
		{"invalid-regexp", fixtureDir + ps + "template-2.2.0-04.json", "validate", true},
	}

	for _, tt := range tests {
		Init()

		t.Run(tt.name, func(t *testing.T) {
			err := Run([]string{tt.cmd, tt.template})
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
