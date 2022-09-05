package main

var usageMsgs = map[string]string{
	"answer-path": "Path to a JSON file containing the values for placeholders (which are the keys) defined by a template.",
	"branch":      "Branch of the template to clone when tmplType=git.",
	"default-val": "Set a default value for any un-set placeholders and skip asking for input from the command line, and useful for un-attended automation.",
	"help":        "(or -h) Prints usage information and exit 0.",
	"out-path":    "Path to output the new project.",
	"tmpl-path":   "URL to a zip or a local path to a directory.",
	"tmpl-type":   "Can be of git|zip|local.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "Print build version information and exit 0.",
}

// messages helpful info to std out
var messages = struct {
	actualArgs            string
	configFileExist       string
	currentVersion        string
	gitCheckout           string
	madeNewConfig         string
	numNonFlagArgs        string
	numParsedFlags        string
	outPathExist          string
	pleaseAnswerQuestions string
	printAllFlags         string
	printFlag             string
	questionAnsweredWith  string
	questionAnswerStat    string
	questionHasAnAnswer   string
	refInfo               string
	runningCommand        string
	subCommands           string
	verboseLevelInfo      string
}{
	actualArgs:            "actual arguments passed in: %v",
	configFileExist:       "config file %q exist",
	currentVersion:        "%v, %v\n",
	gitCheckout:           "git checkout %s",
	madeNewConfig:         "wrote %d bytes to a new config file %q",
	numNonFlagArgs:        "number of non-flag arguments passed in: %d",
	numParsedFlags:        "number of parsed flags = %v",
	outPathExist:          "out-path already exits %q",
	pleaseAnswerQuestions: "please answer the following\nnote: providing no value will render an empty string in its place",
	printAllFlags:         "printing all flags set:",
	printFlag:             "\t%s = %v (default= %v)\n",
	questionAnsweredWith:  "%q was answered with %q",
	questionAnswerStat:    "there are %v placeholders and %v values",
	questionHasAnAnswer:   "question %q already has an answer of %q, so skipping",
	refInfo:               "ref = %v ",
	runningCommand:        "running command %s",
	subCommands:           "sub-commands:\n\n",
	verboseLevelInfo:      "verbose level: %v",
}
