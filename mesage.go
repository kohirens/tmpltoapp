package main

var usageMsgs = map[string]string{
	"answers":     "Path to an answer file.",
	"appPath":     "Path to output the new project.",
	"branch":      "branch to clone.",
	"help":        "print usage information.",
	"subCommands": "\n  sub-commands:\n\n    semver\n\tsee semver --help\n    taggable\n\tsee taggable --help\n    checkConf\n\tsee checkConf --help\n",
	"tmplPath":    "URL to a zip or a local path to a directory.",
	"tmplType":    "Can be of git|zip|local.",
	"verbosity":   "extra detail processing infof.",
	"version":     "print build version information and exit 0.",
}
