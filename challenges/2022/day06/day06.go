package aoc2022_day6

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	const n int = 4

	var result int = -1
	for i := range parsed {
		if !hasDuplicates(parsed, i, n) {
			result = i + n
			break
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	const n int = 14

	var result int = -1
	for i := range parsed {
		if !hasDuplicates(parsed, i, n) {
			result = i + n
			break
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func hasDuplicates(arr []rune, start int, n int) bool {
	seen := make(map[rune]bool)
	for i := start; i < start+n; i++ {
		if _, ok := seen[arr[i]]; ok {
			return true
		}
		seen[arr[i]] = true
	}
	return false
}

func parseInput(input []string) []rune {
	return []rune(input[0])
}

func init() {
	challenges.RegisterChallengeFunc(2022, 6, 1, "day06.txt", part1)
	challenges.RegisterChallengeFunc(2022, 6, 2, "day06.txt", part2)
}
