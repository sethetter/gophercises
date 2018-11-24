package quiz

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		in   [][]string
		want []Problem
	}{
		{
			in: [][]string{
				{"2+2", "4"},
				{"5+2", "7"},
			},
			want: []Problem{
				{q: "2+2", a: "4"},
				{q: "5+2", a: "7"},
			},
		},
	}

	for i, test := range tests {
		quiz := New(test.in, make(chan time.Time))
		for j, p := range test.want {
			if quiz.Problems[j] != p {
				t.Errorf("#%d: expected %q, got %q", i, p, quiz.Problems[j])
			}
		}
	}
}

func TestShowQuestion(t *testing.T) {
	tests := []struct {
		in   Problem
		want string
	}{
		{
			in:   Problem{q: "Sup?", a: "yo"},
			want: "Problem #1: Sup? = ",
		},
		{
			in:   Problem{q: "5 + 5", a: "yo"},
			want: "Problem #2: 5 + 5 = ",
		},
	}

	for i, test := range tests {
		got := ShowQuestion(i, &test.in)
		if got != test.want {
			t.Errorf("#%d: expected %s, got %s", i, test.want, got)
		}
	}
}

func TestAskQuestions(t *testing.T) {
	lines := [][]string{
		{"5+5", "10"},
		{"2+5", "7"},
	}
	timerC := make(chan time.Time)
	quiz := New(lines, timerC)

	output := testAskQuestions(t, quiz, "10\n7\n")

	// check expected correct answers
	expected := 2
	if quiz.Correct != expected {
		t.Errorf("got %d, expected %d", quiz.Correct, expected)
	}

	// check output
	expectedOut := "Problem #1: 5+5 = "
	expectedOut += "Problem #2: 2+5 = "
	expectedOut += "Score: 2/2\n"
	if output != expectedOut {
		t.Errorf("got %v, expected %v", output, expectedOut)
	}
}

func TestAskQuestionsTimeout(t *testing.T) {
	lines := [][]string{
		{"5+5", "10"},
		{"2+5", "7"},
	}
	timerC := make(chan time.Time)
	quiz := New(lines, timerC)

	// Simulate timer end
	go func() { timerC <- time.Now() }()

	output := testAskQuestions(t, quiz, "")

	// check expected correct answers
	expected := 0
	if quiz.Correct != expected {
		t.Errorf("got %d, expected %d", quiz.Correct, expected)
	}

	// check output
	expectedOut := "Problem #1: 5+5 = "
	expectedOut += "Score: 0/2\n"
	if output != expectedOut {
		t.Errorf("got %v, expected %v", output, expectedOut)
	}
}

// Test function that takes an input and then simulates
// asking questions on a quiz struct, returning the output
func testAskQuestions(t *testing.T, quiz *Quiz, input string) string {
	var out bytes.Buffer
	r := bufio.NewReader(strings.NewReader(input))
	w := bufio.NewWriter(&out)

	quiz.AskQuestions(r, w)

	// check expected score
	return out.String()
	// check expected output
}
