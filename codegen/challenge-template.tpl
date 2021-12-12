package aoc{{.Year}}_day{{.Day}}

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	var result int
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	// parsed := parseInput(input)

	var result int
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		nums = append(nums, util.Atoi(s))
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc({{.Year}}, {{.Day}}, 1, "day{{ printf "%.2d" .Day }}.txt", part1)
	challenges.RegisterChallengeFunc({{.Year}}, {{.Day}}, 2, "day{{ printf "%.2d" .Day }}.txt", part2)
}
