package main

//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName appConfig

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: Change name to tmplpress

const (
	PS       = string(os.PathSeparator)
	DIR_MODE = 0774
)

var (
	appConfig = &Config{} // store all settings (including CLI flag values).
)

func init() {
	// Define all flags
	appConfig.define()
}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			fmt.Print("\nfatal error detected: ")
			log.Fatalln(mainErr.Error())
		}
		os.Exit(0)
	}()

	mainErr = flagMain(appConfig)
	if mainErr != nil {
		return
	}

	appDataDir, mainErr := stdlib.AppDataDir()
	if mainErr != nil {
		return
	}

	// Make a directory for tmpltoapp to store data.
	appDataDir = appDataDir + PS + "tmpltoapp"
	mainErr = os.MkdirAll(appDataDir, DIR_MODE)
	if mainErr != nil {
		return
	}

	// Make a configuration file when there is not one.
	configFile := appDataDir + PS + "config.json"
	mainErr = initConfigFile(configFile)
	if mainErr != nil {
		return
	}

	mainErr = settings(configFile, appConfig)

	if mainErr != nil {
		return
	}

	// TODO: Move to configMain
	appConfig.cacheDir = appDataDir + PS + "cache"
	mainErr = os.MkdirAll(appConfig.cacheDir, DIR_MODE)
	if mainErr != nil {
		mainErr = fmt.Errorf("could not make cache directory, error: %s", mainErr.Error())
		return
	}

	appConfig.tmplLocation = getTmplLocation(appConfig.tmplPath)

	if appConfig.tmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		appConfig.tmpl = filepath.Clean(appConfig.tmplPath)
	}

	if appConfig.tmplType == "zip" {
		var zipFile string
		var iErr error
		zipFile = appConfig.tmplPath
		if appConfig.tmplLocation == "remote" {
			client := http.Client{}
			zipFile, iErr = download(appConfig.tmplPath, appConfig.cacheDir, &client)
			if iErr != nil {
				mainErr = iErr
				return
			}
		}

		appConfig.tmpl, iErr = extract(zipFile)
		if iErr != nil {
			mainErr = iErr
			return
		}
	}

	if appConfig.tmplType == "git" {
		var repo, commitHash string
		var err2 error

		repoDir := appConfig.cacheDir + PS + getRepoDir(appConfig.tmplPath)
		infof("repoDir = %q\n", repoDir)

		// Do a pull when the repo already exists. This will fail if it downloaded a zip.
		if stdlib.DirExist(repoDir + "/.git") {
			infof("pulling latest\n")
			repo, commitHash, err2 = gitCheckout(repoDir, appConfig.branch)
		} else {
			repo, commitHash, err2 = gitClone(appConfig.tmplPath, appConfig.cacheDir, appConfig.branch)
		}

		infof("repo = %q; %q", repo, commitHash)
		if err2 != nil {
			mainErr = err2
			return
		}
		appConfig.tmpl = repo
	}

	if !stdlib.DirExist(appConfig.tmpl) {
		mainErr = fmt.Errorf("invalid template directory %q", appConfig.tmpl)
		return
	}

	fec, err1 := stdlib.NewFileExtChecker(&appConfig.ExcludeFileExtensions, &appConfig.IncludeFileExtensions)
	if err1 != nil {
		mainErr = fmt.Errorf("error instantiating file extension checker: %v", err1.Error())
	}

	// Require template directories to have a specific file in order to be processed to prevent processing directories unintentionally.
	tmplManifestFile := appConfig.tmpl + PS + TMPL_MANIFEST
	tmplManifest, errX := readTemplateJson(tmplManifestFile)
	if errX != nil {
		mainErr = fmt.Errorf(errors.missingTmplJson, TMPL_MANIFEST, tmplManifestFile, errX.Error())
		return
	}

	appConfig.TmplJson = *tmplManifest

	if len(appConfig.answersPath) > 0 { // optionally load the answers.json
		appConfig.answersJson, mainErr = loadAnswers(appConfig.answersPath)
		if mainErr != nil {
			return
		}
	}

	// Checks for any missing placeholder values waits for their input from the CLI.
	if e := getPlaceholderInput(&appConfig.TmplJson, &appConfig.answersJson.Placeholders, os.Stdin); e != nil {
		mainErr = fmt.Errorf(errors.gettingAnswers, e.Error())
	}

	// TODO: showAllQuestionsAndAnswer: Output each question with its answer

	mainErr = parseDir(appConfig.tmpl, appConfig.outPath, appConfig.answersJson.Placeholders, fec, tmplManifest.Excludes)
}
