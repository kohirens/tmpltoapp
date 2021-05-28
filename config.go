package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	answers tplVars
	allowedUrls    []string
	answersPath    string
	appPath        string
	cacheDir       string
	tplPath        string
	verbosityLevel int
}

// Load configuration file.
func initConfigFile(file string) (err error) {

	_, er := os.Stat(file)

	if err != nil {
		if !os.IsNotExist(er) {
			verboseF(1, "config file exist %v", er.Error())
			return
		}

		err = er
	}

	f, err := os.Create(file)

	defer f.Close()

	f.WriteString(DEFAULT_CFG)

	return
}

func settings(filename string) (cfg Config, err error) {
	var data map[string]interface{}

	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}

	val, ok := data["allowedUrls"].([]interface{})
	if ok && len(data) > 0 {
		cfg.allowedUrls = make([]string, len(data))
		for i, v := range val {
			cfg.allowedUrls[i] = v.(string)
		}
	}

	return
}

// loadAnswers Loan answers from a JSON file to be used to fill in placeholders when processing Go templates.
func loadAnswers(filename string) (answers tplVars, err error) {
	content, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		err = fmt.Errorf("could not %s", err.Error())
		return
	}

	err = json.Unmarshal(content, &answers)
	if err != nil {
		return
	}

	return
}
