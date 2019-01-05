package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
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
		exit(err.Error())
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		exit(err.Error())
	}

	ret := make([]QuizItem, 0)

	for _, line := range lines {
		data := QuizItem{
			Question: line[0],
			Answer:   strings.TrimSpace(line[1]),
		}
		ret = append(ret, data)
	}
	return ret
}

func exit(msg string) {
	fmt.Printf("%s\n", msg)
	os.Exit(1)
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

func scoreExam(exam *Exam) {
	score, correct := score(*exam)
	fmt.Printf("You answered %d of %d correctly. Score: %.2f%%\n", correct, len(exam.Questions), score)
}

func main() {
	var timeLimit time.Duration
	csvFileName := flag.String("csv", "problems.csv", "a CSV file in 'question,answer' format")
	flag.DurationVar(&timeLimit, "time-limit", time.Second*30, "number of seconds until the exam ends")
	flag.Parse()

	exam := Exam{Questions: readCsv(csvFileName), Score: 100}

	fmt.Printf("The Exam will begin now. There are %d questions, you have %s time remaining.\n", len(exam.Questions), timeLimit)
	_ = prompt("Press [Enter] to begin")

	quit := make(chan struct{})
	complete := make(chan struct{})

	go func() {
		time.Sleep(timeLimit)
		close(quit)
	}()

	go func() {

		for index, question := range exam.Questions {

			str := fmt.Sprintf("%d.  What is %s ?  ", index+1, question.Question)
			exam.Questions[index].Answered = prompt(str)
		}
		close(complete)
	}()

	for {
		select {
		case <-quit:
			fmt.Println("Out of time")
			scoreExam(&exam)
			return
		case <-complete:
			scoreExam(&exam)
			return
		}
	}

}
