package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day14_part1(input []string) (string, error) {
	result := 0
	nums := day14_parseInts(input)
	for i, n := range nums {
		log.Debugf("Line %d: %d\n", i, n)

	}
	return fmt.Sprintf("%d", result), nil
}

func day14_part2(input []string) (string, error) {
	result := 0
	nums := day14_parseInts(input)
	for i, n := range nums {
		log.Debugf("Line %d: %d\n", i, n)

	}
	return fmt.Sprintf("%d", result), nil
}

func day14_parseInts(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	registerChallengeFunc(14, 1, "day14.txt", day14_part1)
	registerChallengeFunc(14, 2, "day14.txt", day14_part2)
}
