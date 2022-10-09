package command

var errors = struct {
	encodingJson   string
	savingManifest string
	parsingFile    string
	pathNotExist   string
}{
	encodingJson:   "could not marshall actions in file %v, error: %v",
	savingManifest: "could not save file %v, error: %v",
	parsingFile:    "could not parse file %v, error: %v",
	pathNotExist:   "could not locate the path %q",
}
