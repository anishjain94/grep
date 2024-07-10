package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func handleError(err error) {
	if err != nil {
		errMsg := err.Error()
		if os.IsExist(err) {
			errMsg = "file already exists. Please delete and try again."
		}
		println(errMsg)
		os.Exit(1)
	}
}

func fileValidations(filepath string) {
	_, err := os.Stat(filepath)

	if err != nil && os.IsNotExist(err) {
		println(filepath + " No such file")
		os.Exit(1)
	} else if err != nil && os.IsPermission(err) {
		println(filepath + " Permission denied")
		os.Exit(1)
	}

}

func readDataAndMatch(r io.Reader, sourceName *string, flagConfig *FlagConfig, searchStr string) []string {
	var outputLines []string
	reader := bufio.NewScanner(r)

	var directory string
	if sourceName != nil {
		directory = (*sourceName) + ": "
	}

	for reader.Scan() {
		inputStr := directory + reader.Text()

		if flagConfig != nil && flagConfig.isFlagIEnabled() {
			if strings.Contains(strings.ToLower(inputStr), strings.ToLower(searchStr)) {
				outputLines = append(outputLines, inputStr)
			}
		} else {
			if strings.Contains(inputStr, searchStr) {
				outputLines = append(outputLines, inputStr)
			}
		}

	}
	return outputLines
}

type FlagConfig struct {
	FlagI bool   //case-inSensitive search
	FlagO string //output file
}

func (flagConfig *FlagConfig) isFlagIEnabled() bool {
	return flagConfig.FlagI
}

func (flagConfig *FlagConfig) isFlagOEnabled() bool {
	return flagConfig.FlagO != ""
}
