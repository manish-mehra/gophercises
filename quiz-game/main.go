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

func main() {

	fileName := flag.String("fName", "problems_1.csv", "a file")
	flag.Parse()

	file, err := os.Open(*fileName)
	check(err, "Error opening csv file")
	defer file.Close()

	// create csv reader
	reader := csv.NewReader(file)

	records, rError := reader.ReadAll()
	check(rError, "Error reading csv")

	totalScore := 0

	inputReader := bufio.NewReader(os.Stdin)

	for _, record := range records {

		question := record[0]
		answer := record[1]
		fmt.Println("question: ", question)
		userInput, err := inputReader.ReadString('\n')
		check(err, "Error reading input")

		userInput = strings.TrimSpace(userInput)

		if userInput == answer {
			totalScore++
		}
	}

	fmt.Printf("%d out of %d correct", totalScore, len(records))
}
