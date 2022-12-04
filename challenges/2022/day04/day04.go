package aoc2022_day4

import (
	"fmt"
	"regexp"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

var (
	parseRegex = regexp.MustCompile(`(\d+)-(\d+),(\d+)-(\d+)`)
)

type pair struct {
	FirstStart  int
	FirstStop   int
	SecondStart int
	SecondStop  int
}

func part1(input []string) (string, error) {
	pairs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", pairs)

	var result int
	for _, p := range pairs {
		if p.FirstStart <= p.SecondStart && p.FirstStop >= p.SecondStop {
			result += 1
		} else if p.SecondStart <= p.FirstStart && p.SecondStop >= p.FirstStop {
			result += 1
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	pairs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", pairs)

	var result int
	for _, p := range pairs {
		if p.FirstStart <= p.SecondStart && p.FirstStop >= p.SecondStart {
			result += 1
		} else if p.SecondStart <= p.FirstStart && p.SecondStop >= p.FirstStart {
			result += 1
		}
	}
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []pair {
	pairs := make([]pair, 0, len(input))
	for _, s := range input {
		matches := parseRegex.FindStringSubmatch(s)

		pairs = append(pairs, pair{
			FirstStart:  util.Atoi(matches[1]),
			FirstStop:   util.Atoi(matches[2]),
			SecondStart: util.Atoi(matches[3]),
			SecondStop:  util.Atoi(matches[4]),
		})
	}
	return pairs
}

func init() {
	challenges.RegisterChallengeFunc(2022, 4, 1, "day04.txt", part1)
	challenges.RegisterChallengeFunc(2022, 4, 2, "day04.txt", part2)
}
