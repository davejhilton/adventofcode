package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day000_part1(input []string) (string, error) {
	nums := day000_parse(input)
	var result int

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day000_part2(input []string) (string, error) {
	nums := day000_parse(input)
	var result int

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day000_parse(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	return nums
}

func init() {
	registerChallengeFunc(000, 1, "day000.txt", day000_part1)
	registerChallengeFunc(000, 2, "day000.txt", day000_part2)
}
