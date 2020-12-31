package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day25_part1(input []string) (string, error) {
	pk1, pk2 := day25_parse(input)
	_, _ = pk1, pk2

	var subject = 7
	var value = 1
	var loopNum = 0
	var divisor = 20201227

	var pk1Loops = -1
	var pk2Loops = -1

	for pk1Loops == -1 && pk2Loops == -1 {
		loopNum++
		value = day25_transform(value, subject, divisor)
		if value == pk1 {
			pk1Loops = loopNum
		} else if value == pk2 {
			pk2Loops = loopNum
		}
	}

	var nLoops int
	var subject2 int
	if pk1Loops != -1 {
		subject2 = pk2
		nLoops = pk1Loops
	} else {
		subject2 = pk1
		nLoops = pk2Loops
	}

	var value2 = 1
	for i := 0; i < nLoops; i++ {
		value2 = day25_transform(value2, subject2, divisor)
	}

	return fmt.Sprintf("%d", value2), nil
}

func day25_transform(value, subject, divisor int) int {
	value *= subject
	value %= divisor
	return value
}

func day25_part2(input []string) (string, error) {
	pk1, pk2 := day25_parse(input)
	_, _ = pk1, pk2

	log.Debugf("Result: %d\n", 0)
	return fmt.Sprintf("%d", 0), nil
}

func day25_parse(input []string) (int, int) {
	n1, _ := strconv.Atoi(input[0])
	n2, _ := strconv.Atoi(input[1])
	return n1, n2
}

func init() {
	registerChallengeFunc(25, 1, "day25.txt", day25_part1)
	registerChallengeFunc(25, 2, "day25.txt", day25_part2)
}
