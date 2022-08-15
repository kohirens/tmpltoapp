package main

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestFlagExitCode(t *testing.T) {

	// This was adapted from https://golang.org/src/flag/flag_test.go; line 596-657 at the time.
	// This is called recursively, because we will have this test call
	// itself in a sub-command. A call to `main()` MUST exit or
	// you'll spin out of control.
	if os.Getenv("GO_CHILD_FLAG") != "" {
		// We re in the test binary, so test flags are set, lets reset it so
		// so that only the program is set
		// and whatever flags we want.
		os.Args = []string{os.Args[0], os.Getenv("GO_CHILD_FLAG")}

		// Anything you print here will be passed back to the cmd.Stderr and
		// cmd.Stdout below, for example:
		//fmt.Printf("os args = %v", os.Args)

		// Strange, I was expecting to have to call the content of init(),
		// but that seem to happen automatically. So yet more I have learn.
		main()
	}

	tests := []struct {
		flag       string
		flagHandle string
		expectExit int
	}{
		{
			flag:       "-h",
			expectExit: 0,
		},
		{
			flag:       "-help",
			expectExit: 0,
		},
		{
			flag:       "-version",
			expectExit: 0,
		},
		{
			flag:       "-undefined",
			expectExit: 2,
		},
	}

	for _, test := range tests {
		// Run the test binary and tell it to run just this test with
		// environment set. Update this string if you change the function name.
		cmd := exec.Command(os.Args[0], "--test.run", "TestFlagExitCode")
		cmd.Env = append(
			os.Environ(),
			"GO_CHILD_FLAG="+test.flag,
			"GO_CHILD_FLAG_HANDLE="+test.flagHandle,
		)

		// Uncomment when you want to debug
		//cmd.Stderr = os.Stderr
		//cmd.Stdout = os.Stdout

		cmd.Run()

		got := cmd.ProcessState.ExitCode()

		// ExitCode is either 0 or 1 on Plan 9.
		if runtime.GOOS == "plan9" && test.expectExit != 0 {
			test.expectExit = 1
		}

		if got != test.expectExit {
			t.Errorf("unexpected exit code for test case %+v \n: got %d, expect %d",
				test, got, test.expectExit)
		}
	}
}
