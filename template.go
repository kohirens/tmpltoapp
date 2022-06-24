package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/kohirens/stdlib"
)

const (
	MAX_TPL_SIZE  = 1e+7
	TMPL_MANIFEST = "template.json"
	EMPTY_FILE    = ".empty"
)

type Client interface {
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
}

type tplVars map[string]string // TODO: change to tmplVars for consistency

var regExpTmplType = regexp.MustCompile("^(zip|git|dir)$")

var regExpRelativePath = regexp.MustCompile(`^(\.\.|\.|~)(/[a-zA-Z/._\-].*)?`)
var regExpAbsolutePath = regexp.MustCompile(`^/([a-zA-Z._\-][a-zA-Z/._\-].*)?`)
var regExpWinDrive = regexp.MustCompile(`^[a-zA-Z]:\\[a-zA-Z/._\\-].*$`)

// getTmplLocation Check if the path is an HTTP or local directory URL.
func getTmplLocation(tmplPath string) string {
	pathType := "remote"

	if regExpAbsolutePath.MatchString(tmplPath) || regExpRelativePath.MatchString(tmplPath) || regExpWinDrive.MatchString(tmplPath) {
		pathType = "local"
	}

	return pathType
}

// download a template from a URL to a local directory.
func download(url, dstDir string, client Client) (zipFile string, err error) {
	// Save to a unique filename in the cache.
	dest := strings.ReplaceAll(url, "https://", "")
	dest = strings.ReplaceAll(dest, "/", "-")
	// HTTP Request
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		err = fmt.Errorf(errors.tmplPath, resp.Status, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	zipFile = dstDir + PS + dest
	// make handle to the file.
	out, err := os.Create(zipFile)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	infof("downloading %v to %v\n", url, dest)

	return
}

func extract(archivePath string) (string, error) {
	tmplDir := ""
	zipParentDir := ""
	dest := strings.ReplaceAll(archivePath, ".zip", "")
	// Get resource to zip archive.
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return tmplDir, fmt.Errorf("could not open archive %q, error: %v", archivePath, err.Error())
	}

	err = os.MkdirAll(dest, DIR_MODE)
	if err != nil {
		return tmplDir, fmt.Errorf("could not write dest %q, error: %v", dest, err.Error())
	}

	infof("extracting %v to %v\n", archivePath, dest)
	for _, file := range archive.File {
		sourceFile, fErr := file.Open()
		if fErr != nil {
			return tmplDir, fmt.Errorf("failed to extract archive %q to dest %q, error: %v", archivePath, dest, file.Name)
		}

		extractionDir := filepath.Join(dest, file.Name)
		// trying to figure out the
		if zipParentDir == "" {
			// TODO: Document the fact that template archives MUST be zip format and contain all template files in a single directory at the root of the zip.
			zipParentDir = extractionDir
		}

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(extractionDir, filepath.Clean(dest)+PS) {
			return tmplDir, fmt.Errorf("illegal file path: %s", extractionDir)
		}

		if file.FileInfo().IsDir() {
			ferr := os.MkdirAll(extractionDir, file.Mode())
			if ferr != nil {
				return tmplDir, ferr
			}
		} else {
			dh, ferr := os.OpenFile(extractionDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if ferr != nil {
				return tmplDir, ferr
			}

			_, ferr = io.Copy(dh, sourceFile)
			if ferr != nil {
				return tmplDir, ferr
			}

			ferr = dh.Close()
			if ferr != nil {
				panic(ferr)
			}
		}

		fErr = sourceFile.Close()
		if fErr != nil {
			return tmplDir, fmt.Errorf("unsuccessful extracting archive %q, error: %v", archivePath, fErr.Error())
		}
	}

	err = archive.Close()
	tmplDir = zipParentDir
	dbugf("zipParentDir = %v", zipParentDir)

	return tmplDir, nil
}

// parse a file as a Go template.
func parse(tplFile, dstDir string, vars tplVars) error {

	parser, err1 := template.ParseFiles(tplFile)

	if err1 != nil {
		return err1
	}

	fileStats, err2 := os.Stat(tplFile)

	if err2 != nil {
		return err2
	}

	dstFile := dstDir + PS + fileStats.Name()
	file, err3 := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, fileStats.Mode())

	if err3 != nil {
		return err3
	}

	if e := parser.Execute(file, vars); e != nil {
		return e
	}

	if e := file.Close(); e != nil {
		return e
	}

	return nil
}

// parseDir Recursively walk a directory parsing all files along the way as Go templates.
func parseDir(tplDir, outDir string, vars tplVars, fec *stdlib.FileExtChecker, excludes []string) (err error) {
	// Normalize the path separator in these 2 variables before comparing them.
	normTplDir := strings.ReplaceAll(tplDir, "/", PS)
	normTplDir = strings.ReplaceAll(normTplDir, "\\", PS)

	// Recursively walk the template directory.
	err = filepath.Walk(normTplDir, func(sourcePath string, fi os.FileInfo, wErr error) (rErr error) {
		if wErr != nil {
			rErr = wErr
			return
		}

		infof("processing: %q\n", sourcePath)

		// Do not parse directories.
		if fi.IsDir() {
			return
		}

		// Stop processing files if a template file is too big.
		if fi.Size() > MAX_TPL_SIZE {
			rErr = fmt.Errorf("template file too big to parse, must be less thatn %v bytes", MAX_TPL_SIZE)
			return
		}

		currFile := filepath.Base(sourcePath)
		// Skip non-text files.
		// TODO: Remove FileExtensionCheck in favor of exclude/include list, once globbing is added.
		if currFile != EMPTY_FILE && !fec.IsValid(sourcePath) { // Use an exclusion list, include every file by default.
			infof("will skipp and not process through template engine; could not detect file type for %v", sourcePath)
			return
		}

		// Normalize the path separator in these 2 variables before comparing them.
		normSourcePath := strings.ReplaceAll(sourcePath, "/", PS)
		normSourcePath = strings.ReplaceAll(normSourcePath, "\\", PS)

		// Get the subdirectory from the template source and append it to the
		// output directory, so that files are placed in the correct
		// subdirectories in the output directory.
		partial := strings.ReplaceAll(normSourcePath, normTplDir, "")
		//partial = strings.ReplaceAll(partial, PS, "/")
		saveDir := filepath.Clean(outDir + filepath.Dir(partial))

		// TODO: Make the subdirectories in the new savePath.
		err = os.MkdirAll(saveDir, DIR_MODE)
		if err != nil || currFile == EMPTY_FILE {
			return
		}

		if excludes != nil {
			// TODO: Replace with better method of comparing files.
			fileToCheck := strings.ReplaceAll(normSourcePath, normTplDir, "")
			fileToCheck = strings.ReplaceAll(fileToCheck, PS, "")
			for _, exclude := range excludes {
				fileToCheckB := strings.ReplaceAll(exclude, "\\", "")
				fileToCheckB = strings.ReplaceAll(exclude, "/", "")
				infof("fileToCheck: %q; exclude: %v\n", fileToCheck, fileToCheckB)
				if fileToCheckB == fileToCheck {
					infof("excluding file %q", sourcePath)
					return
				}
			}
		}

		rErr = parse(sourcePath, saveDir, vars)

		return
	})

	return
}

// TODO Rename to manifest or tmplConfig
type questions struct {
	Version   string   `json:"version"`
	Variables tplVars  `json:"variables"`
	Excludes  []string `json:"excludes"`
}

// readTemplateJson read variables needed from the template.json file.
func readTemplateJson(filePath string) (*questions, error) {
	dbugf("\ntemplate manifest path: %q\n", filePath)

	// Verify the TMPL_MANIFEST file is present.
	if !stdlib.PathExist(filePath) {
		return nil, fmt.Errorf(errs.tmplManifest404, TMPL_MANIFEST)
	}

	content, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		return nil, err1
	}

	infof("content = %s \n", content)

	q := questions{}
	if err2 := json.Unmarshal(content, &q); err2 != nil {
		return nil, err2
	}

	infof("content = %v \n", content)

	return &q, nil
}

// getInput Take user input for template variables.
func getInput(questions *questions, answers *tplVars, r *os.File) error {
	nPut := bufio.NewScanner(r)

	for v, q := range questions.Variables {
		a, isAnswered := (*answers)[v]
		if isAnswered {
			infof("question %q already has an answer of %q, so skipping\n", q, a)
			continue
		}
		infof("\n%q: ", q)
		nPut.Scan()
		(*answers)[v] = nPut.Text()
		infof("%q was answered with %q", q, (*answers)[v])
	}

	return nil
}
