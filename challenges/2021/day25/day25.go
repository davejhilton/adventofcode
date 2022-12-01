package aoc2021_day25

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%s\n", grid)

	step := 0
	for {
		r1 := grid.StepRight()
		r2 := grid.StepDown()
		step++
		log.Debugf("%d steps finished\n", step)
		if !r1 && !r2 {
			break
		}
	}

	log.Debugf("\nAfter %d steps:\n%s\n", step, grid)

	return fmt.Sprintf("%d", step), nil
}

func part2(input []string) (string, error) {
	// parsed := parseInput(input)

	var result int
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *Grid {
	grid := make(Grid, 0, len(input))
	for _, s := range input {
		grid = append(grid, []rune(s))
	}
	return &grid
}

type Grid [][]rune

func (g *Grid) StepRight() bool {
	somethingMoved := false
	for r := range *g {
		row := (*g)[r]
		prevVal, prevIdx := row[0], 0
		zeroVal := prevVal
		for c := len((*g)[0]) - 1; c >= 0; c-- {
			curVal := row[c]
			if c == 0 {
				curVal = zeroVal
			}
			if prevVal == '.' && curVal == '>' {
				row[prevIdx] = '>'
				row[c] = '.'
				somethingMoved = true
			}
			prevVal, prevIdx = curVal, c
		}
	}

	return somethingMoved
}

func (g *Grid) StepDown() bool {
	somethingMoved := false
	for c := range (*g)[0] {
		prevVal, prevIdx := (*g)[0][c], 0
		zeroVal := prevVal
		for r := len(*g) - 1; r >= 0; r-- {
			curVal := (*g)[r][c]
			if r == 0 {
				curVal = zeroVal
			}
			if prevVal == '.' && curVal == 'v' {
				(*g)[prevIdx][c] = 'v'
				(*g)[r][c] = '.'
				somethingMoved = true
			}
			prevVal, prevIdx = curVal, r
		}
	}

	return somethingMoved
}

func (g *Grid) String() string {
	var sb strings.Builder
	for r := range *g {
		if r != 0 {
			sb.WriteString("\n")
		}
		for c := range (*g)[r] {
			sb.WriteRune((*g)[r][c])
		}
	}
	return sb.String()
}

func init() {
	challenges.RegisterChallengeFunc(2021, 25, 1, "day25.txt", part1)
	challenges.RegisterChallengeFunc(2021, 25, 2, "day25.txt", part2)
}
