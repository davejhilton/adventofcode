package aoc2021_day1

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day1_part1(input []string) (string, error) {
	nums, err := day1_formatInput(input)
	if err != nil {
		return "", err
	}
	prevDepth := -1
	increases := 0
	for i, depth := range nums {
		if i == 0 {
			log.Debugf("%d (N/A - no previous measurement)\n", depth)
		} else {
			if depth > prevDepth {
				log.Debugf("%d %s\n", depth, log.Colorize("(increases)", log.Green, 0))
				increases++
			} else if depth < prevDepth {
				log.Debugf("%d %s\n", depth, log.Colorize("(decreases)", log.Red, 0))
			} else {
				log.Debugf("%d %s\n", depth, log.Colorize("(no change)", log.Yellow, 0))
			}
		}
		prevDepth = depth
	}
	return fmt.Sprintf("%d", increases), nil
}

func day1_part2(input []string) (string, error) {
	nums, err := day1_formatInput(input)
	if err != nil {
		return "", err
	}
	label := 'A'
	prevSum := -1
	increases := 0
	for i, depth := range nums {
		if i <= 1 {
			// log.Debugf("%d (N/A - no previous measurement)\n", depth)
		} else {
			p1, p2 := nums[i-1], nums[i-2]
			sum := p1 + p2 + depth
			if i == 2 {
				log.Debugf("%s: %d (N/A - no previous sum)\n", label, sum)
			} else if sum > prevSum {
				log.Debugf("%s: %d %s\n", label, sum, log.Colorize("(increases)", log.Green, 0))
				increases++
			} else if sum < prevSum {
				log.Debugf("%s: %d %s\n", label, sum, log.Colorize("(decreases)", log.Red, 0))
			} else {
				log.Debugf("%s: %d %s\n", label, sum, log.Colorize("(no change)", log.Yellow, 0))
			}
			label = label + 1
			prevSum = sum
		}
	}
	return fmt.Sprintf("%d", increases), nil
}

func day1_formatInput(lines []string) ([]int, error) {
	nums := make([]int, 0, len(lines))
	for _, v := range lines {
		if i, err := strconv.Atoi(v); err == nil {
			nums = append(nums, i)
		} else {
			return nums, fmt.Errorf("unable to parse input - '%s' is not a number", v)
		}
	}
	return nums, nil
}

func init() {
	challenges.RegisterChallengeFunc(2021, 1, 1, "day01.txt", day1_part1)
	challenges.RegisterChallengeFunc(2021, 1, 2, "day01.txt", day1_part2)
}
