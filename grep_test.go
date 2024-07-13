package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"unsafe"
)

// TODO: write seperate test files for there spefici files/modules..
// TODO: handle for condition when file limit opening is restricted by os. make is os independent.
// TODO: How to write test case for permission denied. ex - a directory exists in test_files, which does not have read permission, now because of this..my other test cases were failing..
// TODO: How are you load-testing your implementation of grep

var testCases = map[string]struct {
	FileName     string
	SearchStr    string
	Want         []string
	flagConfigIo GrepConfigIo
}{
	"zeroMatch": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "someRandomString",
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: false,
		},
	},
	"oneMatch": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: false,
		},
	},
	"fileDoesNotExists": {
		FileName:  "fileDoesNotExist.txt",
		SearchStr: "temperature",
		Want:      []string{"lstat fileDoesNotExist.txt: no such file or directory"},
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: false,
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
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: false,
		},
	},
	"oneMatchCaseInsensitive": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "Temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: true,
		},
	},
	"oneMatchOutputFile": {
		FileName:  "test_files/testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
		flagConfigIo: GrepConfigIo{
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
		flagConfigIo: GrepConfigIo{
			CaseInsensitiveSearch: false,
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
			"this is test 15",
			"this is test 16",
		},
		flagConfigIo: GrepConfigIo{
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
			"this is test 15",
			"this is test 16",
			"this is line 17",
			"this is line 18",
		},
		flagConfigIo: GrepConfigIo{
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

				fileResult, _ := readAndMatch(&ReadAndMatchIo{
					Reader:     file,
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
				log.Fatalf(err.Error())
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(value.InputStr)); err != nil {
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
				&ReadAndMatchIo{
					Reader: os.Stdin,
					FlagConfig: &GrepConfigIo{
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
		b.Run(fmt.Sprintf("naive-%s", key), func(b *testing.B) {
			file, err := os.Open(value.FileName)
			log.Fatalf(err.Error())
			defer file.Close()

			for i := 0; i < b.N; i++ {
				readAndMatch(&ReadAndMatchIo{
					Reader: file,
					FlagConfig: &GrepConfigIo{
						CountOfLinesBeforeMatch: 0,
						CountOfLinesAfterMatch:  0,
					},
					Pattern: value.SearchStr,
				})
			}
		})
	}
}

func TestWalk(t *testing.T) {
	subFiles, _, err := listFilesInDir("test_files")
	if err != nil && !os.IsPermission(err) {
		log.Fatalf(err.Error())
	}
	fmt.Println(subFiles)
}

func estimateFileResultMapSize(m FileResultMap) uintptr {
	var size uintptr

	// Size of map structure itself
	size += unsafe.Sizeof(m)

	// Iterate through the map to calculate size of contents
	for k, v := range m {
		// Size of key (string)
		size += unsafe.Sizeof(k)
		size += uintptr(len(k))

		// Size of value ([]string)
		size += unsafe.Sizeof(v)

		// Size of each string in the slice
		for _, s := range v {
			size += unsafe.Sizeof(s)
			size += uintptr(len(s))
		}
	}

	return size
}

func TestMem(t *testing.T) {

	resultChannel := make(chan FileResultMap, 10000)
	fmt.Println(resultChannel)
	fmt.Println(unsafe.Sizeof(resultChannel))

	resultChannel <- FileResultMap{"ads": []string{"temp"}}
	resultChannel <- FileResultMap{"ad1": []string{"temp"}}
	resultChannel <- FileResultMap{"ad2": []string{"temp"}}

	fmt.Println(estimateFileResultMapSize(<-resultChannel))
}

func TestBufferChecking(t *testing.T) {
	tests := []struct {
		name                string
		dataBuffer          []string
		lineNumber          int
		linesPrintedMap     map[int]bool
		expectedMatchResult []string
	}{
		{
			name:                "Empty buffer",
			dataBuffer:          []string{},
			lineNumber:          5,
			linesPrintedMap:     map[int]bool{},
			expectedMatchResult: []string{},
		},
		{
			name:                "Buffer with unprintedlines",
			dataBuffer:          []string{"line1", "line2", "line3"},
			lineNumber:          5,
			linesPrintedMap:     map[int]bool{},
			expectedMatchResult: []string{"line1", "line2", "line3"},
		},
		{
			name:                "Buffer with some printed lines",
			dataBuffer:          []string{"line1", "line2", "line3"},
			lineNumber:          5,
			linesPrintedMap:     map[int]bool{3: true},
			expectedMatchResult: []string{"line1", "line3"},
		},
		{
			name:                "Buffer with all lines printed",
			dataBuffer:          []string{"line1", "line2", "line3"},
			lineNumber:          5,
			linesPrintedMap:     map[int]bool{4: true, 2: true, 3: true},
			expectedMatchResult: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchResult := fetchResultFromBuffer(tt.lineNumber, tt.dataBuffer, tt.linesPrintedMap)

			if !reflect.DeepEqual(matchResult, tt.expectedMatchResult) {
				t.Errorf("matchResult = %v, want %v", matchResult, tt.expectedMatchResult)
			}
		})
	}
}
