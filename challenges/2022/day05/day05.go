package aoc2022_day5

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
)

func part1(input []string) (string, error) {
	stacks, instructions := parseInput(input)

	for _, instr := range instructions {
		for i := 0; i < instr.N; i++ {
			str := stacks[instr.From-1].Pop()
			stacks[instr.To-1].Push(str)
		}
	}
	var result strings.Builder
	for _, s := range stacks {
		if len(s) != 0 {
			result.WriteString(s[len(s)-1])
		}
	}

	return result.String(), nil
}

func part2(input []string) (string, error) {
	stacks, instructions := parseInput(input)

	for _, instr := range instructions {
		letters := stacks[instr.From-1].PopN(instr.N)
		stacks[instr.To-1].PushAll(letters)
	}

	var result strings.Builder
	for _, s := range stacks {
		if len(s) != 0 {
			result.WriteString(s[len(s)-1])
		}
	}

	return result.String(), nil
}

func parseInput(input []string) ([]stack, []instruction) {
	numStacks := int((len(input[0]) + 1) / 4)
	stacks := make([]stack, numStacks) // pre-allocate
	for i, s := range input {
		if len(s) == 0 {
			input = input[i+1:]
			break
		}

		for j := 0; j < numStacks; j++ {
			k := j * 4
			if s[k] == '[' {
				stacks[j].LPush(s[k+1 : k+2])
			}
		}
	}

	instructions := make([]instruction, 0, len(input))
	for _, s := range input {
		var instr instruction
		fmt.Sscanf(s, "move %d from %d to %d", &instr.N, &instr.From, &instr.To)
		instructions = append(instructions, instr)
	}

	return stacks, instructions
}

type instruction struct {
	N    int
	From int
	To   int
}

type stack []string

func (s *stack) Push(str string) {
	*s = append(*s, str)
}

func (s *stack) Pop() string {
	str := ""
	if len(*s) != 0 {
		str = (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
	}
	return str
}

func (s *stack) PushAll(letters []string) {
	*s = append(*s, letters...)
}

func (s *stack) LPush(str string) {
	*s = append(*s, str)
	if len(*s) != 1 {
		*s = append((*s)[len(*s)-1:len(*s)], (*s)[:len(*s)-1]...)
	}
}

func (s *stack) PopN(n int) []string {
	var letters []string
	if len(*s) != 0 {
		letters = (*s)[len(*s)-n : len(*s)]
		*s = (*s)[:len(*s)-n]
	}
	return letters
}

func init() {
	challenges.RegisterChallengeFunc(2022, 5, 1, "day05.txt", part1)
	challenges.RegisterChallengeFunc(2022, 5, 2, "day05.txt", part2)
}
