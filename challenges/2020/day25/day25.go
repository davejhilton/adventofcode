package aoc2020_day25

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	pk1, pk2 := parseInput(input)
	_, _ = pk1, pk2

	subject := 7
	value := 1
	loopNum := 0
	divisor := 20201227

	pk1Loops := -1
	pk2Loops := -1

	for pk1Loops == -1 && pk2Loops == -1 {
		loopNum++
		value = transform(value, subject, divisor)
		if value == pk1 {
			pk1Loops = loopNum
		} else if value == pk2 {
			pk2Loops = loopNum
		}
	}

	var nLoops int
	var subject2 int
	if pk1Loops != -1 {
		subject2 = pk2
		nLoops = pk1Loops
	} else {
		subject2 = pk1
		nLoops = pk2Loops
	}

	value2 := 1
	for i := 0; i < nLoops; i++ {
		value2 = transform(value2, subject2, divisor)
	}

	return fmt.Sprintf("%d", value2), nil
}

func transform(value, subject, divisor int) int {
	value *= subject
	value %= divisor
	return value
}

func part2(input []string) (string, error) {
	pk1, pk2 := parseInput(input)
	_, _ = pk1, pk2

	log.Debugf("Result: %d\n", 0)
	return fmt.Sprintf("%d", 0), nil
}

func parseInput(input []string) (int, int) {
	n1, _ := strconv.Atoi(input[0])
	n2, _ := strconv.Atoi(input[1])
	return n1, n2
}

func init() {
	challenges.RegisterChallengeFunc(2020, 25, 1, "day25.txt", part1)
	challenges.RegisterChallengeFunc(2020, 25, 2, "day25.txt", part2)
}
