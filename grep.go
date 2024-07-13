package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func readAndMatch(dataConfigIo *ReadAndMatchIo) ([]string, error) {
	var matchResult []string
	scanner := bufio.NewScanner(dataConfigIo.Reader)
	linesPrintedMap := map[int]bool{}
	lineNumber := 0
	pattern := dataConfigIo.Pattern
	dataBuffer := make([]string, 0, dataConfigIo.FlagConfig.CountOfLinesBeforeMatch)
	linesToPrintAfterMatches := 0

	if dataConfigIo.FlagConfig.shouldSearchCaseInsensitive() {
		pattern = strings.ToLower(pattern)
	}
	re := regexp.MustCompile(pattern)

	for scanner.Scan() {
		inputStr := scanner.Text()
		lineNumber++

		if dataConfigIo.FlagConfig.shouldSearchCaseInsensitive() {
			inputStr = strings.ToLower(inputStr)
		}

		if linesToPrintAfterMatches > 0 {
			matchResult = append(matchResult, inputStr)
			linesPrintedMap[lineNumber] = true
			linesToPrintAfterMatches--
		}
		if re.MatchString(inputStr) {
			// Storing lines before the match for -b option
			if dataConfigIo.FlagConfig.shouldDisplayLinesBeforeMatch() {
				matchResult = append(matchResult, fetchResultFromBuffer(lineNumber, dataBuffer, linesPrintedMap)...)
			}

			// Storing the matching line if not printed already
			if !linesPrintedMap[lineNumber] {
				matchResult = append(matchResult, inputStr)
				linesPrintedMap[lineNumber] = true
			}

			// Storing lines after the match for -a option
			if dataConfigIo.FlagConfig.shouldDisplayLinesAfterMatch() {
				linesToPrintAfterMatches = dataConfigIo.FlagConfig.CountOfLinesAfterMatch
			}
			dataBuffer = []string{}
		}

		// Maintaining a buffer of lines before the match for -b option
		if dataConfigIo.FlagConfig.shouldDisplayLinesBeforeMatch() {
			if len(dataBuffer) == dataConfigIo.FlagConfig.CountOfLinesBeforeMatch {
				dataBuffer = dataBuffer[1:]
			}
			dataBuffer = append(dataBuffer, inputStr)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matchResult, nil
}

func fetchResultFromBuffer(lineNumber int, dataBuffer []string, linesPrintedMap map[int]bool) []string {
	results := make([]string, 0, len(dataBuffer))
	bufferStartIndex := lineNumber - len(dataBuffer)
	for i, bufferedLine := range dataBuffer {
		if !linesPrintedMap[bufferStartIndex+i] {
			results = append(results, bufferedLine)
			linesPrintedMap[bufferStartIndex+i] = true
		}
	}
	return results
}

func displayResult(dataIo *DisplayResultIo) error {
	var writer io.Writer = os.Stdout

	if dataIo.FlagConfig.shouldStoreOutput() {
		file, err := os.OpenFile(dataIo.FlagConfig.OutputFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		writer = file
	}

	for _, filePath := range dataIo.FilesInDirectory {
		if dataIo.FlagConfig.shouldShowCount() {
			fmt.Fprintln(writer, strconv.Itoa(len(dataIo.matchedResultMap)))
			continue
		}
		for _, value := range dataIo.matchedResultMap[filePath] {
			valueToPrint := value
			if dataIo.IsDirectory {
				valueToPrint = filePath + ": " + value
			}
			fmt.Fprintln(writer, valueToPrint)
		}
	}
	return nil
}
