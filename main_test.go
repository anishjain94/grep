package main

import (
	"fmt"
	"testing"
)

var testCases = map[string]struct {
	FileName  string
	SearchStr string
	Want      []string
}{
	"Zero matches": {
		FileName:  "testfile.txt",
		SearchStr: "someRandomString",
	},
	"One match": {
		FileName:  "testfile.txt",
		SearchStr: "temperature",
		Want:      []string{"this is temperature."},
	},
	"Multiple matches": {
		FileName:  "testfile.txt",
		SearchStr: "anish",
		Want:      []string{"this is anish.", "is this anish.", "this is anish?", "anish"},
	},
}

// func TestNaiveGrep(t *testing.T) {
// 	for key, value := range testCases {
// 		t.Run(key, func(t *testing.T) {
// 			got := naiveGrep(value.FileName, value.SearchStr)

// 			if !reflect.DeepEqual(got, value.Want) {
// 				t.Errorf("got %s \n --- want %s ", got, value.Want)
// 			}
// 		})
// 	}
// }

// func TestRegexGrep(t *testing.T) {
// 	for key, value := range testCases {
// 		t.Run(key, func(t *testing.T) {
// 			got := regexGrep(value.FileName, value.SearchStr)

// 			if !reflect.DeepEqual(got, value.Want) {
// 				t.Errorf("got %s \n --- want %s ", got, value.Want)
// 			}
// 		})
// 	}
// }

func BenchmarkTableRegex(b *testing.B) {
	for key, value := range testCases {
		b.Run(fmt.Sprintf("regex-%s", key), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				regexGrep(value.FileName, value.SearchStr)
			}
		})
		b.Run(fmt.Sprintf("naive-%s", key), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				naiveGrep(value.FileName, value.SearchStr)
			}
		})
	}
}

func BenchmarkGrepString(b *testing.B) {
	fileName := "testfile.txt"
	searchStr := "anish"
	for i := 0; i < b.N; i++ {
		naiveGrep(fileName, searchStr)
	}
}

func BenchmarkGrepRegex(b *testing.B) {
	fileName := "testfile.txt"
	searchStr := "anish"
	for i := 0; i < b.N; i++ {
		regexGrep(fileName, searchStr)
	}
}
