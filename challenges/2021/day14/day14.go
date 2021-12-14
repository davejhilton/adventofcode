package aoc2021_day14

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	head, rules := parseInput(input)

	N_PASSES := 10

	for i := 1; i <= N_PASSES; i++ {
		cur := head
		for cur != nil && cur.Next != nil {
			next := cur.Next
			pair := string([]rune{cur.Char, next.Char})
			cur.Next = &Node{
				Char: rules[pair],
				Next: next,
			}
			cur = next
		}
	}

	counts := make(RuneCounter)
	cur := head
	for cur != nil {
		counts.Incr(cur.Char, 1)
		cur = cur.Next
	}

	max := -1
	min := math.MaxInt
	for _, c := range counts {
		min = util.Min(min, c)
		max = util.Max(max, c)
	}

	result := max - min
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	head, rules := parseInput(input)

	pairCounts := make(PairCounter)
	// first, dump everything from the "template" into a map of pairs -> count
	cur := head
	for cur != nil && cur.Next != nil {
		next := cur.Next
		pairCounts.Incr(string([]rune{cur.Char, next.Char}), 1)
		cur = cur.Next
	}
	tail := cur // we'll need to know what the last element is, for later...

	// run all of the "insertion" passes, just counting the number of "pairs" we encounter after each
	N_PASSES := 40
	for i := 1; i <= N_PASSES; i++ {
		newCounts := make(PairCounter)
		for p, c := range pairCounts {
			insert := rules[p]
			chars := []rune(p)
			newCounts.Incr(string([]rune{chars[0], insert}), c)
			newCounts.Incr(string([]rune{insert, chars[1]}), c)
		}
		pairCounts = newCounts
	}

	// now... count all of occurrences of each Char in the "pairs"
	counts := make(RuneCounter)
	for p, c := range pairCounts {
		chars := []rune(p)
		counts.Incr(chars[0], c)
		counts.Incr(chars[1], c)
	}

	// aaaaaannd divide them by 2, since chars get double-counted (left-side pair, right-side pair)
	// ... except the very first char and last char in the polymer. so add 1 back for those
	for r := range counts {
		counts[r] /= 2
		if r == head.Char {
			counts[r] += 1
		}
		if r == tail.Char {
			counts[r] += 1
		}
	}

	// find the min and max
	var max = int(-1)
	var min = int(math.MaxInt)
	for _, c := range counts {
		min = util.Min(min, c)
		max = util.Max(max, c)
	}

	var result = max - min
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) (head *Node, rules map[string]rune) {
	tplStr := input[0]
	var cur *Node
	for _, r := range tplStr {
		if head == nil {
			head = &Node{
				Char: r,
			}
			cur = head
		} else {
			cur.Next = &Node{
				Char: r,
			}
			cur = cur.Next
		}
	}
	ruleStrs := input[2:]
	rules = make(map[string]rune)
	for _, s := range ruleStrs {
		parts := strings.Split(s, " -> ")
		rules[parts[0]] = ([]rune(parts[1]))[0]
	}
	return head, rules
}

type PairCounter map[string]int

func (pc *PairCounter) Incr(pair string, incrBy int) {
	if _, ok := (*pc)[pair]; !ok {
		(*pc)[pair] = 0
	}
	(*pc)[pair] += incrBy
}

type RuneCounter map[rune]int

func (rc *RuneCounter) Incr(r rune, incrBy int) {
	if _, ok := (*rc)[r]; !ok {
		(*rc)[r] = 0
	}
	(*rc)[r] += incrBy
}

type Node struct {
	Char rune
	Next *Node
}

func (n *Node) String() string {
	var sb strings.Builder
	for n != nil {
		sb.WriteRune(n.Char)
		n = n.Next
	}
	return sb.String()
}

func init() {
	challenges.RegisterChallengeFunc(2021, 14, 1, "day14.txt", part1)
	challenges.RegisterChallengeFunc(2021, 14, 2, "day14.txt", part2)
}
