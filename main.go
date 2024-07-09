package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

//TODO: use cobra

func main() {

	flagi := flag.Bool("i", false, "case insensitive search")

	flag.Parse()

	flagconfig := &FlagConfig{
		FlagI: *flagi,
	}

	args := sanitizeArgs(os.Args)
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

	output := naiveGrep(inputStr, searchStr, flagconfig)

	if len(output) > 0 {
		fmt.Println(output)
	}
}

func sanitizeArgs(args []string) []string {
	var newArgs []string

	for _, val := range args {
		if !strings.HasPrefix(val, "-") {
			newArgs = append(newArgs, val)
		}
	}
	return newArgs

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
