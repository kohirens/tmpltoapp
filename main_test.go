package main

import (
	"fmt"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Only runs when this environment variable is set.
	if _, ok := os.LookupEnv(test.SubCmdFlags); ok {
		runAppMain()
	}
	// call flag.Parse() here if TestMain uses flags
	_ = os.RemoveAll(test.TmpDir)
	// Set up a temporary dir for generate files
	_ = os.Mkdir(test.TmpDir, cli.DirMode) // set up a temporary dir for generate files
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
	if os.Getenv(test.SubCmdFlags) != "" {
		// We're in the test binary, so test flags are set, lets reset it
		// so that only the program is set
		// and whatever flags we want.
		args := strings.Split(os.Getenv(test.SubCmdFlags), " ")
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
				"-answer-path", test.FixturesDir + cli.PS + "answers-parse-dir-02.json",
				"-tmpl-path", test.FixturesDir + cli.PS + "parse-dir-02",
				"-out-path", test.TmpDir + cli.PS + "app-parse-dir-02",
				"-tmpl-type", "dir",
			},
		},
		{
			"downloadZipTemplate",
			0,
			[]string{
				"-tmpl-path", "https://github.com/kohirens/tmpl-go-web/archive/refs/tags/0.3.0.zip",
				"-out-path", test.TmpDir + cli.PS + "tmpl-go-web-02",
				"-answer-path", test.FixturesDir + cli.PS + "answers-tmpl-go-web.json",
				"-tmpl-type", "zip",
			},
		},
		{
			"remoteGitTemplate",
			0,
			[]string{
				"-tmpl-path", "https://github.com/kohirens/tmpl-go-web.git",
				"-out-path", test.TmpDir + cli.PS + "tmpl-go-web-03",
				"-answer-path", test.FixturesDir + cli.PS + "answers-tmpl-go-web.json",
				"-tmpl-type", "git",
				"-branch", "refs/tags/0.3.0",
			},
		},
		{"manifest0", 0, []string{"manifest", "-h"}},
		{"manifest0", 1, []string{"manifest"}},
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

			if testing.Verbose() && sce != nil {
				fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n\n", string(out))
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}
		})
	}
}

// runAppMain run main passing only the fixture flags
func runAppMain() {
	// Get fixture flags set in the unit test.
	args := strings.Split(os.Getenv(test.SubCmdFlags), " ")
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

	subEnvVar := test.SubCmdFlags + "=" + strings.Join(args, " ")
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}
