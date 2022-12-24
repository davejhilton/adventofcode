package aoc2022_day24

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util/set"
)

type Set = set.Set[Coord]

const (
	E     int = -2
	WALL  int = -1
	EMPTY int = 0

	UP    int = 1
	RIGHT int = 2
	DOWN  int = 4
	LEFT  int = 8
)

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	var t int
	var start = Coord{0, 1}
	var end = Coord{len(grid) - 1, len(grid[0]) - 2}
	TraverseGrid(grid, start, end, &t)

	return fmt.Sprintf("%d", t), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	var start = Coord{0, 1}
	var end = Coord{len(grid) - 1, len(grid[0]) - 2}

	var t int
	grid = TraverseGrid(grid, start, end, &t)
	grid = TraverseGrid(grid, end, start, &t)
	_ = TraverseGrid(grid, start, end, &t)
	return fmt.Sprintf("%d", t), nil
}

func parseInput(input []string) Grid {
	grid := NewGrid(len(input), len(input[0]))
	for r, s := range input {
		for c, ch := range s {
			grid[r][c] = INPUT_MAP[ch]
		}
	}
	return grid
}

var INPUT_MAP = map[rune]int{
	'#': WALL,
	'.': EMPTY,
	'>': RIGHT,
	'v': DOWN,
	'<': LEFT,
	'^': UP,
}

func TraverseGrid(grid Grid, start, end Coord, t *int) Grid {
	var prevCoords = Set{}
	prevCoords.Add(start)
	for {
		*(t) += 1
		grid = MoveAllBlizzards(grid)
		coords := Set{}
		for coord := range prevCoords {
			if coords.Has(end) {
				return grid
			}
			for _, c := range GetPossibleMoves(coord.Row, coord.Col) {
				if c.Row >= 0 && c.Row < len(grid) && c.Col >= 0 && c.Col < len(grid[0]) {
					if grid[c.Row][c.Col] == EMPTY {
						coords.Add(c)
					}
				}
			}
		}
		prevCoords = coords
		log.Debugf("After t=%d:\n%s\n%s\n", t, grid, DebugPositions(grid, coords))
	}
}

func MoveAllBlizzards(grid Grid) (newGrid Grid) {
	h, w := len(grid), len(grid[0])
	newGrid = NewGrid(h, w)
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			v := grid[r][c]
			if v == WALL {
				newGrid[r][c] = v
				continue
			} else if v == EMPTY {
				continue
			}
			if v&UP != 0 {
				if r == 1 {
					newGrid[h-2][c] |= UP
				} else {
					newGrid[r-1][c] |= UP
				}
			}
			if v&RIGHT != 0 {
				if c == w-2 {
					newGrid[r][1] |= RIGHT
				} else {
					newGrid[r][c+1] |= RIGHT
				}
			}
			if v&DOWN != 0 {
				if r == h-2 {
					newGrid[1][c] |= DOWN
				} else {
					newGrid[r+1][c] |= DOWN
				}
			}
			if v&LEFT != 0 {
				if c == 1 {
					newGrid[r][w-2] |= LEFT
				} else {
					newGrid[r][c-1] |= LEFT
				}
			}
		}
	}
	return newGrid
}

func GetPossibleMoves(r, c int) []Coord {
	return []Coord{
		{r - 1, c},
		{r, c - 1}, {r, c}, {r, c + 1},
		{r + 1, c},
	}
}

func init() {
	challenges.RegisterChallengeFunc(2022, 24, 1, "day24.txt", part1)
	challenges.RegisterChallengeFunc(2022, 24, 2, "day24.txt", part2)
}

type Coord struct {
	Row int
	Col int
}

type Grid [][]int

func NewGrid(h, w int) Grid {
	grid := make(Grid, 0, h)
	for i := 0; i < h; i++ {
		grid = append(grid, make([]int, w))
	}
	return grid
}

func DebugPositions(g Grid, coords map[Coord]struct{}) string {
	if log.DebugEnabled() {
		var grid Grid = NewGrid(len(g), len(g[0]))
		for r, row := range g {
			copy(grid[r], row)
		}
		for c := range coords {
			grid[c.Row][c.Col] = E
		}
		return grid.String()
	}
	return ""
}

func (g Grid) String() string {
	var b strings.Builder
	for r := 0; r < len(g); r++ {
		for c := 0; c < len(g[r]); c++ {
			v := g[r][c]
			if v == WALL || v == EMPTY || v == E {
				b.WriteString(STRINGS[v])
				continue
			}
			var s string
			var n int
			if v&UP != 0 {
				s = STRINGS[UP]
				n++
			}
			if v&RIGHT != 0 {
				s = STRINGS[RIGHT]
				n++
			}
			if v&DOWN != 0 {
				s = STRINGS[DOWN]
				n++
			}
			if v&LEFT != 0 {
				s = STRINGS[LEFT]
				n++
			}
			if n == 1 {
				b.WriteString(s)
			} else {
				b.WriteString(fmt.Sprintf("%d", n))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

var STRINGS = map[int]string{
	WALL:  "#",
	EMPTY: ".",
	UP:    "^",
	RIGHT: ">",
	DOWN:  "v",
	LEFT:  "<",
	E:     "E",
}
