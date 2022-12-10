package aoc2022_day9

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type instruction struct {
	Direction string
	Magnitude int
}

type coordinate struct {
	X int
	Y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func part1(input []string) (string, error) {
	instrs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", instrs)

	coords := make([]*coordinate, 0)
	for i := 0; i < 2; i++ {
		coords = append(coords, &coordinate{0, 0})
	}

	history := make(map[string]bool)
	history[coords[0].String()] = true

	for _, inst := range instrs {
		move(coords, inst, &history)
		for i, c := range coords {
			log.Debugf("%d: %s\n", i, c.String())
		}
	}

	var result int = len(history)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	instrs := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", instrs)

	coords := make([]*coordinate, 0)
	for i := 0; i < 10; i++ {
		coords = append(coords, &coordinate{0, 0})
	}

	history := make(map[string]bool)
	history[coords[0].String()] = true

	for _, inst := range instrs {
		move(coords, inst, &history)
		for i, c := range coords {
			log.Debugf("%d: %s\n", i, c.String())
		}
	}

	var result int = len(history)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []instruction {
	instr := make([]instruction, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " ")
		instr = append(instr, instruction{
			Direction: parts[0],
			Magnitude: util.Atoi(parts[1]),
		})
	}
	return instr
}

func init() {
	challenges.RegisterChallengeFunc(2022, 9, 1, "day09.txt", part1)
	challenges.RegisterChallengeFunc(2022, 9, 2, "day09.txt", part2)
}

func move(coords []*coordinate, in instruction, history *map[string]bool) {
	var coord *coordinate
	var prev *coordinate

	log.Debugf("\n== MOVING %s %d ==\n", in.Direction, in.Magnitude)

	for i := 0; i < in.Magnitude; i++ {
		prev = nil
		log.Debugln()
		for c := 0; c < len(coords); c++ {
			coord = coords[c]

			if prev == nil {
				switch in.Direction {
				case "R":
					coord.X += 1
					log.Debugf("    coords[%d] MOVED %-2s : POS: %s\n", c, in.Direction, coord.String())
				case "L":
					coord.X -= 1
					log.Debugf("    coords[%d] MOVED %-2s : POS: %s\n", c, in.Direction, coord.String())
				case "U":
					coord.Y += 1
					log.Debugf("    coords[%d] MOVED %-2s : POS: %s\n", c, in.Direction, coord.String())
				case "D":
					coord.Y -= 1
					log.Debugf("    coords[%d] MOVED %-2s : POS: %s\n", c, in.Direction, coord.String())
				}
			} else if !touching(*coord, *prev) {
				dir := ""
				if prev.X-coord.X > 0 {
					coord.X += 1
					dir = "R"
				} else if coord.X-prev.X > 0 {
					coord.X -= 1
					dir = "L"
				}
				if prev.Y-coord.Y > 0 {
					coord.Y += 1
					dir = fmt.Sprintf("%sU", dir)
				} else if coord.Y-prev.Y > 0 {
					coord.Y -= 1
					dir = fmt.Sprintf("%sD", dir)
				}
				log.Debugf("    coords[%d] MOVED %-2s : POS: %s\n", c, dir, coord.String())
			} else {
				break
			}
			if c == len(coords)-1 {
				(*history)[coord.String()] = true
			}
			prev = coord
		}
	}
}

func touching(c1, c2 coordinate) bool {
	t := util.Abs(c1.X-c2.X) <= 1 && util.Abs(c1.Y-c2.Y) <= 1
	log.Debugf("\tX dist: %d, Y dist: %d (touching: %v)\n", util.Abs(c1.X-c2.X), util.Abs(c1.Y-c2.Y), t)
	return t
}
