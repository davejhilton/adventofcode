package aoc2021_day13

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	coords, folds, w, h := parseInput(input)
	p := NewPaper(w, h, &coords)
	p.Fold(folds[0])
	log.Debugln(p)
	result := p.Count()
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	coords, folds, w, h := parseInput(input)
	p := NewPaper(w, h, &coords)
	for _, f := range folds {
		p.Fold(f)
	}
	return fmt.Sprintf("\n%s", p.String()), nil
}

func parseInput(input []string) ([]Coord, []FoldingPoint, int, int) {
	coords := make([]Coord, 0)
	i := 0

	maxX := 0
	maxY := 0
	for ; i < len(input); i++ {
		if input[i] == "" {
			break
		}
		parts := strings.Split(input[i], ",")
		x := util.Atoi(parts[0])
		y := util.Atoi(parts[1])
		coords = append(coords, Coord{
			Col: x,
			Row: y,
		})
		maxX = util.Max(maxX, x)
		maxY = util.Max(maxY, y)
	}

	folds := make([]FoldingPoint, 0)
	for i += 1; i < len(input); i++ {
		parts := strings.Split(input[i], "=")
		folds = append(folds, FoldingPoint{
			Axis:  strings.Replace(parts[0], "fold along ", "", -1),
			Value: util.Atoi(parts[1]),
		})
	}

	return coords, folds, maxX + 1, maxY + 1
}

type Coord struct {
	Row int
	Col int
}

type FoldingPoint struct {
	Axis  string
	Value int
}

func NewPaper(w int, h int, coords *[]Coord) paper {
	p := make(paper, h)
	for r := range p {
		p[r] = make([]bool, w)
	}
	if coords != nil {
		for _, c := range *coords {
			p[c.Row][c.Col] = true
		}
	}
	return p
}

type paper [][]bool

func (p paper) String() string {
	var sb strings.Builder
	for r := range p {
		if r != 0 {
			sb.WriteString("\n")
		}
		for c := range p[r] {
			if p[r][c] {
				sb.WriteRune('#')
			} else {
				sb.WriteRune('.')
			}
		}
	}
	return sb.String()
}

func (p paper) Count() int {
	count := 0
	for r := range p {
		for c := range p[r] {
			if p[r][c] {
				count += 1
			}
		}
	}
	return count
}

func (p *paper) Fold(f FoldingPoint) {
	for r := range *p {
		r2 := r
		if f.Axis == "y" && r == f.Value {
			continue
		} else if f.Axis == "y" && r > f.Value {
			r2 = f.Value - (r - f.Value)
		}
		for c := range (*p)[r] {
			c2 := c
			if f.Axis == "x" && c == f.Value {
				continue
			} else if f.Axis == "x" && c > f.Value {
				c2 = f.Value - (c - f.Value)
			}
			(*p)[r2][c2] = (*p)[r][c] || (*p)[r2][c2]
		}
	}

	if f.Axis == "x" {
		for r := range *p {
			(*p)[r] = (*p)[r][:f.Value]
		}
	} else {
		*p = (*p)[:f.Value]
	}
}

func init() {
	challenges.RegisterChallengeFunc(2021, 13, 1, "day13.txt", part1)
	challenges.RegisterChallengeFunc(2021, 13, 2, "day13.txt", part2)
}
