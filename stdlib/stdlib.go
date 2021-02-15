package stdlib

import (
	"os"

	"golang.org/x/tools/godoc/util"
)

const (
	PS = string(os.PathSeparator)
)

// Try to detect if a file is a text file.
func IsTextFile(tplDir string, fi os.FileInfo) (ret bool) {
	file, err := os.OpenFile(tplDir, os.O_RDONLY, fi.Mode())
	ret = false

	if err != nil {
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return
	}

	return util.IsText(buffer)
}

func PathExist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
