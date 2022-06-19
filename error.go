package main

var errors = struct {
	answerPath       string
	appDataDir       string
	badTmplType      string
	localOutPath     string
	pathNotAllowed   string
	tmplPath         string
	unhandledHttpErr string
}{
	answerPath:       "please specify a path to a valid answer file that exist: given %q",
	appDataDir:       "the following error occurred trying to get the app data directory: %q",
	badTmplType:      "tmplType type was not set, must be zip|git",
	localOutPath:     "enter a local path to output the app",
	pathNotAllowed:   "path/URL to template is not in the allow-list",
	tmplPath:         "please specify a path (or URL) to a template",
	unhandledHttpErr: "template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}
