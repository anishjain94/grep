package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

		inputStr = fetchAllfilesAndGetContent(filePath)

	} else if len(args) < 2 {
		searchStr = args[0]
		inputStr = readDataFromSource(os.Stdin, nil)
	}

	output := naiveGrep(inputStr, searchStr, flagconfig)
	displayResult(output, flagconfig)
}

// TODO: Refactor
func fetchAllfilesAndGetContent(filePath string) []string {
	var inputStr []string
	subFiles, isDirectory := walkDir(filePath)

	for _, subFile := range subFiles {
		file, err := os.Open(subFile)
		handleError(err)
		defer file.Close()

		if isDirectory {
			fileResult := readDataFromSource(file, &subFile)
			inputStr = append(inputStr, fileResult...)
		} else {
			fileResult := readDataFromSource(file, nil)
			inputStr = append(inputStr, fileResult...)
		}
	}
	return inputStr
}

func displayResult(output []string, flagconfig *FlagConfig) {
	var outputDestination io.Writer

	if flagconfig.isFlagOEnabled() {
		file, err := os.OpenFile(flagconfig.FlagO, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		handleError(err)
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

func walkDir(path string) ([]string, bool) {
	var subFiles []string
	var isDirectory bool

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			handleError(err)
		}
		if !d.IsDir() {
			subFiles = append(subFiles, path)
		} else {
			isDirectory = true
		}
		return nil
	})

	return subFiles, isDirectory
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
