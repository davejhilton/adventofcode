package aoc2023_day17

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

// This solution is pretty heavily based on Zack's solutions:
// https://github.com/xathien/adventofcode-2023/blob/main/day17/pt1.rb
// https://github.com/xathien/adventofcode-2023/blob/main/day17/pt2.rb

type Grid [][]int

type Node struct {
	X      int
	Y      int
	DX     int
	DY     int
	InARow int
}

func (n Node) Key() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d", n.X, n.Y, n.DX, n.DY, n.InARow)
}

// Item represents an item in the priority queue
type Item struct {
	value    *Node
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

var nodes = make(map[string]*Node)

var gScores = make(map[string]int)

func getGScore(node *Node) int {
	if g, ok := gScores[node.Key()]; ok {
		return g
	}
	return math.MaxInt
}

func getNode(x, y, dx, dy, inARow int) *Node {
	key := fmt.Sprintf("%d,%d,%d,%d,%d", x, y, dx, dy, inARow)
	if n, ok := nodes[key]; ok {
		return n
	}
	n := &Node{x, y, dx, dy, inARow}
	nodes[key] = n
	return n
}

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	maxX := len(grid[0]) - 1
	maxY := len(grid) - 1

	cameFrom := make(map[*Node]*Node)

	var gScore int
	startNode := getNode(0, 0, 1, 0, 0)
	gScores[startNode.Key()] = 0
	var heatLoss int

	queue := make(PriorityQueue, 0)
	heap.Push(&queue, &Item{value: startNode, priority: 0})

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(*Item).value
		log.Debugf("Current: %v\n", current)
		x, y, dx, dy, inARow := current.X, current.Y, current.DX, current.DY, current.InARow
		gScore = getGScore(current)
		if x == maxX && y == maxY {
			heatLoss = gScore
			break
		}

		neighbors := []*Node{getNode(x+dx, y+dy, dx, dy, inARow+1)}
		if dy == 0 {
			neighbors = append(neighbors, getNode(x, y-1, 0, -1, 1)) // UP
			neighbors = append(neighbors, getNode(x, y+1, 0, 1, 1))  // DOWN
		} else {
			neighbors = append(neighbors, getNode(x-1, y, -1, 0, 1)) // LEFT
			neighbors = append(neighbors, getNode(x+1, y, 1, 0, 1))  // RIGHT
		}

		for _, neighbor := range neighbors {
			nx, ny, nInARow := neighbor.X, neighbor.Y, neighbor.InARow
			if nx < 0 || nx > maxX || ny < 0 || ny > maxY || nInARow > 3 {
				continue
			}
			nextCost := grid[ny][nx]
			nextGScore := gScore + nextCost
			if nextGScore >= getGScore(neighbor) {
				continue
			}
			cameFrom[neighbor] = current
			gScores[neighbor.Key()] = nextGScore
			priority := nextGScore + util.Abs(nx-maxX) + util.Abs(ny-maxY)
			heap.Push(&queue, &Item{value: neighbor, priority: priority * -1})
		}
	}

	var result int = heatLoss
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	maxX := len(grid[0]) - 1
	maxY := len(grid) - 1

	cameFrom := make(map[*Node]*Node)

	var gScore int
	startNode1 := getNode(0, 0, 1, 0, 0)
	startNode2 := getNode(0, 0, 1, 0, 0)
	gScores[startNode1.Key()] = 0
	gScores[startNode2.Key()] = 0
	var heatLoss = math.MaxInt

	queue := make(PriorityQueue, 0)
	heap.Push(&queue, &Item{value: startNode1, priority: 0})
	heap.Push(&queue, &Item{value: startNode2, priority: 0})

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(*Item).value
		log.Debugf("Current: %v\n", current)
		x, y, dx, dy, inARow := current.X, current.Y, current.DX, current.DY, current.InARow
		gScore = getGScore(current)
		if x == maxX && y == maxY {
			if inARow < 4 {
				continue
			}
			heatLoss = gScore
			break
		}

		neighbors := make([]*Node, 0)
		if inARow < 10 {
			neighbors = append(neighbors, getNode(x+dx, y+dy, dx, dy, inARow+1))
		}
		if dy == 0 && inARow > 3 {
			neighbors = append(neighbors, getNode(x, y-1, 0, -1, 1)) // UP
			neighbors = append(neighbors, getNode(x, y+1, 0, 1, 1))  // DOWN
		} else if dx == 0 && inARow > 3 {
			neighbors = append(neighbors, getNode(x-1, y, -1, 0, 1)) // LEFT
			neighbors = append(neighbors, getNode(x+1, y, 1, 0, 1))  // RIGHT
		}

		for _, neighbor := range neighbors {
			nx, ny := neighbor.X, neighbor.Y
			if nx < 0 || nx > maxX || ny < 0 || ny > maxY {
				continue
			}
			nextCost := grid[ny][nx]
			nextGScore := gScore + nextCost
			if nextGScore >= getGScore(neighbor) {
				continue
			}
			cameFrom[neighbor] = current
			gScores[neighbor.Key()] = nextGScore
			priority := nextGScore
			heap.Push(&queue, &Item{value: neighbor, priority: priority * -1})
		}
	}

	var result int = heatLoss
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Grid {
	grid := make(Grid, 0, len(input))
	for _, s := range input {
		grid = append(grid, util.AtoiSplit(s, ""))
	}
	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2023, 17, 1, "day17.txt", part1)
	challenges.RegisterChallengeFunc(2023, 17, 2, "day17.txt", part2)
}
