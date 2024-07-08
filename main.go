package main


import (
    "bufio"
    "fmt"
    "os"
    "strings"
)


func main() {


    filePath := "testfile1.txt"
    toSearch := "Anish jain"


    output := naiveGrep(filePath, toSearch)


    if len(output) > 0 {
        fmt.Println(output)
    }
}


func naiveGrep(filePath string, searchStr string) []string {
    file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)


    printError(err)
    defer file.Close()


    var outputLines []string
    scanner := bufio.NewScanner(file)


    for scanner.Scan() {
        currentLine := scanner.Text()
        if strings.Contains(currentLine, searchStr) {
            outputLines = append(outputLines, currentLine)
        }
    }


    return outputLines
}
