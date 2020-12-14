package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day14_part1(input []string) (string, error) {

	mem := make(map[int]int64)
	mask := strings.TrimPrefix(input[0], "mask = ")
	for _, line := range input[1:] {
		if strings.HasPrefix(line, "mask") {
			mask = strings.TrimPrefix(line, "mask = ")
		} else {
			i, v := day14_parseMemVal(line)

			binString := []rune(day14_intToBinString(v))
			for j := len(mask) - 1; j >= 0; j-- {
				if mask[j] != 'X' {
					binString[j] = rune(mask[j])
				}
			}
			n, _ := strconv.ParseInt(string(binString), 2, 64)
			log.Debugf("%s - %d\n", string(binString), n)
			mem[i] = n
		}
	}

	sum := int64(0)
	for _, n := range mem {
		sum += n
	}
	return fmt.Sprintf("%d", sum), nil
}

func day14_part2(input []string) (string, error) {
	mem := make(map[int64]int)
	mask := strings.TrimPrefix(input[0], "mask = ")
	for _, line := range input[1:] {
		log.Debugln(line)
		if strings.HasPrefix(line, "mask") {
			mask = strings.TrimPrefix(line, "mask = ")
			log.Debugln("  set new mask")
		} else {
			i, v := day14_parseMemVal(line)

			day14_applyPart2Mask(i, mask, &mem, v)
		}
	}

	sum := 0
	for _, n := range mem {
		sum += n
	}
	return fmt.Sprintf("%d", sum), nil
}

func day14_parseMemVal(m string) (int, int) {
	m = strings.TrimPrefix(m, "mem[")
	parts := strings.SplitN(m, "]", 2)
	i, _ := strconv.Atoi(parts[0])
	n, _ := strconv.Atoi(strings.TrimPrefix(parts[1], " = "))

	return i, n
}

func day14_intToBinString(n int) string {
	pad := "000000000000000000000000000000000000"
	str := fmt.Sprintf("%s%s", pad, strconv.FormatInt(int64(n), 2))
	return str[len(str)-36:]
}

func day14_applyPart2Mask(idx int, mask string, mem *map[int64]int, v int) {
	binString := []rune(day14_intToBinString(idx))
	log.Debugf("  %d - %s\n", idx, string(binString))
	for i := len(mask) - 1; i >= 0; i-- {
		if mask[i] != '0' {
			binString[i] = rune(mask[i])
		}
	}

	day14_storeAllPermutations(string(binString), 0, mem, v)
}

func day14_storeAllPermutations(str string, idx int, mem *map[int64]int, v int) {
	input := []rune(str)
	floated := false
	for i := idx; i < len(input); i++ {
		if input[i] == 'X' {
			input[i] = '1'
			day14_storeAllPermutations(string(input), i+1, mem, v)
			input[i] = '0'
			day14_storeAllPermutations(string(input), i+1, mem, v)
			floated = true
		}
	}
	if !floated {
		n, _ := strconv.ParseInt(str, 2, 64)
		(*mem)[n] = v
	}
}

func init() {
	registerChallengeFunc(14, 1, "day14.txt", day14_part1)
	registerChallengeFunc(14, 2, "day14.txt", day14_part2)
}
