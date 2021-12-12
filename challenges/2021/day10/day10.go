package aoc2021_day10

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
)

func part1(input []string) (string, error) {
	lines := parseInput(input)

	pointVals := map[string]int{
		")": 3,
		"]": 57,
		"}": 1197,
		">": 25137,
	}

	score := 0
	for _, line := range lines {
		stack := NewStack()
		for _, delim := range line {
			if IsOpenerDelim(delim) {
				stack.Push(delim)
			} else {
				opener := stack.Pop()
				if delimMatches[opener] != delim {
					score += pointVals[delim]
					break
				}
			}
		}
	}

	return fmt.Sprintf("%d", score), nil
}

func part2(input []string) (string, error) {
	lines := parseInput(input)

	pointVals := map[string]int{
		")": 1,
		"]": 2,
		"}": 3,
		">": 4,
	}

	scores := make([]int, 0)
	for _, line := range lines {
		stack := NewStack()
		invalid := false
		for _, delim := range line {
			if IsOpenerDelim(delim) {
				stack.Push(delim)
			} else {
				opener := stack.Pop()
				if delimMatches[opener] != delim {
					invalid = true
					break
				}
			}
		}
		if invalid {
			continue
		}

		score := 0
		for !stack.IsEmpty() {
			score = score*5 + pointVals[delimMatches[stack.Pop()]]
		}
		scores = append(scores, score)
	}

	sort.Ints(scores)
	middleScore := scores[len(scores)/2]

	return fmt.Sprintf("%d", middleScore), nil
}

func parseInput(input []string) [][]string {
	lines := make([][]string, 0, len(input))
	for _, s := range input {
		lines = append(lines, strings.Split(s, ""))
	}
	return lines
}

type delimStack []string

func NewStack() delimStack {
	return make(delimStack, 0)
}

func (s *delimStack) Push(v string) {
	*s = append(*s, v)
}

func (s *delimStack) Pop() string {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *delimStack) IsEmpty() bool {
	return len(*s) == 0
}

var delimMatches = map[string]string{
	"{": "}",
	"[": "]",
	"(": ")",
	"<": ">",
}

func IsOpenerDelim(delim string) bool {
	_, ok := delimMatches[delim]
	return ok
}

func init() {
	challenges.RegisterChallengeFunc(2021, 10, 1, "day10.txt", part1)
	challenges.RegisterChallengeFunc(2021, 10, 2, "day10.txt", part2)
}
