package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type quizQuestionAnswer struct {
	question string
	answer   string
}

func processQuestions(fileName string) ([]*quizQuestionAnswer, error) {
	var quizQuestions []*quizQuestionAnswer
	csvfile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(fmt.Sprintf("error reading file %v", err))
		return nil, err
	}
	fileReader := csv.NewReader(csvfile)
	for {
		// Read each record from csv
		record, err := fileReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("error reading file %v", err))
		}
		quizQuestions = append(quizQuestions, &quizQuestionAnswer{
			question: record[0],
			answer:   record[1],
		})
	}
	return quizQuestions, nil
}

func runQuiz(quizQuestions []*quizQuestionAnswer, shuffle bool) int {
	score := 0
	if shuffle {
		shuffleQuiz(quizQuestions)
	}
	for index, quizQuestion := range quizQuestions {
		fmt.Print(fmt.Sprintf("%d. %v = ", index+1, quizQuestion.question))
		timer := time.NewTimer(5 * time.Second)
		answers := make(chan string)
		go scanUserkeyedInput(answers)
		select {
		case answer := <-answers:
			timer.Stop()
			if answer == quizQuestion.answer {
				score = score + 1
			}
		case <-timer.C:
			return score
		}
	}
	return score
}

func scanUserkeyedInput(answers chan<- string) {
	answerInput := bufio.NewScanner(os.Stdin)
	answerInput.Scan()
	answers <- answerInput.Text()
}

func shuffleQuiz(quizQuestions []*quizQuestionAnswer) {
	for i := range quizQuestions {
		j := rand.Intn(i + 1)
		quizQuestions[i], quizQuestions[j] = quizQuestions[j], quizQuestions[i]
	}
}
func main() {
	csvFileName := flag.String("csv", "/Users/kiruba.duraisamy/gopath/src/personalPrograms/gophercises/1/resources/problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()
	quizQuestions, _ := processQuestions(*csvFileName)
	score := runQuiz(quizQuestions, true)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("you scored %d ", score))
}
