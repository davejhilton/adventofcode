package aoc2022_day12

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type coord struct {
	Row int
	Col int
}

func (c coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

func part1(input []string) (string, error) {
	grid, start, end := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n%v\n%v\n", grid, start, end)

	var finalPath []coord
	queue := [][]coord{{start}}
	seen := make(map[string]bool)

	for len(queue) > 0 && len(finalPath) == 0 {
		thisPath := queue[0]
		queue = queue[1:]
		cur := thisPath[len(thisPath)-1]

		dirsToMove := []coord{
			{cur.Row + 1, cur.Col},
			{cur.Row, cur.Col + 1},
			{cur.Row - 1, cur.Col},
			{cur.Row, cur.Col - 1},
		}

		for _, d := range dirsToMove {
			if d.Row >= 0 && d.Row < len(grid) && d.Col >= 0 && d.Col < len(grid[d.Row]) && grid[d.Row][d.Col]-grid[cur.Row][cur.Col] <= 1 && !seen[d.String()] {
				nPath := make([]coord, len(thisPath))
				copy(nPath, thisPath)
				nPath = append(nPath, d)

				seen[d.String()] = true
				queue = append(queue, nPath)
				if d.Row == end.Row && d.Col == end.Col {
					finalPath = nPath
					break
				}
			}
		}
	}

	log.Debugf("\n\n%v\n\n", finalPath)

	var result int = len(finalPath) - 1

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid, _, end := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n%v\n", grid, end)

	var finalPath []coord
	queue := [][]coord{{end}}
	seen := make(map[string]bool)

	for len(queue) > 0 && len(finalPath) == 0 {
		thisPath := queue[0]
		queue = queue[1:]
		cur := thisPath[len(thisPath)-1]

		dirsToMove := []coord{
			{cur.Row + 1, cur.Col},
			{cur.Row, cur.Col + 1},
			{cur.Row - 1, cur.Col},
			{cur.Row, cur.Col - 1},
		}

		for _, d := range dirsToMove {
			if d.Row >= 0 && d.Row < len(grid) && d.Col >= 0 && d.Col < len(grid[d.Row]) && grid[cur.Row][cur.Col]-grid[d.Row][d.Col] <= 1 && !seen[d.String()] {
				nPath := make([]coord, len(thisPath))
				copy(nPath, thisPath)
				nPath = append(nPath, d)

				seen[d.String()] = true
				queue = append(queue, nPath)
				if grid[d.Row][d.Col] == 0 {
					finalPath = nPath
					break
				}
			}
		}
	}

	log.Debugf("\n\n%v\n\n", finalPath)

	var result int = len(finalPath) - 1

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) ([][]int, coord, coord) {
	grid := make([][]int, 0, len(input))
	var start, end coord
	for r, row := range input {
		gridRow := make([]int, 0, len(row))
		for c, square := range row {
			var v int
			if square == 'S' {
				start = coord{r, c}
				v = 1
			} else if square == 'E' {
				end = coord{r, c}
				v = 26
			} else {
				v = int(square) - 97
			}
			gridRow = append(gridRow, v)
		}
		grid = append(grid, gridRow)
	}
	return grid, start, end
}

func init() {
	challenges.RegisterChallengeFunc(2022, 12, 1, "day12.txt", part1)
	challenges.RegisterChallengeFunc(2022, 12, 2, "day12.txt", part2)
}
