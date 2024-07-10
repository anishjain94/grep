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
	"sync"
)

func main() {

	flagi := flag.Bool("i", false, "case insensitive search")
	flago := flag.String("o", "", "output file")
	flag.Parse()

	flagConfig := &FlagConfig{
		FlagI: *flagi,
		FlagO: *flago,
	}

	args := flag.Args()
	var inputStr []string
	var output []string
	var searchStr string

	numOfWorkers := 5

	var wg sync.WaitGroup

	if len(args) == 2 {
		searchStr = args[0]
		filePath := args[1]
		fileValidations(filePath)

		subFiles := getAllfileNames(filePath)
		jobs := make(chan string, len(subFiles))
		result := make(chan []string, len(subFiles))

		for range numOfWorkers {
			go workers(jobs, result, flagConfig, searchStr)
		}

		for _, subFileName := range subFiles {
			wg.Add(1)
			jobs <- subFileName
		}
		close(jobs)

		go func() {
			for outputFromFiles := range result {
				output = append(output, outputFromFiles...)
				wg.Done()
			}
		}()
		wg.Wait()
		displayResult(output, flagConfig)

	} else if len(args) < 2 {
		searchStr = args[0]
		inputStr = readDataAndMatch(os.Stdin, nil, nil, searchStr)

		output := naiveGrep(inputStr, searchStr, flagConfig)
		displayResult(output, flagConfig)
	}
}

func workers(jobs chan string, result chan []string, flagConfig *FlagConfig, searchStr string) {
	for job := range jobs {
		fileMatchedLines := executeGrep(job, flagConfig, searchStr)
		result <- fileMatchedLines
	}
}

func executeGrep(subFileName string, flagconfig *FlagConfig, searchStr string) []string {
	file, err := os.Open(subFileName)
	handleError(err)
	defer file.Close()

	fileResult := readDataAndMatch(file, &subFileName, flagconfig, searchStr)
	return fileResult
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

func getAllfileNames(path string) []string {
	var subFiles []string

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			handleError(err)
		}
		if !d.IsDir() {
			subFiles = append(subFiles, path)
		}
		return nil
	})

	return subFiles
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
