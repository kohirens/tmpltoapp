package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kohirens/stdlib"
)

const (
	PS       = string(os.PathSeparator)
	CONF     = "config.json"
	DIR_MODE = 0774
)

var (
	appConfig = &Config{}
	buildVersion string
	flagStore *flagStorage
	errMsgs = [...]string{
		"please specify a path (or URL) to a template",
		"enter a local path to output the app",
		"the following error occurred trying to get the app data directory: %q",
		"path/URL to template is not in the allow-list",
		"template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
		"please specify a path to an answer file",
	}
	programName string
)

func init() {
	// Use `path/filepath.Base` for cross-platform compatibility.
	programName = filepath.Base(os.Args[0])
	fs, err := defineFlags(programName, flag.ContinueOnError)
	if err != nil {
		return
	}

	flagStore = fs
}

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Fatalln(err)
		}
	}()

	err = flagStore.Flags.Parse(os.Args[1:])
	if err != nil {
		return
	}

	help, _ := flagStore.GetBool("help")
	if help {
		flagStore.Flags.SetOutput(os.Stdout)
		flagStore.Flags.Usage()
		os.Exit(0)
	}

	version, _ := flagStore.GetBool("version")
	if version {
		flagStore.Flags.SetOutput(os.Stdout)
		fmt.Printf("\n%v\n", buildVersion)
		os.Exit(0)
	}

	err = extractParsedFlags(flagStore, os.Args, appConfig)
	if err != nil {
		return
	}

	verboseF(2, "running program %q", programName)
	verboseF(1, "verbose level: %v", verbosityLevel)

	appDataDir, err := stdlib.AppDataDir()
	if err != nil {
		return
	}

	appDataDir = appDataDir + PS + "go-gitter"
	err = os.MkdirAll(appDataDir, DIR_MODE)
	if err != nil {
		return
	}

	configFile := appDataDir + PS + "config.json"

	err = initConfigFile(configFile)
	if err != nil {
		return
	}

	options, err := settings(configFile)
	if err != nil {
		return
	}

	verboseF(1, "config location %q", configFile)

	isUrl, isAllowed := urlIsAllowed(options.tplPath, options.allowedUrls)
	if isUrl && !isAllowed {
		err = fmt.Errorf(errMsgs[3])
		return
	}
	verboseF(1, "isUrl %v", isUrl)

	options.cacheDir = appDataDir + PS + "cache"
	err = os.MkdirAll(options.cacheDir, DIR_MODE)
	if err != nil {
		err = fmt.Errorf("could not make cache directory, error: %s", err.Error())
		return
	}

	tmplPathType := getPathType(options.tplPath)

	if tmplPathType == "http" {
		client := http.Client{}
		zipFile, iErr := download(options.tplPath, options.cacheDir, &client)
		if iErr != nil {
			err = iErr
			return
		}

		iErr = extract(zipFile, os.TempDir())
		if iErr != nil {
			err = iErr
			return
		}
	}
}

// Check to see if a URL is in the allowed list to download template from.
func urlIsAllowed(loc string, urls []string) (isUrl, isAllowed bool) {
	isUrl = strings.HasPrefix(loc, "https://")
	isAllowed = false

	if isUrl {
		for _, url := range urls {
			if strings.HasPrefix(loc, url) {
				isAllowed = true
				break
			}
		}
	}

	return
}
