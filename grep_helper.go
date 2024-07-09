package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func printError(err error) {
	if err != nil {
		errMsg := err.Error()
		println(errMsg)
		os.Exit(1)
	}
}

func fileValidations(filepath string) {
	stat, err := os.Stat(filepath)

	if err == nil && stat.IsDir() {
		println(filepath + " is a directory, instead of a file.")
		os.Exit(1)
	} else if err != nil && os.IsNotExist(err) {
		println(filepath + " No such file")
		os.Exit(1)
	} else if err != nil && os.IsPermission(err) {
		println(filepath + " Permission denied")
		os.Exit(1)
	}

}

func readDataFromSource(r io.Reader) []string {
	var inputStr []string
	reader := bufio.NewScanner(r)

	for reader.Scan() {
		inputStr = append(inputStr, reader.Text())
	}
	return inputStr
}

type FlagConfig struct {
	FlagI bool   //case-inSensitive search
	FlagO string //output file
}

func (flagConfig *FlagConfig) isFlagIEnabled() bool {
	return flagConfig.FlagI
}

func (flagConfig *FlagConfig) isFlagOEnabled() bool {
	return flagConfig.FlagO != ""
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
