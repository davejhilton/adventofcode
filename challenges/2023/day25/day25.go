package aoc2023_day25

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

// This is heavily based on the solution from kenthklui:
// https://github.com/kenthklui/adventofcode/blob/master/2023/day25/task1.go

type Node struct {
	Name  string
	Edges map[*Edge]*Node
	Seen  bool
	Group int
}

type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge
}

func (g Graph) ResetNodes() {
	for _, n := range g.Nodes {
		n.Seen = false
	}
}

func (g Graph) ResetEdges() {
	for _, e := range g.Edges {
		e.Seen = false
	}
}

type Edge struct {
	Seen bool
}

type queueItem struct {
	Node *Node
	Edge *Edge
	Prev *queueItem
}

func (g Graph) BFS(source *Node, dest *Node) bool {
	queue := make([]*queueItem, 0)
	queue = append(queue, &queueItem{Node: source})

	found := false
	for len(queue) > 0 {
		log.Debugf("BFS queue len=%d\n", len(queue))
		// "pop" the first item off the queue
		cur := queue[0]
		queue = queue[1:]
		if cur.Prev != nil {
			cur.Node.Group = cur.Prev.Node.Group
		}
		if cur.Node == dest {
			for n := cur; n.Edge != nil; n = n.Prev {
				n.Edge.Seen = true
			}
			found = true
			break
		}
		for edge, node := range cur.Node.Edges {
			if edge.Seen || node.Seen {
				continue
			}
			node.Seen = true
			queue = append(queue, &queueItem{Node: node, Edge: edge, Prev: cur})
		}
	}
	g.ResetNodes()
	return found
}

func (g Graph) CutPaths(source *Node, dest *Node, pathNum int) bool {
	complete := true
	for i := 0; i < pathNum; i++ {
		log.Debugf("CutPaths BFS at i=%d\n", i)
		if !g.BFS(source, dest) {
			complete = false
			break
		}
	}
	return complete
}

func (g Graph) Split() (int, int) {
	var source *Node
	// get the first node
	for _, n := range g.Nodes {
		source = n
		break
	}

	source.Group = 1

	for _, dest := range g.Nodes {
		log.Debugf("Splitting %s and %s\n", source.Name, dest.Name)
		if dest == source {
			continue
		}
		if dest.Group > 0 {
			continue
		}

		if !g.CutPaths(source, dest, 4) {
			log.Debugf("Failed to cut paths between %s and %s\n", source.Name, dest.Name)
			g.BFS(source, nil)
			log.Debugf("BFS #1\n")
			dest.Group = 2
			log.Debugf("BFS #2\n")
			g.BFS(dest, nil)
		}
		log.Debugf("Resetting Edges\n")
		g.ResetEdges()
	}

	disconnected := 0
	for _, n := range g.Nodes {
		if n.Group != 1 {
			disconnected++
		}
	}
	return len(g.Nodes) - disconnected, disconnected
}

func part1(input []string) (string, error) {
	graph := parseInput(input)
	size1, size2 := graph.Split()
	result := size1 * size2
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	var result int
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *Graph {
	nodes := make(map[string]*Node)

	for _, s := range input {
		nodeName := strings.Split(s, ": ")[0]
		nodes[nodeName] = &Node{Name: nodeName, Edges: make(map[*Edge]*Node)}
	}

	edges := make([]*Edge, 0)

	for _, s := range input {
		parts := strings.Split(s, ": ")
		name, otherNames := parts[0], strings.Split(parts[1], " ")
		node := nodes[name]
		for _, otherName := range otherNames {
			if _, ok := nodes[otherName]; !ok {
				nodes[otherName] = &Node{Name: otherName, Edges: make(map[*Edge]*Node)}
			}
			dest := nodes[otherName]
			edge := &Edge{}
			edges = append(edges, edge)
			node.Edges[edge] = dest
			dest.Edges[edge] = node
		}
	}
	return &Graph{Nodes: nodes, Edges: edges}
}

func init() {
	challenges.RegisterChallengeFunc(2023, 25, 1, "day25.txt", part1)
	challenges.RegisterChallengeFunc(2023, 25, 2, "day25.txt", part2)
}
