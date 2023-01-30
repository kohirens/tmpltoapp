package test

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/tmpltoapp/internal"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	fixturesDir = "testdata"
	testTmp     = "tmp"
	// SubCmdFlags space separated list of command line flags.
	SubCmdFlags = "SUB_CMD_FLAGS"
	// subCmdFlags space separated list of command line flags.
	subCmdFlags = "RECURSIVE_TEST_FLAGS"
	testRemotes = testTmp + internal.PS + "remotes"
)

func SetupARepository(bundleName string) string {
	repoPath := testRemotes + internal.PS + bundleName
	// It may have already been unbundled.
	fileInfo, err1 := os.Stat(repoPath)
	if (err1 == nil && fileInfo.IsDir()) || os.IsExist(err1) {
		absPath, e2 := filepath.Abs(repoPath)
		if e2 == nil {
			return absPath
		}
		return repoPath
	}

	wd, e := os.Getwd()
	if e != nil {
		panic(fmt.Sprintf("%v failed to get working directory", e.Error()))
	}

	srcRepo := wd + internal.PS + fixturesDir + internal.PS + bundleName + ".bundle"
	// It may not exist.
	if !stdlib.PathExist(srcRepo) {
		panic(fmt.Sprintf("%v bundle not found", srcRepo))
	}

	cmd := exec.Command("git", "clone", "-b", "main", srcRepo, repoPath)
	_, _ = cmd.CombinedOutput()
	if ec := cmd.ProcessState.ExitCode(); ec != 0 {
		log.Panicf("error un-bundling %q to a temporary repo %q for a unit test", srcRepo, repoPath)
	}

	absPath, e2 := filepath.Abs(repoPath)
	if e2 != nil {
		fmt.Printf("\n\npanicing\n\n")
		panic(e2.Error())
	}

	fmt.Printf("\nabsPath = %v\n", absPath)
	return absPath
}
