package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

func validateFile(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil && os.IsNotExist(err) {
		return err
	} else if os.IsPermission(err) {
		return err
	}

	return nil
}

func readAndMatch(dataConfigIo *ReadAndMatchIo) []string {
	var matchResult []string
	scanner := bufio.NewScanner(dataConfigIo.Reader)
	printed := map[int]bool{}
	lineNumber := 0
	pattern := dataConfigIo.Pattern
	bufferLines := make([]string, 0, dataConfigIo.FlagConfig.CountOfLinesBeforeMatch)

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

		if re.MatchString(inputStr) {
			// Storing lines before the match
			if dataConfigIo.FlagConfig.shouldDisplayLinesBeforeMatch() {
				bufferStartIndex := lineNumber - len(bufferLines)
				for i, bufferedLine := range bufferLines {
					if !printed[bufferStartIndex+i] {
						matchResult = append(matchResult, bufferedLine)
						printed[bufferStartIndex+i] = true
					}
				}
			}

			// Storing the matching line
			matchResult = append(matchResult, inputStr)
			printed[lineNumber] = true

			// Storing lines after the match
			if dataConfigIo.FlagConfig.shouldDisplayLinesAfterMatch() {
				for i := 0; i < dataConfigIo.FlagConfig.CountOfLinesAfterMatch; i++ {
					if scanner.Scan() {
						nextLine := scanner.Text()
						matchResult = append(matchResult, nextLine)
						lineNumber++
						printed[lineNumber] = true
					}
				}
			}

			bufferLines = []string{}
		}

		// Maintaining a buffer of lines before the match
		if dataConfigIo.FlagConfig.shouldDisplayLinesBeforeMatch() {
			if len(bufferLines) == dataConfigIo.FlagConfig.CountOfLinesBeforeMatch {
				bufferLines = bufferLines[1:]
			}
			bufferLines = append(bufferLines, inputStr)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err.Error())
	}

	return matchResult
}
