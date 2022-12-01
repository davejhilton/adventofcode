package aoc2021_day17

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	maxY := tryAllVelocities(input, false)
	return fmt.Sprintf("%d", maxY), nil
}

func part2(input []string) (string, error) {
	count := tryAllVelocities(input, true)
	return fmt.Sprintf("%d", count), nil
}

func tryAllVelocities(input []string, retCount bool) int {
	x1, x2, y1, y2 := parseInput(input)
	log.Debugf("Parsed Input:\n%d %d %d %d\n", x1, x2, y1, y2)

	// map["vX,vY"] -> maxY
	hits := make(map[string]int)

	for vX := 1; vX <= x2; vX++ {
		for vY := y1; vY < 200; vY++ {
			if hit, maxY := tryInitialVelocity(vX, vY, x1, x2, y1, y2, false); hit {
				hits[fmt.Sprintf("%d,%d", vX, vY)] = maxY
			}
		}
	}

	// testX, testY := 16, 5
	// hit, y, _ := tryInitialVelocity(testX, testY, x1, x2, y1, y2, true)
	// if hit {
	// 	hits[fmt.Sprintf("%d,%d", testX, testY)] = y
	// }

	if retCount {
		return len(hits)
	}

	maxY := y1

	log.Debugf("%d hits:\n", len(hits))
	for v, y := range hits {
		if y > 0 {
			log.Debugf(" - %9s: %4d\n", v, y)
		}
		maxY = util.Max(maxY, y)
	}

	return maxY
}

func tryInitialVelocity(vX, vY, tx1, tx2, ty1, ty2 int, logSteps bool) (hits bool, maxY int) {
	_log := func(f string, v ...interface{}) {
		if logSteps {
			log.Debugf(f, v...)
		}
	}

	step, x, y, maxY := 0, 0, 0, 0
	maxSteps := 500
	for x <= tx2 && y >= ty1 && step < maxSteps {
		step++
		x += vX
		y += vY
		maxY = util.Max(maxY, y)
		// check for hit
		_log("Step #%2d: (%d,%d) -> (%d,%d)", step, vX, vY, x, y)
		if x >= tx1 && x <= tx2 && y >= ty1 && y <= ty2 {
			_log(" %s\n", "HIT!")
			// log.Debugf("HIT IT AFTER %d steps -- (%d,%d)\n", step, x, y)
			return true, maxY
		}
		if x > tx2 || y < ty1 {
			// log.Debugf("PASSED IT AFTER %d steps -- (%d,%d)\n", step, x, y)
			return false, maxY
		}
		_log(" %s\n", "miss...")

		// update velocities
		if vX > 0 {
			vX -= 1
		} else if vX < 0 {
			vX += 1
		}
		vY -= 1
	}
	// log.Debugf("TRIED MAX %d steps -- (%d,%d)\n", step, x, y)
	return false, maxY
}

func parseInput(input []string) (int, int, int, int) {
	line := input[0]
	var x1, x2, y1, y2 int
	fmt.Sscanf(line, "target area: x=%d..%d, y=%d..%d", &x1, &x2, &y1, &y2)
	return x1, x2, y1, y2
}

func init() {
	challenges.RegisterChallengeFunc(2021, 17, 1, "day17.txt", part1)
	challenges.RegisterChallengeFunc(2021, 17, 2, "day17.txt", part2)
}
