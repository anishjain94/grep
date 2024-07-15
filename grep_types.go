package main

import "io"

type FileResult map[string][]string

type GrepConfig struct {
	CaseInsensitiveSearch   bool   //case-inSensitive search
	CountOfMatches          bool   //displays count of matches
	CountOfLinesBeforeMatch int    //displays n lines before the match
	CountOfLinesAfterMatch  int    //displays n lines after the match
	OutputFileName          string //output file
	ResursiveSearch         bool   //search in a directory
}

func (flagConfig *GrepConfig) shouldDisplayLinesBeforeMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesBeforeMatch != 0
}

func (flagConfig *GrepConfig) shouldDisplayLinesAfterMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesAfterMatch != 0
}

func (flagConfig *GrepConfig) shouldShowCount() bool {
	return flagConfig != nil && flagConfig.CountOfMatches
}

func (flagConfig *GrepConfig) shouldSearchCaseInsensitive() bool {
	return flagConfig != nil && flagConfig.CaseInsensitiveSearch
}

func (flagConfig *GrepConfig) shouldStoreOutput() bool {
	return flagConfig != nil && flagConfig.OutputFileName != ""
}

// TODO: Ask mohit, do i make first letter lowercase, since we are not exporting or using it in other modules. 
type DisplayResultInput struct {
	MatchedResult    map[string][]string
	FlagConfig       *GrepConfig
	FilesInDirectory []string
	IsDirectory      bool
}

type ReadAndMatchInput struct {
	Reader     io.Reader
	FlagConfig *GrepConfig
	Pattern    string
}
