package aoc2021_day23

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

// hallway: [...........]
//             0 0 0 0
//             1 1 1 1
//             . . . .
// rooms:      A B C D
func part1(input []string) (string, error) {
	burrow := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", burrow)

	var result int
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	// parsed := parseInput(input)

	var result int
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Burrow {
	rooms := make(map[string][]string)
	rooms["A"] = make([]string, 2)
	rooms["A"][0] = string(input[2][3])
	rooms["A"][1] = string(input[3][3])

	rooms["B"] = make([]string, 2)
	rooms["B"][0] = string(input[2][5])
	rooms["B"][1] = string(input[3][5])

	rooms["C"] = make([]string, 2)
	rooms["C"][0] = string(input[2][7])
	rooms["C"][1] = string(input[3][7])

	rooms["D"] = make([]string, 2)
	rooms["D"][0] = string(input[2][9])
	rooms["D"][1] = string(input[3][9])

	hallway := make([]string, 11)
	for i := range hallway {
		hallway[i] = "."
	}

	return Burrow{
		Hallway: hallway,
		Rooms:   rooms,
	}
}

var EnergyCosts = map[string]int{
	"A": 1,
	"B": 10,
	"C": 100,
	"D": 1000,
}

type Burrow struct {
	Hallway []string
	Rooms   map[string][]string
}

func (b Burrow) String() string {
	return fmt.Sprintf(`#############
#%s#
###%s#%s#%s#%s###
  #%s#%s#%s#%s#  
  #########`,
		strings.Join(b.Hallway, ""),
		b.Rooms["A"][0], b.Rooms["B"][0], b.Rooms["C"][0], b.Rooms["D"][0],
		b.Rooms["A"][1], b.Rooms["B"][1], b.Rooms["C"][1], b.Rooms["D"][1],
	)
}

func init() {
	challenges.RegisterChallengeFunc(2021, 23, 1, "day23.txt", part1)
	challenges.RegisterChallengeFunc(2021, 23, 2, "day23.txt", part2)
}
