package challenges

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day7_part1(input []string) (string, error) {
	// return day7_part1_simple(input)
	return day7_part1_efficient(input)
}

func day7_part1_efficient(input []string) (string, error) {
	rules := day7_parseRules(input)

	outer := make(map[string]bool)
	log.Debugln()
	for color, _ := range rules {
		if found, known := outer[color]; !known {
			log.Debugf("CHECKING %s:\n", color)
			containsGold := day7_canContainGold_efficient(rules, color, &outer, "\t")
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

func day7_canContainGold_efficient(rules map[string]map[string]int, color string, known *map[string]bool, indent string) bool {
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
				if ok4 := day7_canContainGold_efficient(rules, k, known, fmt.Sprintf("%s\t", indent)); ok4 {
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

func day7_part1_simple(input []string) (string, error) {
	rules := day7_parseRules(input)

	outer := make(map[string]bool)
	log.Debugln()
	for color, _ := range rules {
		log.Debugf("CHECKING %s:\n", color)
		containsGold := day7_canContainGold_simple(rules, color, "\t")
		if containsGold {
			outer[color] = true
		}
	}
	log.Debugf("%v\n", outer)
	return fmt.Sprintf("%d", len(outer)), nil
}

func day7_canContainGold_simple(rules map[string]map[string]int, color string, indent string) bool {
	if contents, ok := rules[color]; ok {
		if _, ok2 := contents["shiny gold"]; ok2 {
			return true
		} else {
			for k := range contents {
				log.Debugf("%sCHECKING %s:\n", indent, k)
				if ok3 := day7_canContainGold_simple(rules, k, fmt.Sprintf("%s\t", indent)); ok3 {
					return true
				}
			}
		}
	}
	return false
}

func day7_part2(input []string) (string, error) {
	rules := day7_parseRules(input)
	total := 0
	for color, count := range rules["shiny gold"] {
		total += count*day7_countInnerBags(color, rules) + count
	}
	return fmt.Sprintf("%d", total), nil
}

func day7_countInnerBags(color string, rules map[string]map[string]int) int {
	total := 0
	if contents, ok := rules[color]; ok {
		for c, count := range contents {
			total += count*day7_countInnerBags(c, rules) + count
		}
	}
	return total
}

func day7_parseRules(input []string) map[string]map[string]int {
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
	registerChallengeFunc(7, 1, "day07.txt", day7_part1)
	registerChallengeFunc(7, 2, "day07.txt", day7_part2)
}
