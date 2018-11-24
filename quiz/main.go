package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sethetter/gophercises/quiz/quiz"
)

func main() {
	csvFilename := flag.String("csv", "filename", "a csv file in the format of question,answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Could not parse CSV")
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	quiz := quiz.New(lines, timer.C)
	quiz.AskQuestions(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
