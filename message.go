package main

var usageMsgs = map[string]string{
	"answer-path": "Path to a JSON file containing the values for placeholders (which are the keys) defined by a template.",
	"branch":      "Branch of the template to clone when tmplType=git.",
	"default-val": "Set a default value for any un-set placeholders and skip asking for input from the command line, and useful for un-attended automation.",
	"help":        "(or -h) Prints usage information and exit 0.",
	"out-path":    "Path to output the new project.",
	"tmpl-path":   "URL to a zip or a local path to a directory.",
	"tmpl-type":   "Can be of git|zip.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "Print build version information and exit 0.",
}
