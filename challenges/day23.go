package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day23_part1(input []string) (string, error) {
	cups := day23_parse(input)

	for i := 1; i <= 100; i++ {
		log.Debugf("-- move %d --\n", i)
		cups.Move()
	}

	log.Debugf("-- final --\nCups: %s\n", cups)

	var oneIdx int

	for i, v := range cups.Cups {
		if v == 1 {
			oneIdx = i
			break
		}
	}

	var b strings.Builder
	for i := (oneIdx + 1) % len(cups.Cups); i != oneIdx; i = (i + 1) % len(cups.Cups) {
		b.WriteString(fmt.Sprintf("%d", cups.Cups[i]))
	}

	return fmt.Sprintf("%s", b.String()), nil
}

func day23_part2(input []string) (string, error) {
	_ = day23_parse(input)
	var result int

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

type day23_cups struct {
	Cups     []int
	CurIndex int
	MinVal   int
	MaxVal   int
}

func (c *day23_cups) Move() {

	log.Debugf("Cups: %s\n", c)
	curVal := c.Cups[c.CurIndex]
	// tmp := append(make([]int, 0, len(c.Cups)), c.Cups[0:i1+1], c.Cups[c.CurIndex+3:])
	// pickedUp := c.Cups[c.CurIndex+1 : c.CurIndex+4]
	pickedUp := make([]int, 0, 3)
	tmp := make([]int, 0, len(c.Cups))
	tmp = append(tmp, curVal)
	for x := (c.CurIndex + 1) % len(c.Cups); x != c.CurIndex; x = (x + 1) % len(c.Cups) {
		if len(pickedUp) < 3 {
			pickedUp = append(pickedUp, c.Cups[x])
		} else {
			tmp = append(tmp, c.Cups[x])
		}
	}
	log.Debugf("Pick Up: %d, %d, %d\n", pickedUp[0], pickedUp[1], pickedUp[2])

	findVal := curVal - 1
	targIndex := -1
	for targIndex == -1 {
		if findVal < c.MinVal {
			findVal = c.MaxVal
		}
		for j := 0; j < len(tmp); j++ {
			if tmp[j] == findVal {
				targIndex = j
				break
			}
		}
		findVal -= 1
	}
	// log.Debugf("Temp: %v\n", tmp)
	// log.Debugf("Target index: %d\n", targIndex)
	log.Debugf("Target: %d\n", tmp[targIndex])

	c.Cups = append(make([]int, 0, len(c.Cups)), tmp[0:targIndex+1]...)
	c.Cups = append(c.Cups, pickedUp...)
	c.Cups = append(c.Cups, tmp[targIndex+1:]...)

	for i, v := range c.Cups {
		if v == curVal {
			c.CurIndex = (i + 1) % len(c.Cups)
			break
		}
	}
	log.Debugln()
}

func (c day23_cups) String() string {
	var b strings.Builder
	for i, v := range c.Cups {
		if i == c.CurIndex {
			b.WriteString(fmt.Sprintf("(%d) ", v))
		} else {
			b.WriteString(fmt.Sprintf(" %d  ", v))
		}
	}
	return strings.TrimSpace(b.String())
}

func day23_parse(input []string) day23_cups {
	nums := make([]int, 0, len(input[0]))
	min := 10
	max := -1
	for _, s := range input[0] {
		n, _ := strconv.Atoi(string(s))
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
		nums = append(nums, n)
	}
	return day23_cups{
		Cups:   nums,
		MinVal: min,
		MaxVal: max,
	}
}

func init() {
	registerChallengeFunc(23, 1, "day23.txt", day23_part1)
	registerChallengeFunc(23, 2, "day23.txt", day23_part2)
}
