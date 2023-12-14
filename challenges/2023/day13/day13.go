package aoc2023_day13

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type Pattern []string

func (p Pattern) String() string {
	return strings.Join(p, "\n")
}

func (p Pattern) Smudge(row, col int) Pattern {
	newPattern := make(Pattern, len(p))
	for r, s := range p {
		newRow := make([]byte, len(s))
		for c := range s {
			if r == row && c == col {
				if p[r][c] == '.' {
					newRow[c] = '#'
				} else {
					newRow[c] = '.'
				}
			} else {
				newRow[c] = p[r][c]
			}
		}
		newPattern[r] = string(newRow)
	}
	return newPattern
}

func (p Pattern) RotateR() Pattern {
	newPattern := make(Pattern, 0, len(p[0]))

	for c := 0; c < len(p[0]); c++ {
		newRow := make([]byte, 0, len(p))
		for r := len(p) - 1; r >= 0; r-- {
			newRow = append(newRow, p[r][c])
		}
		newPattern = append(newPattern, string(newRow))
	}
	return newPattern
}

func (p Pattern) FindSymmetryAxis(ignoreI int, ignoreAxis string) (index int, axis string, ok bool) {
	vIgnore := -1
	hIgnore := -1
	if ignoreAxis == "v" {
		vIgnore = ignoreI
	} else {
		hIgnore = ignoreI
	}

	if v := p.FindVerticalSymmetryAxis(vIgnore); v != -1 {
		return v, "v", true
	} else {

		v2 := p.FindHorizontalSymmetryAxis(hIgnore)
		if v2 != -1 {
			return v2, "h", true
		}
	}
	return -1, "", false
}

func (p Pattern) FindHorizontalSymmetryAxis(ignoreI int) int {
	for r := range p {
		r2 := r + 1
		if r2 < len(p) && p[r] == p[r2] {
			match := true
			if r2 != ignoreI {
				for i := 1; r2+i < len(p) && r-i >= 0; i++ {
					if p[r-i] != p[r2+i] {
						match = false
						break
					}
				}
				if match {
					return r2
				}
			}
		}
	}
	return -1
}

func (p Pattern) FindVerticalSymmetryAxis(ignoreI int) int {
	np := p.RotateR()
	return np.FindHorizontalSymmetryAxis(ignoreI)
}

func part1(input []string) (string, error) {
	patterns := parseInput(input)

	var result int
	for _, p := range patterns {
		v, a, _ := p.FindSymmetryAxis(-1, "")
		if a == "h" {
			result += (v * 100)
		} else {
			result += v
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	patterns := parseInput(input)

	var result int
	for _, p := range patterns {
		log.Debugf("Checking pattern:\n%s\n", p.String())
	patternloop:
		for r := range p {
			origV, origA, _ := p.FindSymmetryAxis(-1, "")
			for c := range p[r] {
				np := p.Smudge(r, c)
				v, a, ok := np.FindSymmetryAxis(origV, origA)
				if ok && (v != origV || a != origA) {
					if a == "h" {
						result += (v * 100)
					} else {
						result += v
					}
					log.Debugf("Found smudge at %d,%d!\n%s\n\n", r, c, np.String())
					break patternloop
				}
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []Pattern {
	p := make([]Pattern, 0, len(input))

	curPattern := make(Pattern, 0)
	for _, s := range input {
		if s == "" {
			p = append(p, curPattern)
			curPattern = make(Pattern, 0)
			continue
		}
		curPattern = append(curPattern, s)
	}
	p = append(p, curPattern)
	return p
}

func init() {
	challenges.RegisterChallengeFunc(2023, 13, 1, "day13.txt", part1)
	challenges.RegisterChallengeFunc(2023, 13, 2, "day13.txt", part2)
}
