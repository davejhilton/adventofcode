package aoc2022_day20

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Node struct {
	val  int
	next *Node
	prev *Node
}

func part1(input []string) (string, error) {
	nodes := parseInput(input)

	var zeroIdx int = mix(nodes)

	var result = calcSum(nodes, zeroIdx)

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	nodes := parseInput(input)

	for _, n := range nodes {
		n.val *= 811589153
	}
	var zeroIdx int
	for i := 0; i < 10; i++ {
		zeroIdx = mix(nodes)
	}
	var result = calcSum(nodes, zeroIdx)

	return fmt.Sprintf("%d", result), nil
}

func mix(nodes []*Node) (zeroIdx int) {
	for i, n := range nodes {
		iterCount := n.val % (len(nodes) - 1)
		if iterCount == 0 {
			if n.val == 0 {
				zeroIdx = i
			}
			continue
		}
		// cut the node out of the list
		n.prev.next = n.next
		n.next.prev = n.prev

		// find the location to re-insert it
		cur := n
		if iterCount > 0 {
			for j := 0; j < iterCount; j++ {
				cur = cur.next
			}
		} else if iterCount < 0 {
			for j := iterCount; j <= 0; j++ {
				cur = cur.prev
			}
		}

		// re-insert the node
		cur.next.prev = n
		n.next = cur.next
		cur.next = n
		n.prev = cur
	}
	return zeroIdx
}

func calcSum(nodes []*Node, zeroIdx int) int {
	var nums = [3]int{}
	cur := nodes[zeroIdx]
	for i := 0; i < 3; i++ {
		for j := 1; j <= 1000; j++ {
			cur = cur.next
		}
		nums[i] = cur.val
	}
	log.Debugf("A: %d, B: %d, C: %d\n", nums[0], nums[1], nums[2])

	return nums[0] + nums[1] + nums[2]
}

func parseInput(input []string) []*Node {
	nodesInOrder := make([]*Node, len(input))
	var prev *Node
	for i, s := range input {
		nodesInOrder[i] = &Node{val: util.Atoi(s), prev: prev}
		prev = nodesInOrder[i]
	}
	nodesInOrder[0].prev = nodesInOrder[len(nodesInOrder)-1]
	next := nodesInOrder[0]
	for i := len(nodesInOrder) - 1; i >= 0; i-- {
		nodesInOrder[i].next = next
		next = nodesInOrder[i]
	}
	return nodesInOrder
}

func init() {
	challenges.RegisterChallengeFunc(2022, 20, 1, "day20.txt", part1)
	challenges.RegisterChallengeFunc(2022, 20, 2, "day20.txt", part2)
}
