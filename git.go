package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// gitClone Clone a repo from a path/URL to a local directory.
func gitClone(repoUrl, localCache, branchName string) (string, string, error) {
	infof("branch to clone is %q", branchName)
	infof("git clone %s", repoUrl)

	sco, e1 := gitCmd(localCache, "clone", repoUrl)

	if e1 != nil {
		return "", "", fmt.Errorf("error cloning %v: %s", repoUrl, e1.Error())
	}

	infof("clone output \n%s", sco)

	repoDir := localCache + PS + getRepoDir(repoUrl)

	latestCommitHash, e2 := getLastCommitHash(repoDir)
	if e2 != nil {
		return "", "", fmt.Errorf("error getting commit hash %v: %s", repoDir, e2.Error())
	}

	return repoDir, latestCommitHash, nil
}

// gitCheckout Open an existing repo and checkout commit by full ref-name
func gitCheckout(repoLocalPath, ref string) (string, string, error) {
	sco, e1 := gitCmd(repoLocalPath, "fetch", "--all", "-p")
	if e1 != nil {
		infof("failed it this far")
		return "", "", fmt.Errorf("fetch failed on %s and %s; %s", repoLocalPath, ref, e1.Error())
	}

	infof("git fetch output %v", sco)
	infof("ref = %v ", ref)
	infof("git checkout %s", ref)

	sco2, e2 := gitCmd(repoLocalPath, "checkout", ""+ref)
	if e2 != nil {
		return "", "", fmt.Errorf("git checkout failed: %s", e2.Error())
	}

	infof("git fetch output %v", sco2)

	repoDir, e8 := filepath.Abs(repoLocalPath)
	if e8 != nil {
		return "", "", e8
	}

	latestCommitHash, e4 := getLastCommitHash(repoDir)
	if e4 != nil {
		return "", "", e4
	}

	return repoDir, latestCommitHash, nil
}

// gitCmd run a git command.
func gitCmd(repoPath string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()
	cmd.Dir = repoPath
	cmdStr := cmd.String()
	infof("running command %s", cmdStr)
	cmdOut, cmdErr := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()

	if cmdErr != nil {
		return nil, fmt.Errorf("error running git %v: %v", args, cmdErr.Error())
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("git %v returned exit code %q", args, exitCode)
	}

	return cmdOut, nil
}

// getLastCommitHash Returns the HEAD commit hash.
func getLastCommitHash(repoDir string) (string, error) {
	latestCommitHash, e1 := gitCmd(repoDir, "rev-parse", "HEAD")
	if e1 != nil {
		return "", fmt.Errorf("error getting commit hash %v: %s", repoDir, e1.Error())
	}

	return strings.Trim(string(latestCommitHash), "\n"), nil
}

// getRepoDir extract a local dirname from a Git URL.
func getRepoDir(repoLocation string) string {
	if len(repoLocation) < 1 {
		return repoLocation
	}

	isGitUri := regexp.MustCompile("^(git|http|https)://|.+git$")
	baseName := filepath.Base(repoLocation)
	if isGitUri.MatchString(repoLocation) {
		return strings.Replace(baseName, ".git", "", 1)
	}

	return baseName
}
