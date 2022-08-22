package main

var errors = struct {
	answerFile404          string
	appDataDir             string
	badTmplType            string
	cannotReadAnswerFile   string
	cannotDecodeAnswerFile string
	flagOrderErr           string
	gettingAnswers         string
	localOutPath           string
	missingTmplJson        string
	pathNotAllowed         string
	tmplManifest404        string
	tmplOutput             string
	tmplPath               string
	unhandledHttpErr       string
}{
	answerFile404:          "could not find the answer file %q, please specify a path to a valid answer file that exist: given %q",
	appDataDir:             "the following error occurred trying to get the app data directory: %q",
	badTmplType:            "%q is an invalid value for flag tmplType, or it was not set, must be zip|git",
	cannotReadAnswerFile:   "there was an error reading the answer file %q: %s",
	cannotDecodeAnswerFile: "could not decode JSON in answer file %q, because of: %s",
	flagOrderErr:           "flag %v MUST come before any non-flag arguments, a fix would be to move this flag to the left of other input arguments",
	gettingAnswers:         "problem getting answers; error %q",
	localOutPath:           "enter a local path to output the app",
	missingTmplJson:        "%s is a file that is required to be in the template, there was a problem reading %q; error %q",
	pathNotAllowed:         "path/URL to template is not in the allow-list",
	tmplManifest404:        "the required manifest template.json file was not %s found",
	tmplOutput:             "template has NOT been cloned locally",
	tmplPath:               "please specify a path (or URL) to a template",
	unhandledHttpErr:       "template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}
