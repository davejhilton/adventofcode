package challenges2021

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func day09_part1(input []string) (string, error) {
	grid := day09_parse(input)
	sum := 0
	for _, c := range getLowPoints(grid) {
		sum += c.V + 1
	}
	return fmt.Sprintf("%d", sum), nil
}

func getLowPoints(grid [][]int) []coord {
	coords := make([]coord, 0)
	for r := range grid {
		for c := range grid[r] {
			if isLowPoint(r, c, grid) {
				coords = append(coords, coord{R: r, C: c, V: grid[r][c]})
			}
		}
	}
	return coords
}

func isLowPoint(r int, c int, grid [][]int) bool {
	val := grid[r][c]
	if r > 0 && val >= grid[r-1][c] { // check up
		return false
	}
	if r < len(grid)-1 && val >= grid[r+1][c] { // check down
		return false
	}
	if c > 0 && val >= grid[r][c-1] { // check left
		return false
	}
	if c < len(grid[0])-1 && val >= grid[r][c+1] { // check right
		return false
	}
	return true
}

func arrContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

type coord struct {
	R int
	C int
	V int
}

func coordStr(r, c int) string {
	return fmt.Sprintf("%d,%d", r, c)
}

func day09_part2(input []string) (string, error) {
	grid := day09_parse(input)
	lowPoints := getLowPoints(grid)

	sizes := make([]int, 0)
	for _, p := range lowPoints {
		basin := make([]string, 0)
		basin = append(basin, coordStr(p.R, p.C))
		checkBoundaries(p.R, p.C, grid, &basin)
		log.Debugf("Low point: (%s) - basin size: %d\n", coordStr(p.R, p.C), len(basin))
		sizes = append(sizes, len(basin))
	}
	sort.Ints(sizes)

	product := 1
	for i := len(sizes) - 3; i < len(sizes); i++ {
		product *= sizes[i]
	}

	log.Debugf("Result: %d\n", product)
	return fmt.Sprintf("%d", product), nil
}

func checkBoundaries(r int, c int, grid [][]int, basin *[]string) {
	if r > 0 { // check up
		above := coordStr(r-1, c)
		if grid[r-1][c] != 9 && !arrContains(above, *basin) {
			*basin = append(*basin, above)
			checkBoundaries(r-1, c, grid, basin)
		}
	}
	if r < len(grid)-1 { // check down
		below := coordStr(r+1, c)
		if grid[r+1][c] != 9 && !arrContains(below, *basin) {
			*basin = append(*basin, below)
			checkBoundaries(r+1, c, grid, basin)
		}
	}
	if c > 0 { // check left
		left := coordStr(r, c-1)
		if grid[r][c-1] != 9 && !arrContains(left, *basin) {
			*basin = append(*basin, left)
			checkBoundaries(r, c-1, grid, basin)
		}
	}
	if c < len(grid[0])-1 { // check right
		right := coordStr(r, c+1)
		if grid[r][c+1] != 9 && !arrContains(right, *basin) {
			*basin = append(*basin, right)
			checkBoundaries(r, c+1, grid, basin)
		}
	}
}

func day09_parse(input []string) [][]int {
	grid := make([][]int, 0, len(input))
	for _, s := range input {
		strDigits := strings.Split(s, "")
		row := make([]int, 0)
		for _, sd := range strDigits {
			row = append(row, util.Atoi(sd))
		}
		grid = append(grid, row)
	}
	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2021, 9, 1, "day09.txt", day09_part1)
	challenges.RegisterChallengeFunc(2021, 9, 2, "day09.txt", day09_part2)
}
