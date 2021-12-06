package challenges2021

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day05_part1(input []string) (string, error) {
	lines, w, h := day05_parse(input)
	g := NewGrid(w, h)
	log.Debugf("Grid size (w x h): %d x %d\n", w, h)

	for _, l := range lines {
		DrawLine(g, l, false)
	}
	log.Debugf("%s\n\n", *g)
	result := g.CountOverlaps()

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day05_part2(input []string) (string, error) {
	lines, w, h := day05_parse(input)
	g := NewGrid(w, h)
	log.Debugf("Grid size (w x h): %d x %d\n", w, h)

	for _, l := range lines {
		DrawLine(g, l, true)
	}
	log.Debugf("%s\n\n", *g)
	result := g.CountOverlaps()

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func day05_parse(input []string) ([]*line, int, int) {
	lines := make([]*line, 0, len(input))
	maxX := 0
	maxY := 0
	for _, coordStr := range input {
		parts := strings.SplitN(coordStr, " -> ", 2)
		p1 := strings.SplitN(parts[0], ",", 2)
		x1s, y1s := p1[0], p1[1]
		x1, _ := strconv.Atoi(x1s)
		y1, _ := strconv.Atoi(y1s)

		p2 := strings.SplitN(parts[1], ",", 2)
		x2s, y2s := p2[0], p2[1]
		x2, _ := strconv.Atoi(x2s)
		y2, _ := strconv.Atoi(y2s)
		lines = append(lines, &line{
			x1: x1,
			y1: y1,
			x2: x2,
			y2: y2,
		})
		maxX = max(maxX, x1, x2)
		maxY = max(maxY, y1, y2)
	}
	return lines, maxX + 1, maxY + 1
}

func DrawLine(g *grid, l *line, includeDiagonal bool) {
	if l.y2 == l.y1 {
		y := l.y1
		start := min(l.x1, l.x2)
		end := max(l.x1, l.x2)
		for x := start; x <= end; x++ {
			(*g)[y][x] += 1
		}
	} else if l.x2 == l.x1 {
		x := l.x1
		start := min(l.y1, l.y2)
		end := max(l.y1, l.y2)
		for y := start; y <= end; y++ {
			(*g)[y][x] += 1
		}
	} else if includeDiagonal {
		// log.Debugf("Drawing diagonal: %s\n", *l)
		xStart := min(l.x1, l.x2)
		xEnd := max(l.x1, l.x2)
		yStart := l.y1
		yAsc := l.y2 > l.y1
		if xStart == l.x2 {
			yStart = l.y2
			yAsc = l.y1 > l.y2
		}

		for x, y := xStart, yStart; x <= xEnd; x += 1 {
			// log.Debugf("   point: %d,%d\n", x, y)
			(*g)[y][x] += 1

			if yAsc {
				y += 1
			} else {
				y -= 1
			}
		}
	}
}

type line struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func (l line) String() string {
	return fmt.Sprintf("%d,%d -> %d,%d", l.x1, l.y1, l.x2, l.y2)
}

type grid [][]int

func NewGrid(w int, h int) *grid {
	g := make(grid, h)
	for i := range g {
		g[i] = make([]int, w)
	}
	return &g
}

func (g grid) CountOverlaps() int {
	overlaps := 0
	for i := range g {
		for j := range g[i] {
			if g[i][j] > 1 {
				overlaps += 1
			}
		}
	}
	return overlaps
}

func (g grid) String() string {
	var sb strings.Builder
	for i, row := range g {
		if i != 0 {
			sb.WriteString("\n")
		}
		for _, n := range row {
			if n == 0 {
				sb.WriteString(".")
			} else {
				fmt.Fprintf(&sb, "%d", n)
			}
		}
	}
	return sb.String()
}

func max(nums ...int) int {
	max := math.MinInt
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

func min(nums ...int) int {
	min := math.MaxInt
	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	return min
}

func init() {
	challenges.RegisterChallengeFunc(2021, 5, 1, "day05.txt", day05_part1)
	challenges.RegisterChallengeFunc(2021, 5, 2, "day05.txt", day05_part2)
}
