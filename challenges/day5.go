package challenges

import (
	"fmt"

	"github.com/davejhilton/adventofcode2020/log"
)

func day5_part1(input []string, isExample bool) (string, error) {
	maxSeat := 0
	for i, line := range input {
		log.Debugf("Line %d\n", i)
		seat := day5_findSeat(line)
		if seat > maxSeat {
			maxSeat = seat
		}
		log.Debugf("seat = %d\n\n", seat)
	}
	return fmt.Sprintf("%d", maxSeat), nil
}

func day5_part2(input []string, isExample bool) (string, error) {
	seats := make([]bool, 1027, 1027)
	for i, line := range input {
		log.Debugf("Line %d\n", i)
		seat := day5_findSeat(line)
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

func day5_findSeat(line string) int {
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
	registerChallengeFunc(5, 1, "day5.txt", day5_part1)
	registerChallengeFunc(5, 2, "day5.txt", day5_part2)
}
