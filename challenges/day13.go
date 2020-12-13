package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day13_part1(input []string, isExample bool) (string, error) {
	time, buses := day13_parseInput(input, false)
	log.Debugf("TIME: %d, BUSES: ", time)
	for _, b := range buses {
		log.Debugf("%d ", b)
	}
	log.Debugln()

	arrivals := make(map[int]int)

	for _, b := range buses {
		var t int
		for t = b; t < time; t += b {
		}
		arrivals[b] = t
	}

	firstBusId := -1
	firstBusTime := -1
	for b, t := range arrivals {
		if firstBusTime == -1 || t < firstBusTime {
			firstBusTime = t
			firstBusId = b
		}
	}

	wait := firstBusTime - time
	log.Debugf("FIRST BUS (id: %d) ARRIVES: %d (%d mins after %d)\n", firstBusId, firstBusTime, wait, time)
	return fmt.Sprintf("%d", wait*firstBusId), nil
}

func day13_part2(input []string, isExample bool) (string, error) {
	_, buses := day13_parseInput(input, true)

	overlaps := make([][]int, len(buses), len(buses))

	for i, b := range buses {
		if overlaps[i] == nil {
			overlaps[i] = make([]int, 0)
		}
		if b == 0 {
			continue
		}
		overlaps[i] = append(overlaps[i], b)
		for bb := i - b; bb >= 0; bb -= b {
			if overlaps[bb] == nil {
				overlaps[bb] = make([]int, 0)
			}
			overlaps[bb] = append(overlaps[bb], b)
		}
		for bb := i + b; bb < len(buses); bb += b {
			if overlaps[bb] == nil {
				overlaps[bb] = make([]int, 0)
			}
			overlaps[bb] = append(overlaps[bb], b)
		}
	}

	log.Debugln("Finding overlapping stops:")
	largestMultiple, idxOfLargest := 0, 0
	for i, n := range overlaps {
		if len(n) > 0 {
			log.Debugf("\t%2d - [", i)
			p := 1
			for j, m := range n {
				log.Debugf("%d", m)
				if j != len(n)-1 {
					log.Debug(",")
				}
				p *= m
			}
			if len(n) > 1 {
				log.Debugf("] = %d\n", p)
			} else {
				log.Debugln("]")
			}
			if p > largestMultiple {
				largestMultiple, idxOfLargest = p, i
			}
		} else {
			log.Debugf("\t%2d - x\n", i)
		}
	}

	log.Debugf("\nLargest multiple is %d, at index %d\n\n", largestMultiple, idxOfLargest)

	t := largestMultiple - idxOfLargest
	for {
		match := true
		for i, b := range buses {
			t2 := t + i
			if b == 0 {
				continue
			}
			if t2%b != 0 {
				match = false
				break
			}
		}
		if match {
			break
		}
		t += largestMultiple
	}
	return fmt.Sprintf("%d", t), nil
}

func day13_parseInput(input []string, xsAsZeroes bool) (int, []int) {
	time, _ := strconv.Atoi(input[0])
	buses := make([]int, 0)
	nums := strings.Split(input[1], ",")
	for _, n := range nums {
		if n == "x" {
			if xsAsZeroes {
				buses = append(buses, 0)
			}
		} else {
			i, _ := strconv.Atoi(n)
			buses = append(buses, i)
		}
	}
	return time, buses
}

func init() {
	registerChallengeFunc(13, 1, "day13.txt", day13_part1)
	registerChallengeFunc(13, 2, "day13.txt", day13_part2)
}
