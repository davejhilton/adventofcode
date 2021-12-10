package aoc2020_day9

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	nums := parseInts(input)

	var i int
	var window int
	var result int

	if strings.Contains(challenges.CurrentChallenge.InputFileName, "example") {
		window = 5
	} else {
		window = 25
	}
	for i = window; i < len(nums); i++ {
		if !checkSum(nums[i], nums[i-window:i]) {
			result = nums[i]
			break
		}
	}
	log.Debugf("Result: %d (index %d)\nprev %d:\n%v\n\n", result, i, window, nums[i-window:i])

	return fmt.Sprintf("%d", result), nil
}

func checkSum(val int, nums []int) bool {
	for i, n := range nums {
		if n >= val {
			continue
		}
		for _, m := range nums[i+1:] {
			if m+n == val {
				log.Debugf("%d + %d = %d\n", n, m, val)
				return true
			}
		}
	}
	return false
}

func part2(input []string) (string, error) {
	nums := parseInts(input)

	var target int
	if strings.Contains(challenges.CurrentChallenge.InputFileName, "example") {
		target = 127
	} else {
		target = 2089807806
	}
	first, last := 0, 0

	sum := nums[0]
	for last < len(nums) {
		if sum < target {
			last++
			sum += nums[last]
		} else if sum > target {
			sum -= nums[first]
			first++
		} else {
			break
		}
	}

	log.Debugf("RANGE: nums[%d] - nums[%d]\n", first, last)
	log.Debugf("%v\n", nums[first:last+1])

	sm, lg := nums[first], nums[first]
	for i := first + 1; i <= last; i++ {
		if nums[i] < sm {
			sm = nums[i]
		}
		if nums[i] > lg {
			lg = nums[i]
		}
	}

	log.Debugf("SMALLEST: %d, LARGEST: %d\n", sm, lg)
	return fmt.Sprintf("%d", sm+lg), nil
}

func parseInts(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2020, 9, 1, "day09.txt", part1)
	challenges.RegisterChallengeFunc(2020, 9, 2, "day09.txt", part2)
}
