package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var testCases = map[string]struct {
	FileName  string
	SearchStr string
	Want      []string
	Iflag     bool
}{
	"Zero matches": {
		FileName:  "testfile.txt",
		SearchStr: "someRandomString",
		Iflag:     false,
	},
	"One match": {
		FileName:  "testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
		Iflag:     false,
	},
	"Multiple matches": {
		FileName:  "testfile.txt",
		SearchStr: "anish",
		Want:      []string{"this is anish.", "is this anish.", "this is anish?", "anish"},
		Iflag:     false,
	},
	"One match - case insensitive": {
		FileName:  "testfile.txt",
		SearchStr: "Temperature",
		Want:      []string{"this is temperature."},
		Iflag:     true,
	},
}

func TestGrep(t *testing.T) {
	for key, value := range testCases {
		t.Run(key, func(t *testing.T) {

			file, err := os.Open(value.FileName)
			printError(err)
			defer file.Close()

			inputStr := readDataFromSource(file)
			gotContains := naiveGrep(inputStr, value.SearchStr, &FlagConfig{
				FlagI: value.Iflag,
			})

			if !reflect.DeepEqual(gotContains, value.Want) {
				t.Errorf("got %s \n --- want %s ", gotContains, value.Want)
			}
		})
	}
}

var testCasesUserInput = map[string]struct {
	InputStr  string
	SearchStr string
	Want      []string
}{
	"Zero matches": {
		InputStr:  "this does not contain the word.\nthis is empty",
		SearchStr: "someRandomString",
	},
	"One match": {
		InputStr:  "this is temperature.\nthis is one match",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
	},
	"Multiple matches": {
		InputStr:  "this is anish.\nis this anish.\nthis is anish?\nanish",
		SearchStr: "anish",
		Want:      []string{"this is anish.", "is this anish.", "this is anish?", "anish"},
	},
}

func TestUserInput(t *testing.T) {
	for key, value := range testCasesUserInput {
		t.Run(key, func(t *testing.T) {
			file, err := os.CreateTemp("", "tempfile")
			printError(err)
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(value.InputStr)); err != nil {
				printError(err)
			}

			if _, err := file.Seek(0, 0); err != nil {
				printError(err)
			}

			oldStdIn := os.Stdin
			os.Stdin = file

			defer func() {
				os.Stdin = oldStdIn
			}()

			inputStr := readDataFromSource(os.Stdin)
			gotContains := naiveGrep(inputStr, value.SearchStr, nil)

			if !reflect.DeepEqual(gotContains, value.Want) {
				t.Errorf("got %s \n --- want %s ", gotContains, value.Want)
			}
		})
	}
}

func BenchmarkTableRegex(b *testing.B) {
	for key, value := range testCases {

		file, err := os.Open(value.FileName)
		printError(err)
		defer file.Close()
		inputStr := readDataFromSource(file)

		b.Run(fmt.Sprintf("naive-%s", key), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				naiveGrep(inputStr, value.SearchStr, nil)
			}
		})
	}
}

// func BenchmarkGrepString(b *testing.B) {
// 	fileName := "testfile.txt"
// 	searchStr := "anish"
// 	for i := 0; i < b.N; i++ {
// 		naiveGrep(fileName, searchStr)
// 	}
// }

// func BenchmarkGrepRegex(b *testing.B) {
// 	fileName := "testfile.txt"
// 	searchStr := "anish"
// 	for i := 0; i < b.N; i++ {
// 		regexGrep(fileName, searchStr)
// 	}
// }
