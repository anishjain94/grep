package main

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"testing"
)

// TODO: write tests for if files does not exists, permission not exists
// TODO: handle for condition when file limit opening is restricted by os. make is os independent.

var testCases = map[string]struct {
	FileName  string
	SearchStr string
	Want      []string
	Iflag     bool
	Oflag     string
}{
	"zeroMatch": {
		FileName:  "testfile.txt",
		SearchStr: "someRandomString",
		Iflag:     false,
	},
	"oneMatch": {
		FileName:  "testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"testfile.txt: this is temperature."},
		Iflag:     false,
	},
	"fileDoesNotExists": {
		FileName:  "fileDoesNotExist.txt",
		SearchStr: "temperature",
		Want:      []string{"testfile.txt: this is temperature."},
		Iflag:     false,
	},
	"multipleMatch": {
		FileName:  "testfile.txt",
		SearchStr: "anish",
		Want: []string{
			"testfile.txt: this is anish.",
			"testfile.txt: is this anish.",
			"testfile.txt: this is anish?",
			"testfile.txt: anish"},
		Iflag: false,
	},
	"oneMatchCaseInsensitive": {
		FileName:  "testfile.txt",
		SearchStr: "Temperature",
		Want:      []string{"testfile.txt: this is temperature."},
		Iflag:     true,
	},
	"oneMatchOutputFile": {
		FileName:  "testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"testfile.txt: this is temperature."},
		Oflag:     "output.txt",
	},
	"multipleMatchesFileinput": {
		FileName:  "root_dir",
		SearchStr: "anish",
		Want: []string{
			"root_dir/parent_dir1/child_dir2/child_dir2_file.txt: this is anish parent_dir1/child_dir1/child_dir1_file.txt",
			"root_dir/parent_dir1/child_dir2/child_dir2_file.txt: is this anish parent_dir1/child_dir1/child_dir1_file.txt",
			"root_dir/parent_dir1/child_dir2/child_dir2_file.txt: this is anish? parent_dir1/child_dir1/child_dir1_file.txt",
			"root_dir/parent_dir2/parent_dir2_file1.txt: this is anish parent_dir2/parent_dir2_file1.txt",
			"root_dir/parent_dir2/parent_dir2_file1.txt: is this anish parent_dir2/parent_dir2_file1.txt",
			"root_dir/parent_dir2/parent_dir2_file1.txt: this is anish? parent_dir2/parent_dir2_file1.txt",
			"root_dir/parent_dir1/child_dir1/child_dir1_file.txt: this is anish parent_dir1/child_dir1/child_dir1_file.txt",
			"root_dir/parent_dir1/child_dir1/child_dir1_file.txt: is this anish parent_dir1/child_dir1/child_dir1_file.txt",
			"root_dir/parent_dir1/child_dir1/child_dir1_file.txt: this is anish? parent_dir1/child_dir1/child_dir1_file.txt",
		},
		Iflag: false,
	},
}

func TestGrep(t *testing.T) {
	for key, value := range testCases {
		t.Run(key, func(t *testing.T) {

			var got []string
			flagConfig := &FlagConfig{
				FlagI: value.Iflag,
				FlagO: value.Oflag,
			}
			subFiles := getAllfileNames(value.FileName)

			for _, subFileName := range subFiles {
				file, err := os.Open(subFileName)
				handleError(err)
				defer file.Close()

				fileResult := readDataAndMatch(file, &subFileName, flagConfig, value.SearchStr)
				got = append(got, fileResult...)
			}

			slices.Sort(got)
			slices.Sort(value.Want)

			if !reflect.DeepEqual(got, value.Want) {
				t.Errorf("got %s \n --- want %s ", got, value.Want)
			}

			displayResult(got, flagConfig)
		})
	}
}

var testCasesUserInput = map[string]struct {
	InputStr  string
	SearchStr string
	Want      []string
}{
	"zero matches": {
		InputStr:  "this does not contain the word.\nthis is empty",
		SearchStr: "someRandomString",
	},
	"one match": {
		InputStr:  "this is temperature.\nthis is one match",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
	},
	"multiple matches": {
		InputStr:  "this is anish.\nis this anish.\nthis is anish?\nanish",
		SearchStr: "anish",
		Want:      []string{"this is anish.", "is this anish.", "this is anish?", "anish"},
	},
}

func TestUserInput(t *testing.T) {
	for key, value := range testCasesUserInput {
		t.Run(key, func(t *testing.T) {
			file, err := os.CreateTemp("", "tempfile")
			handleError(err)
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(value.InputStr)); err != nil {
				handleError(err)
			}

			if _, err := file.Seek(0, 0); err != nil {
				handleError(err)
			}

			oldStdIn := os.Stdin
			os.Stdin = file

			defer func() {
				os.Stdin = oldStdIn
			}()

			inputStr := readDataAndMatch(os.Stdin, nil, nil, value.SearchStr)
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
		handleError(err)
		defer file.Close()

		b.Run(fmt.Sprintf("naive-%s", key), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				inputStr := readDataAndMatch(file, nil, nil, value.SearchStr)
				naiveGrep(inputStr, value.SearchStr, nil)
			}
		})
	}
}

// // func BenchmarkGrepString(b *testing.B) {
// // 	fileName := "testfile.txt"
// // 	searchStr := "anish"
// // 	for i := 0; i < b.N; i++ {
// // 		naiveGrep(fileName, searchStr)
// // 	}
// // }

// // func BenchmarkGrepRegex(b *testing.B) {
// // 	fileName := "testfile.txt"
// // 	searchStr := "anish"
// // 	for i := 0; i < b.N; i++ {
// // 		regexGrep(fileName, searchStr)
// // 	}
// // }
