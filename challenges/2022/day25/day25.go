package aoc2022_day25

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	ints := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", ints)

	sum := 0
	for _, d := range ints {
		sum += d
	}

	var result = IntToSnafu(sum)
	return result, nil
}

func part2(input []string) (string, error) {
	var result = "There is no part 2 for this day"
	return result, nil
}

func parseInput(input []string) []int {
	nums := make([]int, 0, len(input))
	for _, s := range input {
		nums = append(nums, SnafuToInt(s))
	}
	return nums
}

func SnafuToInt(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		m := int(math.Pow(5, float64(len(s)-1-i)))
		v := INPUT_MAP[[]rune(s)[i]]
		n += m * v
	}
	return n
}

func IntToSnafu(n int) string {
	// find the greatest power of 5 that's less than n
	var pow float64 = 0
	for math.Pow(5, pow+1) < float64(n) {
		pow++
	}

	// divide n by powers of 5, and keep using the remainder
	digits := make([]int, 0)
	remainder := n
	for pow >= 0 {
		divisor := int(math.Pow(5, pow))
		digits = append(digits, remainder/divisor)
		pow--
		remainder = n % divisor
	}

	// if any digit is larger than 2, carry over the extra amount to the digit to its left
	carry := 0
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i] + carry
		carry = 0
		for digit > 2 {
			digit -= 5
			carry++
		}
		digits[i] = digit
	}
	// if we have any left to carry over, prepend a new digit to the front of the slice
	if carry != 0 {
		digits = append([]int{carry}, digits...)
	}
	// convert to string
	var b strings.Builder
	for i, d := range digits {
		ru, ok := OUTPUT_MAP[d]
		if !ok {
			fmt.Printf("UNKNOWN DIGIT: %d at index: %d\n", d, i)
		} else {
			b.WriteRune(ru)
		}
	}
	// ... profit!
	return b.String()
}

var INPUT_MAP = map[rune]int{
	'2': 2,
	'1': 1,
	'0': 0,
	'-': -1,
	'=': -2,
}
var OUTPUT_MAP = map[int]rune{
	2:  '2',
	1:  '1',
	0:  '0',
	-1: '-',
	-2: '=',
}

func init() {
	challenges.RegisterChallengeFunc(2022, 25, 1, "day25.txt", part1)
	challenges.RegisterChallengeFunc(2022, 25, 2, "day25.txt", part2)
}
