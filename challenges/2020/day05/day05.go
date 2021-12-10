package aoc2020_day5

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	maxSeat := 0
	for i, line := range input {
		log.Debugf("Line %d\n", i)
		seat := findSeat(line)
		if seat > maxSeat {
			maxSeat = seat
		}
		log.Debugf("seat = %d\n\n", seat)
	}
	return fmt.Sprintf("%d", maxSeat), nil
}

func part2(input []string) (string, error) {
	seats := make([]bool, 1027)
	for i, line := range input {
		log.Debugf("Line %d\n", i)
		seat := findSeat(line)
		log.Debugf("seat = %d\n\n", seat)
		seats[seat] = true
	}
	foundFirst := false
	mySeat := 0
	for i, v := range seats {
		if v {
			if !foundFirst {
				log.Debugf("First seat: %d\n", i)
			}
			foundFirst = true
		} else {
			if foundFirst {
				mySeat = i
				break
			}
		}
	}
	return fmt.Sprintf("%d", mySeat), nil
}

func findSeat(line string) int {
	fRow := 0
	lRow := 127
	fCol := 0
	lCol := 7
	for _, c := range line {
		switch c {
		case 'F':
			lRow = fRow + int((lRow-fRow)/2)
		case 'B':
			fRow = fRow + int((lRow-fRow)/2) + 1
		case 'L':
			lCol = fCol + int((lCol-fCol)/2)
		case 'R':
			fCol = fCol + int((lCol-fCol)/2) + 1
		}
		log.Debugf("\t c = %s | %d - %d | %d - %d\n", string(c), fRow, lRow, fCol, lCol)
	}
	return lRow*8 + lCol
}

func init() {
	challenges.RegisterChallengeFunc(2020, 5, 1, "day05.txt", part1)
	challenges.RegisterChallengeFunc(2020, 5, 2, "day05.txt", part2)
}
