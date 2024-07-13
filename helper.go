package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func validateFile(filepath string) error {
	if _, err := os.Stat(filepath); err != nil {
		return err
	}
	return nil
}

func listFilesInDir(path string) ([]string, bool, error) {
	var subFiles []string
	var isDir bool

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			subFiles = append(subFiles, path)
		} else {
			isDir = true
		}
		return nil
	})

	if err != nil {
		return subFiles, isDir, err //To not error out and exit completely incase we encounter file permission error
	}

	return subFiles, isDir, nil
}
