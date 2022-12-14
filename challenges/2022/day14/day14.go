package aoc2022_day14

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
	ROCK  = 1
	SAND  = 2
)

func part1(input []string) (string, error) {
	paths, minPoint, maxPoint := parseInput(input)
	grid := NewGrid()

	for _, p := range paths {
		grid.AddPath(p)
	}

	log.Debugln(grid)

	nSand := 0
	done := false
	for !done {
		yPos := 0
		xPos := 500
		for {
			if yPos > maxPoint.Y || xPos < minPoint.X || xPos > maxPoint.X { // infinitely falling
				done = true
				break
			} else if grid.Get(xPos, yPos+1) == EMPTY { // move down
				yPos++
				continue
			} else if grid.Get(xPos-1, yPos+1) == EMPTY { // move down + left
				yPos++
				xPos--
			} else if grid.Get(xPos+1, yPos+1) == EMPTY { // move down + right
				yPos++
				xPos++
			} else { // hit rock/sand below, can't move down further
				grid.Set(xPos, yPos, SAND)
				nSand++
				break
			}
		}
	}

	log.Debugln(grid)

	return fmt.Sprintf("%d", nSand), nil
}

func part2(input []string) (string, error) {
	paths, _, maxPoint := parseInput(input)
	grid := NewGrid()

	for _, p := range paths {
		grid.AddPath(p)
	}

	log.Debugln(grid)

	nSand := 0
	done := false
	for !done {
		yPos := 0
		xPos := 500
		for {
			if grid.Get(500, 0) != EMPTY {
				done = true
				break
			}

			if yPos == maxPoint.Y+1 { // hit the "floor"
				grid.Set(xPos, yPos, SAND)
				nSand++
				break
			} else if grid.Get(xPos, yPos+1) == EMPTY { // move down
				yPos++
			} else if grid.Get(xPos-1, yPos+1) == EMPTY { // move down + left
				yPos++
				xPos--
			} else if grid.Get(xPos+1, yPos+1) == EMPTY { // move down + right
				yPos++
				xPos++
			} else { // hit rock/sand below, can't move down further
				grid.Set(xPos, yPos, SAND)
				nSand++
				break
			}
		}
	}

	log.Debugln(grid)

	return fmt.Sprintf("%d", nSand), nil
}

func parseInput(input []string) ([]Path, Coord, Coord) {
	paths := make([]Path, 0)
	var minX, minY int = math.MaxInt, math.MaxInt
	var maxX, maxY int = math.MinInt, math.MinInt
	for _, s := range input {
		path := Path{
			Coords: make([]Coord, 0),
		}

		strCoords := strings.Split(s, " -> ")

		for _, sc := range strCoords {
			xy := strings.Split(sc, ",")
			x, y := util.Atoi(xy[0]), util.Atoi(xy[1])
			minX = util.Min(minX, x)
			maxX = util.Max(maxX, x)
			minY = util.Min(minY, y)
			maxY = util.Max(maxY, y)
			path.Coords = append(path.Coords, Coord{x, y})
		}
		paths = append(paths, path)
	}
	return paths, Coord{minX, minY}, Coord{maxX, maxY}
}

type Coord struct {
	X int
	Y int
}

type Path struct {
	Coords []Coord
}

type Grid struct {
	points map[int]map[int]int
	MinX   int
	MinY   int
	MaxX   int
	MaxY   int
}

func NewGrid() *Grid {
	return &Grid{
		points: make(map[int]map[int]int),
		MinX:   math.MaxInt,
		MaxX:   math.MinInt,
		MinY:   math.MaxInt,
		MaxY:   math.MinInt,
	}
}

func (g *Grid) Set(x, y, val int) {
	if _, ok := g.points[y]; !ok {
		g.points[y] = make(map[int]int)
	}
	g.points[y][x] = val
	g.MaxX = util.Max(g.MaxX, x)
	g.MaxY = util.Max(g.MaxY, y)
	g.MinX = util.Min(g.MinX, x)
	g.MinY = util.Min(g.MinY, y)
}

func (g Grid) Get(x, y int) int {
	if _, ok := g.points[y]; ok {
		if v, ok2 := g.points[y][x]; ok2 {
			return v
		}
	}
	return EMPTY
}

func (g *Grid) AddPath(p Path) {
	var prev Coord
	for i, c := range p.Coords {
		if i != 0 {
			if c.X == prev.X {
				// vertical line
				x := c.X
				start, end := prev, c
				if c.Y < prev.Y {
					start, end = c, prev
				}
				for y := start.Y; y <= end.Y; y++ {
					g.Set(x, y, ROCK)
				}
			} else {
				// horizontal line
				y := c.Y
				start, end := prev, c
				if c.X < prev.X {
					start, end = c, prev
				}
				for x := start.X; x <= end.X; x++ {
					g.Set(x, y, ROCK)
				}
			}
		}
		prev = c
	}
}

func (g Grid) String() string {
	width := g.MaxX - g.MinX + 1
	header := make([][]string, 3)
	for y := range header {
		header[y] = make([]string, width+3)
	}

	for x := 3; x < width+3; x++ {
		xVal := fmt.Sprintf("%3d", x)
		header[0][x] = string(xVal[0])
		header[1][x] = string(xVal[1])
		header[2][x] = string(xVal[2])
	}

	var b strings.Builder
	for _, s := range header {
		b.WriteString(fmt.Sprintf("   %s\n", strings.Join(s, "")))
	}

	for y := g.MinY; y <= g.MaxY; y++ {
		b.WriteString(fmt.Sprintf("%3d", y))
		var row map[int]int
		if r, ok := g.points[y]; ok {
			row = r
		} else {
			row = make(map[int]int)
		}
		for x := g.MinX; x <= g.MaxX; x++ {
			if y == g.MinY && x == 500 && row[x] == EMPTY {
				b.WriteRune('+')
			} else {
				switch row[x] {
				case EMPTY:
					b.WriteRune('.')
				case ROCK:
					b.WriteRune('#')
				case SAND:
					b.WriteRune('o')
				default:
					b.WriteRune('?')
				}
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func init() {
	challenges.RegisterChallengeFunc(2022, 14, 1, "day14.txt", part1)
	challenges.RegisterChallengeFunc(2022, 14, 2, "day14.txt", part2)
}
