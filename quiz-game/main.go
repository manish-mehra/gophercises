package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func check(e error, warning string) {
	if e != nil {
		fmt.Println(warning, e)
		os.Exit(1)
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

func shuffleQuiz(problems []problem) []problem {

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	min := 0
	max := len(problems) - 1

	var temp problem

	for i := 0; i < len(problems); i++ {
		randomIndex := random.Intn(max-min+1) + min
		temp = problems[i]
		problems[i] = problems[randomIndex]
		problems[randomIndex] = temp
	}

	return problems
}

func main() {

	fileName := flag.String("fName", "problems_1.csv", "a file")
	timeLimit := flag.Int("limit", 3, "time limit in seconds")
	shuffle := flag.Bool("shuffle", false, "flag to shuffle quiz")
	flag.Parse()

	file := openCSVFile(*fileName)
	defer file.Close()

	records := readCSV(file)
	problems := parseLines(records)

	if *shuffle {
		fmt.Println("suffled")
		problems = shuffleQuiz(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemloop:
	for index, problem := range problems {
		fmt.Printf("problem #%d: %s = ", index+1, problem.question)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d. \n", correct, len(records))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for index, line := range lines {
		problems[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

type problem struct {
	question string
	answer   string
}
