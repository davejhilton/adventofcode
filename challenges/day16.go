package challenges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode2020/log"
)

func day16_part1(input []string) (string, error) {
	rules, myTicketVals, nearbyTickets := day16_parseInput(input)

	for _, r := range rules {
		log.Debugf("%s: %d-%d or %d-%d\n", r.Name, r.Min1, r.Max1, r.Min2, r.Max2)
	}

	for _, n := range myTicketVals {
		log.Debugf("%d, ", n)
	}
	log.Debugln("\n")
	for _, t := range nearbyTickets {
		for _, n := range t {
			log.Debugf("%d, ", n)
		}
		log.Debugln()
	}

	invalidSum := 0
	for _, ticket := range nearbyTickets {

		for _, n := range ticket {
			var valid bool
			for _, rule := range rules {
				if rule.CheckValue(n) {
					valid = true
					break
				}
			}
			if !valid {
				log.Debugf("invalid value: %d\n", n)
				invalidSum += n
			}
		}
	}

	return fmt.Sprintf("%d", invalidSum), nil
}

func day16_part2(input []string) (string, error) {
	rules, myTicketVals, nearbyTickets := day16_parseInput(input)

	for _, r := range rules {
		log.Debugf("%s: %d-%d or %d-%d\n", r.Name, r.Min1, r.Max1, r.Min2, r.Max2)
	}

	for _, n := range myTicketVals {
		log.Debugf("%d, ", n)
	}
	log.Debugln("\n")
	for _, t := range nearbyTickets {
		for _, n := range t {
			log.Debugf("%d, ", n)
		}
		log.Debugln()
	}

	validTickets := make([][]int, 0)
	for _, ticket := range nearbyTickets {
		ticketValid := true
	outer:
		for _, n := range ticket {
			var valid bool
			for _, rule := range rules {
				if rule.CheckValue(n) {
					valid = true
					break
				}
			}
			if !valid {
				ticketValid = false
				break outer
			}
		}
		if ticketValid {
			validTickets = append(validTickets, ticket)
		}
	}

	possibleIdxs := make(map[string][]int)

	for i, n := range myTicketVals {
		for _, rule := range rules {
			if match := rule.CheckValue(n); match {
				list, ok := possibleIdxs[rule.Name]
				if !ok {
					list = make([]int, 0)
				}
				list = append(list, i)
				possibleIdxs[rule.Name] = list
			}
		}
	}

	matches := make(map[string][]int)
	for ruleName, idxs := range possibleIdxs {

		stillValidIdxs := make([]int, 0)

		for _, i := range idxs {
			valid := true
			for _, ticket := range validTickets {
				if !rules[ruleName].CheckValue(ticket[i]) {
					valid = false
					break
				}
			}
			if valid {
				stillValidIdxs = append(stillValidIdxs, i)
			}
		}
		matches[ruleName] = stillValidIdxs

	}

	idxMatches := make(map[int][]string)
	for ruleName, idxs := range matches {
		for _, idx := range idxs {
			if _, ok := idxMatches[idx]; !ok {
				idxMatches[idx] = make([]string, 0)
			}
			idxMatches[idx] = append(idxMatches[idx], ruleName)
		}
	}

	idxMap := make(map[int]string)
	for len(idxMatches) > 0 && len(matches) > 0 {
		for idx, names := range idxMatches {
			if len(names) == 1 {
				idxMap[idx] = names[0]
				delete(idxMatches, idx)
				delete(matches, names[0])
				log.Debugf("FOUND ONE: %d = %s\n", idx, names[0])
			} else {
				unmatched := make([]string, 0)
				for _, name := range names {
					found := false
					for _, n := range idxMap {
						if n == name {
							found = true
							break
						}
					}
					if !found {
						unmatched = append(unmatched, name)
					}
				}
				idxMatches[idx] = unmatched
			}
		}
	}

	log.Debugln()

	product := 1
	for idx, name := range idxMap {
		log.Debugf("%s - myTicketVals[%d] = %d\n", name, idx, myTicketVals[idx])
		if strings.HasPrefix(name, "departure") {
			product *= myTicketVals[idx]
		}
	}

	return fmt.Sprintf("%d", product), nil
}

type day16_rule struct {
	Name string
	Min1 int
	Max1 int
	Min2 int
	Max2 int
}

func (r day16_rule) CheckValue(n int) bool {
	if n >= r.Min1 && n <= r.Max1 {
		return true
	} else if n >= r.Min2 && n <= r.Max2 {
		return true
	} else {
		return false
	}
}

func day16_parseInput(input []string) (map[string]day16_rule, []int, [][]int) {
	var i int

	rules := make(map[string]day16_rule)
	for {
		if input[i] == "" {
			i++
			break
		}

		rule := day16_rule{}

		parts := strings.Split(input[i], ": ")
		rule.Name = parts[0]
		fmt.Sscanf(parts[1], "%d-%d or %d-%d", &rule.Min1, &rule.Max1, &rule.Min2, &rule.Max2)
		rules[rule.Name] = rule
		i++
	}

	i++
	myTicketVals := make([]int, 0)
	nums := strings.Split(input[i], ",")
	for j := 0; j < len(nums); j++ {
		n, _ := strconv.Atoi(nums[j])
		myTicketVals = append(myTicketVals, n)
	}
	nearbyTickets := make([][]int, 0)
	for i = i + 3; i < len(input); i++ {
		nums = strings.Split(input[i], ",")
		ticket := make([]int, 0, len(nums))
		for j := 0; j < len(nums); j++ {
			n, _ := strconv.Atoi(nums[j])
			ticket = append(ticket, n)
		}
		nearbyTickets = append(nearbyTickets, ticket)
	}

	return rules, myTicketVals, nearbyTickets
}

func init() {
	registerChallengeFunc(16, 1, "day16.txt", day16_part1)
	registerChallengeFunc(16, 2, "day16.txt", day16_part2)
}
