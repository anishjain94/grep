package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

var grepTestCases = map[string]struct {
	FileName      string
	SearchPattern string
	Want          []string
	flagConfigIo  GrepConfig
}{
	"zeroMatch": {
		FileName:      "test_files/testfile.txt",
		SearchPattern: "someRandomString",
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: false,
		},
	},
	"oneMatch": {
		FileName:      "test_files/testfile.txt",
		SearchPattern: "temperature",
		Want:          []string{"this is temperature."},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: false,
		},
	},
	"fileDoesNotExists": {
		FileName:      "fileDoesNotExist.txt",
		SearchPattern: "temperature",
		Want:          []string{"lstat fileDoesNotExist.txt: no such file or directory"},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: false,
		},
	},
	"multipleMatch": {
		FileName:      "test_files/testfile.txt",
		SearchPattern: "anish",
		Want: []string{
			"this is anish.",
			"is this anish.",
			"this is anish?",
			"anish"},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: false,
		},
	},
	"oneMatchCaseInsensitive": {
		FileName:      "test_files/testfile.txt",
		SearchPattern: "Temperature",
		Want:          []string{"this is temperature."},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: true,
		},
	},
	"oneMatchOutputFile": {
		FileName:      "test_files/testfile.txt",
		SearchPattern: "temperature",
		Want:          []string{"this is temperature."},
		flagConfigIo: GrepConfig{
			OutputFileName: "output.txt",
		},
	},
	"multipleMatchesDirectory": {
		FileName:      "test_files",
		SearchPattern: "anish",
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
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch: false,
			ResurciveSearch:       true,
		},
	},
	"nLinesBefore": {
		FileName:      "test_files/testfile2.txt",
		SearchPattern: "test",
		Want: []string{
			"this is line 6",
			"this is line 7",
			"this is test 8",
			"this is line 12",
			"this is line 13",
			"this is test 14",
			"this is test 15",
			"this is test 16",
		},
		flagConfigIo: GrepConfig{
			CountOfLinesBeforeMatch: 2,
		},
	},
	"nLinesAfer": {
		FileName:      "test_files/testfile2.txt",
		SearchPattern: "test",
		Want: []string{
			"this is test 8",
			"this is line 9",
			"this is line 10",
			"this is test 14",
			"this is test 15",
			"this is test 16",
			"this is line 17",
			"this is line 18",
		},
		flagConfigIo: GrepConfig{
			CountOfLinesAfterMatch: 2,
		},
	},
	"multipleFlagsOne": {
		FileName:      "test_files",
		SearchPattern: "Test",
		Want: []string{
			"this is line 6",
			"this is line 7",
			"this is test 8",
			"this is line 9",
			"this is line 12",
			"this is line 13",
			"this is test 14",
			"this is test 15",
			"this is test 16",
			"this is line 17",
		},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch:   true,
			CountOfLinesBeforeMatch: 2,
			CountOfLinesAfterMatch:  1,
			ResurciveSearch:         true,
		},
	},
	"multipleFlagsTwo": {
		FileName:      "test_files",
		SearchPattern: "test",
		Want: []string{
			"this is line 6",
			"this is line 7",
			"this is test 8",
			"this is line 12",
			"this is line 13",
			"this is test 14",
			"this is test 15",
			"this is test 16",
		},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch:   true,
			CountOfLinesBeforeMatch: 2,
			ResurciveSearch:         true,
		},
	},
	"multipleFlagsThree": {
		FileName:      "test_files",
		SearchPattern: "test",
		Want: []string{
			"this is test 8",
			"this is line 9",
			"this is line 10",
			"this is test 14",
			"this is test 15",
			"this is test 16",
			"this is line 17",
			"this is line 18",
		},
		flagConfigIo: GrepConfig{
			CaseInsensitiveSearch:  true,
			CountOfLinesAfterMatch: 2,
			ResurciveSearch:        true,
		},
	},
}

func TestGrep(t *testing.T) {
	for key, value := range grepTestCases {
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

				fileResult, _ := readAndMatch(&ReadAndMatchInput{
					Reader:     file,
					FlagConfig: &value.flagConfigIo,
					Pattern:    value.SearchPattern,
				})
				got = append(got, fileResult...)
			}

			if !reflect.DeepEqual(got, value.Want) {
				t.Errorf("got %s \n --- want %s ", got, value.Want)
			}

		})
	}
}

var userInputTestCases = map[string]struct {
	Input         string
	SearchPattern string
	Want          []string
}{
	"zeroMatches": {
		Input:         "this does not contain the word.\nthis is empty",
		SearchPattern: "someRandomString",
	},
	"oneMatch": {
		Input:         "this is temperature.\nthis is one match",
		SearchPattern: "temperature",
		Want:          []string{"this is temperature."},
	},
	"multipleMatches": {
		Input:         "this is anish.\nis this anish.\nthis is anish?\nanish",
		SearchPattern: "anish",
		Want:          []string{"this is anish.", "is this anish.", "this is anish?", "anish"},
	},
}

func TestUserInput(t *testing.T) {
	for key, value := range userInputTestCases {
		t.Run(key, func(t *testing.T) {
			file, err := os.CreateTemp("", "tempfile")
			if err != nil {
				log.Fatalf(err.Error())
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(value.Input)); err != nil {
				log.Fatalf(err.Error())
			}

			if _, err := file.Seek(0, 0); err != nil {
				log.Fatalf(err.Error())
			}

			oldStdIn := os.Stdin
			os.Stdin = file

			defer func() {
				os.Stdin = oldStdIn
			}()

			gotResult, _ := readAndMatch(
				&ReadAndMatchInput{
					Reader: os.Stdin,
					FlagConfig: &GrepConfig{
						CountOfLinesBeforeMatch: 0,
						CountOfLinesAfterMatch:  0,
					},
					Pattern: value.SearchPattern,
				},
			)

			if !reflect.DeepEqual(gotResult, value.Want) {
				t.Errorf("got %s \n --- want %s ", gotResult, value.Want)
			}
		})
	}
}

func BenchmarkTableRegex(b *testing.B) {
	for key, value := range grepTestCases {
		b.Run(key, func(b *testing.B) {
			file, err := os.Open(value.FileName)
			if err != nil {
				log.Fatalf(err.Error())
			}
			defer file.Close()

			for i := 0; i < b.N; i++ {
				readAndMatch(&ReadAndMatchInput{
					Reader:     file,
					FlagConfig: &value.flagConfigIo,
					Pattern:    value.SearchPattern,
				})
			}
		})
	}

}

var listDirTestCases = map[string]struct {
	Directory string
	Want      []string
}{
	"test_files": {
		Directory: "test_files",
		Want: []string{
			"test_files/parent_dir1/child_dir1/child_dir1_file.txt",
			"test_files/parent_dir1/child_dir2/child_dir2_file.txt",
			"test_files/parent_dir2/parent_dir2_file1.txt",
			"test_files/testfile.txt",
			"test_files/testfile2.txt",
		},
	},
}

func TestWalk(t *testing.T) {
	for key, value := range listDirTestCases {
		t.Run(key, func(t *testing.T) {
			gotSubFiles, _, err := listFilesInDir(value.Directory)
			if err != nil && !os.IsPermission(err) {
				log.Fatalf(err.Error())
			}

			if !reflect.DeepEqual(gotSubFiles, value.Want) {
				t.Errorf("got %s \n --- want %s ", gotSubFiles, value.Want)
			}
		})
	}
}

var bufferTestCases = map[string]struct {
	DataBuffer          []string
	CurrentLine         int
	LinesPrinted        map[int]bool
	ExpectedMatchResult []string
}{
	"emptyBuffer": {
		DataBuffer:          []string{},
		CurrentLine:         5,
		LinesPrinted:        map[int]bool{},
		ExpectedMatchResult: []string{},
	},
	"bufferWithUnprintedLines": {
		DataBuffer:          []string{"line1", "line2", "line3"},
		CurrentLine:         5,
		LinesPrinted:        map[int]bool{},
		ExpectedMatchResult: []string{"line1", "line2", "line3"},
	},
	"bufferWithSomePrintedLines": {
		DataBuffer:          []string{"line1", "line2", "line3"},
		CurrentLine:         5,
		LinesPrinted:        map[int]bool{3: true},
		ExpectedMatchResult: []string{"line1", "line3"},
	},
	"bufferWithAllLinesPrinted": {
		DataBuffer:          []string{"line1", "line2", "line3"},
		CurrentLine:         5,
		LinesPrinted:        map[int]bool{4: true, 2: true, 3: true},
		ExpectedMatchResult: []string{},
	},
}

func TestBufferForOptionB(t *testing.T) {
	for key, value := range bufferTestCases {
		t.Run(key, func(t *testing.T) {
			matchResult := fetchResultFromBuffer(value.CurrentLine, value.DataBuffer, value.LinesPrinted)

			if !reflect.DeepEqual(matchResult, value.ExpectedMatchResult) {
				t.Errorf("matchResult = %v, want %v", matchResult, value.ExpectedMatchResult)
			}
		})
	}
}
