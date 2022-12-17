package aoc2022_day16

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	start, valves := parseInput(input)

	ctx := &Context{
		Start:     start,
		Time:      30,
		Valves:    valves,
		Keys:      util.Keys(valves),
		Distances: precalcDistances(valves),
	}

	best := depthFirstSearch(ctx, start, 0b00, 30, false)

	return fmt.Sprintf("%d", best), nil
}

func part2(input []string) (string, error) {
	start, valves := parseInput(input)

	ctx := &Context{
		Start:     start,
		Time:      26,
		Valves:    valves,
		Keys:      util.Keys(valves),
		Distances: precalcDistances(valves),
	}

	// log.Debugln(util.ToJSON(ctx.Distances))

	sort.Slice(ctx.Keys, func(i, j int) bool { return ctx.Keys[i] < ctx.Keys[j] })

	best := depthFirstSearch(ctx, start, 0b00, 26, true)

	return fmt.Sprintf("%d", best), nil
}

func precalcDistances(valves map[uint64]*Valve) map[uint64]map[uint64]int {
	distMap := make(map[uint64]map[uint64]int)
	for k, val := range valves {
		distMap[k] = make(map[uint64]int)
		for k2 := range valves {
			distMap[k][k2] = math.MaxInt / 10 // default max distance
		}
		distMap[k][k] = 0 // self
		for _, t := range val.Tunnels {
			distMap[k][t] = 1 // immediate neighbors
		}
	}

	for k := range valves {
		for i := range valves {
			for j := range valves {
				distMap[i][j] = util.Min(distMap[i][j], distMap[i][k]+distMap[k][j])
			}
		}
	}

	return distMap
}

var cache = make(map[string]int)

func depthFirstSearch(ctx *Context, valve uint64, bitmask uint64, t int, first bool) int {
	cacheKey := fmt.Sprintf("%d_%d_%d_%v", valve, bitmask, t, first)
	if res, ok := cache[cacheKey]; ok {
		return res
	}

	targets := make([]uint64, 0)
	for _, k := range ctx.Keys {
		if (bitmask&k == 0) && ctx.Valves[k].FlowRate != 0 && ctx.Distances[valve][k] < t {
			targets = append(targets, k)
		}
	}

	pressures := make([]int, 0)
	for _, k := range targets {
		d := ctx.Distances[valve][k] + 1
		timeLeft := t - d
		pressure := timeLeft*ctx.Valves[k].FlowRate +
			depthFirstSearch(ctx, k, bitmask|k, timeLeft, first)
		pressures = append(pressures, pressure)
	}

	if first {
		pressures = append(pressures, depthFirstSearch(ctx, ctx.Start, bitmask, ctx.Time, false))
	}
	best := util.Max(pressures...)

	cache[cacheKey] = best
	return best
}

type Valve struct {
	Name     string
	Key      uint64
	FlowRate int
	Tunnels  []uint64
}

type Context struct {
	Start     uint64
	Time      int
	Keys      []uint64
	Valves    map[uint64]*Valve
	Distances map[uint64]map[uint64]int
}

var (
	parseRegex = regexp.MustCompile(`^Valve ([A-Z]+) has flow rate=(\d+); tunnels? leads? to valves? (.*)$`)
)

func parseInput(input []string) (start uint64, valves map[uint64]*Valve) {
	valves = make(map[uint64]*Valve)
	tunnels := make(map[uint64][]string)
	names := make(map[string]uint64)
	for i, s := range input {
		parsed := parseRegex.FindStringSubmatch(s)
		var k uint64 = (1 << i)
		names[parsed[1]] = k
		valves[k] = &Valve{
			Key:      k,
			Name:     parsed[1],
			FlowRate: util.Atoi(parsed[2]),
		}
		tunnels[k] = strings.Split(parsed[3], ", ")
	}

	for k, tn := range tunnels {
		v := valves[k]
		v.Tunnels = make([]uint64, len(tn))
		for i, t := range tn {
			v.Tunnels[i] = names[t]
		}
	}

	return names["AA"], valves
}

func init() {
	challenges.RegisterChallengeFunc(2022, 16, 1, "day16.txt", part1)
	challenges.RegisterChallengeFunc(2022, 16, 2, "day16.txt", part2)
}
