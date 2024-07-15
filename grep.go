package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func readAndMatch(dataConfigIo *ReadAndMatchInput) ([]string, error) {
	var matchResult []string
	scanner := bufio.NewScanner(dataConfigIo.Reader)
	linesPrinted := map[int]bool{}
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
			linesPrinted[lineNumber] = true
			linesToPrintAfterMatches--
		}
		if re.MatchString(inputStr) {
			// Storing lines before the match for -b option
			if dataConfigIo.FlagConfig.shouldDisplayLinesBeforeMatch() {
				matchResult = append(matchResult, fetchResultFromBuffer(lineNumber, dataBuffer, linesPrinted)...)
			}

			// Storing the matching line if not printed already
			if !linesPrinted[lineNumber] {
				matchResult = append(matchResult, inputStr)
				linesPrinted[lineNumber] = true
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

func fetchResultFromBuffer(lineNumber int, dataBuffer []string, linesPrinted map[int]bool) []string {
	results := make([]string, 0, len(dataBuffer))
	bufferStartIndex := lineNumber - len(dataBuffer)
	for i, bufferedLine := range dataBuffer {
		if !linesPrinted[bufferStartIndex+i] {
			results = append(results, bufferedLine)
			linesPrinted[bufferStartIndex+i] = true
		}
	}
	return results
}

func displayResult(dataIo *DisplayResultInput) error {
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
			valueToPrint := strconv.Itoa(len(dataIo.MatchedResult[filePath]))
			if dataIo.IsDirectory {
				valueToPrint = filePath + ": " + valueToPrint
			}
			fmt.Fprintln(writer, valueToPrint)
			continue
		}
		for _, value := range dataIo.MatchedResult[filePath] {
			valueToPrint := value
			if dataIo.IsDirectory {
				valueToPrint = filePath + ": " + value
			}
			fmt.Fprintln(writer, valueToPrint)
		}
	}
	return nil
}

func listFilesInDir(path string) ([]string, bool, error) {
	var subFiles []string
	var isDir bool

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			subFiles = append(subFiles, path)
		} else {
			isDir = true
		}
		return nil
	})

	if err != nil {
		return subFiles, isDir, err //To not error out and exit completely incase we encounter file permission error
	}

	return subFiles, isDir, nil
}
