package aoc2021_day15

import (
	"fmt"
	"math"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	nodeGrid, allNodes, dists := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", "nope")

	/* shortestPath */
	_, prevNodes := Dijkstras(allNodes, nodeGrid[0][0], dists)

	// log.Debugln(prevNodes)
	// log.Debugln(shortestPath)
	start := nodeGrid[0][0]
	// log.Debugf("%s\n%#v\n", start, start.Edges)
	end := nodeGrid[len(nodeGrid)-1][len(nodeGrid[0])-1]
	cur := end

	total := cur.Risk
	log.Debugf("%s -->\n", cur)
	for cur != start {
		cur = prevNodes[cur]
		total += cur.Risk
	}
	total -= cur.Risk

	return fmt.Sprintf("%d", total), nil
}

func part2(input []string) (string, error) {
	nodeGrid, allNodes, dists := parseInput2(input)
	log.Debugf("Parsed Input:\n%v\n", "nope")

	var sb strings.Builder
	for r := range nodeGrid {
		if r != 0 {
			sb.WriteString("\n")
		}
		for c := range nodeGrid[r] {
			fmt.Fprintf(&sb, "%d", nodeGrid[r][c].Risk)
		}
	}

	/* shortestPath */
	_, prevNodes := Dijkstras(allNodes, nodeGrid[0][0], dists)

	// log.Debugln(prevNodes)
	// log.Debugln(shortestPath)
	start := nodeGrid[0][0]
	// log.Debugf("%s\n%#v\n", start, start.Edges)
	end := nodeGrid[len(nodeGrid)-1][len(nodeGrid[0])-1]
	cur := end

	total := cur.Risk
	log.Debugf("%s -->\n", cur)
	for cur != start {
		cur = prevNodes[cur]
		total += cur.Risk
	}
	total -= cur.Risk

	return fmt.Sprintf("%d", total), nil
}

func Dijkstras(unvisited []*Node, startNode *Node, dists map[*Node]int) (shortestPath map[*Node]int, prevNodes map[*Node]*Node) {

	shortestPath = dists
	prevNodes = make(map[*Node]*Node)

	for len(unvisited) > 0 {
		if len(unvisited)%1000 == 0 {
			log.Debugf("%d nodes still unvisited...\n", len(unvisited))
		}
		var curNode *Node
		var curIdx int
		for i, n := range unvisited {
			if curNode == nil || shortestPath[n] < shortestPath[curNode] {
				curIdx, curNode = i, n
			}
		}
		// log.Debugf("Visiting %s\n", curNode)
		for _, n := range curNode.Edges {
			if maybe := shortestPath[curNode] + n.Risk; maybe < shortestPath[n] {
				shortestPath[n] = maybe
				prevNodes[n] = curNode
			}
		}

		// remove this node from the unvisited list.
		// to do it efficiently (since we don't care about order)
		// copy the last node in the list into this spot, and then
		// trim the last item off the list
		unvisited[curIdx] = unvisited[len(unvisited)-1]
		unvisited = unvisited[:len(unvisited)-1]
	}

	return
}

func parseInput(input []string) (nodeGrid [][]*Node, allNodes []*Node, dists map[*Node]int) {
	nodeGrid = make([][]*Node, 0, len(input))
	allNodes = make([]*Node, 0)
	dists = make(map[*Node]int)
	for rowNum, s := range input {
		row := make([]*Node, 0, len(s))
		strs := strings.Split(s, "")
		for colNum, v := range strs {
			num := util.Atoi(v)
			node := &Node{
				Coord: Coordinate{
					Row: rowNum,
					Col: colNum,
				},
				Risk:  num,
				Edges: make([]*Node, 0),
			}
			row = append(row, node)
			allNodes = append(allNodes, node)
			if colNum > 0 {
				row[colNum-1].AddEdge(node)
				node.AddEdge(row[colNum-1])
			}
			if rowNum > 0 {
				nodeGrid[rowNum-1][colNum].AddEdge(node)
				node.AddEdge(nodeGrid[rowNum-1][colNum])
			}
			if rowNum == 0 && colNum == 0 {
				dists[node] = 0
			} else {
				dists[node] = math.MaxInt
			}
		}
		nodeGrid = append(nodeGrid, row)
	}
	return
}

func parseInput2(input []string) (nodeGrid [][]*Node, allNodes []*Node, dists map[*Node]int) {
	nodeGrid = make([][]*Node, 0, len(input)*5)
	allNodes = make([]*Node, 0)
	dists = make(map[*Node]int)
	h := len(input)
	w := len(input[0])
	for rRepeat := 0; rRepeat < 5; rRepeat++ {
		for rowNum, s := range input {
			log.Debugf("Start row %d (rRepeat = %d)\n", rowNum, rRepeat)
			row := make([]*Node, 0, len(s)*5)
			strs := strings.Split(s, "")
			for cRepeat := 0; cRepeat < 5; cRepeat++ {
				for colNum, v := range strs {
					modColNum := colNum + (cRepeat * w)
					modRowNum := rowNum + (rRepeat * h)
					num := util.Atoi(v) + rRepeat + cRepeat
					for num > 9 {
						num -= 9
					}
					node := &Node{
						Coord: Coordinate{
							Row: modRowNum,
							Col: modColNum,
						},
						Risk:  num,
						Edges: make([]*Node, 0),
					}
					row = append(row, node)
					allNodes = append(allNodes, node)
					if modColNum > 0 {
						row[modColNum-1].AddEdge(node)
						node.AddEdge(row[modColNum-1])
					}
					if modRowNum > 0 {
						nodeGrid[modRowNum-1][modColNum].AddEdge(node)
						node.AddEdge(nodeGrid[modRowNum-1][modColNum])
					}
					if modRowNum == 0 && modColNum == 0 {
						dists[node] = 0
					} else {
						dists[node] = math.MaxInt
					}
				}
			}
			nodeGrid = append(nodeGrid, row)
		}
	}
	return
}

type Node struct {
	Coord Coordinate
	Risk  int
	Edges []*Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%s [%d]", n.Coord, n.Risk)
}

func (n Node) IsConnectedTo(other *Node) bool {
	for _, e := range n.Edges {
		if e == other {
			return true
		}
	}
	return false
}

func (n *Node) AddEdge(other *Node) {
	n.Edges = append(n.Edges, other)
}

type Coordinate struct {
	Row int
	Col int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d,%d", c.Row, c.Col)
}

func init() {
	challenges.RegisterChallengeFunc(2021, 15, 1, "day15.txt", part1)
	challenges.RegisterChallengeFunc(2021, 15, 2, "day15.txt", part2)
}
