package aoc2022_day11

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Monkey struct {
	Index     int
	Items     []int
	Operation Operation
	Test      int
	ThrowTo   map[bool]int
}

type Operation struct {
	First  string
	Op     string
	Second string
}

func part1(input []string) (string, error) {
	monkeys := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", monkeys)

	inspectCounts := make([]int, len(monkeys))
	for round := 0; round < 20; round++ {
		for monkeyIdx, monkey := range monkeys {
			for _, item := range monkey.Items {
				inspectCounts[monkeyIdx] += 1

				var f, s int
				if monkey.Operation.First == "old" {
					f = item
				} else {
					f = util.Atoi(monkey.Operation.First)
				}
				if monkey.Operation.Second == "old" {
					s = item
				} else {
					s = util.Atoi(monkey.Operation.Second)
				}
				if monkey.Operation.Op == "*" {
					item = f * s
				} else {
					item = f + s
				}
				item = item / 3
				isDiv := item%monkey.Test == 0
				throwTo := monkey.ThrowTo[isDiv]
				monkeys[throwTo].Items = append(monkeys[throwTo].Items, item)
			}
			monkey.Items = make([]int, 0)
		}
	}

	log.Debugf("Inspect counts: %v\n", inspectCounts)

	sort.Ints(inspectCounts)

	biggest := inspectCounts[len(inspectCounts)-2:]

	var result int = biggest[0] * biggest[1]
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	monkeys := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", monkeys)

	primeProduct := 1
	for _, m := range monkeys {
		primeProduct *= m.Test
	}

	inspectCounts := make([]int, len(monkeys))
	for round := 0; round < 10000; round++ {
		for monkeyIdx, monkey := range monkeys {
			for _, item := range monkey.Items {
				inspectCounts[monkeyIdx] += 1

				var f, s int
				if monkey.Operation.First == "old" {
					f = item
				} else {
					f = util.Atoi(monkey.Operation.First)
				}
				if monkey.Operation.Second == "old" {
					s = item
				} else {
					s = util.Atoi(monkey.Operation.Second)
				}
				if monkey.Operation.Op == "*" {
					item = f * s
				} else {
					item = f + s
				}

				isDiv := item%monkey.Test == 0
				throwTo := monkey.ThrowTo[isDiv]
				item = item % primeProduct
				monkeys[throwTo].Items = append(monkeys[throwTo].Items, item)
			}
			monkey.Items = make([]int, 0)
		}
	}

	log.Debugf("Inspect counts: %v\n", inspectCounts)

	sort.Ints(inspectCounts)

	biggest := inspectCounts[len(inspectCounts)-2:]

	var result int = biggest[0] * biggest[1]
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []*Monkey {
	monkeys := make([]*Monkey, 0)
	for i := 0; i < len(input); i += 7 {
		monkey := Monkey{
			Items: make([]int, 0),
		}
		fmt.Sscanf(input[i], "Monkey %d:", &monkey.Index)
		numStrs := strings.Split(strings.Split(input[i+1], ": ")[1], ", ")
		for _, s := range numStrs {
			monkey.Items = append(monkey.Items, util.Atoi(s))
		}
		fmt.Sscanf(input[i+2], "  Operation: new = %s %s %s", &monkey.Operation.First, &monkey.Operation.Op, &monkey.Operation.Second)
		fmt.Sscanf(input[i+3], "  Test: divisible by %d", &monkey.Test)

		var trueVal, falseVal int
		fmt.Sscanf(input[i+4], "    If true: throw to monkey %d", &trueVal)
		fmt.Sscanf(input[i+5], "    If false: throw to monkey %d", &falseVal)
		monkey.ThrowTo = map[bool]int{true: trueVal, false: falseVal}

		monkeys = append(monkeys, &monkey)
	}
	return monkeys
}

func init() {
	challenges.RegisterChallengeFunc(2022, 11, 1, "day11.txt", part1)
	challenges.RegisterChallengeFunc(2022, 11, 2, "day11.txt", part2)
}
