package aoc2020_day15

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	return fmt.Sprintf("%d", iterate(input[0], 2020)), nil
}

func part2(input []string) (string, error) {
	return fmt.Sprintf("%d", iterate(input[0], 30000000)), nil
}

func iterate(numList string, target int) int {
	var i int
	var prev int
	lastSeen := make(map[int]int)

	numStrings := strings.Split(numList, ",")
	for i = 0; i < len(numStrings); i++ {
		n, _ := strconv.Atoi(numStrings[i])
		lastSeen[n] = i + 1
		prev = n
		log.Debugf("%8d: %d\n", i+1, prev)
	}

	for ; i < target; i++ {
		li, ok := lastSeen[prev]
		lastSeen[prev] = i
		if ok {
			prev = i - li
		} else {
			prev = 0
		}
		log.Debugf("%4d: %d\n", i, prev)
	}
	return prev
}

func init() {
	challenges.RegisterChallengeFunc(2020, 15, 1, "day15.txt", part1)
	challenges.RegisterChallengeFunc(2020, 15, 2, "day15.txt", part2)
}
