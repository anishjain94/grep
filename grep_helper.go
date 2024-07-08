package main

import "strings"

func printError(err error) {
	if err != nil {
		errMsg := err.Error()
		displayableErrorMsg := ""
		if strings.Contains(errMsg, "The system cannot find the file specified.") {
			displayableErrorMsg = "No such file or directory"
		}

		print(displayableErrorMsg)
	}
}
