package aoc2023_day16

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Grid struct {
	Cells     [][]string
	Height    int
	Width     int
	Visited   map[string]map[string]bool
	Beams     []Beam
	BeamCache map[string]bool
}

type Beam struct {
	Row       int
	Col       int
	Direction string
}

func (g Grid) String() string {
	s := ""
	for _, row := range g.Cells {
		s += strings.Join(row, "") + "\n"
	}
	return s
}

func (g Grid) Copy() Grid {
	newCells := make([][]string, 0, len(g.Cells))
	for _, row := range g.Cells {
		newCells = append(newCells, append([]string{}, row...))
	}
	return Grid{
		Cells:     newCells,
		Height:    g.Height,
		Width:     g.Width,
		Visited:   make(map[string]map[string]bool),
		BeamCache: make(map[string]bool),
	}
}

func (g Grid) CountVisited() int {
	var result int
	for range g.Visited {
		result++
	}
	return result
}

func (g Grid) Get(r, c int) string {
	if r < 0 || r >= g.Height || c < 0 || c >= g.Width {
		return ""
	}
	return g.Cells[r][c]
}

func (b Beam) String() string {
	return fmt.Sprintf("%d,%d %s", b.Row, b.Col, b.Direction)
}

func (g *Grid) BeamTravel() bool {
	newBeams := make([]Beam, 0)
	for _, b := range g.Beams {
		cell := g.Get(b.Row, b.Col)
		if cell == "" {
			continue // off the grid, beam is done
		}
		if c, ok := g.Visited[fmt.Sprintf("%d,%d", b.Row, b.Col)]; ok {
			if c[b.Direction] {
				continue // already visited this cell in this direction
			} else {
				c[b.Direction] = true
			}
		} else {
			g.Visited[fmt.Sprintf("%d,%d", b.Row, b.Col)] = make(map[string]bool)
			g.Visited[fmt.Sprintf("%d,%d", b.Row, b.Col)][b.Direction] = true
		}

		switch b.Direction {
		case "R":
			if cell == "." || cell == "-" {
				newBeams = append(newBeams, Beam{b.Row, b.Col + 1, b.Direction})
			} else if cell == "|" {
				newBeams = append(newBeams, Beam{b.Row - 1, b.Col, "U"})
				newBeams = append(newBeams, Beam{b.Row + 1, b.Col, "D"})
			} else if cell == "/" {
				newBeams = append(newBeams, Beam{b.Row - 1, b.Col, "U"})
			} else if cell == "\\" {
				newBeams = append(newBeams, Beam{b.Row + 1, b.Col, "D"})
			}
		case "L":
			if cell == "." || cell == "-" {
				newBeams = append(newBeams, Beam{b.Row, b.Col - 1, b.Direction})
			} else if cell == "|" {
				newBeams = append(newBeams, Beam{b.Row - 1, b.Col, "U"})
				newBeams = append(newBeams, Beam{b.Row + 1, b.Col, "D"})
			} else if cell == "/" {
				newBeams = append(newBeams, Beam{b.Row + 1, b.Col, "D"})
			} else if cell == "\\" {
				newBeams = append(newBeams, Beam{b.Row - 1, b.Col, "U"})
			}
		case "U":
			if cell == "." || cell == "|" {
				newBeams = append(newBeams, Beam{b.Row - 1, b.Col, b.Direction})
			} else if cell == "-" {
				newBeams = append(newBeams, Beam{b.Row, b.Col - 1, "L"})
				newBeams = append(newBeams, Beam{b.Row, b.Col + 1, "R"})
			} else if cell == "/" {
				newBeams = append(newBeams, Beam{b.Row, b.Col + 1, "R"})
			} else if cell == "\\" {
				newBeams = append(newBeams, Beam{b.Row, b.Col - 1, "L"})
			}
		case "D":
			if cell == "." || cell == "|" {
				newBeams = append(newBeams, Beam{b.Row + 1, b.Col, b.Direction})
			} else if cell == "-" {
				newBeams = append(newBeams, Beam{b.Row, b.Col - 1, "L"})
				newBeams = append(newBeams, Beam{b.Row, b.Col + 1, "R"})
			} else if cell == "/" {
				newBeams = append(newBeams, Beam{b.Row, b.Col - 1, "L"})
			} else if cell == "\\" {
				newBeams = append(newBeams, Beam{b.Row, b.Col + 1, "R"})
			}
		}
	}
	g.Beams = newBeams
	dedup := make(map[string]Beam)
	for _, b := range g.Beams {
		dedup[b.String()] = b
	}
	g.Beams = make([]Beam, 0)
	keys := make([]string, 0)
	for _, b := range dedup {
		// g.Beams = append(g.Beams, b)
		keys = append(keys, b.String())
	}
	sort.StringSlice(keys).Sort()
	for _, k := range keys {
		g.Beams = append(g.Beams, dedup[k])
	}
	cacheKey := strings.Join(keys, "|")
	if g.BeamCache[cacheKey] {
		return true
	}
	g.BeamCache[strings.Join(keys, "|")] = true
	return false
}

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	grid.Beams = append(grid.Beams, Beam{Row: 0, Col: 0, Direction: "R"})

	for len(grid.Beams) > 0 {
		if grid.BeamTravel() {
			break
		}
	}

	var result = grid.CountVisited()

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	max := 0
	for r := range grid.Cells {
		g2 := grid.Copy()
		g2.Beams = append(g2.Beams, Beam{Row: r, Col: 0, Direction: "R"})
		g3 := grid.Copy()
		g3.Beams = append(g3.Beams, Beam{Row: r, Col: grid.Width - 1, Direction: "L"})
		for len(g2.Beams) > 0 && !g2.BeamTravel() {
		}
		for len(g3.Beams) > 0 && !g3.BeamTravel() {
		}
		max = util.Max(max, g2.CountVisited(), g3.CountVisited())
	}

	for c := range grid.Cells[0] {
		g2 := grid.Copy()
		g2.Beams = append(g2.Beams, Beam{Row: 0, Col: c, Direction: "D"})
		g3 := grid.Copy()
		g3.Beams = append(g3.Beams, Beam{Row: grid.Height - 1, Col: c, Direction: "U"})
		for len(g2.Beams) > 0 && !g2.BeamTravel() {
		}
		for len(g3.Beams) > 0 && !g3.BeamTravel() {
		}
		max = util.Max(max, g2.CountVisited(), g3.CountVisited())
	}

	return fmt.Sprintf("%d", max), nil
}

func parseInput(input []string) Grid {
	cells := make([][]string, 0, len(input))
	for _, s := range input {
		cells = append(cells, strings.Split(s, ""))
	}
	return Grid{
		Cells:     cells,
		Height:    len(cells),
		Width:     len(cells[0]),
		Visited:   make(map[string]map[string]bool),
		BeamCache: make(map[string]bool),
	}
}

func init() {
	challenges.RegisterChallengeFunc(2023, 16, 1, "day16.txt", part1)
	challenges.RegisterChallengeFunc(2023, 16, 2, "day16.txt", part2)
}
