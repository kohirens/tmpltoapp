package template

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var errMsgs = [...]string{
	"path/URL to template is not in the allow-list",
}

type Client interface {
	Get(url string) (*http.Response, error)
}

/* Download a template from a URL. */
func Download(url, dest string, client Client) error {

	if !urlIsAllowed() {
		return fmt.Errorf(errMsgs[0])
	}

	// Request
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// make handle to the file.
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	fmt.Printf("downloading %v to %v\n", url, dest)

	return err
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

func urlIsAllowed() bool {
	return true
}
