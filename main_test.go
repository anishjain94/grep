package main

import (
	"reflect"
	"testing"
)

func TestNaiveGrep(t *testing.T) {

	testCases := map[string]struct {
		FileName  string
		SearchStr string
		Want      []string
	}{
		"Zero matches": {
			FileName:  "testfile.txt",
			SearchStr: "someRandomString",
			// Want:      []string,
		},
		"One matches": {
			FileName:  "testfile.txt",
			SearchStr: "temp",
			Want:      []string{"this is temp."},
		},
		"Multiple matches": {
			FileName:  "testfile.txt",
			SearchStr: "anish",
			Want:      []string{"this is anish.", "is this anish.", "this is anish?"},
		},
	}

	for key, value := range testCases {
		t.Run(key, func(t *testing.T) {
			got := naiveGrep(value.FileName, value.SearchStr)

			if !reflect.DeepEqual(got, value.Want) {
				t.Errorf("%s got %s want", got, value.Want)
			}
		})
	}

}

func Test_naiveGrep(t *testing.T) {
	type args struct {
		filePath  string
		searchStr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := naiveGrep(tt.args.filePath, tt.args.searchStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("naiveGrep() = %v, want %v", got, tt.want)
			}
		})
	}
}
