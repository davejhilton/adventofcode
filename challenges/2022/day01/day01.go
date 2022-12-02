package aoc2022_day1

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	max := 0
	for _, group := range parsed {
		calories := 0
		for _, cals := range group {
			calories += cals
		}
		max = util.Max(max, calories)
	}
	return fmt.Sprintf("%d", max), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)

	top3 := append(make([]int, 0), 0, 0, 0)
	for _, group := range parsed {
		calories := 0
		for _, cals := range group {
			calories += cals
		}
		for i, v := range top3 {
			if calories > v {
				top3 = append(top3[:i+1], top3[i:]...)
				top3[i] = calories
				top3 = top3[0:3]
				break
			}
		}
	}
	total := 0
	for _, v := range top3 {
		total += v
	}
	return fmt.Sprintf("%d", total), nil
}

func parseInput(input []string) [][]int {
	groups := make([][]int, 0)
	nums := make([]int, 0)
	for _, v := range input {
		if v == "" {
			groups = append(groups, nums)
			nums = make([]int, 0)
			continue
		}
		nums = append(nums, util.Atoi(v))
	}
	return append(groups, nums)
}

func init() {
	challenges.RegisterChallengeFunc(2022, 1, 1, "day01.txt", part1)
	challenges.RegisterChallengeFunc(2022, 1, 2, "day01.txt", part2)
}
