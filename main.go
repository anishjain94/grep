package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	toSearch := os.Args[1]
	filePath := os.Args[2]
	fileValidations(filePath)

	output := naiveGrep(filePath, toSearch)

	if len(output) > 0 {
		fmt.Println(output)
	}
}

func naiveGrep(filePath string, searchStr string) []string {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	printError(err)
	defer file.Close()

	var outputLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLine := scanner.Text()
		if strings.Contains(currentLine, searchStr) {
			outputLines = append(outputLines, currentLine)
		}
	}

	return outputLines
}

func regexGrep(filePath string, searchStr string) []string {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	printError(err)
	defer file.Close()

	var outputLines []string
	scanner := bufio.NewScanner(file)

	re := regexp.MustCompile(searchStr)

	for scanner.Scan() {
		currentLine := scanner.Text()
		if re.MatchString(currentLine) {
			outputLines = append(outputLines, currentLine)
		}
	}

	return outputLines
}
