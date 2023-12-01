package aoc2023_day1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

var digits = map[string]int{
	"0":     0,
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", parsed)

	reg := regexp.MustCompile(`[0-9]`)

	result := 0
	for _, line := range parsed {
		matches := reg.FindAllString(line, -1)

		first := -1
		last := -1
		for _, match := range matches {
			if first == -1 {
				first = digits[match]
			}
			last = digits[match]
		}
		lineVal := (first * 10) + last
		log.Debugf("%s: %d\n", line, lineVal)
		result += lineVal
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)

	reg := regexp.MustCompile(`([0-9]|zero|one|two|three|four|five|six|seven|eight|nine)`)

	result := 0

	maxLen := 0
	for _, line := range parsed {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	for _, line := range parsed {
		first := -1
		last := -1
		// find all matches -- even overlapping ones
		origLine := line
		allMatches := make([]string, 0)
		firstMatch := reg.FindStringIndex(line)
		for firstMatch != nil {
			startIdx, endIdx := firstMatch[0], firstMatch[1]
			match := line[startIdx:endIdx]
			allMatches = append(allMatches, match)
			last = digits[match]
			if first == -1 {
				first = last
			}
			line = line[startIdx+1:] // match again, starting after the most recent match
			firstMatch = reg.FindStringIndex(line)
		}
		allMatches[0] = log.Colorize(allMatches[0], log.Yellow, 0)
		allMatches[len(allMatches)-1] = log.Colorize(allMatches[len(allMatches)-1], log.Teal, 0)
		lineVal := (first * 10) + last
		log.Debugf("%s: [%d] %s\n", log.Colorize(origLine, log.Green, -1*(maxLen+1)), lineVal, strings.Join(allMatches, " "))
		result += lineVal
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []string {
	return input
}

func init() {
	challenges.RegisterChallengeFunc(2023, 1, 1, "day01.txt", part1)
	challenges.RegisterChallengeFunc(2023, 1, 2, "day01.txt", part2)
}
