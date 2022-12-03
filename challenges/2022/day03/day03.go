package aoc2022_day3

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	total := 0
	for _, row := range parsed {
		first := row[:(len(row) / 2)]
		second := row[len(row)/2:]
		log.Debugf("First:  %s\n", string(first))
		log.Debugf("Second: %s\n", string(second))

		seen := make(map[rune]bool)
		for _, r := range first {
			seen[r] = true
		}
		matchVal := 0
		for _, r := range second {
			if _, ok := seen[r]; ok {
				matchVal = getPriority(r)
				log.Debugf("COMMON: %s (%d)\n\n", string(r), matchVal)
				break
			}
		}
		total += matchVal
	}

	return fmt.Sprintf("%d", total), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	total := 0
	seenInGroup := make(map[rune]int)
	for i, row := range parsed {
		seenInBag := make(map[rune]bool)
		if i%3 == 0 {
			seenInGroup = make(map[rune]int)
		}

		for _, r := range row {
			if !seenInBag[r] {
				seenInBag[r] = true
				seenInGroup[r] = seenInGroup[r] + 1
				if seenInGroup[r] == 3 {
					total += getPriority(r)
					break
				}
			}
		}
	}

	return fmt.Sprintf("%d", total), nil
}

func parseInput(input []string) [][]rune {
	rows := make([][]rune, 0, len(input))
	for _, s := range input {
		rows = append(rows, []rune(s))
	}
	return rows
}

func getPriority(c rune) int {
	if int(c) > 96 {
		return int(c) - 96
	} else {
		return int(c) - 38
	}
}

func init() {
	challenges.RegisterChallengeFunc(2022, 3, 1, "day03.txt", part1)
	challenges.RegisterChallengeFunc(2022, 3, 2, "day03.txt", part2)
}
