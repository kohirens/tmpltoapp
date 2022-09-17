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
	cloningToCache        string
	configFileExist       string
	currentVersion        string
	currentVersionInfo    string
	gitCheckout           string
	madeNewConfig         string
	numNonFlagArgs        string
	numParsedFlags        string
	outPathExist          string
	provideValues         string
	printAllFlags         string
	printFlag             string
	placeholderAnswer     string
	placeholderAnswerStat string
	placeholderHasAnswer  string
	refInfo               string
	repoDir               string
	runningCommand        string
	skipFile              string
	subCommands           string
	unknownFileType       string
	usingCache            string
	verboseLevelInfo      string
}{
	actualArgs:            "actual arguments passed in: %v",
	cloningToCache:        "no cache; cloning %v to %v",
	configFileExist:       "config file %q exist",
	currentVersion:        "%v, %v\n",
	currentVersionInfo:    "version: %v, %v\n",
	gitCheckout:           "git checkout %s",
	madeNewConfig:         "saved %d bytes to a new config %q",
	numNonFlagArgs:        "number of non-flag arguments passed in: %d",
	numParsedFlags:        "number of parsed flags = %v",
	outPathExist:          "out-path already exits %q",
	provideValues:         "note: entering no value will render the placeholder with an empty string",
	printAllFlags:         "printing all flags set:",
	printFlag:             "\t%s = %v (default= %v)\n",
	placeholderAnswer:     "%v = %q\n",
	placeholderAnswerStat: "please provide values for %v placeholders",
	placeholderHasAnswer:  "placeholder %v has a value of %q, so skipping",
	refInfo:               "ref = %v ",
	repoDir:               "repoDir = %q\n",
	runningCommand:        "running command %s",
	skipFile:              "skipping: %v",
	subCommands:           "sub-commands:\n\n",
	usingCache:            "using cache %v",
	unknownFileType:       "will skip and not process through template engine; could not detect file type for %v",
	verboseLevelInfo:      "verbose level: %v",
}
