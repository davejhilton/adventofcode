package aoc2021_day12

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	edges := parseInput(input)

	checkValidMove := func(p Path, next *Node) bool {
		if next.Value == "start" {
			return false
		} else if IsUpperCase(next.Value) {
			return true
		}
		for _, e := range p {
			if e == next.Value {
				return false
			}
		}
		return true
	}

	nodes := populateGraph(edges)
	start := nodes["start"]
	log.Debugln("start")
	allPaths := findAllPaths(start, Path{start.Value}, checkValidMove, "")
	return fmt.Sprintf("%d", len(allPaths)), nil
}

func part2(input []string) (string, error) {
	edges := parseInput(input)

	checkValidMove := func(p Path, next *Node) bool {
		if next.Value == "start" {
			return false
		} else if IsUpperCase(next.Value) {
			return true
		}
		visited := make(map[string]bool)
		nDoubles := 0
		for _, e := range p {
			if v := visited[e]; v && !IsUpperCase(e) {
				nDoubles++
			}
			visited[e] = true
		}
		return !visited[next.Value] || nDoubles == 0
	}

	nodes := populateGraph(edges)
	start := nodes["start"]
	log.Debugln("start")
	allPaths := findAllPaths(start, Path{start.Value}, checkValidMove, "")
	return fmt.Sprintf("%d", len(allPaths)), nil
}

func findAllPaths(n *Node, p Path, checkValidMove MoveValidator, parentDebugPrefix string) []Path {
	paths := make([]Path, 0)
	for i, e := range n.Edges {
		debugPrefix := generateDebugPrefix(parentDebugPrefix, i == len(n.Edges)-1)
		if checkValidMove(p, e) {
			newPath := append(p[0:], e.Value)
			if e.Value == "end" {
				paths = append(paths, newPath)
				log.Debugf("%s%s,%s\n", debugPrefix, p, log.Colorize(e.Value, log.Green, 0))
			} else {
				log.Debugf("%s%s,%s...\n", debugPrefix, p, e.Value)
				childPaths := findAllPaths(e, newPath, checkValidMove, debugPrefix)
				paths = append(paths, childPaths...)
			}
		} else {
			log.Debugf("%s%s,%s\n", debugPrefix, p, log.Colorize(e.Value, log.Red, 0))
		}
	}
	return paths
}

func populateGraph(edges []Path) map[string]*Node {
	nodes := make(map[string]*Node)
	for _, p := range edges {
		n1 := nodes[p[0]]
		n2 := nodes[p[1]]
		if n1 == nil {
			nodes[p[0]] = &Node{Value: p[0], Edges: make([]*Node, 0)}
			n1 = nodes[p[0]]
		}
		if n2 == nil {
			nodes[p[1]] = &Node{Value: p[1], Edges: make([]*Node, 0)}
			n2 = nodes[p[1]]
		}
		n1.Edges = append(n1.Edges, n2)
		n2.Edges = append(n2.Edges, n1)
	}
	return nodes
}

type Path []string

func (p Path) String() string {
	return strings.Join(p, ",")
}

type MoveValidator func(p Path, next *Node) bool

type Node struct {
	Value string
	Edges []*Node
}

func IsUpperCase(v string) bool {
	return v[0] >= 'A' && v[0] <= 'Z'
}

func generateDebugPrefix(parentPrefix string, isLastEdge bool) string {
	if !log.DebugEnabled() {
		return ""
	}
	var prefix strings.Builder
	for i, r := range []rune(parentPrefix) {
		if i%4 == 0 {
			if r == rune('└') {
				prefix.WriteRune('·')
			} else if r == rune('├') {
				prefix.WriteRune('│')
			} else {
				prefix.WriteRune(r)
			}
		} else {
			prefix.WriteRune(' ')
		}
	}
	if isLastEdge {
		prefix.WriteString("└── ")
	} else {
		prefix.WriteString("├── ")
	}

	return prefix.String()
}

func parseInput(input []string) []Path {
	paths := make([]Path, 0, len(input))
	for _, s := range input {
		paths = append(paths, strings.SplitN(s, "-", 2))
	}
	return paths
}

func init() {
	challenges.RegisterChallengeFunc(2021, 12, 1, "day12.txt", part1)
	challenges.RegisterChallengeFunc(2021, 12, 2, "day12.txt", part2)
}
