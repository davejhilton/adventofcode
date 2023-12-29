package aoc2023_day19

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Rule struct {
	Attribute   string
	Condition   string
	Value       int
	Destination string
}

type Workflow struct {
	Label string
	Rules []Rule
}

type Part struct {
	X int
	M int
	A int
	S int
}

type Range struct {
	Start int
	End   int
}
type PartRanges map[string]Range

func (pr PartRanges) Clone() PartRanges {
	newRange := make(PartRanges)
	for k, v := range pr {
		newRange[k] = Range{v.Start, v.End}
	}
	return newRange
}

func (p Part) String() string {
	return fmt.Sprintf("{x=%d, m=%d, a=%d, s=%d}", p.X, p.M, p.A, p.S)
}

func (w Workflow) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s {\n", w.Label))
	for _, rule := range w.Rules {
		if rule.Attribute == "" {
			sb.WriteString(fmt.Sprintf("  -> %s\n", rule.Destination))
		} else {
			sb.WriteString(fmt.Sprintf("  %s %s %d -> %s\n", rule.Attribute, rule.Condition, rule.Value, rule.Destination))
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (w Workflow) Evaluate(p Part) string {
	for _, rule := range w.Rules {
		if rule.Attribute == "" {
			return rule.Destination
		}
		switch rule.Attribute {
		case "a":
			if rule.Condition == "<" && p.A < rule.Value {
				return rule.Destination
			}
			if rule.Condition == ">" && p.A > rule.Value {
				return rule.Destination
			}
		case "m":
			if rule.Condition == "<" && p.M < rule.Value {
				return rule.Destination
			}
			if rule.Condition == ">" && p.M > rule.Value {
				return rule.Destination
			}
		case "s":
			if rule.Condition == "<" && p.S < rule.Value {
				return rule.Destination
			}
			if rule.Condition == ">" && p.S > rule.Value {
				return rule.Destination
			}
		case "x":
			if rule.Condition == "<" && p.X < rule.Value {
				return rule.Destination
			}
			if rule.Condition == ">" && p.X > rule.Value {
				return rule.Destination
			}
		}
	}
	return ""
}

func EvaluateRanges(ranges PartRanges, curWorkflowName string, workflows map[string]Workflow) int {

	if curWorkflowName == "A" {
		product := 1
		for _, r := range ranges {
			product *= r.End - r.Start + 1
		}
		return product
	} else if curWorkflowName == "R" {
		return 0
	}
	w := workflows[curWorkflowName]
	total := 0
	curRanges := ranges.Clone()
	for _, rule := range w.Rules {
		if rule.Attribute == "" {
			total += EvaluateRanges(curRanges, rule.Destination, workflows)
			break
		}
		switch rule.Condition {
		case "<":
			r := curRanges[rule.Attribute]
			if r.Start > rule.Value {
				// no "true" values for this range
				continue
			}
			if r.End < rule.Value {
				// all values are true
				// recurse
				total += EvaluateRanges(curRanges, rule.Destination, workflows)
				break
			}
			// recurse with true range
			trueRanges := curRanges.Clone()
			trueRanges[rule.Attribute] = Range{r.Start, rule.Value - 1}
			total += EvaluateRanges(trueRanges, rule.Destination, workflows)
			// continue looping with false range
			curRanges[rule.Attribute] = Range{rule.Value, r.End}
		case ">":
			r := curRanges[rule.Attribute]
			if r.End < rule.Value {
				// no "true" values for this range
				continue
			}
			if r.Start > rule.Value {
				// all values are true
				// recurse
				total += EvaluateRanges(curRanges, rule.Destination, workflows)
				break
			}
			// recurse with true range
			trueRanges := curRanges.Clone()
			trueRanges[rule.Attribute] = Range{rule.Value + 1, r.End}
			total += EvaluateRanges(trueRanges, rule.Destination, workflows)
			// continue looping with false range
			curRanges[rule.Attribute] = Range{r.Start, rule.Value}
		}
	}

	return total
}

func part1(input []string) (string, error) {
	workflows, parts := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n%v\n", workflows, parts)

	accepted := make([]Part, 0)
	for _, part := range parts {
		curWf := workflows["in"]
		dest := curWf.Evaluate(part)
		for dest != "A" && dest != "R" {
			curWf = workflows[dest]
			dest = curWf.Evaluate(part)
		}
		if dest == "A" {
			accepted = append(accepted, part)
		}
	}

	var result int

	for _, part := range accepted {
		result += part.X + part.M + part.A + part.S
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	workflows, _ := parseInput(input)

	ranges := PartRanges{
		"x": Range{Start: 1, End: 4000},
		"m": Range{Start: 1, End: 4000},
		"a": Range{Start: 1, End: 4000},
		"s": Range{Start: 1, End: 4000},
	}

	var result = EvaluateRanges(ranges, "in", workflows)
	return fmt.Sprintf("%d", result), nil
}

var wfRegex = regexp.MustCompile(`^([a-zA-Z]+)\{(.*)\}$`)
var ruleTestRegex = regexp.MustCompile(`^([amsx])([<>])([0-9]+)$`)
var partRegex = regexp.MustCompile(`^\{x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)\}$`)

func parseInput(input []string) (map[string]Workflow, []Part) {
	workflows := make(map[string]Workflow)
	parts := make([]Part, 0)
	i := 0
	for i < len(input) && input[i] != "" {
		x := wfRegex.FindAllStringSubmatch(input[i], -1)[0]
		label, rulesStrs := x[1], strings.Split(x[2], ",")
		rules := make([]Rule, 0, len(rulesStrs))
		for _, ruleStr := range rulesStrs {
			ruleParts := strings.Split(ruleStr, ":")
			if len(ruleParts) == 2 {
				testParts := ruleTestRegex.FindAllStringSubmatch(ruleParts[0], -1)[0]
				rules = append(rules, Rule{
					Attribute:   testParts[1],
					Condition:   testParts[2],
					Value:       util.Atoi(testParts[3]),
					Destination: ruleParts[1],
				})
			} else {
				rules = append(rules, Rule{
					Destination: ruleParts[0],
				})
			}
		}
		workflows[label] = Workflow{
			Label: label,
			Rules: rules,
		}
		i++
	}
	i++

	for i < len(input) {
		x := partRegex.FindAllStringSubmatch(input[i], -1)[0]
		parts = append(parts, Part{
			X: util.Atoi(x[1]),
			M: util.Atoi(x[2]),
			A: util.Atoi(x[3]),
			S: util.Atoi(x[4]),
		})
		i++
	}

	return workflows, parts
}

func init() {
	challenges.RegisterChallengeFunc(2023, 19, 1, "day19.txt", part1)
	challenges.RegisterChallengeFunc(2023, 19, 2, "day19.txt", part2)
}
