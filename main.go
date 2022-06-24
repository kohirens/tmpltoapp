package main

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
			log.Fatalln(mainErr)
		}
		os.Exit(0)
	}()

	mainErr = flagMain(appConfig)
	if mainErr != nil {
		return
	}

	infof("verbose level: %v", verbosityLevel)

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

	infof("configured runtime options %v", appConfig)

	appConfig.cacheDir = appDataDir + PS + "cache"
	mainErr = os.MkdirAll(appConfig.cacheDir, DIR_MODE)
	if mainErr != nil {
		mainErr = fmt.Errorf("could not make cache directory, error: %s", mainErr.Error())
		return
	}

	if appConfig.tmplType == "zip" {
		client := http.Client{}
		zipFile, iErr := download(appConfig.tplPath, appConfig.cacheDir, &client)
		if iErr != nil {
			mainErr = iErr
			return
		}

		appConfig.tmpl, iErr = extract(zipFile)
		if iErr != nil {
			mainErr = iErr
			return
		}
	}

	tmplPathType, err3 := getTmplType(appConfig.tmplType)
	if err3 != nil {
		mainErr = err3
		return
	}
	if tmplPathType == "local" {
		appConfig.tmpl = filepath.Clean(appConfig.tplPath)
	}

	if appConfig.tmplType == "git" {
		var repo, commitHash string
		var err2 error

		repoDir := appConfig.cacheDir + PS + getRepoDir(appConfig.tplPath)
		infof("repoDir = %q\n", repoDir)

		// Do a pull when the repo already exists. This will fail if it downloaded a zip.
		if stdlib.DirExist(repoDir) {
			infof("pulling latest\n")
			repo, commitHash, err2 = gitCheckout(repoDir, appConfig.branch)
		} else {
			repo, commitHash, err2 = gitClone(appConfig.tplPath, repoDir, appConfig.branch)
		}

		infof("repo = %q; %q", repo, commitHash)
		if err2 != nil {
			mainErr = err2
			return
		}
		appConfig.tmpl = repo
	}

	errf("appConfig = %v", appConfig)

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
		mainErr = fmt.Errorf(errs.missingTmplJson, TMPL_MANIFEST, tmplManifestFile, errX.Error())
		return
	}

	appConfig.Questions = *tmplManifest
	appConfig.answers, mainErr = loadAnswers(appConfig.answersPath)
	if mainErr != nil {
		return
	}

	if e := getInput(&appConfig.Questions, &appConfig.answers, os.Stdin); e != nil {
		mainErr = fmt.Errorf(errs.gettingAnswers, e.Error())
	}

	//missingAnswers := checkAnswersToQuestions()

	mainErr = parseDir(appConfig.tmpl, appConfig.appPath, appConfig.answers, fec, tmplManifest.Excludes)
}
