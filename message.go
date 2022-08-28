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
	madeNewConfig:         "made a new config file %q exist",
	numNonFlagArgs:        "number of non-flag arguments passed in: %d",
	numParsedFlags:        "number of parsed flags = %v",
	outPathExist:          "out-path already exits %q",
	pleaseAnswerQuestions: "please answer the following\nnote: providing no value will render the placeholder with no valued when processing is done",
	printAllFlags:         "printing all flags set:",
	printFlag:             "\t%s = %v (default= %v)\n",
	questionAnsweredWith:  "%q was answered with %q",
	questionAnswerStat:    "there are %v placeholders and %v values",
	questionHasAnAnswer:   "question %q already has an answer of %q, so skipping",
	refInfo:               "ref = %v ",
	runningCommand:        "running command %s",
	subCommands:           "\n  sub-commands:\n\n    none\n",
	verboseLevelInfo:      "verbose level: %v",
}
