package aoc2022_day13

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func ListAsString(l any) string {
	str, _ := json.Marshal(l)
	return string(str)
}

func part1(input []string) (string, error) {
	lists := parseInput(input)

	var result int

	pairIdx := 1
	for i := 0; i < len(lists); i += 2 {
		v := compare(lists[i], lists[i+1])
		if v == -1 {
			result += pairIdx
		}
		log.Debugf("COMPARING pairs[%d]: %d\n", pairIdx, v)
		pairIdx++
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	var lists [][]any = [][]any{
		{[]any{float64(2)}},
		{[]any{float64(6)}},
	}
	lists = append(lists, parseInput(input)...)

	sort.Slice(lists, func(i, j int) bool {
		return compare(lists[i], lists[j]) == -1
	})

	idxs := make([]int, 0, 2)
	for i, list := range lists {
		s := ListAsString(list)
		if s == "[[2]]" || s == "[[6]]" {
			idxs = append(idxs, i+1)
			if len(idxs) == 2 {
				break
			}
		}
	}

	return fmt.Sprintf("%d", idxs[0]*idxs[1]), nil
}

func compare(a, b any) int {
	_, aIsNumber := a.(float64)
	_, bIsNumber := b.(float64)

	if aIsNumber && bIsNumber {
		if a.(float64) > b.(float64) {
			return 1
		} else if a.(float64) < b.(float64) {
			return -1
		}
		return 0
	} else if !aIsNumber && !bIsNumber {
		aa := a.([]any)
		bb := b.([]any)
		minLen := util.Min(len(aa), len(bb))
		for i := 0; i < minLen; i++ {
			if cmp := compare(aa[i], bb[i]); cmp != 0 {
				return cmp
			}
		}
		if len(aa) > len(bb) {
			return 1
		} else if len(aa) < len(bb) {
			return -1
		} else {
			return 0
		}
	} else {
		var aa, bb []any
		if aIsNumber {
			aa = []any{a}
			bb = b.([]any)
		} else {
			aa = a.([]any)
			bb = []any{b}
		}
		return compare(aa, bb)
	}
}

func parseInput(input []string) [][]any {
	lists := make([][]any, 0)
	for _, s := range input {
		if s != "" {
			var l []any
			json.Unmarshal([]byte(s), &l)
			lists = append(lists, l)
		}
	}
	return lists
}

func init() {
	challenges.RegisterChallengeFunc(2022, 13, 1, "day13.txt", part1)
	challenges.RegisterChallengeFunc(2022, 13, 2, "day13.txt", part2)
}
