package main

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	fixturesDir = "testdata"
	testTmp     = "tmp"
	// SubCmdFlags space separated list of command line flags.
	SubCmdFlags = "SUB_CMD_FLAGS"
	// subCmdFlags space separated list of command line flags.
	subCmdFlags = "RECURSIVE_TEST_FLAGS"
	testRemotes = testTmp + PS + "remotes"
)

func TestMain(m *testing.M) {
	// Only runs when this environment variable is set.
	if _, ok := os.LookupEnv(subCmdFlags); ok {
		runAppMain()
	}
	// call flag.Parse() here if TestMain uses flags
	_ = os.RemoveAll(testTmp)
	// Set up a temporary dir for generate files
	_ = os.Mkdir(testTmp, DirMode) // set up a temporary dir for generate files
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}

func TestCallingMain(tester *testing.T) {
	// This was adapted from https://golang.org/src/flag/flag_test.go; line 596-657 at the time.
	// This is called recursively, because we will have this test call itself
	// in a sub-command with the environment variable `GO_CHILD_FLAG` set.
	// Note that a call to `main()` MUST exit or you'll spin out of control.
	if os.Getenv(SubCmdFlags) != "" {
		// We're in the test binary, so test flags are set, lets reset it
		// so that only the program is set
		// and whatever flags we want.
		args := strings.Split(os.Getenv(SubCmdFlags), " ")
		os.Args = append([]string{os.Args[0]}, args...)

		// Anything you print here will be passed back to the cmd.Stderr and
		// cmd.Stdout below, for example:
		fmt.Printf("os args = %v\n", os.Args)

		// Strange, I was expecting a need to manually call the code in
		// `init()`,but that seem to happen automatically. So yet more I have learn.
		main()
	}

	var tests = []struct {
		name     string
		wantCode int
		args     []string
	}{
		{"versionFlag", 0, []string{"-version"}},
		{"helpFlag", 0, []string{"-help"}},
		{
			"localTemplate",
			0,
			[]string{
				"-answer-path", fixturesDir + PS + "answers-parse-dir-02.json",
				"-tmpl-path", fixturesDir + PS + "parse-dir-02",
				"-out-path", testTmp + PS + "app-parse-dir-02",
				"-tmpl-type", "dir",
			},
		},
		{
			"downloadZipTemplate",
			0,
			[]string{
				"-tmpl-path", "https://github.com/kohirens/tmpl-go-web/archive/refs/tags/0.3.0.zip",
				"-out-path", testTmp + PS + "tmpl-go-web-02",
				"-answer-path", fixturesDir + PS + "answers-tmpl-go-web.json",
				"-tmpl-type", "zip",
			},
		},
		{
			"remoteGitTemplate",
			0,
			[]string{
				"-tmpl-path", "https://github.com/kohirens/tmpl-go-web.git",
				"-out-path", testTmp + PS + "tmpl-go-web-03",
				"-answer-path", fixturesDir + PS + "answers-tmpl-go-web.json",
				"-tmpl-type", "git",
				"-branch", "refs/tags/0.3.0",
			},
		},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			cmd := runMain(tester.Name(), test.args)

			out, sce := cmd.CombinedOutput()

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}

			if sce != nil {
				fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n\n", string(out))
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}
		})
	}
}

// getTestBinCmd return a command to run the test binary in a sub-process, passing it flags as fixtures to produce expected output; `TestMain`, will be run automatically.
func getTestBinCmd(args []string) *exec.Cmd {
	// call the generated test binary directly
	// Have it the function runAppMain.
	cmd := exec.Command(os.Args[0])
	// Run in the context of the source directory.
	_, filename, _, _ := runtime.Caller(0)
	cmd.Dir = path.Dir(filename)
	// Set an environment variable
	// 1. Only exist for the life of the test that calls this function.
	// 2. Passes arguments/flag to your app
	// 3. Lets TestMain know when to run the main function.
	subEnvVar := subCmdFlags + "=" + strings.Join(args, " ")
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}

// runAppMain run main passing only the fixture flags
func runAppMain() {
	// Get fixture flags set in the unit test.
	args := strings.Split(os.Getenv(subCmdFlags), " ")
	// replace the test flaga and arts with the test fixture flags and args.
	os.Args = append([]string{os.Args[0]}, args...)

	// Debug stmt
	//fmt.Printf("\nsub os.Args = %v\n", os.Args)

	main()
}

// runMain execute main in a sub process
func runMain(testFunc string, args []string) *exec.Cmd {
	// Run the test binary and tell it to run just this test with environment set.
	cmd := exec.Command(os.Args[0], "-test.run", testFunc)

	subEnvVar := SubCmdFlags + "=" + strings.Join(args, " ")
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sOut := os.Stdout
	sErr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sOut
		os.Stderr = sErr
		log.SetOutput(os.Stderr)
	}
}

func setupARepository(bundleName string) string {
	repoPath := testRemotes + PS + bundleName

	// It may have already been unbundled.
	fileInfo, err1 := os.Stat(repoPath)
	if (err1 == nil && fileInfo.IsDir()) || os.IsExist(err1) {
		absPath, e2 := filepath.Abs(repoPath)
		if e2 == nil {
			return absPath
		}
		return repoPath
	}

	srcRepo := "." + PS + fixturesDir + PS + bundleName + ".bundle"

	// It may not exist.
	if !stdlib.PathExist(srcRepo) {
		return bundleName
	}

	cmd := exec.Command("git", "clone", "-b", "main", srcRepo, repoPath)
	_, _ = cmd.CombinedOutput()
	if ec := cmd.ProcessState.ExitCode(); ec != 0 {
		log.Panicf("error un-bundling %q to a temporary repo %q for a unit test", srcRepo, repoPath)
	}

	absPath, e2 := filepath.Abs(repoPath)
	if e2 == nil {
		return absPath
	}

	return repoPath
}
