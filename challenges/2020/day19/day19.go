package aoc2020_day19

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	count := countMatches(input, false)
	return fmt.Sprintf("%d", count), nil
}

func part2(input []string) (string, error) {
	count := countMatches(input, true)
	return fmt.Sprintf("%d", count), nil
}

func countMatches(input []string, checkSpecialRules bool) int {
	rules, messages := parseInput(input)

	possibilities := expandRule(rules[0], rules, 0)

	log.Debugln("Sorting...")
	t1 := time.Now()
	sort.Strings(possibilities)
	log.Debugf("Sorted %d items in %v\n", len(possibilities), time.Since(t1))
	inf := ""
	if checkSpecialRules {
		inf = "(well, infinite)"
	}
	log.Debugf("\n---------\nChecking %d messages against %d%s possibilities:\n---------\n\n", len(messages), len(possibilities), inf)

	count := 0
	for _, m := range messages {
		if i := sort.SearchStrings(possibilities, m); i < len(possibilities) && possibilities[i] == m {
			log.Debugf("MATCH! message: %s\n", m)
			count++
		} else if checkSpecialRules && checkSpecial(m) {
			log.Debugln(log.Colorize(fmt.Sprintf("SPECIAL MATCH! message: %s", m), log.Green, 0))
			count++
		} else {
			log.Debugln(log.Colorize(fmt.Sprintf("NO MATCH! message: %s", m), log.Red, 0))
		}
	}
	return count
}

func checkSpecial(message string) bool {
	log.Debugf("CHECKING SPECIAL: %s\n", message)
	if len(message) < 24 {
		return false
	}

	if !check42(string(message[:8])) {
		log.Debugf("\tfirst 8 don't match: %s\n", message[:8])
		return false
	}
	log.Debugf("\tmessage[%d:%d] matches 42: %s\n", 0, 8, message[0:8])
	i := 8
	count42s := 0
	count31s := 0
	// ooo := false
	// matches := make([]int, 0)
	for i+8 <= len(message) {
		j := i + 8
		if check42(string(message[i:j])) {
			log.Debugf("\tmessage[%d:%d] matches 42: %s\n", i, j, message[i:j])
			if count31s > 0 {
				log.Debugln("\t...but already found a 31")
				// ooo = true
				return false
			}
			count42s++
			// matches = append(matches, 42)
		} else if check31(string(message[i:j])) {
			log.Debugf("\tmessage[%d:%d] matches 31: %s\n", i, j, message[i:j])
			count31s++
		} else {
			log.Debugf("\tmessage[%d:%d] doesn't match: %s\n", i, j, message[i:j])
			return false
		}
		i += 8
	}

	log.Debugf("\tfound %d 42s and %d 31s!\n", count42s, count31s)
	return count31s > 0 && count31s <= count42s //&& !ooo
}

func check42(message string) bool {
	if len(message) != 8 {
		return false
	}
	for _, m := range resultsCache[42] {
		if message == m {
			return true
		}
	}
	return false
}

func check31(message string) bool {
	if len(message) != 8 {
		return false
	}
	for _, m := range resultsCache[31] {
		if message == m {
			return true
		}
	}
	return false
}

func parseInput(input []string) (map[int]rule, []string) {
	rules := make(map[int]rule)
	var i int
	for i = 0; i < len(input); i++ {
		if input[i] == "" {
			i++
			break
		}
		rule := parseRule(input[i])
		rules[rule.Id] = rule
	}

	messages := input[i:]

	return rules, messages
}

func expandSequence(sequence sequence, rules map[int]rule, depth int) []string {
	log.Debugf("%*sEXPANDING SEQ:  %s\n", 2*depth, "", sequence.String())
	str := ""
	results := make([]string, 0)
	recursed := false
	for _, t := range sequence {
		if t.Type == "STRING" {
			str += t.StringVal
			if len(results) > 0 {
				for i, r := range results {
					log.Debugf("%*s-appending result: %s%s\n", 2*depth, "", r, t.StringVal)
					results[i] = fmt.Sprintf("%s%s", r, t.StringVal)
				}
			}
		} else {
			options := expandRule(rules[t.Ref], rules, depth+1)
			newResults := make([]string, 0)
			for _, o := range options {
				if len(results) != 0 {
					for _, r := range results {
						// log.Debugf("%*s-adding result: %s%s\n", 2*depth, "", r, o)
						newResults = append(newResults, fmt.Sprintf("%s%s", r, o))
					}
				} else {
					newResults = append(newResults, o)
				}
			}
			results = newResults
			recursed = true
		}
	}
	if !recursed {
		if str != "" {
			results = append(results, str)
			// log.Debugf("%*sadding direct result: %s\n", 2*depth, "", str)
		}
	}
	return results
}

func expandRule(rule rule, rules map[int]rule, depth int) []string {
	if results, ok := resultsCache[rule.Id]; ok {
		log.Debugf("%*sUSING RULE CACHE: %s - \n", 2*depth, "", rule.String())
		return results
	}
	log.Debugf("%*sEXPANDING RULE: %s\n", 2*depth, "", rule.String())
	results := make([]string, 0)
	for _, seq := range rule.Sequences {
		results = append(results, expandSequence(seq, rules, depth+1)...)
	}
	if rule.Id != 0 && rule.Id != 8 && rule.Id != 11 {
		resultsCache[rule.Id] = results
	}
	return results
}

func parseRule(line string) rule {
	parts := strings.Split(line, ": ")
	id, _ := strconv.Atoi(parts[0])
	seqStrs := strings.Split(parts[1], " | ")
	rule := rule{
		Id:        id,
		Sequences: make([]sequence, 0, len(seqStrs)),
	}
	log.Debugf("line: '%s'\n", line)
	for _, str := range seqStrs {
		// log.Debugf("  seq string: '%s'\n", str)
		tokStrs := strings.Split(str, " ")
		seq := make(sequence, 0, len(tokStrs))
		for _, tokStr := range tokStrs {
			if tokStr[0] == '"' {
				// log.Debugf("    tok STRING: '%s'\n", tokStr)
				seq = append(seq, token{
					Type:      "STRING",
					StringVal: string(tokStr[1 : len(tokStr)-1]),
				})
			} else {
				// log.Debugf("    tok REF: '%s'\n", tokStr)
				n, _ := strconv.Atoi(tokStr)
				seq = append(seq, token{
					Type: "REF",
					Ref:  n,
				})
			}
		}
		rule.Sequences = append(rule.Sequences, seq)
	}
	return rule
}

type rule struct {
	Id        int
	Sequences []sequence
}

func (r rule) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d: ", r.Id))

	for i, s := range r.Sequences {
		if i != 0 {
			b.WriteString(" | ")
		}
		b.WriteString(s.String())
	}
	return b.String()
}

type sequence []token

func (s sequence) String() string {
	var b strings.Builder
	for i, t := range s {
		if i != 0 {
			b.WriteString(" ")
		}
		b.WriteString(t.String())
	}
	return b.String()
}

type token struct {
	Type      string
	StringVal string
	Ref       int
}

func (t token) String() string {
	if t.Type == "STRING" {
		return fmt.Sprintf(`"%s"`, t.StringVal)
	} else {
		return fmt.Sprintf("%d", t.Ref)
	}
}

func init() {
	challenges.RegisterChallengeFunc(2020, 19, 1, "day19.txt", part1)
	challenges.RegisterChallengeFunc(2020, 19, 2, "day19.txt", part2)

	resultsCache = make(map[int][]string)
}

var resultsCache map[int][]string
