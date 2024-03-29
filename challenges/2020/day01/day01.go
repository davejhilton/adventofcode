package aoc2020_day1

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	nums, err := parseInput(input)
	if err != nil {
		return "", err
	}
	for i, v1 := range nums {
		for j := i + 1; j < len(nums); j++ {
			v2 := nums[j]
			if v1+v2 == 2020 {
				log.Debugf("%d + %d = 2020\n", v1, v2)
				return fmt.Sprintf("%d", v1*v2), nil
			}
		}
	}
	return "", fmt.Errorf("no solution found%s", "\n")
}

func part2(input []string) (string, error) {
	nums, err := parseInput(input)
	if err != nil {
		return "", err
	}
	for i, v1 := range nums {
		for j := i + 1; j < len(nums); j++ {
			v2 := nums[j]
			if v1+v2 <= 2020 {
				for k := j + 1; k < len(nums); k++ {
					v3 := nums[k]
					if v1+v2+v3 == 2020 {
						log.Debugf("%d (%d) + %d (%d) + %d (%d) = 2020\n", v1, i, v2, j, v3, k)
						return fmt.Sprintf("%d", v1*v2*v3), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no solution found")
}

func parseInput(lines []string) ([]int, error) {
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
	challenges.RegisterChallengeFunc(2020, 1, 1, "day01.txt", part1)
	challenges.RegisterChallengeFunc(2020, 1, 2, "day01.txt", part2)
}
