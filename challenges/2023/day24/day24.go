package aoc2023_day24

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Coord struct {
	X float64
	Y float64
	Z float64
}

func (c Coord) String() string {
	return fmt.Sprintf("(%.0f,%.0f,%.0f)", c.X, c.Y, c.Z)
}

type Velocity struct {
	VX float64
	VY float64
	VZ float64
}

func (v Velocity) String() string {
	return fmt.Sprintf("(%.0f,%.0f,%.0f)", v.VX, v.VY, v.VZ)
}

type Hailstone struct {
	Coord
	Velocity
}

func (h Hailstone) String() string {
	return fmt.Sprintf("{%s, %s}", h.Coord.String(), h.Velocity.String())
}

func (h Hailstone) SlopeXY() float64 {
	if h.VX == 0 {
		return math.Inf(1)
	}
	return float64(h.VY) / float64(h.VX)
}

func (h Hailstone) IntersectXY(h2 Hailstone) *Coord {
	if h.SlopeXY() == h2.SlopeXY() {
		return nil
	}
	if math.IsInf(h.SlopeXY(), 1) {
		return &Coord{h.X, h2.Y, 0}
	}
	if math.IsInf(h2.SlopeXY(), 1) {
		return &Coord{h2.X, h.Y, 0}
	}
	if h.VX == 0 {
		return &Coord{h.X, h2.Y, 0}
	}
	if h2.VX == 0 {
		return &Coord{h2.X, h.Y, 0}
	}
	x := (h.Y - h2.Y - h.X*h.SlopeXY() + h2.X*h2.SlopeXY()) / (h2.SlopeXY() - h.SlopeXY())
	y := h.Y + h.SlopeXY()*(x-h.X)

	// check if the collision happens in the future
	if ((x-h.X) < 0 && h.VX > 0) || ((x-h.X) > 0 && h.VX < 0) {
		return nil
	}
	if ((x-h2.X) < 0 && h2.VX > 0) || ((x-h2.X) > 0 && h2.VX < 0) {
		return nil
	}
	return &Coord{x, y, 0}
}

func part1(input []string) (string, error) {
	stones := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", stones)

	checkRange := []float64{200000000000000, 400000000000000}
	if strings.Contains(challenges.CurrentChallenge.InputFileName, "example") {
		checkRange = []float64{7, 27}
	}
	var result int
	for i := 0; i < len(stones); i++ {
		for j := i + 1; j < len(stones); j++ {
			intersection := stones[i].IntersectXY(stones[j])
			if intersection != nil {
				log.Debugf("Hailstones %s and %s : Intersection: (%.2f, %.2f)\n", stones[i], stones[j], intersection.X, intersection.Y)
				if intersection.X >= checkRange[0] && intersection.X <= checkRange[1] && intersection.Y >= checkRange[0] && intersection.Y <= checkRange[1] {
					result++
				}
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

// Part2's solution is pretty heavily based on rumkugel13's code on github:
// https://github.com/rumkugel13/AdventOfCode2023/blob/main/day24.go
func part2(input []string) (string, error) {
	stones := parseInput(input)

	maybeX, maybeY, maybeZ := []int{}, []int{}, []int{}
	for i := 0; i < len(stones)-1; i++ {
		for j := i + 1; j < len(stones); j++ {
			a := stones[i]
			b := stones[j]
			if a.VX == b.VX {
				nextMaybe := findMatchingVelocities(int(b.X-a.X), int(a.VX))
				if len(maybeX) == 0 {
					maybeX = nextMaybe
				} else {
					maybeX = getIntersection(maybeX, nextMaybe)
				}
			}
			if a.VY == b.VY {
				nextMaybe := findMatchingVelocities(int(b.Y-a.Y), int(a.VY))
				if len(maybeY) == 0 {
					maybeY = nextMaybe
				} else {
					maybeY = getIntersection(maybeY, nextMaybe)
				}
			}
			if a.VZ == b.VZ {
				nextMaybe := findMatchingVelocities(int(b.Z-a.Z), int(a.VZ))
				if len(maybeZ) == 0 {
					maybeZ = nextMaybe
				} else {
					maybeZ = getIntersection(maybeZ, nextMaybe)
				}
			}
		}
	}

	var result int
	if len(maybeX) == 1 && len(maybeY) == 1 && len(maybeZ) == 1 {
		rockVelocity := Velocity{float64(maybeX[0]), float64(maybeY[0]), float64(maybeZ[0])}
		a := stones[0]
		b := stones[1]
		mA := (a.VY - rockVelocity.VY) / (a.VX - rockVelocity.VX)
		mB := (b.VY - rockVelocity.VY) / (b.VX - rockVelocity.VX)
		cA := a.Y - mA*a.X
		cB := b.Y - mB*b.X
		xPos := (cB - cA) / (mA - mB)
		yPos := mA*xPos + cA
		time := (xPos - a.X) / (a.VX - rockVelocity.VX)
		zPos := a.Z + time*(a.VZ-rockVelocity.VZ)
		result = int(xPos + yPos + zPos)
	}

	return fmt.Sprintf("%d", result), nil
}

func findMatchingVelocities(dvel int, pv int) []int {
	match := []int{}
	for v := -1000; v < 1000; v++ {
		if v != pv && dvel%(v-pv) == 0 {
			match = append(match, v)
		}
	}
	return match
}

func getIntersection(a []int, b []int) []int {
	intersection := []int{}
	for _, x := range a {
		if slices.Contains(b, x) {
			intersection = append(intersection, x)
		}
	}
	return intersection
}

func parseInput(input []string) []Hailstone {
	stones := make([]Hailstone, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " @ ")
		cStr := strings.Split(parts[0], ", ")
		vStr := strings.Split(parts[1], ", ")
		c := make([]int, 3)
		v := make([]int, 3)
		for i, s := range cStr {
			c[i] = util.Atoi(strings.TrimSpace(s))
		}
		for i, s := range vStr {
			v[i] = util.Atoi(strings.TrimSpace(s))
		}
		stones = append(stones, Hailstone{
			Coord:    Coord{float64(c[0]), float64(c[1]), float64(c[2])},
			Velocity: Velocity{float64(v[0]), float64(v[1]), float64(v[2])},
		})
	}
	return stones
}

func init() {
	challenges.RegisterChallengeFunc(2023, 24, 1, "day24.txt", part1)
	challenges.RegisterChallengeFunc(2023, 24, 2, "day24.txt", part2)
}
