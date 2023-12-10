package aoc2023_day9

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func calcDiffsLtoR(nums []int) []int {
	diffs := make([]int, len(nums)-1)
	for i := 1; i < len(nums); i++ {
		diffs[i-1] = nums[i] - nums[i-1]
	}
	return diffs
}

func calcDiffsRtoL(nums []int) []int {
	diffs := make([]int, len(nums)-1)
	for i := len(nums) - 2; i >= 0; i-- {
		diffs[i] = nums[i] - nums[i+1]
	}
	return diffs
}

func isAllZeroes(nums []int) bool {
	for _, n := range nums {
		if n != 0 {
			return false
		}
	}
	return true
}

func part1(input []string) (string, error) {
	lists := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", lists)

	var result int
	for _, list := range lists {
		levels := make([][]int, 0)
		levels = append(levels, list)
		curLevel := list
		levelStrings := make([]string, 0)
		for !isAllZeroes(curLevel) {
			curLevel = calcDiffsLtoR(curLevel)
			levels = append(levels, curLevel)
			levelStrings = append(levelStrings, fmt.Sprintf("%v", curLevel))
		}
		log.Debugf("Levels:\n%s\n\n", strings.Join(levelStrings, "\n"))
		prevLast := 0
		for j := len(levels) - 1; j >= 0; j-- {
			// log.Debugf("Level %d: %v\n", j, levels[j])
			levels[j] = append(levels[j], levels[j][len(levels[j])-1]+prevLast)
			prevLast = levels[j][len(levels[j])-1]
		}
		result += prevLast
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	lists := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", lists)

	var result int
	for _, list := range lists {
		levels := make([][]int, 0)
		levels = append(levels, list)
		curLevel := list
		levelStrings := make([]string, 0)
		levelStrings = append(levelStrings, fmt.Sprintf("%v", curLevel))
		for !isAllZeroes(curLevel) {
			curLevel = calcDiffsRtoL(curLevel)
			levels = append(levels, curLevel)
			levelStrings = append(levelStrings, fmt.Sprintf("%v", curLevel))
		}
		log.Debugf("Levels:\n%s\n\n", strings.Join(levelStrings, "\n"))
		prevVal := 0
		for j := len(levels) - 1; j >= 0; j-- {
			// log.Debugf("Level %d: %v\n", j, levels[j])
			levels[j] = append([]int{prevVal + levels[j][0]}, levels[j]...)
			// levels[j] = append(levels[j], levels[j][len(levels[j])-1]+prevVal)
			prevVal = levels[j][0]
		}
		log.Debugf("Levels:\n%s\n\n", strings.Join(levelStrings, "\n"))
		result += prevVal
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) [][]int {
	nums := make([][]int, 0, len(input))
	for _, s := range input {
		nums = append(nums, util.AtoiSplit(s, " "))
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2023, 9, 1, "day09.txt", part1)
	challenges.RegisterChallengeFunc(2023, 9, 2, "day09.txt", part2)
}
