package template

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

var errMsgs = [...]string{
	"template download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
}

type Client interface {
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
}

/* Download a template from a URL. */
func Download(url string, client Client) error {

	dest := path.Base(url)
	// Request
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode > 300 || resp.StatusCode < 200 {
		return fmt.Errorf(errMsgs[0], resp.Status, resp.StatusCode)
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
