package main

import (
	"bufio"
	"os"
	"strings"
)

func fileValidations(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil && os.IsNotExist(err) {
		return err
	} else if os.IsPermission(err) {
		return err
	}

	return nil
}

func readAndMatch(dataConfigIo *ReadAndMatchConfigIo) []string {
	var matchedLines []string
	searchStr := dataConfigIo.SearchStr
	reader := bufio.NewScanner(dataConfigIo.Reader)

	for reader.Scan() {
		inputStr := reader.Text()

		if dataConfigIo.FlagConfig.shouldSearchCaseInsensitive() {
			inputStr = strings.ToLower(inputStr)
			searchStr = strings.ToLower(searchStr)
		}

		if strings.Contains(inputStr, searchStr) {
			matchedLines = append(matchedLines, reader.Text())
		}
	}
	return matchedLines
}
