package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func handleError(err error) {
	if err != nil {
		errMsg := err.Error()
		if os.IsExist(err) {
			errMsg = "file already exists. Please delete and try again."
		}
		log.Print(errMsg)
		os.Exit(1)
	}
}

func fileValidations(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil && os.IsNotExist(err) {
		log.Print(filepath + " No such file")
		return err
	} else if os.IsPermission(err) {
		log.Print(filepath + " Permission denied")
		return err
	}

	return nil
}

func readMatchAndStore(r io.Reader, sourceName *string, flagConfig *FlagConfig, searchStr string) []string {
	var outputLines []string
	reader := bufio.NewScanner(r)

	var directory string
	if sourceName != nil {
		directory = (*sourceName) + ": "
	}

	for reader.Scan() {
		inputStr := directory + reader.Text()

		if flagConfig.isFlagIEnabled() {
			inputStr = strings.ToLower(inputStr)
			searchStr = strings.ToLower(searchStr)
		}

		if strings.Contains(inputStr, searchStr) {
			outputLines = append(outputLines, reader.Text())
		}
	}
	return outputLines
}

type FlagConfig struct {
	FlagI bool   //case-inSensitive search
	FlagC bool   //displays count of matches
	FlagA int    //displays n lines before the match
	FlagB int    //displays n lines after the match
	FlagO string //output file
}

func (flagConfig *FlagConfig) isFlagAEnabled() bool {
	return flagConfig != nil && flagConfig.FlagA != 0
}

func (flagConfig *FlagConfig) isFlagBEnabled() bool {
	return flagConfig != nil && flagConfig.FlagB != 0
}

func (flagConfig *FlagConfig) isFlagCEnabled() bool {
	return flagConfig != nil && flagConfig.FlagC
}

func (flagConfig *FlagConfig) isFlagIEnabled() bool {
	return flagConfig != nil && flagConfig.FlagI
}

func (flagConfig *FlagConfig) isFlagOEnabled() bool {
	return flagConfig != nil && flagConfig.FlagO != ""
}
