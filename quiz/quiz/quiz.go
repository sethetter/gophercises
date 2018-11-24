package quiz

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Quiz struct {
	timerC   <-chan time.Time
	wg       *sync.WaitGroup
	Problems []Problem
	Correct  int
}

type Problem struct {
	q string
	a string
}

func New(lines [][]string, timerC <-chan time.Time) *Quiz {
	correct := 0

	wg := &sync.WaitGroup{}
	wg.Add(1)

	quiz := &Quiz{
		timerC:  timerC,
		wg:      wg,
		Correct: correct,
	}

	quiz.ParseLines(lines)

	return quiz
}

func (q *Quiz) ParseLines(lines [][]string) {
	q.Problems = make([]Problem, len(lines))

	for i, line := range lines {
		q.Problems[i] = Problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
}

func (q *Quiz) ShowScore() string {
	return fmt.Sprintf("Score: %d/%d\n", q.Correct, len(q.Problems))
}

func (q *Quiz) AskQuestions(r *bufio.Reader, w *bufio.Writer) {
	go q.processQuestions(r, w)
	q.wg.Wait()
	w.WriteString(q.ShowScore())
	w.Flush()
}

func (q *Quiz) processQuestions(r *bufio.Reader, w *bufio.Writer) {
	for i, p := range q.Problems {
		w.WriteString(ShowQuestion(i, &p))
		w.Flush()

		answerCh := make(chan string)
		go func() {
			answer, err := r.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			answerCh <- answer
		}()

		select {
		case <-q.timerC:
			q.wg.Done()
			return
		case answer := <-answerCh:
			answer = strings.ToUpper(strings.TrimSpace(answer))
			given := strings.ToUpper(strings.TrimSpace(p.a))
			if answer == given {
				q.Correct++
			}
		}
	}
	q.wg.Done()
}

func ShowQuestion(i int, p *Problem) string {
	return fmt.Sprintf("Problem #%d: %s = ", i+1, p.q)
}
