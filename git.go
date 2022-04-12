package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// gitClone Clone a repo from a path/URL to a local directory.
func gitClone(repoUrl, repoLocalPath, branchName string) (string, string, error) {
	options := &git.CloneOptions{
		URL:          repoUrl,
		Progress:     os.Stdout,
		SingleBranch: false,
		Depth:        1,
	}

	if branchName != "" {
		options.ReferenceName = plumbing.ReferenceName(branchName)
	}

	infof("git clone %s", repoLocalPath)
	repo, e1 := git.PlainClone(repoLocalPath, false, options)
	if e1 != nil {
		return "", "", e1
	}

	// Retrieving the path to the cloned repo.
	w, e4 := repo.Worktree()
	if e4 != nil {
		return "", "", e4
	}

	// Retrieving the commit being pointed by HEAD
	ref, e5 := repo.Head()
	if e5 != nil {
		return "", "", e5
	}

	return w.Filesystem.Root(), ref.Hash().String(), nil
}

func gitCheckout(repoLocalPath, branchName string) (string, string, error) {
	infof("cd %s", repoLocalPath)
	r, e1 := git.PlainOpen(repoLocalPath)
	if e1 != nil {
		return "", "", e1
	}

	wt, e2 := r.Worktree()
	if e2 != nil {
		return "", "", e2
	}

	h, err := r.ResolveRevision(plumbing.Revision(branchName))
	if err != nil {
		return "", "", err
	}
	fmt.Printf("revision = %s\n", h.String())

	infof("git checkout %s", branchName)
	e3 := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
	})
	if e3 != nil {
		return "", "", e3
	}

	ref, e7 := r.Head()
	if e7 != nil {
		return "", "", e7
	}

	retVal, e8 := filepath.Abs(repoLocalPath)
	if e8 != nil {
		return "", "", e8
	}

	return retVal, ref.Hash().String(), nil
}

func getRepoDir(repoLocation string) string {
	if len(repoLocation) < 1 {
		return repoLocation
	}

	isUrl := regexp.MustCompile("(git|http|https)://")
	baseName := filepath.Base(repoLocation)
	if isUrl.MatchString(repoLocation) {
		return strings.Replace(baseName, ".git", "", 1)
	}
	return baseName
}
