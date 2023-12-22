package aoc2023_day18

// Duplicate of day18.go, but with a different Fill() algorithm that doesn't work.
import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Coord = util.Coord

type Dig struct {
	Dir    string
	Meters int
	Color  string
}

type Grid struct {
	PerimSize int
	Vertices  []Coord
}

func NewGrid() Grid {
	return Grid{
		PerimSize: 0,
		Vertices:  []Coord{{Row: 0, Col: 0}},
	}
}

func part1(input []string) (string, error) {
	digs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", digs)
	grid := NewGrid()

	result := grid.GetArea(digs)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	digs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", digs)
	grid := NewGrid()

	for i, dig := range digs {
		digs[i] = ColorToMeters(dig)
	}
	result := grid.GetArea(digs)
	return fmt.Sprintf("%d", result), nil
}

func (g *Grid) GetArea(digs []Dig) int {
	cur := Coord{Row: 0, Col: 0}
	for _, dig := range digs {
		switch dig.Dir {
		case "R":
			cur.Col += dig.Meters
		case "L":
			cur.Col -= dig.Meters
		case "U":
			cur.Row -= dig.Meters
		case "D":
			cur.Row += dig.Meters
		}
		g.Vertices = append(g.Vertices, Coord{Row: cur.Row, Col: cur.Col})
	}

	A := util.ShoelaceArea(g.Vertices)
	b := g.PerimSize
	log.Debugf("Grid Area:\n%d\n", A)
	log.Debugf("PerimeterPoints:\n%d\n", b)

	i := util.PicksTheorem(A, g.PerimSize)
	log.Debugf("InteriorPoints:\n%d\n", i)

	return i + b
}

func ColorToMeters(d Dig) Dig {
	digits := strings.Replace(d.Color, "#", "", -1)
	hex := digits[0:5]
	// 0 means R, 1 means D, 2 means L, and 3 means U
	switch digits[5] {
	case '0':
		d.Dir = "R"
	case '1':
		d.Dir = "D"
	case '2':
		d.Dir = "L"
	case '3':
		d.Dir = "U"
	}
	n, _ := strconv.ParseInt(hex, 16, 64)
	d.Meters = int(n)
	log.Debugf("%s %d (%d)\n", d.Dir, n, d.Meters)
	return d
}

func parseInput(input []string) []Dig {
	digs := make([]Dig, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " ")
		digs = append(digs, Dig{
			Dir:    parts[0],
			Meters: util.Atoi(parts[1]),
			Color:  strings.ReplaceAll(strings.ReplaceAll(parts[2], "(", ""), ")", ""),
		})
	}
	return digs
}

func init() {
	challenges.RegisterChallengeFunc(2023, 18, 1, "day18.txt", part1)
	challenges.RegisterChallengeFunc(2023, 18, 2, "day18.txt", part2)
}
