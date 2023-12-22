package aoc2023_day22

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
	"github.com/davejhilton/adventofcode/util/set"
)

type Coord struct {
	X int
	Y int
	Z int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d,%d)", c.X, c.Y, c.Z)
}

type Brick struct {
	Id     int
	Coords []Coord
}

type Grid [][][]int

func (g Grid) String() string {
	var sb strings.Builder
	for z := len(g) - 1; z >= 0; z-- {
		for y := 0; y < len(g[z]); y++ {
			for x := 0; x < len(g[z][y]); x++ {
				sb.WriteString(fmt.Sprintf("%d", g[z][y][x]))
			}
			if y < len(g[z])-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func NewGrid(maxX, maxY, maxZ int) Grid {
	g := make(Grid, 0, maxZ+1)
	for z := 0; z <= maxZ; z++ {
		g = append(g, make([][]int, 0, maxY+1))
		for y := 0; y <= maxY; y++ {
			g[z] = append(g[z], make([]int, maxX+1))
		}
	}
	return g
}

// for each brick in the grid, move it downward until it hits another brick
// or the bottom of the grid (lowest it can go is z=1)
// return a set of brick ids that moved
func Settle(grid Grid, bricks []Brick) set.Set[int] {
	settledBricks := make(set.Set[int])
	for _, brick := range bricks {
		canMoveDown := true
		for _, c := range brick.Coords {
			if c.Z == 1 || (grid[c.Z-1][c.Y][c.X] != brick.Id && grid[c.Z-1][c.Y][c.X] != 0) {
				canMoveDown = false
				break
			}
		}
		if canMoveDown {
			settledBricks.Add(brick.Id)
			for i, c := range brick.Coords {
				grid[c.Z][c.Y][c.X] = 0
				brick.Coords[i].Z = c.Z - 1
			}
			for _, c := range brick.Coords {
				grid[c.Z][c.Y][c.X] = brick.Id
			}
		}
	}
	return settledBricks
}

func FullySettle(grid Grid, bricks []Brick) set.Set[int] {
	settling := true
	allMovedBricks := make(set.Set[int])
	for settling {
		settling = false
		movedBricks := Settle(grid, bricks)
		if len(movedBricks) > 0 {
			settling = true
			for k := range movedBricks {
				allMovedBricks.Add(k)
			}
		}
	}
	return allMovedBricks
}

func (g Grid) Clone() Grid {
	clone := make(Grid, len(g))
	for z := 0; z < len(g); z++ {
		clone[z] = make([][]int, len(g[z]))
		for y := 0; y < len(g[z]); y++ {
			clone[z][y] = make([]int, len(g[z][y]))
			copy(clone[z][y], g[z][y])
		}
	}
	return clone
}

func CloneBricks(bricks []Brick) []Brick {
	clone := make([]Brick, len(bricks))
	for i, brick := range bricks {
		clone[i] = Brick{
			Id:     brick.Id,
			Coords: make([]Coord, len(brick.Coords)),
		}
		copy(clone[i].Coords, brick.Coords)
	}
	return clone
}

func (g Grid) DisintigrateBrick(brick Brick) {
	for _, c := range brick.Coords {
		g[c.Z][c.Y][c.X] = 0
	}
}

func part1(input []string) (string, error) {
	bricks, max := parseInput(input)

	grid := NewGrid(max.X, max.Y, max.Z)
	for _, brick := range bricks {
		for _, c := range brick.Coords {
			grid[c.Z][c.Y][c.X] = brick.Id
		}
	}

	log.Debugln(grid)
	log.Debugln()

	FullySettle(grid, bricks)

	log.Debugln(grid)
	log.Debugln()

	canDisintigrate := make(set.Set[int])
	grid1 := grid.Clone()
	bricks1 := CloneBricks(bricks)
	for i, brick := range bricks1 {
		grid1.DisintigrateBrick(brick)
		otherBricks := append(CloneBricks(bricks1[:i]), CloneBricks(bricks1[i+1:])...)
		moved := Settle(grid1, otherBricks).Size() > 0
		log.Debugf("Brick %d moved?: %v\n", brick.Id, moved)
		if !moved {
			canDisintigrate.Add(brick.Id)
		}
		grid1 = grid.Clone()
		bricks1 = CloneBricks(bricks)
	}

	result := canDisintigrate.Size()
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	bricks, max := parseInput(input)

	grid := NewGrid(max.X, max.Y, max.Z)
	for _, brick := range bricks {
		for _, c := range brick.Coords {
			grid[c.Z][c.Y][c.X] = brick.Id
		}
	}

	log.Debugln(grid)
	FullySettle(grid, bricks)
	log.Debugln()

	grid1 := grid.Clone()
	bricks1 := CloneBricks(bricks)
	nMoved := 0
	for i, brick := range bricks1 {
		grid1.DisintigrateBrick(brick)
		otherBricks := append(CloneBricks(bricks1[:i]), CloneBricks(bricks1[i+1:])...)
		moved := FullySettle(grid1, otherBricks)
		if moved.Size() > 0 {
			nMoved += moved.Size()
		}
		grid1 = grid.Clone()
		bricks1 = CloneBricks(bricks)
	}

	result := nMoved
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) ([]Brick, Coord) {
	bricks := make([]Brick, 0, len(input))
	var maxX, maxY, maxZ int
	for i, s := range input {
		strs := strings.Split(s, "~")
		c1 := util.AtoiSplit(strs[0], ",")
		c2 := util.AtoiSplit(strs[1], ",")

		coords := make([]Coord, 0)
		for z := c1[2]; z <= c2[2]; z++ {
			for y := c1[1]; y <= c2[1]; y++ {
				for x := c1[0]; x <= c2[0]; x++ {
					maxX = util.Max(maxX, x)
					maxY = util.Max(maxY, y)
					maxZ = util.Max(maxZ, z)
					coords = append(coords, Coord{X: x, Y: y, Z: z})
				}
			}
		}

		bricks = append(bricks, Brick{
			Id:     i + 1,
			Coords: coords,
		})
	}
	return bricks, Coord{X: maxX, Y: maxY, Z: maxZ}
}

func init() {
	challenges.RegisterChallengeFunc(2023, 22, 1, "day22.txt", part1)
	challenges.RegisterChallengeFunc(2023, 22, 2, "day22.txt", part2)
}
