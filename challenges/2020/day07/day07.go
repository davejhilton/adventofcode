package aoc2020_day7

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	// return part1_simple(input)
	return part1_efficient(input)
}

func part1_efficient(input []string) (string, error) {
	rules := parseRules(input)

	outer := make(map[string]bool)
	log.Debugln()
	for color := range rules {
		if found, known := outer[color]; !known {
			log.Debugf("CHECKING %s:\n", color)
			containsGold := canContainGold_efficient(rules, color, &outer, "\t")
			log.Debugf("setting known #0: %s = %v\n", color, containsGold)
			outer[color] = containsGold
		} else {
			log.Debugf("already known %s: %v\n", color, found)
		}
	}
	count := 0
	log.Debugf("\n\n")
	for k, v := range outer {
		log.Debugf("%s: %v\n", k, v)
		if v {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func canContainGold_efficient(rules map[string]map[string]int, color string, known *map[string]bool, indent string) bool {
	if found, ok := (*known)[color]; ok {
		log.Debugf("%salready known %s: %v\n", indent, color, found)
		return found
	}
	if contents, ok := rules[color]; ok {
		if _, ok2 := contents["shiny gold"]; ok2 {
			return true
		} else {
			for k := range contents {
				log.Debugf("%sCHECKING %s:\n", indent, k)
				if found, ok3 := (*known)[k]; ok3 {
					log.Debugf("%s\talready known %s: %v\n", indent, k, found)
					if found {
						return found
					}
				}
				if ok4 := canContainGold_efficient(rules, k, known, fmt.Sprintf("%s\t", indent)); ok4 {
					return true
				}
			}
			log.Debugf("%ssetting known #2: %s = false\n", indent, color)
			(*known)[color] = false
			return false
		}
	}
	log.Debugf("%ssetting known #3: %s = false\n", indent, color)
	(*known)[color] = false
	return false
}

//lint:ignore U1000 this code is here as an example of the "simple" way to do this
func part1_simple(input []string) (string, error) {
	rules := parseRules(input)

	outer := make(map[string]bool)
	log.Debugln()
	for color := range rules {
		log.Debugf("CHECKING %s:\n", color)
		containsGold := canContainGold_simple(rules, color, "\t")
		if containsGold {
			outer[color] = true
		}
	}
	log.Debugf("%v\n", outer)
	return fmt.Sprintf("%d", len(outer)), nil
}

//lint:ignore U1000 this code is here as an example of the "simple" way to do this
func canContainGold_simple(rules map[string]map[string]int, color string, indent string) bool {
	if contents, ok := rules[color]; ok {
		if _, ok2 := contents["shiny gold"]; ok2 {
			return true
		} else {
			for k := range contents {
				log.Debugf("%sCHECKING %s:\n", indent, k)
				if ok3 := canContainGold_simple(rules, k, fmt.Sprintf("%s\t", indent)); ok3 {
					return true
				}
			}
		}
	}
	return false
}

func part2(input []string) (string, error) {
	rules := parseRules(input)
	total := 0
	for color, count := range rules["shiny gold"] {
		total += count*countInnerBags(color, rules) + count
	}
	return fmt.Sprintf("%d", total), nil
}

func countInnerBags(color string, rules map[string]map[string]int) int {
	total := 0
	if contents, ok := rules[color]; ok {
		for c, count := range contents {
			total += count*countInnerBags(c, rules) + count
		}
	}
	return total
}

func parseRules(input []string) map[string]map[string]int {
	rules := make(map[string]map[string]int)
	re := regexp.MustCompile(`^(no|[0-9]+)\s(.*)\sbags?\.?`)
	for _, rule := range input {
		halves := strings.Split(rule, " bags contain ")
		rules[halves[0]] = make(map[string]int)
		contains := strings.Split(halves[1], ", ")
		for _, c := range contains {
			matches := re.FindStringSubmatch(c)
			if len(matches) > 2 && matches[1] != "no" {
				n, _ := strconv.Atoi(matches[1])
				bag := matches[2]
				rules[halves[0]][bag] = n
			}
		}
		log.Debugf("%s - %q\n", halves[0], rules[halves[0]])
	}
	return rules
}

func init() {
	challenges.RegisterChallengeFunc(2020, 7, 1, "day07.txt", part1)
	challenges.RegisterChallengeFunc(2020, 7, 2, "day07.txt", part2)
}
