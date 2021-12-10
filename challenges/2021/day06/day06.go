package aoc2021_day6

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
)

const FISH_TIMER = 6
const NEW_FISH_TIMER = 8

// map of "day number the fish was created on" --> "total number of descendents fish will produce by the last day"
var cache = make(map[int]int)

func part1(input []string) (string, error) {
	fishList := parse(input)
	result := len(fishList)
	for _, n := range fishList {
		result += beFruitfulAndMultiply(n, 0, 80)
	}
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	fishList := parse(input)
	result := len(fishList)
	for _, n := range fishList {
		result += beFruitfulAndMultiply(n, 0, 256)
	}
	return fmt.Sprintf("%d", result), nil
}

func beFruitfulAndMultiply(num int, startDay int, totalDays int) int {
	day := startDay
	count := 0
	for ; day < totalDays; day += 1 {
		if num == 0 {
			if c, ok := cache[day]; ok {
				count += c
			} else {
				c := 1 + beFruitfulAndMultiply(NEW_FISH_TIMER, day+1, totalDays)
				cache[day] = c
				count += c
			}
			num = FISH_TIMER
		} else {
			num = num - 1
		}
	}
	return count
}

func parse(input []string) []int {
	strs := strings.Split(input[0], ",")
	nums := make([]int, 0, len(strs))
	for _, s := range strs {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2021, 6, 1, "day06.txt", part1)
	challenges.RegisterChallengeFunc(2021, 6, 2, "day06.txt", part2)
}
