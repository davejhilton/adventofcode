package aoc2021_day3

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	binLists := parse(input)

	γ, ε := 0, 0
	count := len(binLists)
	bitCount := len(binLists[0])
	sums := make([]int, bitCount)
	for _, l := range binLists {
		for i, b := range l {
			sums[i] += b
		}
	}

	for i, s := range sums {
		if s > count/2 {
			γ += 1 << (bitCount - i - 1)
		} else {
			ε += 1 << (bitCount - i - 1)
		}
	}

	var result = γ * ε

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	binLists := parse(input)

	o2gen, co2scrub := 0, 0
	bitCount := len(binLists[0])
	o2list, co2list := make([][]int, 0), make([][]int, 0)
	for _, l := range binLists {
		o2list = append(o2list, l)
		co2list = append(co2list, l)
	}

	log.Debugf("Possibilities are:\n%s\n\n\n", doubleSliceToString(o2list))
	o2soFar := ""
	for i := 0; i < bitCount; i += 1 {
		sum := 0
		zeros := make([][]int, 0)
		ones := make([][]int, 0)
		for _, bin := range o2list {
			sum += bin[i]
			if bin[i] == 0 {
				zeros = append(zeros, bin)
			} else {
				ones = append(ones, bin)
			}
		}
		var bit int
		n := len(o2list)
		if float64(sum) >= float64(n)/float64(2.0) {
			bit = 1
			o2list = ones
		} else {
			bit = 0
			o2list = zeros
		}
		o2soFar = fmt.Sprintf("%s%d", o2soFar, bit)
		log.Debugf("Bit [%d]: %d/%d have 1 bit, so most common is '%d'.\n%s\nPossibilities are:\n%s\n\n\n", i, sum, n, bit, o2soFar, doubleSliceToString(o2list))

		if len(o2list) == 1 {
			o2gen = binSliceToDec(o2list[0])
			log.Debugf("O2 Gen value: %d\n", o2gen)
			break
		}
	}

	co2soFar := ""
	for i := 0; i < bitCount; i += 1 {
		sum := 0
		zeros := make([][]int, 0)
		ones := make([][]int, 0)
		for _, bin := range co2list {
			sum += bin[i]
			if bin[i] == 0 {
				zeros = append(zeros, bin)
			} else {
				ones = append(ones, bin)
			}
		}
		var bit int
		n := len(co2list)
		if float64(sum) >= float64(n)/float64(2.0) {
			bit = 0
			co2list = zeros
		} else {
			bit = 1
			co2list = ones
		}
		co2soFar = fmt.Sprintf("%s%d", co2soFar, bit)
		log.Debugf("Bit [%d]: %d/%d have 1 bit, so least common is '%d'.\n%s\nPossibilities are:\n%s\n\n\n", i, sum, n, bit, co2soFar, doubleSliceToString(co2list))

		if len(co2list) == 1 {
			co2scrub = binSliceToDec(co2list[0])
			log.Debugf("CO2 Gen value: %d\n", co2scrub)
			break
		}
	}

	result := o2gen * co2scrub

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func binSliceToDec(bin []int) int {
	r := 0
	bitCount := len(bin)
	for i, d := range bin {
		if d == 1 {
			r += 1 << (bitCount - i - 1)
		}
	}
	return r
}

func doubleSliceToString(bl [][]int) string {
	str := ""
	for i, l := range bl {
		if i > 0 {
			str = fmt.Sprintf("%s\n", str)
		}
		for _, b := range l {
			str = fmt.Sprintf("%s%d", str, b)
		}
	}
	return str
}

func parse(input []string) [][]int {
	digitLists := make([][]int, 0, len(input))
	for _, s := range input {
		dl := make([]int, 0)
		for _, d := range s {
			if d == '1' {
				dl = append(dl, 1)
			} else {
				dl = append(dl, 0)
			}
		}
		digitLists = append(digitLists, dl)
	}
	return digitLists
}

func init() {
	challenges.RegisterChallengeFunc(2021, 3, 1, "day03.txt", part1)
	challenges.RegisterChallengeFunc(2021, 3, 2, "day03.txt", part2)
}
