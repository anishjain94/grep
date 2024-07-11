package main

import "io"

type DisplayResultIo struct {
	matchedResultMap map[string][]string
	FlagConfig       *FlagConfigIo
	FilesInDirectory []string
	IsDirectory      bool
}

type ReadAndMatchIo struct {
	Reader     io.Reader
	Source     *string
	FlagConfig *FlagConfigIo
	Pattern    string
}

type FlagConfigIo struct {
	CaseInSensitiveSearch   bool   //case-inSensitive search
	CountOfMatches          bool   //displays count of matches
	CountOfLinesBeforeMatch int    //displays n lines before the match
	CountOfLinesAfterMatch  int    //displays n lines after the match
	OutputFileName          string //output file
}

type FileResultMap map[string][]string

func (flagConfig *FlagConfigIo) shouldDisplayLinesBeforeMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesBeforeMatch != 0
}

func (flagConfig *FlagConfigIo) shouldDisplayLinesAfterMatch() bool {
	return flagConfig != nil && flagConfig.CountOfLinesAfterMatch != 0
}

func (flagConfig *FlagConfigIo) shouldShowCount() bool {
	return flagConfig != nil && flagConfig.CountOfMatches
}

func (flagConfig *FlagConfigIo) shouldSearchCaseInsensitive() bool {
	return flagConfig != nil && flagConfig.CaseInSensitiveSearch
}

func (flagConfig *FlagConfigIo) shouldStoreOutput() bool {
	return flagConfig != nil && flagConfig.OutputFileName != ""
}
