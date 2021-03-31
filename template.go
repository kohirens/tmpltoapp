package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kohirens/stdlib"
)

const (
	MAX_TPL_SIZE = 1e+7
)

type Client interface {
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
}

// Copy a source directory to another destination directory.
func CopyDir(srcDir, dstDir string) (err error) {
	// TODO: Why not just use the OS to copy the files over!?
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return
	}

	err = os.MkdirAll(dstDir, 0774)
	if err != nil {
		return
	}

	for _, file := range files {
		srcPath := srcDir + PS + file.Name()

		if file.IsDir() {
			ferr := CopyDir(srcPath, dstDir+PS+file.Name())
			if ferr != nil {
				err = ferr
				return
			}

			continue
		}

		srcR, ferr := os.Open(srcPath)
		if ferr != nil {
			err = ferr
			break
		}

		dstPath := dstDir + PS + file.Name()
		fmt.Printf("copy %q  to %q ", srcPath, dstPath)
		dstW, ferr := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, file.Mode())
		if ferr != nil {
			err = ferr
			break
		}

		written, ferr := io.Copy(dstW, srcR)
		if ferr != nil {
			err = ferr
			break
		}

		// check file written matches the original file size.
		if written != file.Size() {
			err = fmt.Errorf("failed to copy file correctly, wrote %v, should have written %v", written, file.Size())
		}

		ferr = srcR.Close()
		if ferr != nil {
			err = fmt.Errorf("CopyDir could not close the source file: %q", srcPath)
			break
		}

		ferr = dstW.Close()
		if ferr != nil {
			err = fmt.Errorf("CopyDir could not close the destination file: %q", dstPath)
			break
		}
	}

	return
}

// Download a template from a URL to a local directory.
func Download(url, dstDir string, client Client) (err error) {
	dest := path.Base(url)
	// HTTP Request
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		err = fmt.Errorf(errMsgs[0], resp.Status, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	// make handle to the file.
	out, err := os.Create(dstDir + PS + dest)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	fmt.Printf("downloading %v to %v\n", url, dest)

	return
}

func Extract(archivePath, dest string) (err error) {
	// Get resource to zip archive.
	archive, err := zip.OpenReader(archivePath)

	if err != nil {
		err = fmt.Errorf("could not open archive %q, error: %v", archivePath, err.Error())
		return
	}

	err = os.MkdirAll(dest, 0774)
	if err != nil {
		err = fmt.Errorf("could not write dest %q, error: %v", dest, err.Error())
		return
	}

	for _, file := range archive.File {
		sourceFile, ferr := file.Open()

		if ferr != nil {
			err = fmt.Errorf("failed to extract archive %q to dest %q, error: %v", archivePath, dest, file.Name)
			break
		}

		deflateFilePath := filepath.Join(dest, file.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(deflateFilePath, filepath.Clean(dest)+PS) {
			err = fmt.Errorf("illegal file path: %s", deflateFilePath)
			return
		}

		if file.FileInfo().IsDir() {
			ferr := os.MkdirAll(deflateFilePath, file.Mode())
			if ferr != nil {
				err = ferr
				return
			}
		} else {
			dh, ferr := os.OpenFile(deflateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if ferr != nil {
				err = ferr
				return
			}

			_, ferr = io.Copy(dh, sourceFile)
			if ferr != nil {
				err = ferr
				return
			}

			ferr = dh.Close()
			if ferr != nil {
				panic(ferr)
			}
		}

		ferr = sourceFile.Close()
		if ferr != nil {
			err = fmt.Errorf("unsuccessful extracting archive %q, error: %v", archivePath, ferr.Error())
		}
	}

	archive.Close()

	return
}

type tplVars map[string]string

func Parse(tplFile, dstDir string, vars tplVars) (err error) {

	parser, err := template.ParseFiles(tplFile)

	if err != nil {
		return
	}

	fileStats, err := os.Stat(tplFile)

	if err != nil {
		return
	}

	dstFile := dstDir + PS + fileStats.Name()
	file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, fileStats.Mode())

	if err != nil {
		return
	}

	err = parser.Execute(file, vars)

	return
}

func ParseDir(tplDir, outDir string, vars tplVars) (err error) {
	// Recursively walk the template directory.
	filepath.Walk(tplDir, func(path string, fi os.FileInfo, wErr error) (rErr error) {
		if wErr != nil {
			rErr = wErr
			return
		}

		// Skip directories.
		if fi.IsDir() {
			return
		}

		// Stop processing files if a template file is too big.
		if fi.Size() > MAX_TPL_SIZE {
			rErr = fmt.Errorf("Template file too big to parse, must be less thatn %v bytes.", MAX_TPL_SIZE)
			return
		}

		// Skip non-text files.
		if stdlib.IsTextFile(path) {
			rErr = fmt.Errorf("could not detect file type for %v", path)
			return
		}

		rErr = Parse(path, outDir, vars)
		return
	})

	return
}
