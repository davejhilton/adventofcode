package aoc2020_day3

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
)

func part1(input []string) (string, error) {
	nTrees := traverseSlopes(input, 3, 1)
	return fmt.Sprintf("%d", nTrees), nil
}

func part2(input []string) (string, error) {
	slopes := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	result := 1
	for _, slope := range slopes {
		result = result * traverseSlopes(input, slope[0], slope[1])
	}
	return fmt.Sprintf("%d", result), nil
}

func traverseSlopes(input []string, r int, d int) int {
	width := len(input[0])
	nTrees := 0
	idx := 0
	i := 0
	for i < len(input) {
		if idx >= width {
			idx = idx % width
		}
		c := input[i][idx]
		if c == '#' {
			nTrees++
		}
		idx = idx + r
		i = i + d
	}
	return nTrees
}

func init() {
	challenges.RegisterChallengeFunc(2020, 3, 1, "day03.txt", part1)
	challenges.RegisterChallengeFunc(2020, 3, 2, "day03.txt", part2)
}
