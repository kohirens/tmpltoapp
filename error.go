package main

var errors = struct {
	answerPath       string
	appDataDir       string
	localOutPath     string
	pathNotAllowed   string
	subCmdMissing    string
	tmplPath         string
	unhandledHttpErr string
}{
	answerPath:       "please specify a path to an answer file that exist",
	appDataDir:       "the following error occurred trying to get the app data directory: %q",
	localOutPath:     "enter a local path to output the app",
	pathNotAllowed:   "path/URL to template is not in the allow-list",
	subCmdMissing:    "missing a sub command, either \"semver\" or \"taggable\" or \"checkConf\"",
	tmplPath:         "please specify a path (or URL) to a template",
	unhandledHttpErr: "template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}
