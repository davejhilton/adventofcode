package aoc2022_day10

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type instruction struct {
	Operator string
	Operand  int
}

var (
	cycleCounts = map[string]int{
		"noop": 1,
		"addx": 2,
	}
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	cycle := 0
	nCycles := 0
	x := 1

	sum := 0
	for _, instr := range parsed {
		nCycles = cycleCounts[instr.Operator]

		for i := 0; i < nCycles; i++ {
			cycle++
			if (cycle-20)%40 == 0 {
				sum += (cycle * x)
			}
		}
		if instr.Operator == "addx" {
			x += instr.Operand
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	cycle := 0
	x := 1

	var out strings.Builder

	nCycles := 0
	for _, instr := range parsed {
		nCycles = cycleCounts[instr.Operator]

		for i := 0; i < nCycles; i++ {
			if cycle%40 == 0 {
				out.WriteString("\n")
			}

			if util.Abs((cycle%40)-x) <= 1 {
				out.WriteString("#")
			} else {
				out.WriteString(".")
			}
			cycle++
		}
		if instr.Operator == "addx" {
			x += instr.Operand
		}
	}

	return out.String(), nil
}

var (
	addxRegex = regexp.MustCompile(`^addx\s(-?\d+)$`)
)

func parseInput(input []string) []instruction {
	instrs := make([]instruction, 0, len(input))
	for _, s := range input {
		instr := instruction{Operator: s[:4]}
		if s != "noop" {
			instr.Operand = util.Atoi(addxRegex.FindStringSubmatch(s)[1])
		}
		instrs = append(instrs, instr)
	}
	return instrs
}

func init() {
	challenges.RegisterChallengeFunc(2022, 10, 1, "day10.txt", part1)
	challenges.RegisterChallengeFunc(2022, 10, 2, "day10.txt", part2)
}
