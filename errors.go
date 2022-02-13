package main

var errs = struct {
	gettingAnswers,
	missingTmplJson,
	tmplManifest404 string
}{
	gettingAnswers:  "problem getting answers; error %q",
	missingTmplJson: "%s is a file that is required to be in the template, there was a problem reading %q; error %q",
	tmplManifest404: "not %s found",
}
