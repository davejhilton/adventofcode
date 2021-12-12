package aoc2020_day18

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	expressions := parseInput(input)
	sum := 0
	log.Debugln("---------------")
	for _, e := range expressions {
		v := e.Value()
		log.Debugln("---------------")
		sum += v
	}
	return fmt.Sprintf("%d", sum), nil
}

func part2(input []string) (string, error) {
	expressions := parseInput(input)
	log.Debugln("*****************")
	sum := 0
	for _, e := range expressions {
		v := e.HandleAltPrecedence().Value()
		log.Debugln("---------------")
		sum += v
	}
	return fmt.Sprintf("%d", sum), nil
}

func parseInput(input []string) []token {
	expressions := make([]token, 0, len(input))
	for _, line := range input {
		exprStack := make([]token, 0)
		exprStack = append(exprStack, token{
			Type:      ROOT_EXPRESSION,
			SubTokens: make([]token, 0),
		})
		curIdx := 0
		for i := 0; i < len(line); i++ {
			c := rune(line[i])
			switch c {
			case ' ':
				// log.Debugln("SPACE")
			case '+':
				// log.Debugln("+")
				t := token{Type: ADD}
				exprStack[curIdx].SubTokens = append(exprStack[curIdx].SubTokens, t)
			case '*':
				// log.Debugln("*")
				t := token{Type: MULTIPLY}
				exprStack[curIdx].SubTokens = append(exprStack[curIdx].SubTokens, t)
			case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
				n, err := strconv.Atoi(string(c))
				if err != nil {
					fmt.Printf("PARSE ERROR: %s\n", err)
				}
				// log.Debugf("INT(%d)\n", n)
				t := token{
					Type:   NUMBER,
					IntVal: n,
				}
				exprStack[curIdx].SubTokens = append(exprStack[curIdx].SubTokens, t)
			case '(':
				// log.Debugln("(")
				t := token{
					Type:      PAREN_GROUP,
					SubTokens: make([]token, 0),
				}
				exprStack = append(exprStack, t)
				curIdx++
			case ')':
				// log.Debugln(")")
				t := exprStack[curIdx]
				exprStack = exprStack[0 : len(exprStack)-1]
				curIdx--
				exprStack[curIdx].SubTokens = append(exprStack[curIdx].SubTokens, t)
			}
			log.Debugf("CURRENT EXPR: %s %v\n", exprStack[curIdx].Type, exprStack[curIdx].SubTokens)
		}
		expressions = append(expressions, exprStack[curIdx])
	}
	return expressions
}

type token struct {
	Type      string
	IntVal    int
	SubTokens []token
}

func (t token) String() string {
	var str string
	switch t.Type {
	case ROOT_EXPRESSION:
		var b strings.Builder
		for i, e := range t.SubTokens {
			b.WriteString(e.String())
			if i != len(t.SubTokens)-1 {
				b.WriteString(" ")
			}
		}
		if len(t.SubTokens) == 0 {
			b.WriteString("[EMPTY ROOT]")
		}
		str = b.String()
		b.Reset()
	case PAREN_GROUP:
		var b strings.Builder
		b.WriteString("(")
		for i, e := range t.SubTokens {
			b.WriteString(e.String())
			if i != len(t.SubTokens)-1 {
				b.WriteString(" ")
			}
		}
		if len(t.SubTokens) == 0 {
			b.WriteString("[EMPTY PAREN GROUP]")
		}
		b.WriteString(")")
		str = b.String()
		b.Reset()
	case NUMBER:
		str = fmt.Sprintf("%d", t.IntVal)
	case ADD:
		str = "+"
	case MULTIPLY:
		str = "*"
	default:
		str = fmt.Sprintf("????%s", t.Type)
	}
	return str
}

func (t token) Value() int {
	var value int
	switch t.Type {
	case NUMBER:
		value = t.IntVal
	case ADD, MULTIPLY:
		fmt.Printf("WTF: tried getting Value() of '%s' token\n", t.Type)
	case ROOT_EXPRESSION, PAREN_GROUP:
		value = t.SubTokens[0].Value()
		var op string
		for _, subT := range t.SubTokens[1:] {
			switch subT.Type {
			case ADD, MULTIPLY:
				op = subT.Type
			case NUMBER, PAREN_GROUP:
				if op == ADD {
					value += subT.Value()
				} else if op == MULTIPLY {
					value *= subT.Value()
				} else {
					fmt.Printf("WTF: got a %s (%d) without a preceeding Op?!?\n", subT.Type, subT.Value())
					break
				}
				op = ""
			default:
				fmt.Printf("WTF: got a %s (%s)?!?\n", subT.Type, subT)
			}
		}
		log.Debugf("SOLVED:     %s  ==>  %d\n", t, value)
	}
	return value
}

func (t token) HandleAltPrecedence() token {
	log.Debugf("SIMPLIFYING: %s\n", t)
	var newT token
	switch t.Type {
	case NUMBER:
		newT = t
	case ADD, MULTIPLY:
		newT = t
	case ROOT_EXPRESSION, PAREN_GROUP:
		newT = token{
			Type:      t.Type,
			SubTokens: make([]token, 0),
		}

		var i int
		var pending *token
		for ; i < len(t.SubTokens); i++ {
			curT := t.SubTokens[i]
			// log.Debugf("i=%d : %s\n", i, curT)
			switch curT.Type {
			case NUMBER:
				pending = &curT
			case PAREN_GROUP:
				curT = curT.HandleAltPrecedence()
				pending = &curT
			case MULTIPLY:
				if pending == nil {
					fmt.Println("WTF: MULTIPLY without pending!!!")
				} else {
					newT.SubTokens = append(newT.SubTokens, *pending)
					pending = nil
					newT.SubTokens = append(newT.SubTokens, curT)
				}
			case ADD:
				if pending == nil {
					fmt.Println("WTF: ADD without pending!!!")
				} else {
					nextT := t.SubTokens[i+1]
					// log.Debugf("ADDING: %d * %d\n", *pending, nextT)
					if nextT.Type == PAREN_GROUP {
						nextT = nextT.HandleAltPrecedence()
					}
					pending = &token{
						Type:   NUMBER,
						IntVal: pending.Value() + nextT.Value(),
					}
					i++ // skip the next token, since we consumed it here
				}
			}
		}
		if pending != nil {
			newT.SubTokens = append(newT.SubTokens, *pending)
		}
	}
	log.Debugf("  SIMPLIFIED: %s  ==>  %s\n", t, newT)
	return newT
}

const (
	ROOT_EXPRESSION = "root"
	NUMBER          = "#"
	ADD             = "+"
	MULTIPLY        = "*"
	OPEN_PAREN      = "("
	CLOSE_PAREN     = ")"
	PAREN_GROUP     = "()"
)

func init() {
	challenges.RegisterChallengeFunc(2020, 18, 1, "day18.txt", part1)
	challenges.RegisterChallengeFunc(2020, 18, 2, "day18.txt", part2)
}
