package aoc2022_day8

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	var result int
	var b strings.Builder
	for x := range grid {
		for y := range grid[x] {
			var visible bool = visibleFromLeft(grid, x, y) ||
				visibleFromRight(grid, x, y) ||
				visibleFromTop(grid, x, y) ||
				visibleFromBottom(grid, x, y)

			if visible {
				b.WriteString(log.Colorize(grid[x][y], log.Green, 0))
				result += 1
			} else {
				b.WriteString(log.Colorize(grid[x][y], log.Red, 0))
			}
		}
		b.WriteString("\n")
	}

	log.Debugln(b.String())

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	var result int
	for x := range grid {
		for y := range grid[x] {
			if x == 0 || y == 0 || x == len(grid)-1 || y == len(grid[x])-1 {
				continue
			}
			score := lookLeft(grid, x, y) * lookRight(grid, x, y) * lookUp(grid, x, y) * lookDown(grid, x, y)
			log.Debugf("row %d col %d (%d) has total visibility score of %d\n\n", x, y, grid[x][y], score)
			result = util.Max(score, result)
		}
	}
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) [][]int {
	nums := make([][]int, 0, len(input))
	for _, s := range input {
		row := make([]int, 0, len(s))
		for _, h := range s {
			row = append(row, util.Atoi(string(h)))
		}
		nums = append(nums, row)
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2022, 8, 1, "day08.txt", part1)
	challenges.RegisterChallengeFunc(2022, 8, 2, "day08.txt", part2)
}

// ===========================
//       PART 1 HELPERS
// ===========================

func visibleFromLeft(grid [][]int, x, y int) bool {
	var h int = grid[x][y]
	for i := 0; i < y; i++ {
		if grid[x][i] >= h {
			return false
		}
	}
	log.Debugf("row %d col %d (%d) is visible from the LEFT\n", x, y, grid[x][y])
	return true
}

func visibleFromRight(grid [][]int, x, y int) bool {
	var h int = grid[x][y]
	for i := len(grid[x]) - 1; i > y; i-- {
		if grid[x][i] >= h {
			return false
		}
	}
	log.Debugf("row %d col %d (%d) is visible from the RIGHT\n", x, y, grid[x][y])
	return true
}

func visibleFromTop(grid [][]int, x, y int) bool {
	var h int = grid[x][y]
	for i := 0; i < x; i++ {
		if grid[i][y] >= h {
			return false
		}
	}
	log.Debugf("row %d col %d (%d) is visible from the TOP\n", x, y, grid[x][y])
	return true
}

func visibleFromBottom(grid [][]int, x, y int) bool {
	var h int = grid[x][y]
	for i := len(grid) - 1; i > x; i-- {
		if grid[i][y] >= h {
			return false
		}
	}
	log.Debugf("row %d col %d (%d) is visible from the BOTTOM\n", x, y, grid[x][y])
	return true
}

// ===========================
//       PART 2 HELPERS
// ===========================

func lookLeft(grid [][]int, x, y int) int {
	dist := 0
	h := grid[x][y]
	for k := y - 1; k >= 0; k-- {
		dist += 1
		if grid[x][k] >= h {
			break
		}
	}
	log.Debugf("row %d col %d (%d) has vis dist to the LEFT of %d\n", x, y, grid[x][y], dist)
	return dist
}

func lookRight(grid [][]int, x, y int) int {
	dist := 0
	h := grid[x][y]
	for k := y + 1; k < len(grid[x]); k++ {
		dist += 1
		if grid[x][k] >= h {
			break
		}
	}
	log.Debugf("row %d col %d (%d) has vis dist to the RIGHT of %d\n", x, y, grid[x][y], dist)
	return dist
}

func lookUp(grid [][]int, x, y int) int {
	dist := 0
	h := grid[x][y]
	for k := x - 1; k >= 0; k-- {
		dist += 1
		if grid[k][y] >= h {
			break
		}
	}
	log.Debugf("row %d col %d (%d) has vis dist to the TOP of %d\n", x, y, grid[x][y], dist)
	return dist
}

func lookDown(grid [][]int, x, y int) int {
	dist := 0
	h := grid[x][y]
	for k := x + 1; k < len(grid); k++ {
		dist += 1
		if grid[k][y] >= h {
			break
		}
	}
	log.Debugf("row %d col %d (%d) has vis dist to the BOTTOM of %d\n", x, y, grid[x][y], dist)
	return dist
}
