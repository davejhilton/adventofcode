package aoc2023_day14

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type Grid [][]string

func (g Grid) TiltN() {
	for r, row := range g {
		for c, cell := range row {
			if cell == "O" {
				newR := r
				for i := r - 1; i >= 0; i-- {
					if g[i][c] == "#" || g[i][c] == "O" {
						break
					} else {
						newR = i
					}
				}
				if newR != r {
					g[newR][c] = "O"
					g[r][c] = "."
				}
			}
		}
	}
}

func (g Grid) TiltS() {
	for r := len(g) - 1; r >= 0; r-- {
		for c, cell := range g[r] {
			if cell == "O" {
				newR := r
				for i := r + 1; i < len(g); i++ {
					if g[i][c] == "#" || g[i][c] == "O" {
						break
					} else {
						newR = i
					}
				}
				if newR != r {
					g[newR][c] = "O"
					g[r][c] = "."
				}
			}
		}
	}
}

func (g Grid) TiltE() {
	for r, row := range g {
		for c := len(row) - 1; c >= 0; c-- {
			if row[c] == "O" {
				newC := c
				for i := c + 1; i < len(row); i++ {
					if g[r][i] == "#" || g[r][i] == "O" {
						break
					} else {
						newC = i
					}
				}
				if newC != c {
					g[r][newC] = "O"
					g[r][c] = "."
				}
			}
		}
	}
}

func (g Grid) TiltW() {
	for r, row := range g {
		for c, cell := range row {
			if cell == "O" {
				newC := c
				for i := c - 1; i >= 0; i-- {
					if g[r][i] == "#" || g[r][i] == "O" {
						break
					} else {
						newC = i
					}
				}
				if newC != c {
					g[r][newC] = "O"
					g[r][c] = "."
				}
			}
		}
	}
}

var cache = make(map[string]int)

func SpinGrid(g Grid, n int) string {
	g.TiltN()
	g.TiltW()
	g.TiltS()
	g.TiltE()
	str := g.String()
	if _, ok := cache[str]; ok {
		// log.Debugf("Found a cycle at %d:\n%s\n", n, str)
		return str
	}
	cache[g.String()] = n
	return ""
}

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		for _, cell := range row {
			sb.WriteString(cell)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g Grid) Weight() int {
	var result int
	weight := 0
	for i := len(g) - 1; i >= 0; i-- {
		weight++
		for j := 0; j < len(g[i]); j++ {
			if g[i][j] == "O" {
				result += weight
			}
		}
	}
	return result
}

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	grid.TiltN()
	log.Debugf("Tilted N:\n%s\n", grid)

	result := grid.Weight()
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	totalSpins := 1000000000
	var firstMatchStr string
	var firstMatchI int
	var i int
	for i = 0; i < totalSpins; i++ {
		r := SpinGrid(grid, i)
		if firstMatchStr == "" && r != "" {
			firstMatchStr = r
			firstMatchI = i
		} else if firstMatchStr != "" && r == firstMatchStr {
			log.Debugf("Found a cycle between i = %d and i = %d:\n", firstMatchI, i)
			diff := i - firstMatchI
			remainingSpins := (totalSpins - firstMatchI) % diff
			i = totalSpins - remainingSpins
		}
	}

	log.Debugf("Spun %d times:\n%s\n", i, grid)

	// diff := i - firstMatchI
	// remainingSpins := (totalSpins - firstMatchI) % diff
	// log.Debugf("Found %d remaining spins\n", remainingSpins)
	// for j := 0; j < remainingSpins-1; j++ {
	// 	SpinGrid(grid, i+j)
	// 	log.Debugf("Spun %d times: (%d)\n%s\n", totalSpins-(remainingSpins-j-1), grid.Weight(), grid)
	// }

	result := grid.Weight()
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
	challenges.RegisterChallengeFunc(2023, 14, 1, "day14.txt", part1)
	challenges.RegisterChallengeFunc(2023, 14, 2, "day14.txt", part2)
}
