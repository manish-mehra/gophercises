package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func check(e error, warning string) {
	if e != nil {
		fmt.Println(warning, e)
		return
	}
}

func openCSVFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	check(err, "Error opening CSV file")
	return file
}

func readCSV(file *os.File) [][]string {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err, "Error reading CSV file")
	return records
}

func getUserInput() string {
	inputReader := bufio.NewReader(os.Stdin)
	userInput, err := inputReader.ReadString('\n')
	check(err, "Error reading input")
	return strings.TrimSpace(userInput)
}

func processQuestions(records [][]string) int {
	totalScore := 0

	for _, record := range records {
		question := record[0]
		answer := record[1]
		fmt.Println("question: ", question)
		userInput := getUserInput()
		if userInput == answer {
			totalScore++
		}
	}
	return totalScore
}

func main() {

	fileName := flag.String("fName", "problems_1.csv", "a file")
	flag.Parse()

	file := openCSVFile(*fileName)
	defer file.Close()

	records := readCSV(file)

	totalScore := processQuestions(records)
	fmt.Printf("%d out of %d correct", totalScore, len(records))
}
