package main

// open file
// read lines
// create csv objects

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

type QuizItem struct {
	Question string
	Answer   string
	Answered string
}

type Exam struct {
	Questions []QuizItem
	Score     float64
}

func readCsv(filename *string) []QuizItem {
	path := *filename
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	ret := make([]QuizItem, 0)

	for _, line := range lines {
		data := QuizItem{
			Question: line[0],
			Answer:   line[1],
		}
		ret = append(ret, data)
	}
	return ret
}

func prompt(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s", question)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func score(exam Exam) (float64, int) {
	var correct = 0
	for _, question := range exam.Questions {
		if question.Answer == question.Answered {
			correct += 1
		}
	}
	var score float64
	score = float64(correct) / float64(len(exam.Questions)) * 100.0
	return math.Round(score*100) / 100, correct
}

func main() {

	csvFileName := flag.String("csv", "problems.csv", "a CSV file in 'question,answer' format")
	timeLimit := flag.Int("time-limit", 0, "number of seconds until the exam ends")
	flag.Parse()

	_ = timeLimit

	exam := Exam{Questions: readCsv(csvFileName), Score: 100}

	fmt.Printf("The Exam will begin now. There are %d questions\n", len(exam.Questions))
	for index, question := range exam.Questions {
		str := fmt.Sprintf("%d.  What is %s ?  ", index+1, question.Question)
		exam.Questions[index].Answered = prompt(str)
	}
	score, correct := score(exam)
	fmt.Printf("You answered %d of %d correctly. Score: %.2f%%\n", correct, len(exam.Questions), score)
}
