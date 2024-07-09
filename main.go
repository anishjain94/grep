package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {

	flagi := flag.Bool("i", false, "case insensitive search")
	flago := flag.String("o", "", "output file")
	flag.Parse()

	flagconfig := &FlagConfig{
		FlagI: *flagi,
		FlagO: *flago,
	}

	args := flag.Args()

	var inputStr []string
	var searchStr string

	if len(args) == 2 {
		searchStr = args[0]
		filePath := args[1]
		fileValidations(filePath)

		file, err := os.Open(filePath)
		printError(err)
		defer file.Close()

		inputStr = readDataFromSource(file)

	} else if len(args) < 2 {
		searchStr = args[0]
		inputStr = readDataFromSource(os.Stdin)
	}

	output := naiveGrep(inputStr, searchStr, flagconfig)
	displayResult(output, flagconfig)
}

func displayResult(output []string, flagconfig *FlagConfig) {
	var outputDestination io.Writer

	if flagconfig.isFlagOEnabled() {
		file, err := os.OpenFile(flagconfig.FlagO, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		printError(err)
		defer file.Close()

		outputDestination = file
	} else {
		outputDestination = os.Stdout
	}

	for _, value := range output {
		fmt.Fprint(outputDestination, value+"\n")
	}
}

func naiveGrep(inputStr []string, searchStr string, flagconfig *FlagConfig) []string {
	var outputLines []string
	for _, str := range inputStr {
		if flagconfig != nil && flagconfig.isFlagIEnabled() {
			if strings.Contains(strings.ToLower(str), strings.ToLower(searchStr)) {
				outputLines = append(outputLines, str)
			}
		} else {
			if strings.Contains(str, searchStr) {
				outputLines = append(outputLines, str)
			}
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
