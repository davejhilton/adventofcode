package aoc2020_day6

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	sum := 0
	groups := getGroups(input)
	for i, group := range groups {
		log.Debugf("Group: %d\n", i)
		allQs := make(map[rune]bool)
		for _, line := range group {
			log.Debugf("\t%s\n", line)
			for _, q := range line {
				allQs[q] = true
			}
		}
		log.Debug("questions: ")
		for q := range allQs {
			log.Debugf("%s", string(q))
		}
		log.Debugf(" - %d\n\n", len(allQs))
		sum += len(allQs)
	}
	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	sum := 0
	groups := getGroups(input)
	for g, group := range groups {
		log.Debugf("Group: %d\n", g)
		commonQs := make(map[rune]bool)
		for i, line := range group {
			if i == 0 {
				for _, q := range line {
					commonQs[q] = true
				}
				log.Debugf("\t%s\n", line)
			} else {
				for cq := range commonQs {
					found := false
					for _, c := range line {
						if c == cq {
							found = true
							break
						}
					}
					if !found {
						delete(commonQs, cq)
					}
				}
			}
		}

		log.Debug("questions: ")
		for q := range commonQs {
			log.Debugf("%s", string(q))
		}
		log.Debugf(" - %d\n\n", len(commonQs))
		sum += len(commonQs)
	}
	return fmt.Sprintf("%d", sum), nil
}

func getGroups(input []string) [][]string {
	input = append(input, "")
	groups := make([][]string, 0)
	cur := make([]string, 0)
	for _, line := range input {
		if line == "" {
			groups = append(groups, cur)
			cur = make([]string, 0)
		} else {
			cur = append(cur, line)
		}
	}
	return groups
}

func init() {
	challenges.RegisterChallengeFunc(2020, 6, 1, "day06.txt", part1)
	challenges.RegisterChallengeFunc(2020, 6, 2, "day06.txt", part2)
}
