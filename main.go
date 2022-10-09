package main

//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName appConfig

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/tmpltoapp/internal/command"
	"log"
	"net/http"
	"os"
)

// TODO: Change name to tmplpress

const (
	PS         = string(os.PathSeparator)
	DirMode    = 0774
	AppName    = "tmpltoapp"
	gitConfDir = ".git"
)

var (
	// appConfig Runtime settings used throughout the application.
	appConfig = &Config{
		usrOpts: &userOptions{
			ExcludeFileExtensions: &[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"},
		},
	}
)

func init() {
	// Define all flags
	appConfig.defineFlags()
}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			logf(errors.fatalHeader)
			log.Fatalln(mainErr.Error())
		}
		os.Exit(0)
	}()

	mainErr = appConfig.parseFlags()
	if mainErr != nil {
		return
	}

	mainErr = appConfig.setup(AppName, PS, DirMode)
	if mainErr != nil {
		return
	}

	// Exit if we are just printing help usage
	if appConfig.help {
		mainErr = Usage(appConfig)
		return
	}

	infof(messages.currentVersionInfo, appConfig.CurrentVersion, appConfig.CommitHash)

	// process sub-commands
	switch appConfig.subCmd {
	case cmdConfig:
		// store or get the key and return
		mainErr = updateUserSettings(appConfig, DirMode)
		return
	case cmdManifest:
		// store or get the key and return
		fec, _ := stdlib.NewFileExtChecker(&[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"}, &[]string{})
		_, mainErr = command.GenerateATemplateManifest(appConfig.subCmdManifest.path, fec, []string{})
		return
	}

	if appConfig.tmplType == "zip" {
		var zipFile string
		var iErr error
		zipFile = appConfig.tmplPath
		if appConfig.tmplLocation == "remote" {
			client := http.Client{}
			zipFile, iErr = download(appConfig.tmplPath, appConfig.usrOpts.CacheDir, &client)
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

		if appConfig.branch == "latest" {
			latestTag, e3 := getLatestTag(appConfig.tmplPath)
			// This error is informative, but not worth stopping the program.
			logf(e3.Error())
			if latestTag != "" {
				appConfig.branch = latestTag
			}
		}

		// Determine the cache location
		repoDir := appConfig.usrOpts.CacheDir + PS + getRepoDir(appConfig.tmplPath, appConfig.branch)
		infof(messages.outRepoDir, repoDir)

		// Do a pull when the repo already exists. This will fail if it downloaded a zip.
		if stdlib.DirExist(repoDir + PS + gitConfDir) {
			infof(messages.usingCache, repoDir)
			repo, commitHash, err2 = gitCheckout(repoDir, appConfig.branch)
		} else {
			infof(messages.cloningToCache, repoDir)
			repo, commitHash, err2 = gitClone(appConfig.tmplPath, repoDir, appConfig.branch)
		}

		infof(messages.repoInfo, repo, commitHash)
		if err2 != nil {
			mainErr = err2
			return
		}
		appConfig.tmpl = repo
	}

	if !stdlib.DirExist(appConfig.tmpl) {
		mainErr = fmt.Errorf(errors.invalidTmplDir, appConfig.tmpl)
		return
	}

	fec, err1 := stdlib.NewFileExtChecker(appConfig.usrOpts.ExcludeFileExtensions, &[]string{})
	if err1 != nil {
		mainErr = fmt.Errorf(errors.cannotInitFileChecker, err1.Error())
	}

	// Require template directories to have a specific file in order to be processed to prevent processing directories unintentionally.
	tmplManifestFile := appConfig.tmpl + PS + TmplManifest
	tmplManifest, errX := readTemplateJson(tmplManifestFile)
	if errX != nil {
		mainErr = fmt.Errorf(errors.missingTmplJson, TmplManifest, tmplManifestFile, errX.Error())
		return
	}

	appConfig.TmplJson = tmplManifest
	appConfig.answersJson = newAnswerJson()

	if stdlib.PathExist(appConfig.answersPath) {
		appConfig.answersJson, mainErr = loadAnswers(appConfig.answersPath)
		if mainErr != nil {
			return
		}
	}

	// Checks for any missing placeholder values waits for their input from the CLI.
	if e := getPlaceholderInput(appConfig.TmplJson, &appConfig.answersJson.Placeholders, os.Stdin, appConfig.defaultVal); e != nil {
		mainErr = fmt.Errorf(errors.gettingAnswers, e.Error())
	}

	showAllPlaceholderValues(appConfig.TmplJson, &appConfig.answersJson.Placeholders)

	mainErr = parseDir(appConfig.tmpl, appConfig.outPath, appConfig.answersJson.Placeholders, fec, tmplManifest.Excludes)
}
