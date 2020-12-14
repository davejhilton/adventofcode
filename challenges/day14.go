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
	mem := make(map[string]int)
	mask := strings.TrimPrefix(input[0], "mask = ")
	for _, line := range input[1:] {
		log.Debugln(line)
		if strings.HasPrefix(line, "mask") {
			mask = strings.TrimPrefix(line, "mask = ")
			log.Debugln("  set new mask")
		} else {
			i, v := day14_parseMemVal(line)

			idxs := day14_applyPart2Mask(i, mask)
			for _, s := range idxs {
				x, _ := strconv.ParseInt(s, 2, 64)
				log.Debugf("\t%s - %d\n", s, x)
				mem[s] = v
			}
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

func day14_applyPart2Mask(n int, mask string) []string {
	binString := []rune(day14_intToBinString(n))
	log.Debugf("  %d - %s\n", n, string(binString))
	for i := len(mask) - 1; i >= 0; i-- {
		if mask[i] != '0' {
			binString[i] = rune(mask[i])
		}
	}

	return day14_findAllPermutations(string(binString), 0)
}

func day14_findAllPermutations(str string, idx int) []string {
	input := []rune(str)
	res := make([]string, 0)
	floated := false
	for ; idx < len(input); idx++ {
		if input[idx] == 'X' {
			input[idx] = '1'
			res = append(res, day14_findAllPermutations(string(input), idx+1)...)
			input[idx] = '0'
			res = append(res, day14_findAllPermutations(string(input), idx+1)...)
			floated = true
		}
	}
	if !floated {
		res = append(res, str)
	}
	return res
}

func init() {
	registerChallengeFunc(14, 1, "day14.txt", day14_part1)
	registerChallengeFunc(14, 2, "day14.txt", day14_part2)
}
