package challenges{{.Year}}

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func day{{.Day}}_part1(input []string) (string, error) {
	parsed := day{{.Day}}_parse(input)
	log.Debugf("Parsed Input:\n%v\n", parsed)
	var result int

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day{{.Day}}_part2(input []string) (string, error) {
	// parsed := day{{.Day}}_parse(input)
	var result int

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day{{.Day}}_parse(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		nums = append(nums, util.Atoi(s))
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc({{.Year}}, {{.Day}}, 1, "day{{.Day}}.txt", day{{.Day}}_part1)
	challenges.RegisterChallengeFunc({{.Year}}, {{.Day}}, 2, "day{{.Day}}.txt", day{{.Day}}_part2)
}
