package aoc2020_day23

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	cups := parseInput(input)

	for i := 1; i <= 100; i++ {
		log.Debugf("-- move %d --\n", i)
		cups.Move()
	}

	log.Debugf("-- final --\nCups: %s\n", cups)

	cup := cups.Current
	for cup.Label != 1 {
		cup = cup.Next
	}

	var b strings.Builder
	cup = cup.Next
	for cup.Label != 1 {
		b.WriteString(fmt.Sprintf("%d", cup.Label))
		cup = cup.Next
	}

	return b.String(), nil
}

func part2(input []string) (string, error) {
	cups := parseInput(input)

	// find the "last" node
	c := cups.Current
	for c.Next != cups.Current {
		c = c.Next
	}

	// append the rest of the million cups
	last := cups.MaxVal + 1000000 - len(cups.CupMap)
	for i := cups.MaxVal + 1; i <= last; i++ {
		c.Next = &cup{
			Label: i,
		}
		c = c.Next
		cups.CupMap[i] = c
		cups.MaxVal = i
	}
	c.Next = cups.Current

	debugWasEnabled := log.DebugEnabled()
	for i := 1; i <= 10000000; i++ {
		if (i < 5 || i%1000000 == 0 || i > 9999995) && debugWasEnabled {
			log.EnableDebugLogs(true)
			log.Debugf("-- move %8d --\n", i)
		}

		cups.Move()

		if (i < 5 || i%1000000 == 0 || i > 9999995) && debugWasEnabled {
			log.EnableDebugLogs(false)
		}
	}
	log.EnableDebugLogs(debugWasEnabled)

	c = cups.Current
	for c.Label != 1 {
		c = c.Next
	}

	result := c.Next.Label * c.Next.Next.Label
	log.Debugf("Max value: %d\n", cups.MaxVal)
	log.Debugf("Next two values: %d * %d = %d\n", c.Next.Label, c.Next.Next.Label, result)

	return fmt.Sprintf("%d", result), nil
}

type cup struct {
	Label int
	Next  *cup
}

type cups struct {
	Current *cup
	CupMap  map[int]*cup
	MinVal  int
	MaxVal  int
}

func (c *cups) Move() {
	log.Debugf("Cups: %s\n", c)
	curVal := c.Current.Label

	// cut out the next 3 nodes
	pickedUp := c.Current.Next
	c.Current.Next = pickedUp.Next.Next.Next

	log.Debugf("Pick Up: %d, %d, %d\n", pickedUp.Label, pickedUp.Next.Label, pickedUp.Next.Next.Label)

	// find the insertion target node
	var targetNode *cup
	targetVal := curVal

	for targetNode == nil {
		targetVal -= 1
		if targetVal < c.MinVal {
			targetVal = c.MaxVal
		}

		if n, ok := c.CupMap[targetVal]; ok {
			if n != pickedUp && n != pickedUp.Next && n != pickedUp.Next.Next {
				targetNode = n
				break
			}
		}
	}
	log.Debugf("Target: %d\n", targetVal)

	// insert the "picked up" nodes after the insertion target
	targetNext := targetNode.Next
	targetNode.Next = pickedUp
	pickedUp.Next.Next.Next = targetNext

	// update the "Current" node pointer
	c.Current = c.Current.Next
	log.Debugln()
}

func (c cups) String() string {
	var b strings.Builder

	first := c.Current
	b.WriteString(fmt.Sprintf("(%d) ", first.Label))
	cur := first.Next
	MAX := 10
	n := 1
	for cur != first && n < MAX {
		b.WriteString(fmt.Sprintf(" %d  ", cur.Label))
		cur = cur.Next
		n++
	}
	if MAX < len(c.CupMap) {
		b.WriteString(fmt.Sprintf("...  (and %d more)", len(c.CupMap)-MAX))
	}
	return strings.TrimSpace(b.String())
}

func parseInput(input []string) cups {
	min := 10
	max := -1
	var first *cup
	var prev *cup
	cupMap := make(map[int]*cup)
	for _, s := range input[0] {
		n, _ := strconv.Atoi(string(s))
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}

		c := &cup{
			Label: n,
		}
		if prev == nil {
			first = c
		} else {
			prev.Next = c
		}
		prev = c
		cupMap[n] = c
	}
	prev.Next = first
	return cups{
		Current: first,
		CupMap:  cupMap,
		MinVal:  min,
		MaxVal:  max,
	}
}

func init() {
	challenges.RegisterChallengeFunc(2020, 23, 1, "day23.txt", part1)
	challenges.RegisterChallengeFunc(2020, 23, 2, "day23.txt", part2)
}
