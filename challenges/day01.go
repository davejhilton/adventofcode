package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day1_part1(input []string) (string, error) {
	nums, err := day1_formatInput(input)
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
	return "", fmt.Errorf("No solution found\n")
}

func day1_part2(input []string) (string, error) {
	nums, err := day1_formatInput(input)
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
	return "", fmt.Errorf("No solution found")
}

func day1_formatInput(lines []string) ([]int, error) {
	nums := make([]int, 0, len(lines))
	for _, v := range lines {
		if i, err := strconv.Atoi(v); err == nil {
			nums = append(nums, i)
		} else {
			return nums, fmt.Errorf("Unable to parse input - '%s' is not a number!", v)
		}
	}
	return nums, nil
}

func init() {
	registerChallengeFunc(1, 1, "day01.txt", day1_part1)
	registerChallengeFunc(1, 2, "day01.txt", day1_part2)
}
