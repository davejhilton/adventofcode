package aoc2022_day1

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	groups := parseInput(input)
	max := 0
	for _, group := range groups {
		calories := 0
		for _, cals := range group {
			calories += cals
		}
		if calories > max {
			max = calories
		}
	}
	return fmt.Sprintf("%d", max), nil
}

func part2(input []string) (string, error) {
	groups := parseInput(input)
	top3 := append(make([]int, 0), 0, 0, 0)
	for _, group := range groups {
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
		if i, err := strconv.Atoi(v); err == nil {
			nums = append(nums, i)
		} else {
			log.Printf("unable to parse input - '%s' is not a number", v)
			return nil
		}
	}
	return append(groups, nums)
}

func init() {
	challenges.RegisterChallengeFunc(2022, 1, 1, "day01.txt", part1)
	challenges.RegisterChallengeFunc(2022, 1, 2, "day01.txt", part2)
}
