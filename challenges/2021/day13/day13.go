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
	log.Debugf("-------------\n%s\n", p)
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

func parseInput(input []string) (coords []Coord, folds []FoldingPoint, w int, h int) {
	coords = make([]Coord, 0)
	i := 0

	w = 0
	h = 0
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
		w = util.Max(w, x)
		h = util.Max(h, y)
	}
	w += 1
	h += 1

	folds = make([]FoldingPoint, 0)
	for i += 1; i < len(input); i++ {
		parts := strings.Split(input[i], "=")
		folds = append(folds, FoldingPoint{
			Axis:  strings.Replace(parts[0], "fold along ", "", -1),
			Value: util.Atoi(parts[1]),
		})
	}

	return
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
	if len(p) > 100 || len(p[0]) > 100 {
		return fmt.Sprintf("Paper - w: %d, h: %d", len(p[0]), len(p))
	}
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
	log.Debugf("Folding at %s=%d\n", f.Axis, f.Value)
	if f.Axis == "y" {
		p.foldY(f.Value)
	} else {
		p.foldX(f.Value)
	}
}

func (p *paper) foldX(xVal int) {
	for r := range *p {
		for c := range (*p)[r] {
			c2 := c
			if c == xVal {
				continue
			} else if c > xVal {
				c2 = xVal - (c - xVal)
			}
			(*p)[r][c2] = (*p)[r][c] || (*p)[r][c2]
		}
	}
	for r := range *p {
		(*p)[r] = (*p)[r][0:xVal]
	}
}

func (p *paper) foldY(yVal int) {
	for r := range *p {
		r2 := r
		if r == yVal {
			continue
		} else if r > yVal {
			r2 = yVal - (r - yVal)
		}
		for c := range (*p)[r] {
			(*p)[r2][c] = (*p)[r][c] || (*p)[r2][c]
		}
	}
	*p = (*p)[0:yVal]
}

func init() {
	challenges.RegisterChallengeFunc(2021, 13, 1, "day13.txt", part1)
	challenges.RegisterChallengeFunc(2021, 13, 2, "day13.txt", part2)
}
