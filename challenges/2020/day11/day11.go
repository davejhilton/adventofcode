package aoc2020_day11

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	seating := parseSeating(input)

	iterCount := 0
	numOccupied := 0
	for {
		newSeating := make([][]rune, len(seating))
		hasChanges := false
		numOccupied = 0
		iterCount++
		log.Debugf("ITERATION #%d:\n", iterCount)
		for r, row := range seating {
			newSeating[r] = make([]rune, len(row))
			log.Debug("    ")
			for s, seat := range row {
				log.Debugf("%s", string(seat))
				if seat == '.' {
					newSeating[r][s] = '.'
					continue
				} else if seat == '#' {
					numOccupied++
				}

				numOccupiedAdjacent := checkAdjacentOccupied(seating, r, s)

				nextVal := seat
				if seat == 'L' && numOccupiedAdjacent == 0 {
					nextVal = '#'
					hasChanges = true
				} else if seat == '#' && numOccupiedAdjacent >= 4 {
					nextVal = 'L'
					hasChanges = true
				}
				newSeating[r][s] = nextVal
			}
			log.Debugln()
		}
		if !hasChanges {
			break
		}
		seating = newSeating
		log.Debugln()
	}

	log.Debugf("\n%d seat(s) occupied\n", numOccupied)
	return fmt.Sprintf("%d", numOccupied), nil
}

func checkAdjacentOccupied(seating [][]rune, r int, c int) int {
	numOccupied := 0
	row := seating[r]
	if r > 0 {
		// LOOK UP AND LEFT
		if c > 0 && seating[r-1][c-1] == '#' {
			numOccupied++
		}
		// LOOK UP
		if seating[r-1][c] == '#' {
			numOccupied++
		}
		// LOOK UP AND RIGHT
		if c < len(row)-1 && seating[r-1][c+1] == '#' {
			numOccupied++
		}
	}
	// LOOK LEFT
	if c > 0 && row[c-1] == '#' {
		numOccupied++
	}
	// LOOK RIGHT
	if c < len(row)-1 && row[c+1] == '#' {
		numOccupied++
	}
	if r < len(seating)-1 {
		// LOOK DOWN AND LEFT
		if c > 0 && seating[r+1][c-1] == '#' {
			numOccupied++
		}
		// LOOK DOWN
		if seating[r+1][c] == '#' {
			numOccupied++
		}
		// LOOK DOWN AND RIGHT
		if c < len(row)-1 && seating[r+1][c+1] == '#' {
			numOccupied++
		}
	}
	return numOccupied
}

func part2(input []string) (string, error) {
	seating := parseSeating(input)

	iterCount := 0
	numOccupied := 0
	for {
		nextSeating := make([][]rune, len(seating))
		hasChanges := false
		numOccupied = 0
		iterCount++
		log.Debugf("ITERATION #%d:\n", iterCount)
		for i, row := range seating {
			nextSeating[i] = make([]rune, len(row))
			log.Debug("    ")
			for j, seat := range row {
				log.Debugf("%s", string(seat))
				if seat == '.' {
					nextSeating[i][j] = '.'
					continue
				} else if seat == '#' {
					numOccupied++
				}

				adjacentOccupied := checkLineOfSightOccupied(seating, i, j)

				nextVal := seat
				if seat == 'L' && adjacentOccupied == 0 {
					nextVal = '#'
					hasChanges = true
				} else if seat == '#' && adjacentOccupied >= 5 {
					nextVal = 'L'
					hasChanges = true
				}
				nextSeating[i][j] = nextVal
			}
			log.Debugln()
		}
		log.Debugln()

		if !hasChanges {
			break
		}
		seating = nextSeating
	}

	log.Debugf("%d seat(s) occupied\n", numOccupied)
	return fmt.Sprintf("%d", numOccupied), nil
}

func checkLineOfSightOccupied(seating [][]rune, r int, c int) int {
	TOPLEFT := 0
	TOP := 1
	TOPRIGHT := 2
	LEFT := 3
	RIGHT := 4
	BOTTOMLEFT := 5
	BOTTOM := 6
	BOTTOMRIGHT := 7

	//  0 1 2
	//  3 * 4
	//  5 6 7

	lineOfSight := make([]rune, 8)
	for x := range lineOfSight {
		lineOfSight[x] = '.'
	}

	for curRow := r - 1; curRow >= 0; curRow-- {
		diff := r - curRow
		// LOOK UP AND LEFT
		curCol := c - diff
		if curCol >= 0 && lineOfSight[TOPLEFT] == '.' {
			lineOfSight[TOPLEFT] = seating[curRow][curCol]
		}
		// LOOK UP
		if lineOfSight[TOP] == '.' {
			lineOfSight[TOP] = seating[curRow][c]
		}
		// LOOK UP AND RIGHT
		curCol = c + diff
		if curCol < len(seating[curRow]) && lineOfSight[TOPRIGHT] == '.' {
			lineOfSight[TOPRIGHT] = seating[curRow][curCol]
		}
	}

	for curCol := c - 1; curCol >= 0; curCol-- {
		// LOOK LEFT
		if lineOfSight[LEFT] == '.' {
			lineOfSight[LEFT] = seating[r][curCol]
		} else {
			break
		}
	}
	for curCol := c + 1; curCol < len(seating[r]); curCol++ {
		// LOOK RIGHT
		if lineOfSight[RIGHT] == '.' {
			lineOfSight[RIGHT] = seating[r][curCol]
		} else {
			break
		}
	}

	for curRow := r + 1; curRow < len(seating); curRow++ {
		diff := curRow - r
		// LOOK DOWN AND LEFT
		curCol := c - diff
		if curCol >= 0 && lineOfSight[BOTTOMLEFT] == '.' {
			lineOfSight[BOTTOMLEFT] = seating[curRow][curCol]
		}
		// LOOK DOWN
		if lineOfSight[BOTTOM] == '.' {
			lineOfSight[BOTTOM] = seating[curRow][c]
		}
		// LOOK DOWN AND RIGHT
		curCol = c + diff
		if curCol < len(seating[curRow]) && lineOfSight[BOTTOMRIGHT] == '.' {
			lineOfSight[BOTTOMRIGHT] = seating[curRow][curCol]
		}
	}

	numOccupied := 0
	for _, s := range lineOfSight {
		if s == '#' {
			numOccupied++
		}
	}

	return numOccupied
}

func parseSeating(input []string) [][]rune {
	rows := make([][]rune, 0, len(input))
	for _, s := range input {
		row := make([]rune, 0, len(s))
		for _, c := range s {
			row = append(row, c)
		}
		rows = append(rows, row)
	}
	return rows
}

func init() {
	challenges.RegisterChallengeFunc(2020, 11, 1, "day11.txt", part1)
	challenges.RegisterChallengeFunc(2020, 11, 2, "day11.txt", part2)
}
