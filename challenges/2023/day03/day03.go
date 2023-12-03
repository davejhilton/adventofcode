package aoc2023_day3

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

var numRegex = regexp.MustCompile(`[0-9]`)

const (
	BLANK = "."
	STAR  = "*"
)

type Grid struct {
	Width  int
	Height int
	Rows   [][]string
}

func isSymbol(grid Grid, row, col int) bool {
	if row < 0 || row >= grid.Height || col < 0 || col >= grid.Width {
		return false
	}
	val := grid.Rows[row][col]
	return val != BLANK && !numRegex.MatchString(val)
}

func isNumber(grid Grid, row, col int) bool {
	if row < 0 || row >= grid.Height || col < 0 || col >= grid.Width {
		return false
	}
	val := grid.Rows[row][col]
	return numRegex.MatchString(val)
}

func getFullNumber(grid Grid, row, col int) int {
	if !isNumber(grid, row, col) {
		return 0
	}
	if isNumber(grid, row, col-1) {
		// move left
		return getFullNumber(grid, row, col-1)
	}
	// we're at the leftmost digit, so start with this one and move right
	num := util.Atoi(grid.Rows[row][col])
	for isNumber(grid, row, col+1) {
		num = (num * 10) + util.Atoi(grid.Rows[row][col+1])
		col++
	}
	return num
}

func getNumbersLeftCenterRight(grid Grid, row, col int) []int {
	nums := make([]int, 0)
	if row < 0 || row > grid.Height-1 {
		return nums
	}
	var left = isNumber(grid, row, col-1)
	var center = isNumber(grid, row, col)
	var right = isNumber(grid, row, col+1)

	if left && !center && right {
		// two different numbers!
		nums = append(nums, getFullNumber(grid, row, col-1))
		nums = append(nums, getFullNumber(grid, row, col+1))
	} else if left {
		nums = append(nums, getFullNumber(grid, row, col-1))
	} else if center {
		nums = append(nums, getFullNumber(grid, row, col))
	} else if right {
		nums = append(nums, getFullNumber(grid, row, col+1))
	}
	return nums
}

func part1(input []string) (string, error) {
	grid := parseInput(input)

	var isAdjacent bool
	var result int
	for row := 0; row < grid.Height; row++ {
		curNum := 0
		isAdjacent = false
		for col := 0; col < grid.Width; col++ {
			v := grid.Rows[row][col]
			if isNumber(grid, row, col) {
				curNum = (curNum * 10) + util.Atoi(v)
				if !isAdjacent { // only check if we're not already adjacent
					if isSymbol(grid, row-1, col-1) || isSymbol(grid, row-1, col) || isSymbol(grid, row-1, col+1) ||
						isSymbol(grid, row, col-1) || isSymbol(grid, row, col+1) ||
						isSymbol(grid, row+1, col-1) || isSymbol(grid, row+1, col) || isSymbol(grid, row+1, col+1) {
						isAdjacent = true
					}
				}
			} else {
				if isAdjacent {
					result += curNum
				}
				curNum = 0
				isAdjacent = false
			}
			if col == grid.Width-1 && isAdjacent {
				// end of the row and we're adjacent, so add the number
				result += curNum
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	var result int
	for row := 0; row < grid.Height; row++ {
		for col := 0; col < grid.Width; col++ {
			v := grid.Rows[row][col]
			if v == STAR {
				adjacentNums := make([]int, 0)

				// check row above
				adjacentNums = append(adjacentNums, getNumbersLeftCenterRight(grid, row-1, col)...)
				// check left and right
				adjacentNums = append(adjacentNums, getNumbersLeftCenterRight(grid, row, col)...)
				// check row below
				adjacentNums = append(adjacentNums, getNumbersLeftCenterRight(grid, row+1, col)...)

				// if we have EXACTLY TWO adjacent numbers, multiply them and add to result
				if len(adjacentNums) == 2 {
					result += (adjacentNums[0] * adjacentNums[1])
				}

				log.Debugf("Found star at (row %d, col %d) with %d adjacent num(s) %v\n", row, col, len(adjacentNums), adjacentNums)
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Grid {
	grid := Grid{
		Rows:   make([][]string, 0, len(input)),
		Height: len(input),
	}
	for _, s := range input {
		grid.Rows = append(grid.Rows, strings.Split(s, ""))
		grid.Width = util.Max(grid.Width, len(s))
	}

	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2023, 3, 1, "day03.txt", part1)
	challenges.RegisterChallengeFunc(2023, 3, 2, "day03.txt", part2)
}
