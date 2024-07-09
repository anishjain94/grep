package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	args := os.Args
	var inputStr []string
	var searchStr string

	if len(args) > 2 {
		searchStr = args[1]
		filePath := args[2]
		fileValidations(filePath)

		file, err := os.Open(filePath)
		printError(err)
		defer file.Close()

		inputStr = readDataFromSource(file)

	} else if len(args) == 2 {
		inputStr = readDataFromSource(os.Stdin)
	}

	output := naiveGrep(inputStr, searchStr)

	if len(output) > 0 {
		fmt.Println(output)
	}
}

func naiveGrep(inputStr []string, searchStr string) []string {
	var outputLines []string
	for _, str := range inputStr {
		if strings.Contains(str, searchStr) {
			outputLines = append(outputLines, str)
		}
	}

	return outputLines
}

func regexGrep(inputStr []string, searchStr string) []string {
	var outputLines []string
	re := regexp.MustCompile(searchStr)

	for _, str := range inputStr {
		if re.MatchString(str) {
			outputLines = append(outputLines, str)
		}
	}

	return outputLines
}
