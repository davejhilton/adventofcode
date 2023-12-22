package aoc2023_day21

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Coord = util.Coord
type Grid [][]string

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		sb.WriteString(fmt.Sprintf("%s\n", strings.Join(row, "")))
	}
	return sb.String()
}

func NewEmptyGrid(rows, cols int) Grid {
	grid := make(Grid, 0, rows)
	for r := 0; r < rows; r++ {
		row := make([]string, 0, cols)
		for c := 0; c < cols; c++ {
			row = append(row, ".")
		}
		grid = append(grid, row)
	}
	return grid
}

func AvailableNeighbors(g Grid, c Coord) []Coord {
	neighbors := make([]Coord, 0)
	// case "U":
	if c.Row > 0 && g[c.Row-1][c.Col] != "#" {
		neighbors = append(neighbors, Coord{Row: c.Row - 1, Col: c.Col})
	}
	// case "D":
	if c.Row < len(g)-1 && g[c.Row+1][c.Col] != "#" {
		neighbors = append(neighbors, Coord{Row: c.Row + 1, Col: c.Col})
	}
	// case "L":
	if c.Col > 0 && g[c.Row][c.Col-1] != "#" {
		neighbors = append(neighbors, Coord{Row: c.Row, Col: c.Col - 1})
	}
	// case "R":
	if c.Col < len(g[c.Row])-1 && g[c.Row][c.Col+1] != "#" {
		neighbors = append(neighbors, Coord{Row: c.Row, Col: c.Col + 1})
	}
	return neighbors
}

func (g Grid) StartingPoint() Coord {
	return Coord{Row: len(g) / 2, Col: len(g[0]) / 2}
}

func cacheKey(c Coord, moves int) string {
	return fmt.Sprintf("%s:%d", c.String(), moves)
}

var seen map[string]bool
var resultCache map[string]int
var resultGrid Grid

func Move(g Grid, c Coord, moves int) {
	seen[cacheKey(c, moves)] = true
	if moves == 0 {
		resultCache[cacheKey(c, moves)] = 1
		resultGrid[c.Row][c.Col] = "O"
		return
	}
	for _, n := range AvailableNeighbors(g, c) {
		if !seen[cacheKey(n, moves-1)] {
			Move(g, n, moves-1)
		}
	}
}

func polynomial(x int, yVals []int) int {
	a := yVals[0]
	a1 := yVals[1]
	a2 := yVals[2]

	b := a1 - a
	b1 := a2 - a1

	c := b1 - b

	return a + b*x + c*x*(x-1)/2
}

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	resultCache = make(map[string]int)
	resultGrid = NewEmptyGrid(len(grid), len(grid[0]))
	seen = make(map[string]bool)
	start := grid.StartingPoint()
	Move(grid, start, 64)

	var result int = len(resultCache)

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)

	STEPS := 26501365
	n := 2

	yValues := make([]int, 0, n+1)
	var result int
	for i := 0; i <= n; i++ {
		size := 2*i + 1
		resultGrid = NewEmptyGrid(len(grid)*size, len(grid[0])*size)
		superGrid := NewEmptyGrid(len(grid)*size, len(grid[0])*size)

		for r := 0; r < len(grid)*size; r++ {
			for c := 0; c < len(grid[0])*size; c++ {
				superGrid[r][c] = grid[r%len(grid)][c%len(grid[0])]
			}
		}
		resultCache = make(map[string]int)
		seen = make(map[string]bool)
		start := superGrid.StartingPoint()
		log.Debugf("Starting at %s\n", start)
		Move(superGrid, start, (len(grid)*i)+65)

		result = len(resultCache)
		yValues = append(yValues, result)
		log.Debugf("Result for %d: %d\n", i, result)
		log.Debugf("GRID:\n%s\n", resultGrid)
	}

	center := int(len(grid) / 2)
	x := (STEPS - center) / len(grid)
	result = polynomial(x, yValues)

	// fmt.Println(resultGrid)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Grid {
	grid := make(Grid, 0, len(input))
	for _, s := range input {
		grid = append(grid, strings.Split(s, ""))
	}
	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2023, 21, 1, "day21.txt", part1)
	challenges.RegisterChallengeFunc(2023, 21, 2, "day21.txt", part2)
}
