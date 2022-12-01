package aoc2021_day20

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	alg, grid := parseInput(input)
	log.Debugf("INPUT:\n%s\n", grid)

	NUM_ITERATIONS := 2
	for i := 0; i < NUM_ITERATIONS; i++ {
		grid = grid.Enhance(alg, i)
	}

	result := grid.CountLight()
	log.Debugf("\n\nRESULT:\n%s\n", grid)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	alg, grid := parseInput(input)
	// log.Debugf("INPUT:\n%s\n", grid)

	NUM_ITERATIONS := 50
	for i := 0; i < NUM_ITERATIONS; i++ {
		grid = grid.Enhance(alg, i)
	}

	result := grid.CountLight()
	// log.Debugf("\n\nRESULT:\n%s\n", grid)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) ([]int, Grid) {
	alg := make([]int, 0, len(input[0]))
	for _, r := range input[0] {
		v := 0
		if r == '#' {
			v = 1
		}
		alg = append(alg, v)
	}

	grid := make(Grid, 0, len(input))
	for _, s := range input[2:] {
		row := make([]int, 0, len(s))
		for _, r := range s {
			v := 0
			if r == '#' {
				v = 1
			}
			row = append(row, v)
		}
		grid = append(grid, row)
	}
	return alg, grid
}

type Grid [][]int

func (g Grid) String() string {
	var sb strings.Builder
	for i := range g {
		if i > 0 {
			sb.WriteString("\n")
		}
		for j := range g[i] {
			if g[i][j] == 0 {
				sb.WriteRune('.')
			} else {
				sb.WriteRune('#')
			}
		}
	}
	return sb.String()
}

func (g Grid) Expand(defaultVal int) Grid {
	oldW := len(g[0])
	oldH := len(g)

	newG := make(Grid, oldH+2)
	newG[0] = make([]int, oldW+2)
	if defaultVal == 1 {
		for i := range newG[0] {
			newG[0][i] = 1
		}
	}
	for i := range g {
		row := make([]int, oldW+2)
		for j := range g[i] {
			row[j+1] = g[i][j]
		}
		if defaultVal == 1 {
			row[0] = 1
			row[len(row)-1] = 1
		}
		newG[i+1] = row
	}
	newG[oldH+1] = make([]int, oldW+2)
	if defaultVal == 1 {
		for i := range newG[oldH+1] {
			newG[oldH+1][i] = 1
		}
	}

	return newG
}

func (g Grid) Enhance(alg []int, iter int) Grid {
	defaultValue := 0
	if alg[0] == 1 && iter%2 == 1 {
		defaultValue = 1
	}
	g2 := g.Expand(defaultValue)
	// log.Debugf("AFTER EXPAND:\n%s\n", g2)
	// g2 := g
	newGr := make(Grid, len(g2))
	// newGr[0] = make([]int, len(g2[0]))
	for i := 0; i < len(g2); i++ {
		newGr[i] = make([]int, len(g2[i]))
		for j := 0; j < len(g2[i]); j++ {
			str := ""

			// up & left
			if i == 0 || j == 0 {
				str = fmt.Sprintf("%d", defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i-1][j-1])
			}
			// up
			if i == 0 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i-1][j])
			}
			// up & right
			if i == 0 || j == len(g2[i])-1 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i-1][j+1])
			}
			// left
			if j == 0 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i][j-1])
			}
			// center
			str = fmt.Sprintf("%s%d", str, g2[i][j])
			// right
			if j == len(g2[i])-1 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i][j+1])
			}
			// down & left
			if i == len(g2)-1 || j == 0 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i+1][j-1])
			}
			// down
			if i == len(g2)-1 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i+1][j])
			}
			// down & right
			if i == len(g2)-1 || j == len(g2[i])-1 {
				str = fmt.Sprintf("%s%d", str, defaultValue)
			} else {
				str = fmt.Sprintf("%s%d", str, g2[i+1][j+1])
			}

			num, _ := strconv.ParseUint(str, 2, 16)
			// log.Debugf("PIXEL VALUE of (%d,%d): %d [%s]\n", i, j, num, str)
			newGr[i][j] = alg[num]
		}
	}
	// newGr[len(g2)-1] = make([]int, len(g2[0]))
	// log.Debugf("\nAFTER ENHANCE:\n%s\n\n", newGr)
	return newGr
}

func (g Grid) CountLight() int {
	count := 0
	for i := range g {
		for j := range g[i] {
			count += g[i][j]
		}
	}
	return count
}

func init() {
	challenges.RegisterChallengeFunc(2021, 20, 1, "day20.txt", part1)
	challenges.RegisterChallengeFunc(2021, 20, 2, "day20.txt", part2)
}
