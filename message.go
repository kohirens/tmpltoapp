package main

var usageMsgs = map[string]string{
	"answers":     "Path to an answer file.",
	"out-path":    "Path to output the new project.",
	"branch":      "branch to clone.",
	"help":        "print usage information.",
	"subCommands": "\n  sub-commands:\n\n    none\n",
	"tmplPath":    "URL to a zip or a local path to a directory.",
	"tmplType":    "Can be of git|zip|local.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "print build version information and exit 0.",
}

var messages = struct {
	questionAnsweredWith  string
	questionAnswerStat    string
	questionHasAnAnswer   string
	pleaseAnswerQuestions string
}{
	questionAnsweredWith:  "%q was answered with %q",
	questionAnswerStat:    "there are %v placeholders and %v values",
	questionHasAnAnswer:   "question %q already has an answer of %q, so skipping",
	pleaseAnswerQuestions: "please answer the following\nnote: providing no value will render the placeholder with no valued when processing is done",
}
