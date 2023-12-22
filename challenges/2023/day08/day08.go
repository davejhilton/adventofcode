package aoc2023_day8

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
)

type Node struct {
	Name  string
	Left  *Node
	Right *Node
}

func part1(input []string) (string, error) {
	instructions, nodes := parseInput(input)
	curNode := nodes["AAA"]
	var i, nSteps int
	for curNode != nil && curNode.Name != "ZZZ" {
		nSteps++
		if instructions[i] == "L" {
			curNode = curNode.Left
		} else {
			curNode = curNode.Right
		}
		i = (i + 1) % len(instructions)
	}

	return fmt.Sprintf("%d", nSteps), nil
}

func part2(input []string) (string, error) {
	instructions, nodes := parseInput(input)

	var minSteps uint64 = 1
	for _, node := range nodes {
		if node.Name[2] == 'A' {
			minSteps = util.LCM(minSteps, node.StepsToZ(instructions, 0))
		}
	}

	return fmt.Sprintf("%d", minSteps), nil
}

func parseInput(input []string) ([]string, map[string]*Node) {
	instr := strings.Split(input[0], "")

	nodes := make(map[string]*Node)
	for i := 2; i < len(input); i++ {
		var name, lName, rName string
		fmt.Sscanf(input[i], "%3s = (%3s, %3s)", &name, &lName, &rName)
		if nodes[name] == nil {
			nodes[name] = &Node{name, nil, nil}
		}
		if nodes[lName] == nil {
			nodes[lName] = &Node{lName, nil, nil}
		}
		if nodes[rName] == nil {
			nodes[rName] = &Node{rName, nil, nil}
		}
		nodes[name].Left = nodes[lName]
		nodes[name].Right = nodes[rName]
	}

	return instr, nodes
}

var memo = make(map[string]*uint64)

func (n *Node) StepsToZ(instructions []string, i int) uint64 {
	key := fmt.Sprintf("%s:%d", n.Name, i)
	if memo[key] != nil {
		return *memo[key]
	}

	var steps uint64

	if n.Name[2] == 'Z' {
		memo[key] = &steps
		return steps
	}

	var nextNode *Node
	if instructions[i] == "L" {
		nextNode = n.Left
	} else {
		nextNode = n.Right
	}
	steps = 1 + nextNode.StepsToZ(instructions, (i+1)%len(instructions))
	memo[key] = &steps
	return steps
}

func init() {
	challenges.RegisterChallengeFunc(2023, 8, 1, "day08.txt", part1)
	challenges.RegisterChallengeFunc(2023, 8, 2, "day08.txt", part2)
}
