package aoc2020_day21

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	foods := parseInput(input)

	log.Debugln("FOODS: ")
	for _, f := range foods {
		log.Debugf("  - %v\n", f)
	}

	log.Debug("\n\n")
	suspects := make(map[string][]string)
	culprits := make(map[string]string)

	var loop int
	for loop < 10 && (len(suspects) > 0 || len(culprits) == 0) {
		log.Debugf("LOOP: %d\n", loop)
		log.Debugln("  SUSPECTS: ")
		for a, in := range suspects {
			log.Debugf("    - %-10s: %s\n", a, in)
		}
		log.Debugln("\n  CULPRITS: ")
		for i, a := range culprits {
			log.Debugf("    - %-10s: %s\n", a, i)
		}
		log.Debugln("\n------------------")
		for _, f := range foods {
			for _, a := range f.Allergens {
				if _, ok := suspects[a]; !ok {
					known := false
					for _, al := range culprits {
						if al == a {
							known = true
							break
						}
					}
					if !known {
						suspects[a] = make([]string, len(f.Ingredients))
						copy(suspects[a], f.Ingredients)
					}
				} else {
					if len(suspects[a]) == 1 {
						culprits[suspects[a][0]] = a
						delete(suspects, a)
					} else {
						stillSuspect := make([]string, 0)
						for _, i := range suspects[a] {
							if _, ok := culprits[i]; !ok {
								if f.HasIngredient(i) {
									stillSuspect = append(stillSuspect, i)
								}
							}
						}
						suspects[a] = stillSuspect
					}
				}
			}
		}
		loop++
	}

	harmless := make(map[string]int)
	sum := 0
	for _, f := range foods {
		for _, in := range f.Ingredients {
			if _, ok := culprits[in]; !ok {
				harmless[in] += 1
				sum++
			}
		}
	}

	log.Debugln("HARMLESS: ")
	for i, c := range harmless {
		log.Debugf("  - %-10s: %d\n", i, c)
	}
	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	foods := parseInput(input)

	log.Debugln("FOODS: ")
	for _, f := range foods {
		log.Debugf("  - %v\n", f)
	}

	log.Debug("\n\n")
	suspects := make(map[string][]string)
	culprits := make(map[string]string)

	var loop int
	for loop < 100 && (len(suspects) > 0 || len(culprits) == 0) {
		log.Debugf("LOOP: %d\n", loop)
		log.Debugln("  SUSPECTS: ")
		for a, in := range suspects {
			log.Debugf("    - %-10s: %s\n", a, in)
		}
		log.Debugln("\n  CULPRITS: ")
		for a, i := range culprits {
			log.Debugf("    - %-10s: %s\n", a, i)
		}
		log.Debugln("\n------------------")
		for _, f := range foods {
			for _, a := range f.Allergens {
				if _, ok := suspects[a]; !ok {
					if _, ok2 := culprits[a]; !ok2 {
						suspects[a] = make([]string, len(f.Ingredients))
						copy(suspects[a], f.Ingredients)
					}
				} else {
					if len(suspects[a]) == 1 {
						culprits[a] = suspects[a][0]
						delete(suspects, a)
					} else {
						stillSuspect := make([]string, 0)
						for _, i := range suspects[a] {
							culprit := false
							for _, in := range culprits {
								if in == i {
									culprit = true
									break
								}
							}
							if !culprit && f.HasIngredient(i) {
								stillSuspect = append(stillSuspect, i)
							}
						}
						suspects[a] = stillSuspect
					}
				}
			}
		}
		loop++
	}

	allergens := make([]string, 0, len(culprits))
	for a := range culprits {
		allergens = append(allergens, a)
	}
	sort.Strings(allergens)

	var b strings.Builder
	for i, a := range allergens {
		b.WriteString(culprits[a])
		if i != len(allergens)-1 {
			b.WriteString(",")
		}
	}
	return b.String(), nil
}

func parseInput(input []string) []food {
	foods := make([]food, 0, len(input))

	for _, line := range input {
		halves := strings.Split(line, " (contains ")
		halves[1] = strings.TrimSuffix(halves[1], ")")
		foods = append(foods, food{
			Ingredients: strings.Split(halves[0], " "),
			Allergens:   strings.Split(halves[1], ", "),
		})
	}
	return foods
}

type food struct {
	Ingredients []string
	Allergens   []string
}

func (f food) ListsAllergen(a string) bool {
	for _, ag := range f.Allergens {
		if ag == a {
			return true
		}
	}
	return false
}

func (f food) HasIngredient(in string) bool {
	for _, ing := range f.Ingredients {
		if ing == in {
			return true
		}
	}
	return false
}

func init() {
	challenges.RegisterChallengeFunc(2020, 21, 1, "day21.txt", part1)
	challenges.RegisterChallengeFunc(2020, 21, 2, "day21.txt", part2)
}
