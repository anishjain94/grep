package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

func main() {
	var pattern string

	fileResultMap := FileResultMap{}
	wg := &sync.WaitGroup{}

	flagConfig := parseFlags()
	args := flag.Args()
	numOfWorkers := 5

	if len(args) == 0 || len(args) > 2 {
		log.Fatalf("incorrect number of args")
	}

	// input from stdin
	if len(args) == 1 {
		pattern = args[0]
		sourceName := "stdin"

		output, err := readAndMatch(&ReadAndMatchIo{
			Reader:     os.Stdin,
			FlagConfig: flagConfig,
			Pattern:    pattern,
		})
		if err != nil {
			log.Fatalf(err.Error())
		}

		fileResultMap[sourceName] = output

		displayResult(&DisplayResultIo{
			matchedResultMap: fileResultMap,
			FlagConfig:       flagConfig,
			FilesInDirectory: []string{sourceName},
			IsDirectory:      false,
		})
	}

	// input from file/directory
	if len(args) == 2 {
		pattern = args[0]
		filePath := args[1]

		if err := validateFile(filePath); err != nil {
			log.Fatalf(err.Error())
		}

		filesToBeSearched, isDirectory, err := listFilesInDir(filePath)

		//we do not what our program to error out and exit completely incase we encounter file permission error
		if err != nil && !os.IsPermission(err) {
			log.Fatalf(err.Error())
		}

		jobsChannel := make(chan string, numOfWorkers)
		resultChannel := make(chan FileResultMap, numOfWorkers) // Keeping a map here to make the output consisent with multiple goroutines.

		// if files to search is less than numOfWorkers then only spin up that no of goroutines..
		if len(filesToBeSearched) < numOfWorkers {
			numOfWorkers = len(filesToBeSearched)
		}

		for range numOfWorkers {
			go worker(jobsChannel, resultChannel, flagConfig, pattern)
		}

		wg.Add(len(filesToBeSearched))
		for _, file := range filesToBeSearched {
			jobsChannel <- file
		}

		close(jobsChannel)
		go func() {
			wg.Wait()
			close(resultChannel)
		}()

		fileResultMap = collectResult(resultChannel, wg)
		displayResult(&DisplayResultIo{
			matchedResultMap: fileResultMap,
			FlagConfig:       flagConfig,
			FilesInDirectory: filesToBeSearched,
			IsDirectory:      isDirectory,
		})
	}
}

func parseFlags() *GrepConfigIo {
	caseInsensitiveSearch := flag.Bool("i", false, "case insensitive search")
	outputFileName := flag.String("o", "", "output file")
	countOfMatches := flag.Bool("c", false, "displays count of matches instead of actual matched lines")
	countOfLinesBeforeMatch := flag.Int("B", 0, "displays n lines before the match")
	countOfLinesAfterMatch := flag.Int("A", 0, "displays n lines after the match")

	flag.Parse()
	flagConfig := &GrepConfigIo{
		CaseInsensitiveSearch:   *caseInsensitiveSearch,
		OutputFileName:          *outputFileName,
		CountOfMatches:          *countOfMatches,
		CountOfLinesBeforeMatch: *countOfLinesBeforeMatch,
		CountOfLinesAfterMatch:  *countOfLinesAfterMatch,
	}
	return flagConfig
}

func collectResult(result <-chan FileResultMap, wg *sync.WaitGroup) FileResultMap {
	fileResultMap := FileResultMap{}
	for outputFromFiles := range result {
		for key, value := range outputFromFiles {
			fileResultMap[key] = value
		}
		wg.Done()
	}

	return fileResultMap
}

func worker(filePaths <-chan string, result chan<- FileResultMap, flagConfig *GrepConfigIo, searchStr string) {
	for filePath := range filePaths {
		matchedLines, err := executeGrep(filePath, flagConfig, searchStr)
		if err != nil {
			log.Println(err.Error())
		}
		result <- FileResultMap{filePath: matchedLines}
	}
}

func executeGrep(subFileName string, flagconfig *GrepConfigIo, searchStr string) ([]string, error) {
	file, err := os.Open(subFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileResult, err := readAndMatch(&ReadAndMatchIo{
		Reader:     file,
		FlagConfig: flagconfig,
		Pattern:    searchStr,
	})

	if err != nil {
		return nil, err
	}

	return fileResult, nil
}
