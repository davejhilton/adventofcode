package aoc2022_day21

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	monkeys := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", monkeys)

	unresolved := make(map[string]*Monkey)
	for name, m := range monkeys {
		if !m.Resolved {
			unresolved[name] = m
		}
	}

	var root *Monkey = monkeys["root"]

	for !root.Resolved {
		for name, m := range unresolved {
			if m.DependsOn[0].Resolved && m.DependsOn[1].Resolved {

				v1 := m.DependsOn[0].Value
				v2 := m.DependsOn[1].Value
				switch m.Operator {
				case "+":
					m.Value = v1 + v2
				case "-":
					m.Value = v1 - v2
				case "*":
					m.Value = v1 * v2
				case "/":
					m.Value = v1 / v2
				default:
					fmt.Printf("UNRECOGNIZED OPERATOR: %s\n", m.Operator)
				}
				m.Resolved = true
				//
				delete(unresolved, name)
			}
		}
	}

	var result int = root.Value
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	monkeys := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", monkeys)

	var root *Monkey = monkeys["root"]
	var humn *Monkey = monkeys["humn"]

	root.Operator = "="
	root.Resolved = false
	humn.Resolved = false
	humn.Operator = "?"
	humn.DependsOn = []*Monkey{}

	unresolved := make(map[string]*Monkey)
	for name, m := range monkeys {
		if !m.Resolved {
			unresolved[name] = m
		}
	}

	var numU int
	for !humn.Resolved {
		numU = len(unresolved)
		for name, m := range unresolved {
			if len(m.DependsOn) > 1 && m.DependsOn[0].Resolved && m.DependsOn[1].Resolved {
				v1 := m.DependsOn[0].Value
				v2 := m.DependsOn[1].Value
				switch m.Operator {
				case "+":
					m.Value = v1 + v2
				case "-":
					m.Value = v1 - v2
				case "*":
					m.Value = v1 * v2
				case "/":
					m.Value = v1 / v2
				case "=":
					if v1 == v2 {
						m.Value = v1
					} else {
						continue
					}
				case "?":
					continue
				default:
					fmt.Printf("UNRECOGNIZED OPERATOR: %s\n", m.Operator)
				}
				m.Resolved = true
				delete(unresolved, name)
			}
		}
		if len(unresolved) == numU {
			log.Debugf("No monkeys resolved in this loop. %d/%d are still unresolved!\n", numU, len(monkeys))
			break
		}
	}

	n := 0
	rootExpression, ok := root.ToExpression()
	for ok {
		log.Debugln(rootExpression.ToString(false))
		rootExpression, ok = rootExpression.Simplify()
		n++
	}

	log.Debugf("\n========\n\nSimplified %d times\n", n)
	log.Debugf("ROOT (simplified): %s\n", rootExpression.ToString(false))

	var result int = rootExpression.GetValue()
	return fmt.Sprintf("%d", result), nil
}

type Monkey struct {
	Name      string
	Resolved  bool
	Operator  string
	Value     int
	DependsOn []*Monkey
	Value1    *int
	Value2    *int
}

func (m *Monkey) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s: ", m.Name))
	if m.Resolved {
		b.WriteString(fmt.Sprintf("%d", m.Value))
	} else if len(m.DependsOn) > 0 {
		n1, n2 := m.DependsOn[0].Name, m.DependsOn[1].Name
		b.WriteString(fmt.Sprintf("%s %s %s", n1, m.Operator, n2))
	} else {
		b.WriteString("???")
	}
	return b.String()
}

func (m *Monkey) ToExpression() (AbstractExpression, bool) {
	if len(m.DependsOn) > 1 {
		var v1 AbstractExpression
		var v2 AbstractExpression
		var v1ok, v2ok bool
		if m.DependsOn[0].Resolved {
			v := m.DependsOn[0].Value
			v1 = (*ConcreteValue)(&v)
			v1ok = true
		} else {
			v1, v1ok = m.DependsOn[0].ToExpression()
		}
		if m.DependsOn[1].Resolved {
			v := m.DependsOn[1].Value
			v2 = (*ConcreteValue)(&v)
			v2ok = true
		} else {
			v2, v2ok = m.DependsOn[1].ToExpression()
		}
		return &Expression{
			Value1:   v1,
			Value2:   v2,
			Operator: m.Operator,
		}, v1ok || v2ok || true
	} else if m.Operator == "?" {
		return &UnknownValue{}, false
	}

	return nil, false
}

type Value struct {
	Resolved  bool
	Value     int
	DependsOn AbstractExpression
}

func (m *Value) IsResolved() bool {
	return m.Resolved
}
func (m *Value) GetValue() int {
	return m.Value
}
func (m *Value) Resolve() (int, bool) {
	return 0, false
}
func (m *Value) Simplify() (AbstractExpression, bool) {
	if m.IsResolved() {
		return (*ConcreteValue)(&m.Value), true
	} else if m.DependsOn != nil && m.DependsOn.IsResolved() {
		m.Value = m.DependsOn.GetValue()
		m.Resolved = true
		return m, true
	} else if m.DependsOn != nil {
		return m.DependsOn, true
	} else {
		return nil, false
	}
}
func (m *Value) ToString(parens bool) string {
	if m.IsResolved() {
		return fmt.Sprintf("%d", m.GetValue())
	} else {
		return m.DependsOn.ToString(parens)
	}
}
func (m *Value) ContainsVariable() bool {
	if m.IsResolved() {
		return false
	} else {
		return m.DependsOn.ContainsVariable()
	}
}
func (m *Value) Rebalance() (AbstractExpression, string, AbstractExpression) {
	if !m.IsResolved() && m.DependsOn.ContainsVariable() {
		return m.DependsOn.Rebalance()
	}
	return nil, "", nil
}
func (m *Value) CanRebalance() bool {
	if !m.IsResolved() {
		return m.DependsOn.CanRebalance()
	} else {
		return false
	}
}
func (m *Value) GetName() string {
	return "Value"
}

type UnknownValue struct{}

func (m *UnknownValue) IsResolved() bool {
	return false
}
func (m *UnknownValue) GetValue() int {
	return -1
}
func (m *UnknownValue) Resolve() (int, bool) {
	return -1, false
}
func (m *UnknownValue) Simplify() (AbstractExpression, bool) {
	return m, false
}
func (m *UnknownValue) ToString(parens bool) string {
	return "?"
}
func (m *UnknownValue) ContainsVariable() bool {
	return true
}
func (m *UnknownValue) Rebalance() (AbstractExpression, string, AbstractExpression) {
	return nil, "", nil
}
func (m *UnknownValue) CanRebalance() bool {
	return false
}
func (m *UnknownValue) GetName() string {
	return "UnknownValue"
}

type ConcreteValue int

func (c *ConcreteValue) IsResolved() bool {
	return true
}
func (c *ConcreteValue) GetValue() int {
	return int(*c)
}
func (c *ConcreteValue) Resolve() (int, bool) {
	return int(*c), true
}
func (c *ConcreteValue) Simplify() (AbstractExpression, bool) {
	return c, false
}
func (m *ConcreteValue) ToString(parens bool) string {
	return fmt.Sprintf("%d", m.GetValue())
}
func (m *ConcreteValue) ContainsVariable() bool {
	return false
}
func (m *ConcreteValue) Rebalance() (AbstractExpression, string, AbstractExpression) {
	return nil, "", nil
}
func (m *ConcreteValue) CanRebalance() bool {
	return false
}
func (m *ConcreteValue) GetName() string {
	return "ConcreteValue"
}

type Expression struct {
	Value1   AbstractExpression
	Value2   AbstractExpression
	Operator string
}

func (m *Expression) IsResolved() bool {
	if m.Value1.IsResolved() && m.Value2.IsResolved() {
		return true
	} else if m.Operator == "=" {
		if _, ok := m.Value1.(*UnknownValue); ok && m.Value2.IsResolved() {
			return true
		}
		if _, ok := m.Value2.(*UnknownValue); ok && m.Value1.IsResolved() {
			return true
		}
	}
	return false
}
func (m *Expression) GetValue() int {
	var value int
	switch m.Operator {
	case "+":
		value = m.Value1.GetValue() + m.Value2.GetValue()
	case "-":
		value = m.Value1.GetValue() - m.Value2.GetValue()
	case "*":
		value = m.Value1.GetValue() * m.Value2.GetValue()
	case "/":
		value = m.Value1.GetValue() / m.Value2.GetValue()
	case "=":
		if m.Value1.GetValue() == m.Value2.GetValue() {
			value = m.Value1.GetValue()
		} else {
			if _, ok := m.Value1.(*UnknownValue); ok && m.Value2.IsResolved() {
				return m.Value2.GetValue()
			}
			if _, ok := m.Value2.(*UnknownValue); ok && m.Value1.IsResolved() {
				return m.Value1.GetValue()
			}
		}
	case "?":
		value = -1
	default:
		fmt.Printf("UNRECOGNIZED OPERATOR: %s\n", m.Operator)
	}
	return value
}
func (m *Expression) Resolve() (int, bool) {
	if m.IsResolved() {
		return m.GetValue(), true
	} else if m.Operator == "=" {
		if _, ok := m.Value1.(*UnknownValue); ok && m.Value2.IsResolved() {
			return m.Value2.GetValue(), true
		}
		if _, ok := m.Value2.(*UnknownValue); ok && m.Value1.IsResolved() {
			return m.Value1.GetValue(), true
		}
	}
	return 0, false
}
func (m *Expression) Simplify() (AbstractExpression, bool) {
	if m.IsResolved() {
		v := m.GetValue()
		return (*ConcreteValue)(&v), true
	}
	v1, ok1 := m.Value1.Simplify()
	v2, ok2 := m.Value2.Simplify()
	m.Value1 = v1
	m.Value2 = v2
	if ok1 || ok2 {
		return m, true
	}

	if m.Operator == "=" {
		if m.Value1.CanRebalance() {
			a, op, b := m.Value1.Rebalance()
			m.Value1 = b
			m.Value2 = &Expression{
				Value1:   m.Value2,
				Operator: op,
				Value2:   a,
			}
			return m, true
		} else if m.Value2.CanRebalance() {
			a, op, b := m.Value2.Rebalance()
			m.Value2 = b
			m.Value1 = &Expression{
				Value1:   m.Value1,
				Operator: op,
				Value2:   a,
			}
			return m, true
		}
	}
	return m, false
}
func (m *Expression) ToString(parens bool) string {
	if m.Operator == "=" {
		return fmt.Sprintf("%s %s %s", m.Value1.ToString(false), m.Operator, m.Value2.ToString(false))
	} else if parens {
		return fmt.Sprintf("(%s %s %s)", m.Value1.ToString(true), m.Operator, m.Value2.ToString(true))
	} else {
		return fmt.Sprintf("%s %s %s", m.Value1.ToString(true), m.Operator, m.Value2.ToString(true))
	}
}
func (m *Expression) ContainsVariable() bool {
	return m.Value1.ContainsVariable() || m.Value2.ContainsVariable()
}
func (m *Expression) Rebalance() (AbstractExpression, string, AbstractExpression) {
	if m.Value1.ContainsVariable() {
		if m.Operator == "+" {
			return m.Value2, "-", m.Value1
		} else if m.Operator == "-" {
			return m.Value2, "+", m.Value1
		} else if m.Operator == "/" {
			return m.Value2, "*", m.Value1
		} else if m.Operator == "*" {
			return m.Value2, "/", m.Value1
		}
	} else {
		if m.Operator == "+" {
			return m.Value1, "-", m.Value2
		} else if m.Operator == "-" {
			v := -1
			e1 := &Expression{Value1: m.Value1, Operator: "*", Value2: (*ConcreteValue)(&v)}
			e2 := &Expression{Value1: m.Value2, Operator: "*", Value2: (*ConcreteValue)(&v)}
			return e1, "+", e2
		} else if m.Operator == "/" {
			v := 1
			e := &Expression{Value1: (*ConcreteValue)(&v), Operator: "/", Value2: m.Value1}
			return e, "*", m.Value2
		} else if m.Operator == "*" {
			return m.Value1, "/", m.Value2
		}
	}
	return nil, "", nil
}
func (m *Expression) CanRebalance() bool {
	return m.ContainsVariable() && m.Operator != ""
}
func (m *Expression) GetName() string {
	return "Expression"
}

type AbstractExpression interface {
	IsResolved() bool
	GetValue() int
	Resolve() (int, bool)
	Simplify() (AbstractExpression, bool)
	ToString(bool) string
	ContainsVariable() bool
	Rebalance() (AbstractExpression, string, AbstractExpression)
	CanRebalance() bool
	GetName() string
}

var (
	numberRegex = regexp.MustCompile(`^([a-z]{4}):\s(\d+)$`)
	mathRegex   = regexp.MustCompile(`^([a-z]{4}):\s([a-z]{4})\s(.)\s([a-z]{4})$`)
)

func parseInput(input []string) map[string]*Monkey {
	monkeys := make(map[string]*Monkey)
	deps := make(map[string][]string)
	for _, s := range input {
		var m *Monkey
		if matches := numberRegex.FindStringSubmatch(s); matches != nil {
			m = &Monkey{
				Name:      matches[1],
				Resolved:  true,
				Value:     util.Atoi(matches[2]),
				DependsOn: []*Monkey{},
			}
		} else if matches := mathRegex.FindStringSubmatch(s); matches != nil {
			m = &Monkey{
				Name:      matches[1],
				Resolved:  false,
				Operator:  matches[3],
				DependsOn: []*Monkey{},
			}
			deps[matches[1]] = []string{matches[2], matches[4]}
		}
		monkeys[m.Name] = m
	}

	for name, d := range deps {
		m1 := monkeys[d[0]]
		m2 := monkeys[d[1]]
		monkeys[name].DependsOn = []*Monkey{m1, m2}
	}

	return monkeys
}

func init() {
	challenges.RegisterChallengeFunc(2022, 21, 1, "day21.txt", part1)
	challenges.RegisterChallengeFunc(2022, 21, 2, "day21.txt", part2)
}
