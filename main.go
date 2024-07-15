package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

func main() {
	var pattern string

	fileResult := FileResult{}
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

		output, err := readAndMatch(&ReadAndMatchInput{
			Reader:     os.Stdin,
			FlagConfig: flagConfig,
			Pattern:    pattern,
		})
		if err != nil {
			log.Fatalf(err.Error())
		}

		fileResult[sourceName] = output

		displayResult(&DisplayResultInput{
			MatchedResult:    fileResult,
			FlagConfig:       flagConfig,
			FilesInDirectory: []string{sourceName},
			IsDirectory:      false,
		})
	}

	// input from file/directory
	if len(args) == 2 {
		pattern = args[0]
		filePath := args[1]

		if _, err := os.Stat(filePath); err != nil {
			log.Fatalf(err.Error())
		}

		filesToBeSearched, isDirectory, err := listFilesInDir(filePath)

		if isDirectory && !flagConfig.ResurciveSearch {
			log.Fatalf("%s Is a directory", filePath)
		}

		//we do not what our program to error out and exit completely incase we encounter file permission error
		if err != nil && !os.IsPermission(err) {
			log.Fatalf(err.Error())
		}

		jobs := make(chan string, numOfWorkers)
		results := make(chan FileResult, numOfWorkers) // Keeping a map here to make the output consisent with multiple goroutines.

		// if files to search is less than numOfWorkers then only spin up that no of goroutines..
		if len(filesToBeSearched) < numOfWorkers {
			numOfWorkers = len(filesToBeSearched)
		}

		for range numOfWorkers {
			go worker(jobs, results, flagConfig, pattern)
		}

		wg.Add(len(filesToBeSearched))
		for _, file := range filesToBeSearched {
			jobs <- file
		}

		close(jobs)
		go func() {
			wg.Wait()
			close(results)
		}()

		fileResult = collectResult(results, wg)
		displayResult(&DisplayResultInput{
			MatchedResult:    fileResult,
			FlagConfig:       flagConfig,
			FilesInDirectory: filesToBeSearched,
			IsDirectory:      isDirectory,
		})
	}
}

func parseFlags() *GrepConfig {
	countOfLinesAfterMatch := flag.Int("A", 0, "displays n lines after the match")
	countOfLinesBeforeMatch := flag.Int("B", 0, "displays n lines before the match")
	countOfMatches := flag.Bool("c", false, "displays count of matches instead of actual matched lines")
	caseInsensitiveSearch := flag.Bool("i", false, "case insensitive search")
	outputFileName := flag.String("o", "", "output file")
	recurviceSearch := flag.Bool("r", false, "search in a directory")

	flag.Parse()
	flagConfig := &GrepConfig{
		CaseInsensitiveSearch:   *caseInsensitiveSearch,
		OutputFileName:          *outputFileName,
		CountOfMatches:          *countOfMatches,
		CountOfLinesBeforeMatch: *countOfLinesBeforeMatch,
		CountOfLinesAfterMatch:  *countOfLinesAfterMatch,
		ResurciveSearch:         *recurviceSearch,
	}
	return flagConfig
}

func collectResult(result <-chan FileResult, wg *sync.WaitGroup) FileResult {
	fileResult := FileResult{}
	for outputFromFiles := range result {
		for key, value := range outputFromFiles {
			fileResult[key] = value
		}
		wg.Done()
	}

	return fileResult
}

func worker(filePaths <-chan string, result chan<- FileResult, flagConfig *GrepConfig, searchStr string) {
	for filePath := range filePaths {
		matchedLines, err := executeGrep(filePath, flagConfig, searchStr)
		if err != nil {
			log.Println(err.Error())
		}
		result <- FileResult{filePath: matchedLines}
	}
}

func executeGrep(subFileName string, flagconfig *GrepConfig, searchStr string) ([]string, error) {
	file, err := os.Open(subFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileResult, err := readAndMatch(&ReadAndMatchInput{
		Reader:     file,
		FlagConfig: flagconfig,
		Pattern:    searchStr,
	})

	if err != nil {
		return nil, err
	}

	return fileResult, nil
}
