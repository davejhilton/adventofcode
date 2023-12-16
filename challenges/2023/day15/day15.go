package aoc2023_day15

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Lens struct {
	Label         string
	FocusingPower int
	Next          *Lens
}

func (l *Lens) String() string {
	return fmt.Sprintf("%s_%d", l.Label, l.FocusingPower)
}

type Step struct {
	Label string
	Op    string
	Fp    int
}

type Hashmap map[int]*Lens

func HASH(s string) int {
	cur := 0
	for _, c := range s {
		// fmt.Printf("%d\n", c)
		cur = ((cur + int(c)) * 17) % 256
	}
	return cur
}

var reg = regexp.MustCompile(`([a-z]+)([=-])(\d?)`)

func parseStep(s string) Step {
	matches := reg.FindStringSubmatch(s)
	// log.Debugf("Matches: %v\n", strings.Join(matches, ","))
	return Step{
		Label: matches[1],
		Op:    matches[2],
		Fp:    util.Atoi(matches[3]),
	}
}

func part1(input []string) (string, error) {
	steps := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", steps)

	var result int

	for _, s := range steps {
		result += HASH(s)
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	steps := parseInput(input)

	hashmap := make(Hashmap)

	for _, s := range steps {

		step := parseStep(s)
		hash := HASH(step.Label)
		log.Debugf("Step hash: %d\n", hash)
		if boxLens, ok := hashmap[hash]; ok {
			if step.Op == "-" {
				var prev *Lens
				for boxLens != nil {
					if boxLens.Label == step.Label {
						if prev != nil {
							prev.Next = boxLens.Next
						} else {
							hashmap[hash] = boxLens.Next
						}
						break
					}
					prev = boxLens
					boxLens = boxLens.Next
				}
			} else if step.Op == "=" {
				var replaced = false
				var prev *Lens
				for boxLens != nil {
					if boxLens.Label == step.Label {
						boxLens.FocusingPower = step.Fp
						replaced = true
						break
					}
					prev = boxLens
					boxLens = boxLens.Next
				}
				if !replaced {
					if prev == nil {
						hashmap[hash] = &Lens{step.Label, step.Fp, nil}
					} else {
						prev.Next = &Lens{step.Label, step.Fp, nil}
					}
				}
			} else {
				boxLens.FocusingPower += step.Fp
			}
		} else {
			if step.Op == "=" {
				hashmap[hash] = &Lens{step.Label, step.Fp, nil}
			}
			// else do nothing
		}
	}

	var result int

	for boxNum, boxLens := range hashmap {
		i := 1
		for boxLens != nil {
			result += (1 + boxNum) * boxLens.FocusingPower * i
			boxLens = boxLens.Next
			i++
		}
	}

	log.Debugf("Hashmap: %v\n", hashmap)

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []string {
	return strings.Split(strings.Join(input, ""), ",")
}

func init() {
	challenges.RegisterChallengeFunc(2023, 15, 1, "day15.txt", part1)
	challenges.RegisterChallengeFunc(2023, 15, 2, "day15.txt", part2)
}
