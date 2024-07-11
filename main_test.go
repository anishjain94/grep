package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

// TODO: write tests for if files does not exists, permission not exists
// TODO: handle for condition when file limit opening is restricted by os. make is os independent.

var testCases = map[string]struct {
	FileName     string
	SearchStr    string
	Want         []string
	flagConfigIo FlagConfigIo
}{
	"zeroMatch": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "someRandomString",
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: false,
		},
	},
	"oneMatch": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: false,
		},
	},
	"fileDoesNotExists": {
		FileName:  "fileDoesNotExist.txt",
		SearchStr: "temperature",
		Want:      []string{"lstat fileDoesNotExist.txt: no such file or directory"},
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: false,
		},
	},
	"multipleMatch": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "anish",
		Want: []string{
			"this is anish.",
			"is this anish.",
			"this is anish?",
			"anish"},
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: false,
		},
	},
	"oneMatchCaseInsensitive": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "Temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: true,
		},
	},
	"oneMatchOutputFile": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: FlagConfigIo{
			OutputFileName: "output.txt",
		},
	},
	"multipleMatchesDirectory": {
		FileName:  "test_files",
		SearchStr: "anish",
		Want: []string{
			"this is anish parent_dir1/child_dir1/child_dir1_file.txt",
			"is this anish parent_dir1/child_dir1/child_dir1_file.txt",
			"this is anish? parent_dir1/child_dir1/child_dir1_file.txt",
			"this is anish parent_dir1/child_dir2/child_dir2_file.txt",
			"is this anish parent_dir1/child_dir2/child_dir2_file.txt",
			"this is anish? parent_dir1/child_dir2/child_dir2_file.txt",
			"this is anish parent_dir2/parent_dir2_file1.txt",
			"is this anish parent_dir2/parent_dir2_file1.txt",
			"this is anish? parent_dir2/parent_dir2_file1.txt",
			"this is anish.",
			"is this anish.",
			"this is anish?",
			"anish",
		},
		flagConfigIo: FlagConfigIo{
			CaseInSensitiveSearch: false,
		},
	},
	"NLinesBefore": {
		FileName:  "test_files/testfile2.txt",
		SearchStr: "test",
		Want: []string{
			"this is line 6",
			"this is line 7",
			"this is test 8",
			"this is line 12",
			"this is line 13",
			"this is test 14",
		},
		flagConfigIo: FlagConfigIo{
			CountOfLinesBeforeMatch: 2,
		},
	},
	"NLinesAfer": {
		FileName:  "test_files/testfile2.txt",
		SearchStr: "test",
		Want: []string{
			"this is test 8",
			"this is line 9",
			"this is line 10",
			"this is test 14",
			"this is line 15",
			"this is line 16",
		},
		flagConfigIo: FlagConfigIo{
			CountOfLinesAfterMatch: 2,
		},
	},
}

func TestGrep(t *testing.T) {
	for key, value := range testCases {
		t.Run(key, func(t *testing.T) {

			var got []string
			subFiles, _, err := listFilesInDir(value.FileName)
			if err != nil {
				if err.Error() != value.Want[0] {
					t.Errorf("got %s \n --- want %s ", err.Error(), value.Want[0])
				}
				return
			}

			for _, subFileName := range subFiles {
				file, err := os.Open(subFileName)
				if err != nil {
					panic(err.Error())
				}
				defer file.Close()

				fileResult := readAndMatch(&ReadAndMatchIo{
					Reader:     file,
					Source:     &subFileName,
					FlagConfig: &value.flagConfigIo,
					Pattern:    value.SearchStr,
				})
				got = append(got, fileResult...)
			}

			if !reflect.DeepEqual(got, value.Want) {
				t.Errorf("got %s \n --- want %s ", got, value.Want)
			}

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
			if err != nil {
				log.Panic(err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(value.InputStr)); err != nil {
				log.Panic(err)
			}

			if _, err := file.Seek(0, 0); err != nil {
				log.Panic(err)
			}

			oldStdIn := os.Stdin
			os.Stdin = file

			defer func() {
				os.Stdin = oldStdIn
			}()

			sourceName := "stdin"
			gotResult := readAndMatch(
				&ReadAndMatchIo{
					Reader: os.Stdin,
					Source: &sourceName,
					FlagConfig: &FlagConfigIo{
						CountOfLinesBeforeMatch: 0,
						CountOfLinesAfterMatch:  0,
					},
					Pattern: value.SearchStr,
				},
			)

			if !reflect.DeepEqual(gotResult, value.Want) {
				t.Errorf("got %s \n --- want %s ", gotResult, value.Want)
			}
		})
	}
}

func BenchmarkTableRegex(b *testing.B) {
	for key, value := range testCases {
		file, err := os.Open(value.FileName)
		log.Panic(err)
		defer file.Close()

		b.Run(fmt.Sprintf("naive-%s", key), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				readAndMatch(&ReadAndMatchIo{
					Reader: file,
					Source: nil,
					FlagConfig: &FlagConfigIo{
						CountOfLinesBeforeMatch: 0,
						CountOfLinesAfterMatch:  0,
					},
					Pattern: value.SearchStr,
				})
			}
		})
	}
}
