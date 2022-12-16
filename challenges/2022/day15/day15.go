package aoc2022_day15

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

const (
	EMPTY  = 0
	SENSOR = 1
	BEACON = 2
)

func part1(input []string) (string, error) {
	parsed := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", parsed)

	grid := NewGrid()

	var rowNum int
	if strings.Contains(challenges.CurrentChallenge.InputFileName, "example") {
		rowNum = 10
	} else {
		rowNum = 2000000
	}

	for _, sbp := range parsed {
		dist := sbp.Dist()
		log.Debugf("SENSOR: %v :: BEACON: %v :: dist: %d\n", sbp.Sensor, sbp.Beacon, dist)
		grid.Set(sbp.Sensor.X, sbp.Sensor.Y, SENSOR)
		grid.Set(sbp.Beacon.X, sbp.Beacon.Y, BEACON)

		x, y := sbp.Sensor.X, sbp.Sensor.Y
		for i := 0; i <= dist; i++ {
			if y+i == rowNum || y-i == rowNum {
				for j := 0; j <= dist-i; j++ {
					grid.MarkEmpty(x+j, y+i)
					grid.MarkEmpty(x+j, y-i)
					grid.MarkEmpty(x-j, y+i)
					grid.MarkEmpty(x-j, y-i)
				}
			}
		}
		// log.Debugln(grid)
	}
	// log.Debugln(grid)

	var result int
	for _, val := range grid.Row(rowNum) {
		if val == EMPTY {
			result++
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	parsed := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", parsed)

	var limit int
	if strings.Contains(challenges.CurrentChallenge.InputFileName, "example") {
		limit = 20
	} else {
		limit = 4000000
	}

	ranges := make(map[int][][2]int)

	var dist, x, y, d, r1, r2 int
	var r, rng, rng1 [2]int
	var r1rngs, r2rngs [][2]int
	var ok bool
	for _, sbp := range parsed {
		dist = sbp.Dist()
		// log.Debugf("SENSOR: %v :: BEACON: %v :: dist: %d\n", sbp.Sensor, sbp.Beacon, dist)
		x, y = sbp.Sensor.X, sbp.Sensor.Y
		for i := 0; i <= dist; i++ {
			d = dist - i // + 1
			r1 = y - d
			r2 = y + d
			rng = [2]int{x - i, x + i}

			if !(rng[0] < 0 && rng[1] < 0) && !(rng[0] > limit && rng[1] > limit) {
				adjusted := false
				rng1 = [2]int{rng[0], rng[1]}
				if rng[0] < 0 {
					rng[0] = 0
					adjusted = true
				}
				if rng[1] > limit {
					rng[1] = limit
					adjusted = true
				}
				if adjusted {
					log.Debugf("range [%d,%d] ADJUSTED to [%d,%d]\n", rng1[0], rng1[1], rng[0], rng[1])
				}

				rng1 = [2]int{rng[0], rng[1]}
				if r1 >= 0 && r1 <= limit {
					r1rngs, ok = ranges[r1]
					if !ok {
						log.Debugf("row %d: Solo range [%d,%d]\n", r1, rng1[0], rng1[1])
						ranges[r1] = [][2]int{rng1}
					} else {
						newRanges := make([][2]int, 0)
						for _, r = range r1rngs {
							if (rng1[0] <= r[0] && rng1[1] >= r[0]) || (rng1[0] <= r[1] && rng1[1] >= r[1]) || (r[0] <= rng1[0] && r[1] >= rng1[0]) || (r[0] <= rng1[1] && r[1] >= rng1[1]) || (util.Abs(r[0]-rng1[1]) == 1) || (util.Abs(r[1]-rng1[0]) == 1) {
								log.Debugf("row %d: Ranges [%d,%d] and [%d,%d] overlap! ", r1, rng1[0], rng1[1], r[0], r[1])
								rng1[0] = util.Min(rng1[0], r[0])
								rng1[1] = util.Max(rng1[1], r[1])
								log.Debugf(" :: after merge: [%d,%d]\n", rng1[0], rng1[1])
							} else {
								log.Debugf("row %d: Ranges [%d,%d] and [%d,%d] do NOT overlap\n", r1, rng1[0], rng1[1], r[0], r[1])
								newRanges = append(newRanges, r)
							}
						}
						newRanges = append(newRanges, rng1)
						ranges[r1] = newRanges
					}
				}

				if r2 != r1 && r2 >= 0 && r2 <= limit {
					r2rngs, ok = ranges[r2]
					if !ok {
						log.Debugf("row %d: Solo range [%d,%d]\n", r2, rng[0], rng[1])
						ranges[r2] = [][2]int{rng}
					} else {
						newRanges := make([][2]int, 0)
						for _, r = range r2rngs {
							if (rng[0] <= r[0] && rng[1] >= r[0]) || (rng[0] <= r[1] && rng[1] >= r[1]) || (r[0] <= rng[0] && r[1] >= rng[0]) || (r[0] <= rng[1] && r[1] >= rng[1]) || (util.Abs(r[0]-rng[1]) == 1) || (util.Abs(r[1]-rng[0]) == 1) {
								log.Debugf("row %d: Ranges [%d,%d] and [%d,%d] overlap! ", r2, rng[0], rng[1], r[0], r[1])
								rng[0] = util.Min(rng[0], r[0])
								rng[1] = util.Max(rng[1], r[1])
								log.Debugf(" :: after merge: [%d,%d]\n", rng[0], rng[1])
							} else {
								log.Debugf("row %d: Ranges [%d,%d] and [%d,%d] do NOT overlap\n", r2, rng[0], rng[1], r[0], r[1])
								newRanges = append(newRanges, r)
							}
						}
						newRanges = append(newRanges, rng)
						ranges[r2] = newRanges
					}
				}
			} else {
				log.Debugf("EXCLUDED: [%d,%d]\n", rng[0], rng[1])
			}
		}
	}

	var resX, resY int
	// var b strings.Builder
	for r := 0; r <= limit; r++ {
		// log.Debugf("row %d: %v\n", r, ranges[r])
		if len(ranges[r]) > 1 {
			resY = r
			if ranges[r][0][0] < ranges[r][1][0] && ranges[r][1][0]-ranges[r][0][1] == 2 {
				resX = ranges[r][0][1] + 1
			} else {
				resX = ranges[r][1][1] + 1
			}
			log.Debugf("xx disjoint at row %d\n\t%v\n", r, ranges[r])
			break
		} else if ranges[r][0][0] != 0 {
			resY = r
			resX = 0
			log.Debugf("yy disjoint at row %d\n\t%v\n", r, ranges[r])
			break
		} else if ranges[r][0][1] != limit {
			resY = r
			resX = limit
			log.Debugf("zz disjoint at row %d\n\t%v\n", r, ranges[r])
			break
		}
		// str := make([]rune, limit+1)
		// for i := range str {
		// 	str[i] = ' '
		// }
		// for _, rng := range ranges[r] {
		// 	for i := rng[0]; i <= rng[1]; i++ {
		// 		str[i] = '#'
		// 	}
		// }
		// b.WriteString(fmt.Sprintf("%2d %s\n", r, string(str)))
	}
	// log.Debugln(b.String())
	fmt.Printf("x: %d, y: %d\n", resX, resY)

	var result int = (resX * 4000000) + resY
	// for _, val := range grid.Row() {
	// 	if val == EMPTY {
	// 		result++
	// 	}
	// }

	return fmt.Sprintf("%d", result), nil
}

var (
	lineRegex = regexp.MustCompile(`^Sensor at x=([0-9-]+), y=([0-9-]+): closest beacon is at x=([0-9-]+), y=([0-9-]+)$`)
)

func parseInput(input []string) []SensorBeaconPair {
	pairs := make([]SensorBeaconPair, 0, len(input))
	for _, s := range input {
		matches := lineRegex.FindStringSubmatch(s)
		p := SensorBeaconPair{
			Sensor: Coord{util.Atoi(matches[1]), util.Atoi(matches[2])},
			Beacon: Coord{util.Atoi(matches[3]), util.Atoi(matches[4])},
		}
		pairs = append(pairs, p)
	}
	return pairs
}

type SensorBeaconPair struct {
	Sensor Coord
	Beacon Coord
}

func (sbp SensorBeaconPair) Dist() int {
	x := util.Abs(sbp.Sensor.X - sbp.Beacon.X)
	y := util.Abs(sbp.Sensor.Y - sbp.Beacon.Y)
	return x + y
}

type Coord struct {
	X int
	Y int
}

type Grid struct {
	points map[int]map[int]int
	MinX   int
	MinY   int
	MaxX   int
	MaxY   int
}

func NewGrid() *Grid {
	return &Grid{
		points: make(map[int]map[int]int),
		MinX:   math.MaxInt,
		MaxX:   math.MinInt,
		MinY:   math.MaxInt,
		MaxY:   math.MinInt,
	}
}

func (g *Grid) Set(x, y, val int) {
	if _, ok := g.points[y]; !ok {
		g.points[y] = make(map[int]int)
	}
	// log.Debugf("Setting (%d,%d) to %d\n", y, x, val)
	g.points[y][x] = val
	g.MaxX = util.Max(g.MaxX, x)
	g.MaxY = util.Max(g.MaxY, y)
	g.MinX = util.Min(g.MinX, x)
	g.MinY = util.Min(g.MinY, y)
}

func (g *Grid) MarkEmpty(x, y int) {
	if _, ok := g.points[y]; ok {
		if _, ok2 := g.points[y][x]; !ok2 {
			g.Set(x, y, EMPTY)
		}
	} else {
		g.Set(x, y, EMPTY)
	}
}

func (g Grid) Get(x, y int) int {
	if _, ok := g.points[y]; ok {
		if v, ok2 := g.points[y][x]; ok2 {
			return v
		}
	}
	return 0
}

func (g Grid) Row(y int) map[int]int {
	return g.points[y]
}

func (g Grid) String() string {
	width := g.MaxX - g.MinX + 1
	log.Debugf("MAX X: %d, MIN X: %d\n", g.MaxX, g.MinX)
	header := make([][]string, 3)
	for y := range header {
		header[y] = make([]string, width+3)
	}

	for x := 3; x < width+3; x++ {
		xVal := fmt.Sprintf("%3d", g.MinX+x-3)
		header[0][x] = string(xVal[0])
		header[1][x] = string(xVal[1])
		header[2][x] = string(xVal[2])
	}

	var b strings.Builder
	for _, s := range header {
		b.WriteString(fmt.Sprintf("   %s\n", strings.Join(s, "")))
	}

	for y := g.MinY; y <= g.MaxY; y++ {
		b.WriteString(fmt.Sprintf("%3d", y))
		var row map[int]int
		if r, ok := g.points[y]; ok {
			row = r
		} else {
			row = make(map[int]int)
		}
		for x := g.MinX; x <= g.MaxX; x++ {
			v, ok := row[x]
			switch v {
			case EMPTY:
				if ok {
					b.WriteRune('#')
				} else {
					b.WriteRune('.')
				}
			case SENSOR:
				b.WriteRune('S')
			case BEACON:
				b.WriteRune('B')
			default:
				b.WriteRune('?')
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func init() {
	challenges.RegisterChallengeFunc(2022, 15, 1, "day15.txt", part1)
	challenges.RegisterChallengeFunc(2022, 15, 2, "day15.txt", part2)
}
