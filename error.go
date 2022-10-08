package main

var errors = struct {
	answerFile404          string
	appDataDir             string
	badExcludeFileExt      string
	badTmplType            string
	cannotDecodeAnswerFile string
	cannotInitFileChecker  string
	cannotReadAnswerFile   string
	cloning                string
	couldNot               string
	couldNotCloseFile      string
	couldNotDecode         string
	couldNotEncodeConfig   string
	couldNotMakeCacheDir   string
	couldNotSaveConf       string
	couldNotWriteFile      string
	fatalHeader            string
	flagOrderErr           string
	fileTooBig             string
	gettingAnswers         string
	gettingCommitHash      string
	gitCheckoutFailed      string
	gitFetchFailed         string
	gitExitErrCode         string
	getLatestTag           string
	getRemoteTags          string
	invalidNoArgs          string
	invalidTmplDir         string
	localOutPath           string
	missingTmplJson        string
	noGitTagFound          string
	parsingConfigArgs      string
	pathNotAllowed         string
	runGitFailed           string
	tmplManifest404        string
	tmplOutput             string
	tmplPath               string
	unhandledHttpErr       string
}{
	answerFile404:          "could not find the answer file %q, please specify a path to a valid answer file that exist: given %q",
	appDataDir:             "the following error occurred trying to get the app data directory: %q",
	badExcludeFileExt:      "invalid ExcludeFileExtensions, check format, for example: item1,item2,item3",
	badTmplType:            "%q is an invalid value for flag tmplType, or it was not set, must be zip|git",
	cannotDecodeAnswerFile: "could not decode JSON in answer file %q, because of: %s",
	cannotInitFileChecker:  "cannot instantiate file extension checker: %v",
	cannotReadAnswerFile:   "there was an error reading the answer file %q: %s",
	cloning:                "error cloning %v: %s",
	couldNot:               "could not %s",
	couldNotCloseFile:      "could not close file %v, %v",
	couldNotDecode:         "could not decode %q, error: %s",
	couldNotEncodeConfig:   "could not JSON encode user configuration settings, %v",
	couldNotMakeCacheDir:   "could not make cache directory, error: %s",
	couldNotSaveConf:       "could not save a config file, reason: %v",
	couldNotWriteFile:      "could not write file %v, reason: %v",
	fatalHeader:            "\nfatal error detected: ",
	flagOrderErr:           "flag %v MUST come before any non-flag arguments, a fix would be to move this flag to the left of other input arguments",
	fileTooBig:             "template file too big to parse, must be less thatn %v bytes",
	gettingAnswers:         "problem getting answers; error %q",
	gettingCommitHash:      "error getting commit hash %v: %s",
	gitCheckoutFailed:      "git checkout failed: %s",
	getLatestTag:           "failed to get latest tag from %v: %v",
	getRemoteTags:          "could not get remote tags, please check for a typo, it exist, and is readable: %v",
	gitExitErrCode:         "git %v returned exit code %q",
	gitFetchFailed:         "fetch failed on %s and %s; %s",
	invalidNoArgs:          "invalid number of arguments passed to config sub-command, please try config -h for usage",
	invalidTmplDir:         "invalid template directory %q",
	localOutPath:           "enter a local path to output the app",
	missingTmplJson:        "%s is a file that is required to be in the template, there was a problem reading %q; error %q",
	noGitTagFound:          "no tag found in %v",
	parsingConfigArgs:      "error parsing config command args: %v",
	pathNotAllowed:         "path/URL to template is not in the allow-list",
	runGitFailed:           "error running git %v: %v\n%s",
	tmplManifest404:        "the required manifest template.json file was not %s found",
	tmplOutput:             "template has NOT been cloned locally",
	tmplPath:               "please specify a path (or URL) to a template",
	unhandledHttpErr:       "template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}
