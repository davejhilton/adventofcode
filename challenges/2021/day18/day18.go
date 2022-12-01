package aoc2021_day18

import (
	"fmt"
	"math"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed[0])

	var curN *Number
	for i, n := range parsed {
		if i == 0 {
			curN = n
			continue
		}

		curN = addAndReduce(curN, n)
	}

	result := magnitude(curN)
	return fmt.Sprintf("%d", result), nil
}

func addAndReduce(n1 *Number, n2 *Number) *Number {

	curN := addNumbers(n1, n2)

	needsReducing := true
	for needsReducing {
		if expl, path := checkExplode(curN, curN, make([]int, 0), 0); expl {
			explode(curN, path)
		} else if spl, path := checkSplit(curN, curN, make([]int, 0)); spl {
			split(curN, path)
		} else {
			needsReducing = false
			log.Debugln(curN)
		}
	}
	return curN
}

func checkExplode(n *Number, topLevelN *Number, pathFromTop []int, depth int) (bool, []int) {
	if n.IsLiteral {
		return false, pathFromTop
	}
	if depth == 4 {
		return true, pathFromTop
	}
	if expl, path2 := checkExplode(n.Left, topLevelN, append(pathFromTop[0:], 0), depth+1); expl {
		log.Debugf("NEEDS TO EXPLODE (left): %s\n", n)
		return expl, path2
	}
	if expl, path2 := checkExplode(n.Right, topLevelN, append(pathFromTop[0:], 1), depth+1); expl {
		log.Debugf("NEEDS TO EXPLODE (right): %s\n", n)
		return expl, path2
	}
	return false, pathFromTop
}

func addNumbers(n1 *Number, n2 *Number) *Number {
	return &Number{
		IsLiteral: false,
		Left:      n1,
		Right:     n2,
	}
}

func magnitude(n *Number) int {
	if n.IsLiteral {
		return n.Val
	}
	return 3*magnitude(n.Left) + 2*magnitude(n.Right)
}

func explode(topLevelNumber *Number, pathFromTop []int) {
	var parentOfDeepestRight *Number
	var parentOfDeepestLeft *Number
	curN := topLevelNumber
	for i, v := range pathFromTop {
		if i == len(pathFromTop)-1 {
			// do the explode
			var expl *Number
			if v == 0 {
				expl = curN.Left
				parentOfDeepestLeft = curN
			} else {
				expl = curN.Right
				parentOfDeepestRight = curN
			}
			log.Debugf("EXPLODING! path: %#v\nparent of expl: %s\nexpl: %s\n", pathFromTop, curN, expl)
			log.Debugf("ParentOfDeepestRight: %s\n", parentOfDeepestRight)
			log.Debugf("ParentOfDeepestLeft: %s\n", parentOfDeepestLeft)
			if parentOfDeepestRight != nil {
				n := parentOfDeepestRight.Left
				for !n.IsLiteral {
					n = n.Right
				}
				n.Val += expl.Left.Val
			}
			if parentOfDeepestLeft != nil {
				n := parentOfDeepestLeft.Right
				for !n.IsLiteral {
					n = n.Left
				}
				n.Val += expl.Right.Val
			}
			if v == 0 {
				curN.Left = &Number{
					IsLiteral: true,
					Val:       0,
				}
			} else {
				curN.Right = &Number{
					IsLiteral: true,
					Val:       0,
				}
			}
			break
		} else {
			if v == 0 {
				parentOfDeepestLeft = curN
				curN = curN.Left
			}
			if v == 1 {
				parentOfDeepestRight = curN
				curN = curN.Right
			}
		}
	}

	log.Debugf("AFTER EXPLODE: %s\n", topLevelNumber)
}

func split(topLevelNumber *Number, pathFromTop []int) {
	curN := topLevelNumber
	for i, v := range pathFromTop {
		if i == len(pathFromTop)-1 {
			// do the split
			var spl *Number
			if v == 0 {
				spl = curN.Left
			} else {
				spl = curN.Right
			}
			log.Debugf("SPLITTING! path: %#v\nparent of spl: %s\nspl: %s\n", pathFromTop, curN, spl)
			/*

			 */
			splLeft := int(float64(float64(spl.Val) / float64(2)))
			splRight := int(math.Round(float64(float64(spl.Val) / float64(2))))
			splResult := &Number{
				IsLiteral: false,
				Left: &Number{
					IsLiteral: true,
					Val:       splLeft,
				},
				Right: &Number{
					IsLiteral: true,
					Val:       splRight,
				},
			}
			if v == 0 {
				curN.Left = splResult
			} else {
				curN.Right = splResult
			}
			break
		} else {
			if v == 0 {
				curN = curN.Left
			}
			if v == 1 {
				curN = curN.Right
			}
		}
	}

	log.Debugf("AFTER SPLIT: %s\n", topLevelNumber)
}

func checkSplit(n *Number, topLevelN *Number, pathFromTop []int) (bool, []int) {
	if n.IsLiteral {
		return n.Val >= 10, pathFromTop
	}
	if expl, path2 := checkSplit(n.Left, topLevelN, append(pathFromTop[0:], 0)); expl {
		return expl, path2
	}
	if expl, path2 := checkSplit(n.Right, topLevelN, append(pathFromTop[0:], 1)); expl {
		return expl, path2
	}
	return false, pathFromTop
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", parsed[0])

	tried := make([]string, 0)
	maxMag := -1
	for i := 0; i < len(parsed); i++ {
		for j := 0; j < len(parsed); j++ {
			if i == j {
				continue
			}

			x, y := i, j
			reversed := false
			for {
				strP := fmt.Sprintf("%d,%d", x, y)
				alreadyTried := false
				for _, p := range tried {
					if p == strP {
						alreadyTried = true
						break
					}
				}
				if !alreadyTried {
					fmt.Printf("TRYING PAIR: %s\n", strP)
					num := addAndReduce(cloneNumber(parsed[x]), cloneNumber(parsed[y]))
					maxMag = util.Max(maxMag, magnitude(num))
				}
				if reversed {
					break
				} else {
					reversed = true
					x, y = y, x
				}
			}
		}
	}

	result := maxMag
	return fmt.Sprintf("%d", result), nil
}

func cloneNumber(n *Number) *Number {
	if n.IsLiteral {
		return &Number{
			IsLiteral: true,
			Val:       n.Val,
		}
	}
	return &Number{
		IsLiteral: false,
		Left:      cloneNumber(n.Left),
		Right:     cloneNumber(n.Right),
	}

}

func parseInput(input []string) []*Number {
	nums := make([]*Number, 0, len(input))
	for _, s := range input {
		n, _ := parseNumber(s)
		nums = append(nums, n)
	}
	return nums
}

func parseNumber(s string) (*Number, string) {
	runes := []rune(s)
	if len(s) == 0 {
		log.Debugf("Parsing empty string.... %s\n", s)
		return nil, ""
	} else if len(runes) == 0 {
		log.Debugf("Parsing empty runes?.... '%s'\n", len(runes), s)
	}
	if runes[0] != '[' {
		numStr := string(runes[0])
		numLen := 1
		if runes[1] != ']' && runes[1] != ',' {
			// it's a two-digit num
			numStr = string(runes[0:2])
			log.Debugf("LENGTH 2! %s\n", numStr)
			numLen = 2
		}
		log.Debugf("parsing literal: %s\n", numStr)
		n := &Number{
			IsLiteral: true,
			Val:       util.Atoi(numStr),
		}
		return n, s[numLen:]
	}

	// skip the '[', parse the left number
	log.Debugf("Parsing Left: '%s'\n", s[1:])
	left, s := parseNumber(s[1:])
	// skip the ',', parse the right number
	log.Debugf("Parsing Right: '%s'\n", s[1:])
	right, s := parseNumber(s[1:])

	n := &Number{
		IsLiteral: false,
		Left:      left,
		Right:     right,
	}

	return n, s[1:]
}

/*
start -> literal   -> end
start -> LLBracket -> <number> -> LRBracket -> comma -> RLBracket -> <number> -> RRBracket -> end
*/

type State struct {
	Name string
}

type Number struct {
	IsLiteral bool
	Val       int
	Left      *Number
	Right     *Number
}

func (n *Number) String() string {
	if n.IsLiteral {
		return fmt.Sprintf("%d", n.Val)
	}
	return fmt.Sprintf("[%s,%s]", n.Left, n.Right)
}

func init() {
	challenges.RegisterChallengeFunc(2021, 18, 1, "day18.txt", part1)
	challenges.RegisterChallengeFunc(2021, 18, 2, "day18.txt", part2)
}
