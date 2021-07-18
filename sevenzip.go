package main

import (
	"os"
	"os/exec"
	"runtime"
)

func isSevenZipInstalled() (string, error) {
	cmdPath := ""
	cmd := exec.Command("7z", "-i")
	out1, err1 := cmd.CombinedOutput()
	// get exit code.
	ec := cmd.ProcessState.ExitCode()

	if err1 != nil {
		verboseF(verboseLvlDgb, "stdout: %v\n\nstderr:\n%v\n", out1, err1.Error())
	}

	verboseF(verboseLvlDgb, "7zip check exit code : %v\n\n", ec)

	if ec != 0 && runtime.GOOS == "windows" {
		winPath := os.Getenv("ProgramFiles") + PS + "7-zip" + PS + "7z.exe"
		// On Windows, passing it through CMD works, otherwise you get an error.
		cmd := exec.Command("cmd", cmdPath, "-i")
		out2, err2 := cmd.CombinedOutput()
		ec := cmd.ProcessState.ExitCode()

		if err2 != nil {
			verboseF(verboseLvlDgb, "stdout: %v\n\nstderr:\n%v\n", out2, err2.Error())
		}

		if ec == 0 {
			cmdPath = winPath
		}
	}

	verboseF(verboseLvlDgb, "7zip check exit code : %v\n\n", ec)

	if ec == 0 {
		verboseF(verboseLvlInfo, "7zip installation found.")
		return "", nil
	}

	verboseF(verboseLvlInfo, "7zip is not installed.")

	return "", nil
}
