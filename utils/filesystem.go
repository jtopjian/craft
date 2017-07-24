package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

func FileGetLines(fileName string) (lines []string, err error) {
	_, err = os.Stat(fileName)
	if err != nil {
		return
	}

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	lines = strings.Split(string(fileContent), "\n")
	return
}

func WriteFile(fileName, content string) (err error) {
	fi, err := os.Stat(fileName)
	if err != nil {
		return
	}
	mode := fi.Mode()

	err = ioutil.WriteFile(fileName, []byte(content), mode)
	if err != nil {
		return
	}

	return
}
