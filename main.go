package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

func main() {

	flagi := flag.Bool("i", false, "case insensitive search")
	flago := flag.String("o", "", "output file")
	flagC := flag.Bool("c", false, "displays count of matches instead of actual matched lines")
	flagA := flag.Int("a", 0, "displays n lines before the match")
	flagB := flag.Int("b", 0, "displays n lines after the match")

	flag.Parse()

	flagConfig := &FlagConfig{
		FlagI: *flagi,
		FlagO: *flago,
		FlagC: *flagC,
		FlagA: *flagA,
		FlagB: *flagB,
	}

	args := flag.Args()
	var output []string
	var searchStr string

	numOfWorkers := 5

	var wg sync.WaitGroup
	var fileResultMap = map[string][]string{}

	if len(args) == 0 {
		log.Println("Usage: ./grep search_query filename.txt")
		os.Exit(1)
	}

	if len(args) == 1 {
		searchStr = args[0]
		output = readMatchAndStore(os.Stdin, nil, nil, searchStr)

		fileResultMap[""] = output
		displayResult(fileResultMap, flagConfig, nil, false)
	}

	if len(args) == 2 {
		searchStr = args[0]
		filePath := args[1]

		err := fileValidations(filePath)
		handleError(err)

		filesToBeSearched, isDirectory := listFilesInDir(filePath)
		jobs := make(chan string, len(filesToBeSearched))
		result := make(chan map[string][]string, len(filesToBeSearched)) //NOTE: Keeping a map here to make the output consisent with multiple goroutines.

		for range numOfWorkers {
			go workers(jobs, result, flagConfig, searchStr)
		}

		wg.Add(len(filesToBeSearched))
		for _, fileToBeSearched := range filesToBeSearched {
			fileResultMap[fileToBeSearched] = []string{} //Adding file path to map
			jobs <- fileToBeSearched
		}
		close(jobs)

		gatherResult(result, fileResultMap, &wg)
		displayResult(fileResultMap, flagConfig, filesToBeSearched, isDirectory)
	}

}

func gatherResult(result chan map[string][]string, fileResultMap map[string][]string, wg *sync.WaitGroup) {
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

func workers(jobs chan string, result chan map[string][]string, flagConfig *FlagConfig, searchStr string) {
	for job := range jobs {
		fileMatchedLines := executeGrep(job, flagConfig, searchStr)

		fileResultMap := map[string][]string{
			job: fileMatchedLines,
		}
		result <- fileResultMap
	}
}

func executeGrep(subFileName string, flagconfig *FlagConfig, searchStr string) []string {
	file, err := os.Open(subFileName)
	handleError(err)
	defer file.Close()

	fileResult := readMatchAndStore(file, &subFileName, flagconfig, searchStr)
	return fileResult
}

func displayResult(output map[string][]string, flagconfig *FlagConfig, fileOrder []string, isDirectory bool) {
	var outputDestination io.Writer = os.Stdout

	if flagconfig.isFlagOEnabled() {
		file, err := os.OpenFile(flagconfig.FlagO, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		handleError(err)
		defer file.Close()

		outputDestination = file
	}

	if len(fileOrder) == 0 {
		fileOrder = []string{""}
	}

	for _, files := range fileOrder {
		if flagconfig.isFlagCEnabled() {
			fmt.Fprint(outputDestination, strconv.Itoa(len(output))+"\n")
			continue
		}
		for _, value := range output[files] {
			valueToPrint := value
			if isDirectory {
				valueToPrint = files + ": " + value
			}

			fmt.Fprint(outputDestination, valueToPrint+"\n")
		}
	}
}

func listFilesInDir(path string) ([]string, bool) {
	var subFiles []string
	var isDir bool
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			handleError(err)
		}
		if !d.IsDir() {
			subFiles = append(subFiles, path)
		} else {
			isDir = true
		}
		return nil
	})

	return subFiles, isDir
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
