package aoc2023_day12

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type SpringRecord struct {
	Springs string
	Counts  []int
}

func (sr SpringRecord) String() string {
	return fmt.Sprintf("%s %s", sr.Springs, util.JoinInts(sr.Counts, ","))
}

func part1(input []string) (string, error) {
	parsed := parseInput(input)

	var result int
	for _, sr := range parsed {
		n := countPossibilitiesMemoized(sr)
		log.Debugf("Found %d possibilities for record: %s\n", n, sr)
		result += n
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)

	for _, sr := range parsed {
		sr.FiveFold()
	}

	var result int
	for _, sr := range parsed {
		n := countPossibilitiesMemoized(sr)
		log.Debugf("Found %d possibilities for record: %s\n", n, sr)
		result += n
	}
	return fmt.Sprintf("%d", result), nil
}

var resultCache = make(map[string]int)

func countPossibilitiesMemoized(sr *SpringRecord) int {
	key := sr.String()
	if _, ok := resultCache[key]; !ok {
		resultCache[key] = countPossibilities(sr)
	}
	return resultCache[key]
}

func countPossibilities(sr *SpringRecord) int {
	if len(sr.Springs) == 0 {
		if len(sr.Counts) == 0 {
			return 1
		}
		return 0
	}

	if len(sr.Counts) == 0 {
		for _, s := range sr.Springs {
			if s == '#' {
				return 0
			}
		}
		return 1
	}

	if len(sr.Springs) < util.Sum(sr.Counts)+len(sr.Counts)-1 {
		return 0
	}

	switch sr.Springs[0] {
	case '.':
		return countPossibilitiesMemoized(&SpringRecord{sr.Springs[1:], sr.Counts})
	case '#':
		currentCount := sr.Counts[0]
		for i := 1; i < currentCount; i++ {
			if i >= len(sr.Springs) || sr.Springs[i] == '.' {
				return 0
			}
		}
		if len(sr.Springs) == currentCount && len(sr.Counts) == 1 {
			return 1 // we've reached the end of the springs and the counts
		}
		if len(sr.Springs) > currentCount && sr.Springs[currentCount] == '#' {
			return 0 // we've reached the end of the current count, but the next spring is still bad
		}
		return countPossibilitiesMemoized(&SpringRecord{sr.Springs[currentCount+1:], sr.Counts[1:]})
	case '?':
		result1 := countPossibilitiesMemoized(&SpringRecord{string('#') + sr.Springs[1:], sr.Counts})
		result2 := countPossibilitiesMemoized(&SpringRecord{string('.') + sr.Springs[1:], sr.Counts})
		return result1 + result2
	default:
		panic(fmt.Sprintf("WTF: Unknown spring type: %s", string(sr.Springs[0])))
	}
}

func (sr *SpringRecord) FiveFold() *SpringRecord {
	sr.Springs = fmt.Sprintf("%s?%s?%s?%s?%s", sr.Springs, sr.Springs, sr.Springs, sr.Springs, sr.Springs)
	arr := make([]int, 0, len(sr.Counts)*5)
	for i := 0; i < 5; i++ {
		arr = append(arr, sr.Counts...)
	}
	sr.Counts = arr
	log.Debugln(sr)
	return sr
}

func parseInput(input []string) []*SpringRecord {
	records := make([]*SpringRecord, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " ")
		records = append(records, &SpringRecord{
			Springs: parts[0],
			Counts:  util.AtoiSplit(parts[1], ","),
		})
	}
	return records
}

func init() {
	challenges.RegisterChallengeFunc(2023, 12, 1, "day12.txt", part1)
	challenges.RegisterChallengeFunc(2023, 12, 2, "day12.txt", part2)
}
