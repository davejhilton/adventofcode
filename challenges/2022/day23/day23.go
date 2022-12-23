package aoc2022_day23

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

const (
	EMPTY = 0
	ELF   = 1

	N int = 0
	S int = 1
	W int = 2
	E int = 3
)

var (
	DIRECTIONS = []int{N, S, W, E}

	curDirectionIdx = 0
)

func part1(input []string) (string, error) {
	grid := parseInput(input)

	MoveElves(grid, 10)

	var result = 0
	for r := grid.MinRow; r <= grid.MaxRow; r++ {
		for c := grid.MinCol; c <= grid.MaxCol; c++ {
			if v := grid.Get(r, c); v != ELF {
				result++
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)

	rounds := MoveElves(grid, 10_000)

	return fmt.Sprintf("%d", rounds), nil
}

func MoveElves(grid *Grid, maxIterations int) int {
	log.Debugf("\n== START ==\n%s\n", grid)

	var round int
	for round = 1; round <= maxIterations; round++ {
		proposalCounts := NewGrid()
		proposalByElf := make(map[Coord]*Coord)
		for r := range grid.points {
			for c := range grid.points[r] {
				proposedMove := false
				if AnyNeighbors(r, c, grid) {
					for di := 0; di < 4; di++ {
						dIdx := (curDirectionIdx + di) % len(DIRECTIONS)
						if !NeighborsInDirection(DIRECTIONS[dIdx], r, c, grid) {
							r1, c1 := GetCoordsInDirection(DIRECTIONS[dIdx], r, c)
							proposalCounts.Incr(r1, c1, 1)
							proposalByElf[Coord{r, c}] = &Coord{r1, c1}
							proposedMove = true
							break
						}
					}
				}
				if !proposedMove {
					proposalByElf[Coord{r, c}] = nil
				}
			}
		}

		grid = NewGrid()
		nMoves := 0
		for elf, prop := range proposalByElf {
			if prop != nil && proposalCounts.Get(prop.Row, prop.Col) == 1 {
				grid.Set(prop.Row, prop.Col, ELF)
				nMoves++
			} else {
				grid.Set(elf.Row, elf.Col, ELF)
			}
		}

		curDirectionIdx = (curDirectionIdx + 1) % len(DIRECTIONS)

		log.Debugf("\n== END OF ROUND %d ==\n%s\n", round, grid)

		if nMoves == 0 {
			break
		}
	}

	return round
}

func AnyNeighbors(r int, c int, grid *Grid) (occupied bool) {
	coords := []Coord{
		{r - 1, c - 1},
		{r - 1, c},
		{r - 1, c + 1},
		{r, c - 1},
		{r, c + 1},
		{r + 1, c - 1},
		{r + 1, c},
		{r + 1, c + 1},
	}
	for _, c := range coords {
		if grid.Get(c.Row, c.Col) != EMPTY {
			return true
		}
	}
	return false
}

func NeighborsInDirection(d int, r int, c int, grid *Grid) (occupied bool) {
	switch d {
	case N:
		if grid.Get(r-1, c-1) != EMPTY || grid.Get(r-1, c) != EMPTY || grid.Get(r-1, c+1) != EMPTY {
			return true
		}
	case S:
		if grid.Get(r+1, c-1) != EMPTY || grid.Get(r+1, c) != EMPTY || grid.Get(r+1, c+1) != EMPTY {
			return true
		}
	case W:
		if grid.Get(r-1, c-1) != EMPTY || grid.Get(r, c-1) != EMPTY || grid.Get(r+1, c-1) != EMPTY {
			return true
		}
	case E:
		if grid.Get(r-1, c+1) != EMPTY || grid.Get(r, c+1) != EMPTY || grid.Get(r+1, c+1) != EMPTY {
			return true
		}
	}
	return false
}

func GetCoordsInDirection(d int, r int, c int) (r1 int, c1 int) {
	switch d {
	case N:
		return r - 1, c
	case S:
		return r + 1, c
	case W:
		return r, c - 1
	case E:
		return r, c + 1
	}
	return r, c
}

func parseInput(input []string) *Grid {
	grid := NewGrid()
	for r, s := range input {
		for c, ch := range []rune(s) {
			if ch == '#' {
				grid.Set(r, c, ELF)
			}
		}
	}
	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2022, 23, 1, "day23.txt", part1)
	challenges.RegisterChallengeFunc(2022, 23, 2, "day23.txt", part2)
}

type Coord struct {
	Row int
	Col int
}

type Grid struct {
	points map[int]map[int]int
	MinCol int
	MinRow int
	MaxCol int
	MaxRow int
}

func NewGrid() *Grid {
	return &Grid{
		points: make(map[int]map[int]int),
		MinCol: math.MaxInt,
		MaxCol: math.MinInt,
		MinRow: math.MaxInt,
		MaxRow: math.MinInt,
	}
}

func (g *Grid) Set(row, col, val int) {
	if _, ok := g.points[row]; !ok {
		g.points[row] = make(map[int]int)
	}
	g.points[row][col] = val
	g.MaxCol = util.Max(g.MaxCol, col)
	g.MaxRow = util.Max(g.MaxRow, row)
	g.MinCol = util.Min(g.MinCol, col)
	g.MinRow = util.Min(g.MinRow, row)
}

func (g *Grid) Incr(row, col, inc int) {
	if _, ok := g.points[row]; !ok {
		g.points[row] = make(map[int]int)
	}
	g.points[row][col] = g.points[row][col] + inc
}

func (g Grid) Get(row, col int) int {
	if _, ok := g.points[row]; ok {
		if v, ok2 := g.points[row][col]; ok2 {
			return v
		}
	}
	return EMPTY
}

func (g *Grid) String() string {
	var b strings.Builder
	var val int
	for r := g.MinRow - 1; r <= g.MaxRow+1; r++ {
		for c := g.MinCol - 2; c <= g.MaxCol+2; c++ {
			val = g.Get(r, c)
			if val == ELF {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
