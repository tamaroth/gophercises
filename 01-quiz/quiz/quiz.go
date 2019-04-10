package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"../args"
)

// A Quiz instance.
type Quiz struct {
	questions        map[int]string
	answers          map[int]string
	order            []int
	questionsTotal   int
	questionsAsked   int
	questionsCorrect int
	timeout          int
}

func (q *Quiz) addQuestion(id int, question, answer string) {
	q.order = append(q.order, id)
	q.questions[id] = question
	q.answers[id] = answer
	q.questionsTotal++
}

func (q *Quiz) askQuestion(id int, answered chan bool) {
	var answer string
	fmt.Printf("%s\n", q.questions[q.order[id]])
	fmt.Scanln(&answer)
	if q.verifyAnswer(id, answer) {
		q.questionsCorrect++
	}
	q.questionsAsked++
	answered <- true
}

func (q *Quiz) loadQuizFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	id := 0
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		q.addQuestion(id, row[0], row[1])
		id++
	}
}

func (q *Quiz) randomizeQuestions() {
	rand.Seed(time.Now().UTC().UnixNano())
	tmp := make([]int, len(q.order))
	perm := rand.Perm(len(q.order))
	for i, v := range perm {
		tmp[v] = q.order[i]
	}
	q.order = tmp
}

func (q *Quiz) verifyAnswer(id int, answer string) bool {
	if strings.EqualFold(
		strings.TrimSpace(q.answers[q.order[id]]),
		strings.TrimSpace(answer),
	) {
		return true
	}
	return false
}

// NewQuiz returns a new quiz.
func NewQuiz(arguments args.Args) (Quiz, error) {
	q := Quiz{
		map[int]string{},
		map[int]string{},
		[]int{},
		0,
		0,
		0,
		arguments.Timeout,
	}
	err := q.loadQuizFromFile(arguments.Filename)
	if arguments.Randomize {
		q.randomizeQuestions()
	}
	return q, err
}

// RunQuiz runs the quiz with the loaded questions.
func (q *Quiz) RunQuiz() {
	fmt.Print("Press 'Enter' to start quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	timesUp := time.After(time.Second * time.Duration(q.timeout))
	answered := make(chan bool)
	for id := 0; id < q.questionsTotal; id++ {
		go q.askQuestion(id, answered)
		select {
		case <-timesUp:
			fmt.Println("Time's up!")
			return
		case <-answered:
			// Do nothing, just consume.
		}
	}
}

// GetResults returns the result of the quiz.
func (q *Quiz) GetResults() {
	fmt.Printf(
		"Results:\n\tTotal questions: %d\n\tQuestions asked: %d\n\tCorrect answers: %d\n",
		q.questionsTotal,
		q.questionsAsked,
		q.questionsCorrect,
	)
}
