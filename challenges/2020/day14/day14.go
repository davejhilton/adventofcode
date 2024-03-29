package aoc2020_day14

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	mem := make(map[uint64]uint64)
	var mask string
	for _, line := range input {
		log.Debugln(line)
		if strings.HasPrefix(line, "mask") {
			mask = strings.TrimPrefix(line, "mask = ")
			log.Debugf("\tset mask: %s\n", mask)
		} else {
			addr, value := parseMemLine(line)
			log.Debugf("\tbinary:   %.*b (before mask)\n", len(mask), value)
			for i := len(mask) - 1; i >= 0; i-- {
				switch mask[i] {
				case '0':
					// use bitwise 'AND NOT' &^ to clear the bit
					value = value &^ (1 << (len(mask) - i - 1))
				case '1':
					// use bitwise 'OR' | to set the bit
					value = value | (1 << (len(mask) - i - 1))
				case 'X':
				default:
				}
			}
			mem[addr] = value
			log.Debugf("\tstoring:  %.*b (%d) at addr %d\n", len(mask), value, value, addr)
		}
	}

	var sum uint64
	for _, n := range mem {
		sum += n
	}
	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	mem := make(map[uint64]uint64)
	var mask string
	for _, line := range input {
		log.Debugln(line)
		if strings.HasPrefix(line, "mask") {
			mask = strings.TrimPrefix(line, "mask = ")
			log.Debugf("\tset new mask: %s\n", mask)
		} else {
			addr, value := parseMemLine(line)
			storeAllPermutations(addr, mask, value, &mem)
		}
	}

	var sum uint64
	for _, n := range mem {
		sum += n
	}
	return fmt.Sprintf("%d", sum), nil
}

func parseMemLine(line string) (uint64, uint64) {
	line = strings.TrimPrefix(line, "mem[")
	parts := strings.SplitN(line, "]", 2)
	addr, _ := strconv.ParseUint(parts[0], 10, 64)
	value, _ := strconv.ParseUint(strings.TrimPrefix(parts[1], " = "), 10, 64)

	return addr, value
}

func storeAllPermutations(addr uint64, mask string, value uint64, mem *map[uint64]uint64) {
	log.Debugln("\trecursing:")
	log.Debugf("\t\taddr: %.8b - %d\n", addr, addr)
	log.Debugf("\t\tmask: %8s\n", mask)
	floated := false
	for i := 0; i < len(mask); i++ {
		if mask[i] == 'X' {
			// bitwise 'OR' | to set the bit
			addr = addr | (1 << (len(mask) - i - 1))
			storeAllPermutations(addr, mask[i+1:], value, mem)
			// bitwise 'AND NOT' &^ to clear the bit
			addr = addr &^ (1 << (len(mask) - i - 1))
			storeAllPermutations(addr, mask[i+1:], value, mem)
			floated = true
		} else if mask[i] == '1' {
			// bitwise 'OR' | to set the bit
			addr = addr | (1 << (len(mask) - i - 1))
		}
	}

	if !floated {
		log.Debugf("\t\tstoring value %d at addr %.8b (%d)\n", value, addr, addr)
		(*mem)[addr] = value
	}
}

func init() {
	challenges.RegisterChallengeFunc(2020, 14, 1, "day14.txt", part1)
	challenges.RegisterChallengeFunc(2020, 14, 2, "day14.txt", part2)
}
