package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day15_part1(input []string) (string, error) {
	return fmt.Sprintf("%d", day15_iterate(input[0], 2020)), nil
}

func day15_part2(input []string) (string, error) {
	return fmt.Sprintf("%d", day15_iterate(input[0], 30000000)), nil
}

func day15_iterate(numList string, target int) int {

	strs := strings.Split(numList, ",")
	startingNums := make([]int, 0, len(strs))
	for _, s := range strs {
		n, _ := strconv.Atoi(s)
		startingNums = append(startingNums, n)
	}

	var i int
	var prev int
	lastSeen := make(map[int]int)

	for i = 0; i < len(startingNums); i++ {
		lastSeen[startingNums[i]] = i + 1
		prev = startingNums[i]
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
	registerChallengeFunc(15, 1, "day15.txt", day15_part1)
	registerChallengeFunc(15, 2, "day15.txt", day15_part2)
}
