package aoc2020_day2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	valid := 0
	for _, line := range input {
		parts := strings.SplitN(line, " ", 3)
		minMaxParts := strings.SplitN(parts[0], "-", 2)
		min, _ := strconv.Atoi(minMaxParts[0])
		max, _ := strconv.Atoi(minMaxParts[1])
		char := strings.Replace(parts[1], ":", "", -1)
		pwd := parts[2]
		n := 0
		for _, c := range pwd {
			if string(c) == char {
				n++
			}
		}
		log.Debugf("line: '%s', min: '%d', max: '%d', char: '%s', pwd: '%s', matches: '%d', valid: '%v'", line, min, max, char, pwd, n, min <= n && n <= max)
		if min <= n && n <= max {
			valid++
		}
	}

	return fmt.Sprintf("%d", valid), nil
}

func part2(input []string) (string, error) {
	valid := 0
	for _, line := range input {
		parts := strings.SplitN(line, " ", 3)
		indexParts := strings.SplitN(parts[0], "-", 2)
		idx1, _ := strconv.Atoi(indexParts[0])
		idx2, _ := strconv.Atoi(indexParts[1])
		char := strings.Replace(parts[1], ":", "", -1)
		pwd := parts[2]
		matches := 0
		if len(pwd) > idx1-1 {
			if string(pwd[idx1-1]) == char {
				matches++
			}
		}
		if len(pwd) > idx2-1 {
			if string(pwd[idx2-1]) == char {
				matches++
			}
		}

		log.Debugf("line: '%s', idx1: '%d', idx2: '%d', char: '%s', pwd: '%s', valid: '%#v'\n", line, idx1, idx2, char, pwd, matches == 1)
		if matches == 1 {
			valid++
		}
	}
	return fmt.Sprintf("%d", valid), nil
}

func init() {
	challenges.RegisterChallengeFunc(2020, 2, 1, "day02.txt", part1)
	challenges.RegisterChallengeFunc(2020, 2, 2, "day02.txt", part2)
}
