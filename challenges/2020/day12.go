package challenges2020

import (
	"fmt"
	"math"
	"strconv"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day12_part1(input []string) (string, error) {
	directions := day12_parseDirections(input)

	state := day12_state{
		NorthVal: 0,
		EastVal:  0,
		Facing:   EAST,
	}

	log.Debugf("START: %d N, %d E, Facing: %s\n", state.NorthVal, state.EastVal, day12_getDirectionName(state.Facing))

	for _, dir := range directions {
		switch dir.Action {
		case 'N':
			state.NorthVal += dir.Value
		case 'S':
			state.NorthVal -= dir.Value
		case 'E':
			state.EastVal += dir.Value
		case 'W':
			state.EastVal -= dir.Value
		case 'R':
			units := int(dir.Value / 90)
			state.Facing = (state.Facing + units) % 4
		case 'L':
			units := int(dir.Value / 90)
			state.Facing = (state.Facing + 4 - units) % 4
		case 'F':
			if state.Facing == EAST {
				state.EastVal += dir.Value
			} else if state.Facing == WEST {
				state.EastVal -= dir.Value
			} else if state.Facing == NORTH {
				state.NorthVal += dir.Value
			} else if state.Facing == SOUTH {
				state.NorthVal -= dir.Value
			}
		}
		log.Debugf("%s%d --> %d N, %d E, Facing: %s\n", string(dir.Action), dir.Value, state.NorthVal, state.EastVal, day12_getDirectionName(state.Facing))
	}

	log.Debugf("Ended at: North %d, East %d\n", state.NorthVal, state.EastVal)
	result := int(math.Abs(float64(state.NorthVal))) + int(math.Abs(float64(state.EastVal)))
	return fmt.Sprintf("%d", result), nil
}

func day12_part2(input []string) (string, error) {
	directions := day12_parseDirections(input)

	ship := day12_state{
		NorthVal: 0,
		EastVal:  0,
	}

	wp := day12_state{
		NorthVal: 1,
		EastVal:  10,
	}

	log.Debugf("START: SHIP(%d N, %d E), WAYPOINT(%d N, %d E)\n", ship.NorthVal, ship.EastVal, wp.NorthVal, wp.EastVal)

	for _, dir := range directions {
		switch dir.Action {
		case 'N':
			wp.NorthVal += dir.Value
		case 'S':
			wp.NorthVal -= dir.Value
		case 'E':
			wp.EastVal += dir.Value
		case 'W':
			wp.EastVal -= dir.Value
		case 'R':
			units := int(dir.Value / 90)
			newState := day12_state{}

			if wp.NorthVal >= 0 {
				nd := (NORTH + units) % 4
				if nd == EAST {
					newState.EastVal = wp.NorthVal
				} else if nd == SOUTH {
					newState.NorthVal = wp.NorthVal * -1
				} else if nd == WEST {
					newState.EastVal = wp.NorthVal * -1
				}
			} else {
				nd := (SOUTH + units) % 4
				if nd == EAST {
					newState.EastVal = wp.NorthVal * -1
				} else if nd == NORTH {
					newState.NorthVal = wp.NorthVal * -1
				} else if nd == WEST {
					newState.EastVal = wp.NorthVal
				}
			}

			if wp.EastVal >= 0 {
				ed := (EAST + units) % 4
				if ed == SOUTH {
					newState.NorthVal = wp.EastVal * -1
				} else if ed == WEST {
					newState.EastVal = wp.EastVal * -1
				} else if ed == NORTH {
					newState.NorthVal = wp.EastVal
				}
			} else {
				ed := (WEST + units) % 4
				if ed == NORTH {
					newState.NorthVal = wp.EastVal * -1
				} else if ed == EAST {
					newState.EastVal = wp.EastVal * -1
				} else if ed == SOUTH {
					newState.NorthVal = wp.EastVal
				}
			}
			wp = newState

		case 'L':

			// N = -2, E = 7
			units := int(dir.Value / 90)
			newState := day12_state{}
			if wp.NorthVal >= 0 {
				nd := (NORTH + 4 - units) % 4
				if nd == EAST {
					newState.EastVal = wp.NorthVal
				} else if nd == SOUTH {
					newState.NorthVal = wp.NorthVal * -1
				} else if nd == WEST {
					newState.EastVal = wp.NorthVal * -1
				}
			} else {
				nd := (SOUTH + 4 - units) % 4
				if nd == EAST {
					newState.EastVal = wp.NorthVal * -1
				} else if nd == NORTH {
					newState.NorthVal = wp.NorthVal * -1
				} else if nd == WEST {
					newState.EastVal = wp.NorthVal
				}
			}

			if wp.EastVal >= 0 {
				ed := (EAST + 4 - units) % 4
				if ed == SOUTH {
					newState.NorthVal = wp.EastVal * -1
				} else if ed == WEST {
					newState.EastVal = wp.EastVal * -1
				} else if ed == NORTH {
					newState.NorthVal = wp.EastVal
				}
			} else {
				ed := (WEST + 4 - units) % 4
				if ed == NORTH {
					newState.NorthVal = wp.EastVal * -1
				} else if ed == EAST {
					newState.EastVal = wp.EastVal * -1
				} else if ed == SOUTH {
					newState.NorthVal = wp.EastVal
				}
			}
			wp = newState

		case 'F':
			ship.NorthVal += (dir.Value * wp.NorthVal)
			ship.EastVal += (dir.Value * wp.EastVal)
		}
		log.Debugf("%s%d --> SHIP(%d N, %d E), WAYPOINT(%d N, %d E)\n", string(dir.Action), dir.Value, ship.NorthVal, ship.EastVal, wp.NorthVal, wp.EastVal)
	}

	log.Debugf("Ended at: SHIP(%d N, %d E), WAYPOINT(%d N, %d E)\n", ship.NorthVal, ship.EastVal, wp.NorthVal, wp.EastVal)
	result := int(math.Abs(float64(ship.NorthVal))) + int(math.Abs(float64(ship.EastVal)))
	return fmt.Sprintf("%d", result), nil
}

const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

func day12_getDirectionName(d int) string {
	switch d {
	case NORTH:
		return "NORTH"
	case EAST:
		return "EAST"
	case SOUTH:
		return "SOUTH"
	case WEST:
		return "WEST"
	default:
		return fmt.Sprintf("??? (%d)", d)
	}
}

type day12_direction struct {
	Action rune
	Value  int
}

type day12_state struct {
	NorthVal int
	EastVal  int
	Facing   int
}

func day12_parseDirections(input []string) []day12_direction {
	directions := make([]day12_direction, 0, len(input))
	for _, s := range input {
		d := rune(s[0])
		v, _ := strconv.Atoi(s[1:])
		directions = append(directions, day12_direction{
			Action: d,
			Value:  v,
		})
	}
	return directions
}

func init() {
	challenges.RegisterChallengeFunc(2020, 12, 1, "day12.txt", day12_part1)
	challenges.RegisterChallengeFunc(2020, 12, 2, "day12.txt", day12_part2)
}
