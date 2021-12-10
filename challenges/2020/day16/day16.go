package aoc2020_day16

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	fields, _, tickets := parseInput(input)

	log.Debugln("Checking tickets for impossible field values:")
	var debugStr strings.Builder
	sum := 0
	for i, ticket := range tickets {
		for _, value := range ticket {
			var valid bool
			for _, field := range fields {
				if field.RulesApply(value) {
					valid = true
					break
				}
			}
			if !valid {
				sum += value
				fmt.Fprintf(&debugStr, "%s,", log.Colorize(value, log.Red, 0))
			} else {
				fmt.Fprintf(&debugStr, "%d,", value)
			}
		}
		log.Debugf("  %3d:   %s\n", i, debugStr.String()[:debugStr.Len()-1])
		debugStr.Reset()
	}

	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	fields, myTicketVals, tickets := parseInput(input)

	// First, weed out all of the invalid tickets
	validTickets := make([][]int, 0)
	for _, ticket := range tickets {
		invalidIndex := -1
		for i, value := range ticket {
			var valid bool
			for _, field := range fields {
				if field.RulesApply(value) {
					valid = true
					break
				}
			}
			if !valid {
				invalidIndex = i
				break
			}
		}
		if invalidIndex == -1 {
			validTickets = append(validTickets, ticket)
		} else {
			log.Debug("Eliminating invalid ticket: [")
			for i, v := range ticket {
				if i == invalidIndex {
					log.Debug(log.Colorize(v, log.Red, 0))
				} else {
					log.Debugf("%d", v)
				}
				if i != len(ticket)-1 {
					log.Debug(",")
				}
			}
			log.Debugln("]")
		}
	}
	log.Debugln()

	// next, create and populate a map of field => possible ticket indexes
	fieldMap := make(map[string]map[int]bool)
	for field := range fields {
		fieldMap[field] = make(map[int]bool)
		for i := 0; i < len(myTicketVals); i++ {
			fieldMap[field][i] = true
		}
	}

	// next, for each field, eliminate ticket indexes that have any values that break that field's rules
	for field, idxMap := range fieldMap {
		for i := range idxMap {
			for _, ticket := range validTickets {
				if !fields[field].RulesApply(ticket[i]) {
					delete(fieldMap[field], i)
					break
				}
			}
		}
	}

	// lastly, assign a ticket index to each field, based on an iterative process of elimination
	final := make([]string, len(myTicketVals))
	for len(fieldMap) > 0 {
		for field, idxMap := range fieldMap {
			if len(idxMap) == 1 {
				for idx := range idxMap {
					final[idx] = field
					log.Debugf("Found a 1:1 match: Field %-20s is at index: %2d\n", fmt.Sprintf("'%s'", field), idx)
				}
				delete(fieldMap, field)
			} else {
				for idx := range idxMap {
					for i, name := range final {
						if i == idx && name != "" {
							delete(fieldMap[field], i)
							break
						}
					}
				}
			}
		}
	}

	// now, multiply together all the values on "my ticket" for fields that start with "departure"
	log.Debugln("\nTicket Fields, in order:\n------------------------")
	product := 1
	for idx, name := range final {
		if strings.HasPrefix(name, "departure") {
			log.Debugf("%2d. %-31s (value = %s)\n", idx, strings.ReplaceAll(name, "departure", log.Colorize("departure", log.Green, 0)), log.Colorize(myTicketVals[idx], log.Green, 3))
			product *= myTicketVals[idx]
		} else {
			log.Debugf("%2d. %-20s (value = %3d)\n", idx, name, myTicketVals[idx])
		}
	}

	return fmt.Sprintf("%d", product), nil
}

type field struct {
	Name  string
	Rules []rule
}

func (f field) RulesApply(n int) bool {
	for _, rule := range f.Rules {
		if rule.CheckValue(n) {
			return true
		}
	}
	return false
}

type rule struct {
	Min int
	Max int
}

func (r rule) CheckValue(n int) bool {
	return n >= r.Min && n <= r.Max
}

func parseInput(input []string) (map[string]field, []int, [][]int) {
	var i int
	fields := make(map[string]field)
	log.Debugln("Ticket Fields and Their Rules:")
	for {
		if input[i] == "" {
			i++
			break
		}

		rule1, rule2 := rule{}, rule{}
		field := field{}

		parts := strings.Split(input[i], ": ")
		field.Name = parts[0]
		fmt.Sscanf(parts[1], "%d-%d or %d-%d", &rule1.Min, &rule1.Max, &rule2.Min, &rule2.Max)
		field.Rules = append(field.Rules, rule1, rule2)
		fields[field.Name] = field
		i++
		log.Debugf("\t%-20s : %d-%d or %d-%d\n", field.Name, field.Rules[0].Min, field.Rules[0].Max, field.Rules[1].Min, field.Rules[1].Max)
	}
	log.Debugln()

	i++ // skip the "your ticket:" header

	myTicketVals := make([]int, 0)
	nums := strings.Split(input[i], ",")
	for j := 0; j < len(nums); j++ {
		n, _ := strconv.Atoi(nums[j])
		myTicketVals = append(myTicketVals, n)
	}

	i += 3 // skip the empty line and the "nearby tickets:" header

	tickets := make([][]int, 0)
	for ; i < len(input); i++ {
		nums = strings.Split(input[i], ",")
		ticket := make([]int, 0, len(nums))
		for j := 0; j < len(nums); j++ {
			n, _ := strconv.Atoi(nums[j])
			ticket = append(ticket, n)
		}
		tickets = append(tickets, ticket)
	}

	return fields, myTicketVals, tickets
}

func init() {
	challenges.RegisterChallengeFunc(2020, 16, 1, "day16.txt", part1)
	challenges.RegisterChallengeFunc(2020, 16, 2, "day16.txt", part2)
}
