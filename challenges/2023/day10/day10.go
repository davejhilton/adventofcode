package aoc2023_day10

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type Grid [][]string

type Coord struct {
	Row int
	Col int
}

type PseudoCoord struct {
	Row float32
	Col float32
}

func (pc PseudoCoord) String() string {
	return fmt.Sprintf("(%f,%f)", pc.Row, pc.Col)
}

func (g Grid) String() string {
	var s strings.Builder
	for _, row := range g {
		for _, c := range row {
			s.WriteString(c)
		}
		s.WriteString("\n")
	}
	return s.String()
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

func (c Coord) Equals(other Coord) bool {
	return c.Row == other.Row && c.Col == other.Col
}

func (g Grid) StringAt(c Coord) string {
	var s strings.Builder
	for i, row := range g {
		for j, col := range row {
			if c.Row == i && c.Col == j {
				s.WriteString(log.Colorize(col, log.Yellow, 0))
			} else {
				s.WriteString(col)
			}
		}
		s.WriteString("\n")
	}
	return s.String()
}

func (g Grid) StringAtPseudo(pc PseudoCoord) string {
	var s strings.Builder
	for i, row := range g {
		for j, col := range row {
			if (int(pc.Row+0.5) == i || int(pc.Row-0.5) == i) && (int(pc.Col+0.5) == j || int(pc.Col-0.5) == j) {
				s.WriteString(log.Colorize(col, log.Red, 0))
			} else {
				s.WriteString(col)
			}
		}
		s.WriteString("\n")
	}
	return s.String()
}

func (g Grid) Get(Row, Col int) string {
	if Row < 0 || Row >= len(g) {
		return EMPTY
	}
	if Col < 0 || Col >= len(g[Row]) {
		return EMPTY
	}
	return g[Row][Col]
}

func (g Grid) GetNeighbor(c Coord, dir string) string {
	switch dir {
	case DIR_UP:
		return g.Get(c.Row-1, c.Col)
	case DIR_DOWN:
		return g.Get(c.Row+1, c.Col)
	case DIR_LEFT:
		return g.Get(c.Row, c.Col-1)
	case DIR_RIGHT:
		return g.Get(c.Row, c.Col+1)
	}
	return PIPE_EMPTY
}

func (g Grid) GetNeighborCoord(c Coord, dir string) Coord {
	switch dir {
	case DIR_UP:
		return Coord{c.Row - 1, c.Col}
	case DIR_DOWN:
		return Coord{c.Row + 1, c.Col}
	case DIR_LEFT:
		return Coord{c.Row, c.Col - 1}
	case DIR_RIGHT:
		return Coord{c.Row, c.Col + 1}
	}
	return Coord{-1, -1}
}

func (g *Grid) Set(Row, Col int, val string) {
	if Row < 0 || Row >= len(*g) {
		return
	}
	if Col < 0 || Col >= len((*g)[Row]) {
		return
	}
	(*g)[Row][Col] = val
}

func (g Grid) GetCoord(c Coord) string {
	if c.Row < 0 || c.Row >= len(g) {
		return EMPTY
	}
	if c.Col < 0 || c.Col >= len(g[c.Row]) {
		return EMPTY
	}
	return g[c.Row][c.Col]
}

func NewEmptyGrid(other Grid) Grid {
	grid := make(Grid, len(other))
	for i := range grid {
		grid[i] = make([]string, len(other[i]))
		for j := range grid[i] {
			grid[i][j] = PIPE_EMPTY
		}
	}
	return grid
}

func (g Grid) GetStartingPosition() Coord {
	for r, row := range g {
		for c, v := range row {
			if v == STARTING_POSITION {
				return Coord{r, c}
			}
		}
	}
	return Coord{-1, -1}
}

const (
	STARTING_POSITION = "S"
	PIPE_HORIZONTAL   = "-"
	PIPE_VERTICAL     = "|"
	PIPE_TOP_LEFT     = "F"
	PIPE_TOP_RIGHT    = "7"
	PIPE_BOTTOM_LEFT  = "L"
	PIPE_BOTTOM_RIGHT = "J"
	PIPE_EMPTY        = "."
	EMPTY             = ""

	OUTSIDE = "O"

	DIR_UP    = "UP"
	DIR_DOWN  = "DOWN"
	DIR_LEFT  = "LEFT"
	DIR_RIGHT = "RIGHT"
)

func (g Grid) Connects(c Coord, dir string) bool {
	var this string = g.GetCoord(c)
	var other = g.GetNeighbor(c, dir)
	if this == PIPE_EMPTY || this == EMPTY || other == PIPE_EMPTY || other == EMPTY {
		return false
	}
	switch dir {
	case DIR_UP:
		if this == PIPE_HORIZONTAL || this == PIPE_TOP_LEFT || this == PIPE_TOP_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false // can't go up from these
		}
		if other == PIPE_HORIZONTAL || other == PIPE_BOTTOM_LEFT || other == PIPE_BOTTOM_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false
		}
	case DIR_DOWN:
		if this == PIPE_HORIZONTAL || this == PIPE_BOTTOM_LEFT || this == PIPE_BOTTOM_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false // can't go down from these
		}
		if other == PIPE_HORIZONTAL || other == PIPE_TOP_LEFT || other == PIPE_TOP_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false
		}
	case DIR_LEFT:
		if this == PIPE_VERTICAL || this == PIPE_TOP_LEFT || this == PIPE_BOTTOM_LEFT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false // can't go right from these
		}
		if other == PIPE_VERTICAL || other == PIPE_TOP_RIGHT || other == PIPE_BOTTOM_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false
		}
	case DIR_RIGHT:
		if this == PIPE_VERTICAL || this == PIPE_TOP_RIGHT || this == PIPE_BOTTOM_RIGHT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false // can't go left from these
		}
		if other == PIPE_VERTICAL || other == PIPE_TOP_LEFT || other == PIPE_BOTTOM_LEFT {
			log.DebugColor(fmt.Sprintf("Can't go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Red)
			return false
		}
	}
	log.DebugColor(fmt.Sprintf("Can go %s from %v. This: %s, Other: %s", dir, c, this, other), log.Green)
	return true
}

var OppositeDir = map[string]string{
	DIR_UP:    DIR_DOWN,
	DIR_DOWN:  DIR_UP,
	DIR_LEFT:  DIR_RIGHT,
	DIR_RIGHT: DIR_LEFT,
}

func (g *Grid) FindLoop(orig Grid, c Coord, sourceDir string) bool {
	dirs := make([]string, 0, 3)
	for _, d := range []string{DIR_UP, DIR_DOWN, DIR_LEFT, DIR_RIGHT} {
		if d != sourceDir {
			dirs = append(dirs, d)
		}
	}

	printed := false
	isInLoop := false
	log.DebugColor(fmt.Sprintf("Checking %v", c), log.Yellow)
	for _, dir := range dirs {
		if !orig.Connects(c, dir) {
			continue
		}
		if !printed {
			g.Set(c.Row, c.Col, orig.Get(c.Row, c.Col)) // set it as part of the path (for now)
			log.Debugln(orig.StringAt(c))
			printed = true
		}
		if g.GetNeighbor(c, dir) == STARTING_POSITION {
			return true
		}
		if g.FindLoop(orig, g.GetNeighborCoord(c, dir), OppositeDir[dir]) {
			isInLoop = true
			break
		}
	}
	if !isInLoop && g.GetCoord(c) != STARTING_POSITION {
		g.Set(c.Row, c.Col, PIPE_EMPTY) // set it back to empty
	}
	return isInLoop
}

func (g *Grid) FindFurthest(c Coord) int {
	var dir string
	for _, d := range []string{DIR_UP, DIR_DOWN, DIR_LEFT, DIR_RIGHT} {
		if g.Connects(c, d) {
			dir = d
			break
		}
	}

	pos := g.GetNeighborCoord(c, dir)
	n := 0
	for !pos.Equals(c) {
		n++
		for _, d := range []string{DIR_UP, DIR_DOWN, DIR_LEFT, DIR_RIGHT} {
			if d != OppositeDir[dir] && g.Connects(pos, d) {
				dir = d
				break
			}
		}
		pos = g.GetNeighborCoord(pos, dir)
	}
	return n/2 + 1
}

var seen = make(map[string]bool)

func (g *Grid) FloodTravel(pc PseudoCoord, dir string) (PseudoCoord, bool) {
	switch dir {
	case DIR_DOWN:
		// if seen[fmt.Sprintf("(%f,%f)%s", pc.Row+1, pc.Col, DIR_UP)] {
		// 	return PseudoCoord{}, false
		// }
		if int(pc.Row+0.5) >= len(*g) {
			return PseudoCoord{}, false
		}
		n1 := g.Get(int(pc.Row+0.5), int(pc.Col-0.5))
		if n1 == PIPE_HORIZONTAL || n1 == PIPE_TOP_LEFT || n1 == PIPE_BOTTOM_LEFT {
			return PseudoCoord{}, false
		}
		if n1 == STARTING_POSITION {
			n2 := g.Get(int(pc.Row+0.5), int(pc.Col+0.5))
			if n2 == PIPE_HORIZONTAL || n2 == PIPE_TOP_RIGHT || n2 == PIPE_BOTTOM_RIGHT || n2 == EMPTY {
				return PseudoCoord{}, false
			}
		}
		return PseudoCoord{pc.Row + 1, pc.Col}, true
	case DIR_UP:
		// if seen[fmt.Sprintf("(%f,%f)%s", pc.Row-1, pc.Col, DIR_DOWN)] {
		// 	return PseudoCoord{}, false
		// }
		if int(pc.Row-0.5) < 0 {
			return PseudoCoord{}, false
		}
		n1 := g.Get(int(pc.Row-0.5), int(pc.Col-0.5))
		if n1 == PIPE_HORIZONTAL || n1 == PIPE_BOTTOM_LEFT || n1 == PIPE_TOP_LEFT {
			return PseudoCoord{}, false
		}
		if n1 == STARTING_POSITION {
			n2 := g.Get(int(pc.Row-0.5), int(pc.Col+0.5))
			if n2 == PIPE_HORIZONTAL || n2 == PIPE_BOTTOM_RIGHT || n2 == PIPE_TOP_RIGHT || n2 == EMPTY {
				return PseudoCoord{}, false
			}
		}
		return PseudoCoord{pc.Row - 1, pc.Col}, true
	case DIR_RIGHT:
		// if seen[fmt.Sprintf("(%f,%f)%s", pc.Row, pc.Col+1, DIR_LEFT)] {
		// 	return PseudoCoord{}, false
		// }
		if int(pc.Col+0.5) >= len((*g)[0]) {
			return PseudoCoord{}, false
		}
		n1 := g.Get(int(pc.Row-0.5), int(pc.Col+0.5))
		if n1 == PIPE_VERTICAL || n1 == PIPE_TOP_LEFT || n1 == PIPE_TOP_RIGHT {
			return PseudoCoord{}, false
		}
		if n1 == STARTING_POSITION {
			n2 := g.Get(int(pc.Row+0.5), int(pc.Col+0.5))
			if n2 == PIPE_VERTICAL || n2 == PIPE_BOTTOM_LEFT || n2 == PIPE_BOTTOM_RIGHT || n2 == EMPTY {
				return PseudoCoord{}, false
			}
		}
		return PseudoCoord{pc.Row, pc.Col + 1}, true
	case DIR_LEFT:
		// if seen[fmt.Sprintf("(%f,%f)%s", pc.Row, pc.Col-1, DIR_RIGHT)] {
		// 	return PseudoCoord{}, false
		// }
		if int(pc.Col-0.5) < 0 {
			return PseudoCoord{}, false
		}
		n1 := g.Get(int(pc.Row-0.5), int(pc.Col-0.5))
		if n1 == PIPE_VERTICAL || n1 == PIPE_TOP_RIGHT || n1 == PIPE_TOP_LEFT {
			return PseudoCoord{}, false
		}
		if n1 == STARTING_POSITION {
			n2 := g.Get(int(pc.Row+0.5), int(pc.Col-0.5))
			if n2 == PIPE_VERTICAL || n2 == PIPE_BOTTOM_RIGHT || n2 == PIPE_BOTTOM_LEFT || n2 == EMPTY {
				return PseudoCoord{}, false
			}
		}
		return PseudoCoord{pc.Row, pc.Col - 1}, true
	default:
		return PseudoCoord{}, false
	}
}

func (g *Grid) FloodOutsideLoop(pc PseudoCoord, sourceDir string) {
	if seen[fmt.Sprintf("%s%s", pc.String(), sourceDir)] {
		return
	}
	log.Debugln(g.StringAtPseudo(pc))
	seen[fmt.Sprintf("%s%s", pc.String(), sourceDir)] = true
	neighborsToCheck := make([]Coord, 0)

	neighborsToCheck = append(neighborsToCheck,
		Coord{int(pc.Row + 0.5), int(pc.Col - 0.5)},
		Coord{int(pc.Row + 0.5), int(pc.Col + 0.5)},
		Coord{int(pc.Row - 0.5), int(pc.Col - 0.5)},
		Coord{int(pc.Row - 0.5), int(pc.Col + 0.5)},
	)

	// switch sourceDir {
	// case DIR_UP:
	// 	neighborsToCheck = append(neighborsToCheck, Coord{int(pc.Row + 0.5), int(pc.Col - 0.5)}, Coord{int(pc.Row + 0.5), int(pc.Col + 0.5)})
	// case DIR_DOWN:
	// 	neighborsToCheck = append(neighborsToCheck, Coord{int(pc.Row - 0.5), int(pc.Col - 0.5)}, Coord{int(pc.Row - 0.5), int(pc.Col + 0.5)})
	// case DIR_LEFT:
	// 	neighborsToCheck = append(neighborsToCheck, Coord{int(pc.Row - 0.5), int(pc.Col + 0.5)}, Coord{int(pc.Row + 0.5), int(pc.Col + 0.5)})
	// case DIR_RIGHT:
	// 	neighborsToCheck = append(neighborsToCheck, Coord{int(pc.Row - 0.5), int(pc.Col - 0.5)}, Coord{int(pc.Row + 0.5), int(pc.Col - 0.5)})
	// default:
	// 	neighborsToCheck = append(neighborsToCheck, Coord{int(pc.Row + 0.5), int(pc.Col + 0.5)})
	// }

	for _, c := range neighborsToCheck {
		v := g.GetCoord(c)
		if v == PIPE_EMPTY {
			g.Set(c.Row, c.Col, OUTSIDE) // nice
		}
	}

	nextHops := make(map[string]PseudoCoord)
	if sourceDir != DIR_DOWN && pc.Row < float32(len(*g)) {
		if newPC, ok := g.FloodTravel(pc, DIR_DOWN); ok {
			log.DebugColor(fmt.Sprintf("Can go %s ", DIR_DOWN), log.Green)
			nextHops[DIR_UP] = newPC
		} else {
			log.DebugColor(fmt.Sprintf("Can't go %s ", DIR_DOWN), log.Red)
		}
	}
	if sourceDir != DIR_UP && pc.Row > 0.0 {
		if newPC, ok := g.FloodTravel(pc, DIR_UP); ok {
			log.DebugColor(fmt.Sprintf("Can go %s ", DIR_UP), log.Green)
			nextHops[DIR_DOWN] = newPC
		} else {
			log.DebugColor(fmt.Sprintf("Can't go %s ", DIR_UP), log.Red)
		}
	}
	if sourceDir != DIR_LEFT && pc.Col > 0.0 {
		if newPC, ok := g.FloodTravel(pc, DIR_LEFT); ok {
			log.DebugColor(fmt.Sprintf("Can go %s ", DIR_LEFT), log.Green)
			nextHops[DIR_RIGHT] = newPC
		} else {
			log.DebugColor(fmt.Sprintf("Can't go %s ", DIR_LEFT), log.Red)
		}
	}
	if sourceDir != DIR_RIGHT && pc.Col < float32(len((*g)[0])) {
		if newPC, ok := g.FloodTravel(pc, DIR_RIGHT); ok {
			log.DebugColor(fmt.Sprintf("Can go %s ", DIR_RIGHT), log.Green)
			nextHops[DIR_LEFT] = newPC
		} else {
			log.DebugColor(fmt.Sprintf("Can't go %s ", DIR_RIGHT), log.Red)
		}
	}

	// for _, c := range neighborsToCheck {
	// 	v := g.GetCoord(c)
	// 	if v == PIPE_EMPTY {
	// 		g.Set(c.Row, c.Col, OUTSIDE) // nice
	// 	}
	// }

	for nDir, newPC := range nextHops {
		g.FloodOutsideLoop(newPC, nDir)
	}

	// printed := false
	// // log.DebugColor(fmt.Sprintf("Checking %v", c), log.Yellow)
	// for _, dir := range dirs {

	// 	if !printed {
	// 		g.Set(c.Row, c.Col, orig.Get(c.Row, c.Col)) // set it as part of the path (for now)
	// 		log.Debugln(orig.StringAt(c))
	// 		printed = true
	// 	}
	// 	if g.GetNeighbor(c, dir) == STARTING_POSITION {
	// 		return true
	// 	}
	// 	if g.FindLoop(orig, g.GetNeighborCoord(c, dir), OppositeDir[dir]) {
	// 		isInLoop = true
	// 		break
	// 	}
	// }

}

func part1(input []string) (string, error) {
	g := parseInput(input)
	log.Debugf("Parsed Input:\n%s\n", g)

	pos := g.GetStartingPosition()
	g2 := NewEmptyGrid(g)
	g2.Set(pos.Row, pos.Col, STARTING_POSITION)
	g2.FindLoop(g, pos, DIR_UP)
	result := g2.FindFurthest(pos)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	g := parseInput(input)
	log.Debugf("Parsed Input:\n%s\n", g)
	g2 := NewEmptyGrid(g)
	pos := g.GetStartingPosition()
	g2.Set(pos.Row, pos.Col, STARTING_POSITION)
	g2.FindLoop(g, pos, "")

	pc := PseudoCoord{float32(-0.5), float32(-0.5)}
	g2.FloodOutsideLoop(pc, DIR_UP)
	log.Debugln(g2.String())
	var result int
	for _, row := range g2 {
		for _, col := range row {
			if col == PIPE_EMPTY {
				result++
			}
		}
	}
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
	challenges.RegisterChallengeFunc(2023, 10, 1, "day10.txt", part1)
	challenges.RegisterChallengeFunc(2023, 10, 2, "day10.txt", part2)
}
