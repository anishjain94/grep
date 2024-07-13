package main

import "io"

type DisplayResultIo struct {
	matchedResultMap map[string][]string
	FlagConfig       *GrepConfigIo
	FilesInDirectory []string
	IsDirectory      bool
}

type ReadAndMatchIo struct {
	Reader     io.Reader
	FlagConfig *GrepConfigIo
	Pattern    string
}

type GrepConfigIo struct {
	CaseInsensitiveSearch   bool   //case-inSensitive search
	CountOfMatches          bool   //displays count of matches
	CountOfLinesBeforeMatch int    //displays n lines before the match
	CountOfLinesAfterMatch  int    //displays n lines after the match
	OutputFileName          string //output file
}

type FileResultMap map[string][]string

func (flagConfig *GrepConfigIo) shouldDisplayLinesBeforeMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesBeforeMatch != 0
}

func (flagConfig *GrepConfigIo) shouldDisplayLinesAfterMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesAfterMatch != 0
}

func (flagConfig *GrepConfigIo) shouldShowCount() bool {
	return flagConfig != nil && flagConfig.CountOfMatches
}

func (flagConfig *GrepConfigIo) shouldSearchCaseInsensitive() bool {
	return flagConfig != nil && flagConfig.CaseInsensitiveSearch
}

func (flagConfig *GrepConfigIo) shouldStoreOutput() bool {
	return flagConfig != nil && flagConfig.OutputFileName != ""
}
