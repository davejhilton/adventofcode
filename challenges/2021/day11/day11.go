package aoc2021_day11

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	grid := parseInput(input)
	count := 0
	for i := 0; i < 100; i++ {
		count += runStep(&grid)
	}
	return fmt.Sprintf("%d", count), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	stepNum := 1
	for {
		count := runStep(&grid)
		log.Debugf("STEP %d - COUNT: %d\n", stepNum, count)
		if count == 100 {
			break
		}
		stepNum += 1
	}
	return fmt.Sprintf("%d", stepNum), nil
}

type coord struct {
	R int
	C int
}

func incrAndCheckFlash(grid *[][]int, flashed *[]coord, r int, c int) {
	(*grid)[r][c] += 1
	if (*grid)[r][c] == 10 {
		*flashed = append(*flashed, coord{r, c})
	}
}

func runStep(g *[][]int) int {
	grid := *g
	flashed := make([]coord, 0)
	for r := range grid {
		for c := range grid[r] {
			grid[r][c] += 1
			if grid[r][c] == 10 {
				flashed = append(flashed, coord{r, c})
			}
		}
	}
	for len(flashed) != 0 {
		newFlashed := make([]coord, 0)
		for _, crd := range flashed {
			r, c := crd.R, crd.C
			if r > 0 { // handle up
				incrAndCheckFlash(&grid, &newFlashed, r-1, c)
			}
			if r > 0 && c < len(grid[0])-1 { // handle up/right
				incrAndCheckFlash(&grid, &newFlashed, r-1, c+1)
			}
			if c < len(grid[0])-1 { // handle right
				incrAndCheckFlash(&grid, &newFlashed, r, c+1)
			}
			if c < len(grid[0])-1 && r < len(grid)-1 { // handle down/right
				incrAndCheckFlash(&grid, &newFlashed, r+1, c+1)
			}
			if r < len(grid)-1 { // handle down
				incrAndCheckFlash(&grid, &newFlashed, r+1, c)
			}
			if r < len(grid)-1 && c > 0 { // handle down/left
				incrAndCheckFlash(&grid, &newFlashed, r+1, c-1)
			}
			if c > 0 { // handle left
				incrAndCheckFlash(&grid, &newFlashed, r, c-1)
			}
			if c > 0 && r > 0 { // handle up/left
				incrAndCheckFlash(&grid, &newFlashed, r-1, c-1)
			}
		}
		flashed = newFlashed
	}

	count := 0
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] > 9 {
				grid[r][c] = 0
				count += 1
			}
		}
	}
	return count
}

func parseInput(input []string) [][]int {
	nums := make([][]int, 0, len(input))
	for _, s := range input {
		rowStrs := strings.Split(s, "")
		row := make([]int, 0, len(rowStrs))
		for _, n := range rowStrs {
			row = append(row, util.Atoi(n))
		}
		nums = append(nums, row)
	}
	return nums
}

func init() {
	challenges.RegisterChallengeFunc(2021, 11, 1, "day11.txt", part1)
	challenges.RegisterChallengeFunc(2021, 11, 2, "day11.txt", part2)
}
