package stdlib

import (
	"fmt"
	"os"
	"runtime"
)

var errMsgs = [...]string{
	"cannot get user home dir",
	"cannot get app data dir",
}

/* Home directory of the user your application runs under. */
func HomeDir() (homeDir string, err error) {

	// Linux/Mac
	homeDir = os.Getenv("HOME")

	if runtime.GOOS == "windows" {
		homeDir = fmt.Sprintf("%s%v%s", os.Getenv("HOMEDRIVE"), os.PathSeparator, os.Getenv("HOMEPATH"))
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE")
		}
	}

	if homeDir == "" {
		err = fmt.Errorf(errMsgs[0])
	}

	return
}

/* Base location where you store application configuration data. */
func AppDataDir() (dataDir string, err error) {

	// Linux/Mac
	dataDir, err = HomeDir()

	if err != nil {
		err = fmt.Errorf("%s, reason: %s", errMsgs[1], err.Error())
		return
	}

	if runtime.GOOS == "windows" {
		dataDir = dataDir + "\\AppData\\Local"
	}

	return
}
