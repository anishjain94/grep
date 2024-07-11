package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	caseInSensitiveSearch := flag.Bool("i", false, "case insensitive search")
	outputFileName := flag.String("o", "", "output file")
	countOfMatches := flag.Bool("c", false, "displays count of matches instead of actual matched lines")
	countOfLinesBeforeMatch := flag.Int("a", 0, "displays n lines before the match")
	countOfLinesAfterMatch := flag.Int("b", 0, "displays n lines after the match")

	flag.Parse()
	flagConfig := &FlagConfigIo{
		CaseInSensitiveSearch:   *caseInSensitiveSearch,
		OutputFileName:          *outputFileName,
		CountOfMatches:          *countOfMatches,
		CountOfLinesBeforeMatch: *countOfLinesBeforeMatch,
		CountOfLinesAfterMatch:  *countOfLinesAfterMatch,
	}

	var output []string
	var searchStr string

	fileResultMap := FileResultMap{}
	wg := &sync.WaitGroup{}

	args := flag.Args()
	numOfWorkers := 5

	if len(args) == 0 || len(args) > 2 {
		log.Panic("incorrect number of args")
	}

	if len(args) == 1 {
		searchStr = args[0]
		sourceName := "stdin"

		output = readAndMatch(&ReadAndMatchIo{
			Reader:     os.Stdin,
			Source:     &sourceName,
			FlagConfig: flagConfig,
			Pattern:    searchStr,
		})

		fileResultMap[sourceName] = output

		displayResult(&DisplayResultIo{
			matchedResultMap: fileResultMap,
			FlagConfig:       flagConfig,
			FilesInDirectory: []string{sourceName},
			IsDirectory:      false,
		})
	}

	if len(args) == 2 {
		searchStr = args[0]
		filePath := args[1]

		err := validateFile(filePath)
		if err != nil {
			log.Panic(err.Error())
		}

		filesToBeSearched, isDirectory, err := listFilesInDir(filePath)
		if err != nil {
			log.Panic(err.Error())
		}

		jobs := make(chan string, len(filesToBeSearched))
		result := make(chan FileResultMap, len(filesToBeSearched)) //NOTE: Keeping a map here to make the output consisent with multiple goroutines.

		for range numOfWorkers {
			go workers(jobs, result, flagConfig, searchStr)
		}

		(wg).Add(len(filesToBeSearched))
		for _, fileToBeSearched := range filesToBeSearched {
			fileResultMap[fileToBeSearched] = []string{} //Adding file path to map
			jobs <- fileToBeSearched
		}
		close(jobs)

		collectResult(result, fileResultMap, wg)
		displayResult(&DisplayResultIo{
			matchedResultMap: fileResultMap,
			FlagConfig:       flagConfig,
			FilesInDirectory: filesToBeSearched,
			IsDirectory:      isDirectory,
		})
	}
}

func collectResult(result chan FileResultMap, fileResultMap FileResultMap, wg *sync.WaitGroup) {
	go func() {
		for outputFromFiles := range result {
			for key, value := range outputFromFiles {
				fileResultMap[key] = value
			}
			wg.Done()
		}
	}()
	wg.Wait()
}

func workers(filePaths chan string, result chan FileResultMap, flagConfig *FlagConfigIo, searchStr string) {
	for filePath := range filePaths {
		fileMatchedLines, err := executeGrep(filePath, flagConfig, searchStr)
		if err != nil {
			log.Panic(err.Error())
		}

		result <- FileResultMap{filePath: fileMatchedLines}
	}
}

func executeGrep(subFileName string, flagconfig *FlagConfigIo, searchStr string) ([]string, error) {
	file, err := os.Open(subFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileResult := readAndMatch(&ReadAndMatchIo{
		Reader:     file,
		Source:     &subFileName,
		FlagConfig: flagconfig,
		Pattern:    searchStr,
	})

	return fileResult, nil
}

func displayResult(dataIo *DisplayResultIo) error {
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
			fmt.Fprintln(writer, strconv.Itoa(len(dataIo.matchedResultMap)))
			continue
		}
		for _, value := range dataIo.matchedResultMap[filePath] {
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
		return nil, false, err
	}

	return subFiles, isDir, nil
}
