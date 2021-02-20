package stdlib

import (
	"os"
	"path"
)

const (
	PS = string(os.PathSeparator)
)

var textFileTypes = [4]string{
	".json",
	".md",
	".txt",
	".xml",
}

//  Returns true for files that match the text extensions.
func IsTextFile(file string) (ret bool) {
	ret = false

	ext := path.Ext(file)

	if ext != "" {
	txtcompare: // sorry, I just wanted to play with this so I get used to it. Even though this is single loop or I could just use return. I like to be explicit.
		for _, t := range textFileTypes {
			if t == ext {
				ret = true
				break txtcompare
			}
		}
	}

	return
}

func PathExist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
