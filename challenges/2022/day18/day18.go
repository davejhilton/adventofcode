package aoc2022_day18

import (
	"fmt"
	"math"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
	"github.com/davejhilton/adventofcode/util/set"
)

func part1(input []string) (string, error) {
	coords := parseInput(input)

	openSides := 0
	for c := range coords {
		for n := range c.Neighbors() {
			if !coords.Has(n) {
				openSides++
			}
		}
	}

	return fmt.Sprintf("%d", openSides), nil
}

func part2(input []string) (string, error) {
	coords := parseInput(input)

	var minX, maxX int = math.MaxInt, 0
	var minY, maxY int = math.MaxInt, 0
	var minZ, maxZ int = math.MaxInt, 0
	for c := range coords {
		minX, maxX = util.Min(minX, c.X-1), util.Max(maxX, c.X+1)
		minY, maxY = util.Min(minY, c.Y-1), util.Max(maxY, c.Y+1)
		minZ, maxZ = util.Min(minZ, c.Z-1), util.Max(maxZ, c.Z+1)
	}

	start := Coord{X: maxX, Y: maxY, Z: maxZ}
	outside := set.New(start)
	for {
		newPoints := set.Set[Coord]{}
		for c := range outside {
			for n := range c.Neighbors() {
				if !coords.Has(n) && !outside.Has(n) {
					if n.X >= minX && n.X <= maxX &&
						n.Y >= minY && n.Y <= maxY &&
						n.Z >= minZ && n.Z <= maxZ {
						newPoints.Add(n)
					}
				}
			}
		}
		if newPoints.Size() > 0 {
			outside.AddAll(newPoints)
		} else {
			break
		}
	}

	outerFaces := 0
	for c := range coords {
		for n := range c.Neighbors() {
			if outside.Has(n) {
				outerFaces++
			}
		}
	}

	return fmt.Sprintf("%d", outerFaces), nil
}

type Coord struct {
	X int
	Y int
	Z int
}

func (c Coord) Neighbors() set.Set[Coord] {
	return set.FromSlice([]Coord{
		{c.X + 1, c.Y, c.Z},
		{c.X - 1, c.Y, c.Z},
		{c.X, c.Y + 1, c.Z},
		{c.X, c.Y - 1, c.Z},
		{c.X, c.Y, c.Z + 1},
		{c.X, c.Y, c.Z - 1},
	})
}

func parseInput(input []string) set.Set[Coord] {
	coords := set.Set[Coord]{}
	for _, s := range input {
		n := util.AtoiSplit(s, ",")
		coords.Add(Coord{n[0], n[1], n[2]})
	}
	return coords
}

func init() {
	challenges.RegisterChallengeFunc(2022, 18, 1, "day18.txt", part1)
	challenges.RegisterChallengeFunc(2022, 18, 2, "day18.txt", part2)
}
