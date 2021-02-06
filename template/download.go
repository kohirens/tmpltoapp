package template

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	PS = string(os.PathSeparator)
)

var errMsgs = [...]string{
	"template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}

type Client interface {
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
}

// Download a template from a URL to a local directory.
func Download(url, dstDir string, client Client) error {
	dest := path.Base(url)
	// HTTP Request
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		return fmt.Errorf(errMsgs[0], resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()

	// make handle to the file.
	out, err := os.Create(dstDir + PS + dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	fmt.Printf("downloading %v to %v\n", url, dest)

	return err
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
		sourceFile, err := file.Open()

		if err != nil {
			err = fmt.Errorf("failed to extract archive %q to dest %q, error: %v", archivePath, dest, file.Name)
			break
		}

		deflateFilePath := filepath.Join(dest, file.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(deflateFilePath, filepath.Clean(dest)+PS) {
			return fmt.Errorf("illegal file path: %s", deflateFilePath)
		}

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(deflateFilePath, file.Mode())
			if err != nil {
				return err
			}
		} else {
			dh, err := os.OpenFile(deflateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if err != nil {
				return err
			}

			_, err = io.Copy(dh, sourceFile)
			if err != nil {
				return err
			}

			err = dh.Close()
			if err != nil {
				panic(err)
			}
		}

		sourceFile.Close()
		if err != nil {
			err = fmt.Errorf("unsuccessful extracting archive %q, error: %v", archivePath, err.Error())
		}
	}

	archive.Close()

	return
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetTempDir() {
	os.TempDir()
}
