package aoc2023_day23

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Grid [][]string

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		sb.WriteString(strings.Join(row, ""))
		sb.WriteString("\n")
	}
	return sb.String()
}

type Coord = util.Coord

type Node struct {
	Coord Coord
	Cost  int
}

func GetNode(row, col, cost int) *Node {
	key := fmt.Sprintf("%d,%d", row, col)
	if _, ok := nodes[key]; !ok {
		nodes[key] = &Node{Coord{Row: row, Col: col}, cost}
	}
	return nodes[key]
}

var nodes = make(map[string]*Node)

// Item represents an item in the priority queue
type Item struct {
	value    *Coord
	path     []string
	priority int
	index    int
}

// PriorityQueue is a min-heap implementation
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// // a* search, but looking for the longest path instead of the shortest
// func ReverseAStar(grid Grid, start *Coord, end *Coord) int {
// 	queue := &PriorityQueue{}
// 	heap.Push(queue, &Item{value: start, priority: 0})
// 	cameFrom := make(map[*Coord]*Coord)
// 	cameFrom[start] = nil
// 	costSoFar := make(map[*Coord]int)
// 	costSoFar[start] = 0

// 	for queue.Len() > 0 {
// 		current := heap.Pop(queue).(*Item).value
// 		if current.Row == len(g.Grid)-1 && current.Col == len(g.Grid[current.Row])-1 {
// 			break
// 		}
// 		for _, edge := range g.Edges[current] {
// 			newCost := costSoFar[current] + edge.To.Cost
// 			if _, ok := costSoFar[edge.To]; !ok || newCost < costSoFar[edge.To] {
// 				costSoFar[edge.To] = newCost
// 				priority := newCost + ManhattanDistance(edge.To, end)
// 				heap.Push(queue, &Item{value: edge.To, priority: priority})
// 				cameFrom[edge.To] = current
// 			}
// 		}
// 	}

// 	return costSoFar[end]
// }

// func Dijkstra2(g Grid) int {
// 	nRows := len(g)
// 	nCols := len(g[0])
// 	start := &Node{
// 		Coord{Row: 0, Col: 1},
// 		0,
// 	}
// 	end := &Node{
// 		Coord{Row: nRows-1, Col: nCols-2},
// 		0,
// 	}
// 	visited := make(map[*Node]bool)
// 	queue := &PriorityQueue{}
// 	heap.Push(queue, &Item{value: &start.Coord, priority: 0})

// 	// queue = append(queue, &State{
// 	// 	Node: GetNode2(0, 1, "R", 1), // &Node2{0, 1, "R", 1},
// 	// 	Cost: g.Grid[0][1],
// 	// 	Prev: nil,
// 	// })
// 	// queue = append(queue, &State{
// 	// 	// Node: &Node2{1, 0, "D", 1},
// 	// 	Node: GetNode2(1, 0, "D", 1),
// 	// 	Cost: g.Grid[1][0],
// 	// 	Prev: nil,
// 	// })

// 	for queue.Len() > 0 {
// 		cur := heap.Pop(queue).(*Item).value
// 		if visited[cur] {
// 			continue
// 		}
// 		visited[cur] = true

// 	for len(queue) > 0 {
// 		cur := queue[0]
// 		queue = queue[1:]
// 		if visited[cur.Node] {
// 			continue
// 		}

// 		neighbors := g.GetNeighbors(cur.Node)
// 		visited[cur.Node] = true
// 		if cur.Node.Row == nRows-1 && cur.Node.Col == nCols-1 {
// 			return cur.Cost
// 		}
// 		queue = cur.AddNext(queue, &g)
// 	}
// 	return -1
// }

type Entry struct {
	Cost int
	Path []string
}

func GetNeighbors(g Grid, c *Coord, ignoreSlopes bool) []*Coord {
	neighbors := make([]*Coord, 0)
	v := g[c.Row][c.Col]
	if v == "." || ignoreSlopes {
		if c.Row > 0 && g[c.Row-1][c.Col] != "#" && (g[c.Row-1][c.Col] != "v" || ignoreSlopes) {
			neighbors = append(neighbors, GetCoord(c.Row-1, c.Col))
		}
		if c.Row < len(g)-1 && g[c.Row+1][c.Col] != "#" && (g[c.Row+1][c.Col] != "^" || ignoreSlopes) {
			neighbors = append(neighbors, GetCoord(c.Row+1, c.Col))
		}
		if c.Col > 0 && g[c.Row][c.Col-1] != "#" && (g[c.Row][c.Col-1] != ">" || ignoreSlopes) {
			neighbors = append(neighbors, GetCoord(c.Row, c.Col-1))
		}
		if c.Col < len(g[c.Row])-1 && g[c.Row][c.Col+1] != "#" && (g[c.Row][c.Col+1] != "<" || ignoreSlopes) {
			neighbors = append(neighbors, GetCoord(c.Row, c.Col+1))
		}
	} else if v == ">" {
		neighbors = append(neighbors, GetCoord(c.Row, c.Col+1))
	} else if v == "<" {
		neighbors = append(neighbors, GetCoord(c.Row, c.Col-1))
	} else if v == "^" {
		neighbors = append(neighbors, GetCoord(c.Row-1, c.Col))
	} else if v == "v" {
		neighbors = append(neighbors, GetCoord(c.Row+1, c.Col))
	}
	// log.Debugf("Neighbors of %v: %v\n", c, neighbors)
	return neighbors
}

var allCoords = make(map[string]*Coord)

func GetCoord(row, col int) *Coord {
	key := fmt.Sprintf("%d,%d", row, col)
	if _, ok := allCoords[key]; !ok {
		allCoords[key] = &Coord{Row: row, Col: col}
	}
	return allCoords[key]
}

func GetPaths(c *Coord, paths map[*Coord][][]*Coord) [][]*Coord {
	if p, ok := paths[c]; ok {
		return p
	}
	return make([][]*Coord, 0)
}

func PathEquals(a []*Coord, b []*Coord) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ContainsPath(paths [][]*Coord, path []*Coord) bool {
	for _, p := range paths {
		if PathEquals(p, path) {
			return true
		}
	}
	return false
}

func (g Grid) Dijkstra(start *Coord, end *Coord, part2 bool) int {
	paths := make(map[*Coord][][]*Coord)
	paths[start] = [][]*Coord{{start}}
	queue := &PriorityQueue{}

	heap.Push(queue, &Item{value: start, priority: 0})

	maxSolution := 0
	for queue.Len() > 0 {
		current := queue.Pop().(*Item).value
		curPaths := GetPaths(current, paths)
		neighbors := GetNeighbors(g, current, part2)
		for _, cp := range curPaths {
			for _, neighbor := range neighbors {
				if !util.Contains(cp, neighbor) {
					nPaths := GetPaths(neighbor, paths)
					newPath := make([]*Coord, len(cp))
					copy(newPath, cp)
					newPath = append(newPath, neighbor)
					if !ContainsPath(nPaths, newPath) {
						nPaths = append(nPaths, newPath)
						paths[neighbor] = nPaths
						heap.Push(queue, &Item{value: neighbor, priority: len(newPath)})
						if neighbor.Row == end.Row && neighbor.Col == end.Col {
							if len(newPath) > maxSolution {
								fmt.Printf("Found longer path: %d\n", len(newPath)-1)
								maxSolution = len(newPath)
							}
						}
					}
				}
			}
		}
	}

	endPaths := GetPaths(end, paths)

	maxPath := 0
	for _, p := range endPaths {
		maxPath = util.Max(maxPath, len(p))
	}

	return maxPath - 1
}

func GetPath(c *Coord, paths map[*Coord][]*Coord) []*Coord {
	if p, ok := paths[c]; ok {
		return p
	}
	return make([]*Coord, 0)
}

// }
// func (g Grid) DijkstraOptimized(start *Coord, end *Coord, part2 bool) int {
// 	paths := make(map[*Coord][]*Coord)
// 	paths[start] = []*Coord{start}
// 	queue := &PriorityQueue{}

// 	heap.Push(queue, &Item{value: start, priority: 0})

// 	for queue.Len() > 0 {
// 		current := queue.Pop().(*Item2).value
// 		cp := GetPath(current, paths)
// 		neighbors := GetNeighbors(g, current.Cur, part2)
// 		for _, neighbor := range neighbors {
// 			if !util.Contains(cp, neighbor) {
// 				nPath := GetPath(neighbor, paths)
// 				newPath := make([]*Coord, len(cp))
// 				copy(newPath, cp)
// 				newPath = append(newPath, neighbor)
// 				if len(newPath) > len(nPath) {
// 					paths[neighbor] = newPath
// 					heap.Push(queue, &Item{value: neighbor, priority: len(newPath)})
// 				}
// 			}
// 		}
// 	}

// 	endPath := GetPath(end, paths)

// 	return len(endPath) - 1
// }

type Edge struct {
	To   *Coord
	Dist int
}

type Graph map[*Coord]map[*Coord]int

type PathInfo struct {
	Cur        *Coord
	PrevCoord  *Coord
	PrevVertex *Coord
	Steps      int
}

// Item represents an item in the priority queue
type Item2 struct {
	value    *PathInfo
	priority int
	index    int
}

// PriorityQueue is a min-heap implementation
type PriorityQueue2 []*Item2

func (pq PriorityQueue2) Len() int { return len(pq) }

func (pq PriorityQueue2) Less(i, j int) bool { return pq[i].priority < pq[j].priority }

func (pq PriorityQueue2) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue2) Push(x interface{}) {
	item := x.(*Item2)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue2) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func SimplifyGraph(g Grid, isPart2 bool) Graph {
	start := GetCoord(1, 1)
	end := GetCoord(len(g)-1, len(g[0])-2)
	graph := make(Graph)
	// graph[start] = make([]Edge, 0)
	// graph[start] = append(graph[start], Edge{To: start, Dist: 1})
	queue := &PriorityQueue2{}
	heap.Push(queue, &Item2{
		priority: 0,
		value: &PathInfo{
			Cur:        start,
			PrevCoord:  GetCoord(0, 1),
			PrevVertex: GetCoord(0, 1),
			Steps:      1,
		},
	})

	for queue.Len() > 0 {
		cur := queue.Pop().(*Item2).value
		if cur.Cur.Row == end.Row && cur.Cur.Col == end.Col || graph[cur.Cur] != nil {
			if graph[cur.Cur] == nil {
				graph[cur.Cur] = make(map[*Coord]int)
			}
			graph[cur.Cur][cur.PrevVertex] = cur.Steps
			if graph[cur.PrevVertex] == nil {
				graph[cur.PrevVertex] = make(map[*Coord]int)
			}
			graph[cur.PrevVertex][cur.Cur] = cur.Steps
			continue
		}

		newNeighbors := make([]*Coord, 0)
		directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range directions {
			newCoord := GetCoord(cur.Cur.Row+dir[0], cur.Cur.Col+dir[1])
			if g[newCoord.Row][newCoord.Col] != "#" && newCoord != cur.PrevCoord {
				newNeighbors = append(newNeighbors, newCoord)
			}
		}
		// for _, neighbor := range GetNeighbors(g, cur.Cur, isPart2) {
		// 	if neighbor != cur.PrevCoord {
		// 		newNeighbors = append(newNeighbors, neighbor)
		// 	}
		// }
		if len(newNeighbors) > 1 {
			if graph[cur.Cur] == nil {
				graph[cur.Cur] = make(map[*Coord]int)
			}
			graph[cur.Cur][cur.PrevVertex] = cur.Steps
			if graph[cur.PrevVertex] == nil {
				graph[cur.PrevVertex] = make(map[*Coord]int)
			}
			graph[cur.PrevVertex][cur.Cur] = cur.Steps
			cur.PrevVertex = cur.Cur
			cur.Steps = 0
		}
		for _, neighbor := range newNeighbors {
			heap.Push(queue, &Item2{
				priority: cur.Steps,
				value: &PathInfo{
					Cur:        neighbor,
					PrevCoord:  cur.Cur,
					PrevVertex: cur.PrevVertex,
					Steps:      cur.Steps + 1,
				},
			})
		}
	}
	return graph
}

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", grid)

	start := GetCoord(0, 1)
	end := GetCoord(len(grid)-1, len(grid[0])-2)
	result := grid.Dijkstra(start, end, false)

	// for k, v := range distances {
	// 	log.Debugf("%s: %d\n", k, v.Cost)
	// }
	// log.Debugf("------------\n%s\n", end)
	// var result = distances[end.String()].Cost
	return fmt.Sprintf("%d", result), nil
}

type PathEntry struct {
	Coord    *Coord
	CameFrom *Coord
	Dist     int
}

func GetGraphPaths(c *Coord, paths map[*Coord][][]*PathEntry) [][]*PathEntry {
	if p, ok := paths[c]; ok {
		return p
	}
	return make([][]*PathEntry, 0)
}

func IsInPath(c *Coord, path []*PathEntry) bool {
	for _, p := range path {
		if p.Coord == c {
			return true
		}
	}
	return false
}

func PathEquals2(a []*PathEntry, b []*PathEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Coord != b[i].Coord {
			return false
		}
	}
	return true
}

func SeenPath(paths [][]*PathEntry, path []*PathEntry) bool {
	for _, p := range paths {
		if PathEquals2(p, path) {
			return true
		}
	}
	return false
}

func PathDistance(path []*PathEntry) int {
	dist := 0
	for _, e := range path {
		dist += e.Dist
	}

	return dist
}

type QueueEntry struct {
	Coord *Coord
	Steps int
	Seen  map[*Coord]bool
}

func (g Graph) Zack(start *Coord, end *Coord) int {
	queue := make([]*QueueEntry, 0)
	queue = append(queue, &QueueEntry{
		Coord: start,
		Steps: 0,
		Seen:  make(map[*Coord]bool),
	})
	maxSteps := 0
	// longestPath := make(map[*Coord]bool)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.Coord == end {
			if cur.Steps > maxSteps {
				maxSteps = cur.Steps
				// longestPath = cur.Seen
			}
			continue
		}

		for neighbor, dist := range g[cur.Coord] {
			if !cur.Seen[neighbor] {
				newSeen := make(map[*Coord]bool)
				for k, v := range cur.Seen {
					newSeen[k] = v
				}
				newSeen[neighbor] = true
				queue = append(queue, &QueueEntry{
					Coord: neighbor,
					Steps: cur.Steps + dist,
					Seen:  newSeen,
				})
			}
		}
	}
	return maxSteps
}

func (g Graph) Dijkstra(start *Coord, end *Coord, part2 bool) int {
	paths := make(map[*Coord][][]*PathEntry)
	paths[start] = [][]*PathEntry{{&PathEntry{start, nil, 0}}}
	queue := &PriorityQueue2{}
	heap.Push(queue, &Item2{
		priority: 0,
		value: &PathInfo{
			Cur:        start,
			PrevCoord:  nil,
			PrevVertex: nil,
			Steps:      0,
		},
	})

	maxSolution := 0
	for queue.Len() > 0 {
		current := queue.Pop().(*Item2).value
		curPaths := GetGraphPaths(current.Cur, paths)
		neighbors := g[current.Cur]
		for _, cp := range curPaths {
			for neighbor, dist := range neighbors {
				if neighbor == current.PrevVertex {
					continue
				}
				if !IsInPath(neighbor, cp) {
					nPaths := GetGraphPaths(neighbor, paths)
					newPath := make([]*PathEntry, len(cp))
					copy(newPath, cp)
					newPath = append(newPath, &PathEntry{neighbor, current.Cur, dist})
					if !SeenPath(nPaths, newPath) {
						nPaths = append(nPaths, newPath)
						paths[neighbor] = nPaths
						heap.Push(queue, &Item2{
							priority: PathDistance(newPath),
							value: &PathInfo{
								Cur:        neighbor,
								PrevCoord:  current.Cur,
								PrevVertex: current.PrevVertex,
								Steps:      PathDistance(newPath),
							},
						})
						if neighbor.Row == end.Row && neighbor.Col == end.Col {
							d := PathDistance(newPath)
							if d > maxSolution {
								fmt.Printf("Found longer path: %d\n", d)
								maxSolution = d
							}
						}
					}
				}
			}
		}
	}

	endPaths := GetGraphPaths(end, paths)

	maxPath := 0
	for _, p := range endPaths {
		maxPath = util.Max(maxPath, PathDistance(p))
	}

	return maxPath
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", grid)

	g := SimplifyGraph(grid, true)

	for k, v := range g {
		log.Debugf("%s: %v\n", k, v)
	}
	// log.Debugf("Graph:\n%v\n", g)

	start := GetCoord(0, 1)
	end := GetCoord(len(grid)-1, len(grid[0])-2)
	// result := g.Dijkstra(start, end, true)
	result := g.Zack(start, end)

	// for k, v := range distances {
	// 	log.Debugf("%s: %d\n", k, v.Cost)
	// }
	// log.Debugf("------------\n%s\n", end)
	// var result = distances[end.String()].Cost
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) Grid {
	grid := make(Grid, 0, len(input))
	for _, s := range input {
		grid = append(grid, strings.Split(s, ""))
	}
	return grid
}

func init() {
	challenges.RegisterChallengeFunc(2023, 23, 1, "day23.txt", part1)
	challenges.RegisterChallengeFunc(2023, 23, 2, "day23.txt", part2)
}
