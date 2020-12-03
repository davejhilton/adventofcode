package challenges

import (
	"fmt"
	"strconv"
	"strings"
)

func day2_part1(input []string) (string, error) {
	valid := 0
	for _, line := range input {
		parts := strings.SplitN(line, " ", 3)
		minMaxParts := strings.SplitN(parts[0], "-", 2)
		min, _ := strconv.Atoi(minMaxParts[0])
		max, _ := strconv.Atoi(minMaxParts[1])
		char := strings.Replace(parts[1], ":", "", -1)
		pwd := parts[2]
		// fmt.Printf("min: '%d', max: '%d', char: '%s', pwd: '%s'\n", min, max, char, pwd)
		n := 0
		for _, c := range pwd {
			if string(c) == char {
				n++
			}
		}
		if min <= n && n <= max {
			valid++
		}
	}
	// fmt.Printf("====================\n%d valid password(s)\n====================\n", valid)
	return fmt.Sprintf("%d", valid), nil
}

func day2_part2(input []string) (string, error) {
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

		// fmt.Printf("idx1: '%d', idx2: '%d', char: '%s', pwd: '%s', valid: '%#v'\n", idx1, idx2, char, pwd, matches == 1)
		if matches == 1 {
			valid++
		}
	}
	// fmt.Printf("====================\n%d valid password(s)\n====================\n", valid)
	return fmt.Sprintf("%d", valid), nil
}

func init() {
	registerChallengeFunc(2, 1, "day2.txt", day2_part1)
	registerChallengeFunc(2, 2, "day2.txt", day2_part2)
}
