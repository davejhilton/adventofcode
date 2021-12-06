package challenges2021

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

const day06_FISH_TIMER = 6
const day06_NEW_FISH_TIMER = 8

// map of "day number the fish was created on" --> "total number of descendents fish will produce by the last day"
var day06_cache = make(map[int]int)

func day06_part1(input []string) (string, error) {
	fishList := day06_parse(input)
	result := day06_simulateFishLife(fishList, 80)
	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day06_part2(input []string) (string, error) {
	fishList := day06_parse(input)
	result := day06_simulateFishLife(fishList, 256)
	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day06_simulateFishLife(fishList []int, numDays int) int {
	result := len(fishList)
	for _, n := range fishList {
		result += day06_beFruitfulAndMultiply(n, 0, numDays)
	}
	return result
}

func day06_beFruitfulAndMultiply(num int, startDay int, totalDays int) int {
	day := startDay
	count := 0
	for ; day < totalDays; day += 1 {
		if num == 0 {
			if c, ok := day06_cache[day]; ok {
				count += c
			} else {
				c := 1 + day06_beFruitfulAndMultiply(day06_NEW_FISH_TIMER, day+1, totalDays)
				day06_cache[day] = c
				count += c
			}
			num = day06_FISH_TIMER
		} else {
			num = num - 1
		}
	}
	return count
}

func day06_parse(input []string) []int {
	strs := strings.Split(input[0], ",")
	nums := make([]int, 0, len(strs))
	for _, s := range strs {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2021, 06, 1, "day06.txt", day06_part1)
	challenges.RegisterChallengeFunc(2021, 06, 2, "day06.txt", day06_part2)
}
