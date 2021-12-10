package aoc2021_day7

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
)

func minDistance(positions []int, max int, adder func(n int) int) int {
	min := math.MaxInt
	var sum int
	for i := 0; i <= max; i++ {
		sum = 0
		for _, n := range positions {
			diff := util.Abs(n - i)
			sum += adder(diff)
		}
		min = util.Min(sum, min)
	}
	return min
}

func part1(input []string) (string, error) {
	positions, max := parse(input)
	result := minDistance(positions, max, func(n int) int {
		return n
	})
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	positions, max := parse(input)
	result := minDistance(positions, max, func(n int) int {
		return n * (n + 1) / 2
	})
	return fmt.Sprintf("%d", result), nil
}

func parse(input []string) ([]int, int) {
	strs := strings.Split(input[0], ",")
	nums := make([]int, 0, len(strs))
	var max int
	for _, s := range strs {
		n := util.Atoi(s)
		nums = append(nums, n)
		max = util.Max(n, max)
	}
	return nums, max
}

func init() {
	challenges.RegisterChallengeFunc(2021, 7, 1, "day07.txt", part1)
	challenges.RegisterChallengeFunc(2021, 7, 2, "day07.txt", part2)
}
