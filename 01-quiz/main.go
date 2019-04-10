package main

import (
	"fmt"

	"./args"
	"./quiz"
)

func main() {
	arguments := args.ParseCommandlineArgs()
	q, e := quiz.NewQuiz(arguments)
	if e != nil {
		fmt.Printf("Could not load data from file %s.\n", arguments.Filename)
		return
	}
	q.RunQuiz()
	q.GetResults()
}
