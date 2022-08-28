package main

var usageMsgs = map[string]string{
	"answer-path": "Path to a JSON file containing the values for placeholders (which are the keys) defined by a template.",
	"branch":      "Branch of the template to clone when tmplType=git.",
	"help":        "Prints usage information and exit 0.",
	"out-path":    "Path to output the new project.",
	"tmpl-path":   "URL to a zip or a local path to a directory.",
	"tmpl-type":   "Can be of git|zip|local.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "Print build version information and exit 0.",
}

// messages helpful info to std out
var messages = struct {
	pleaseAnswerQuestions string
	questionAnsweredWith  string
	questionAnswerStat    string
	questionHasAnAnswer   string
	subCommands           string
}{
	pleaseAnswerQuestions: "please answer the following\nnote: providing no value will render the placeholder with no valued when processing is done",
	questionAnsweredWith:  "%q was answered with %q",
	questionAnswerStat:    "there are %v placeholders and %v values",
	questionHasAnAnswer:   "question %q already has an answer of %q, so skipping",
	subCommands:           "\n  sub-commands:\n\n    none\n",
}
