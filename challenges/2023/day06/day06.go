package aoc2023_day6

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type Race struct {
	Time     int
	Distance int
}

func part1(input []string) (string, error) {
	races := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", races)

	var result int = 1
	for _, race := range races {
		nWinners := 0
		for holdSpeed := 0; holdSpeed < race.Time; holdSpeed++ {
			if (race.Time-holdSpeed)*holdSpeed > race.Distance {
				nWinners++
			}
		}
		result *= nWinners
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	realTime, realDist := reParseInputForPart2(parseInput(input))

	var result int = 0
	for holdSpeed := 0; holdSpeed < realTime; holdSpeed++ {
		if (realTime-holdSpeed)*holdSpeed > realDist {
			result++
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func reParseInputForPart2(races []Race) (time, dist int) {
	timeStr := "" // ghetto-concatenate all the times into one big number
	distStr := "" // ghetto-concatenate all the distances into one big number
	for _, race := range races {
		timeStr = fmt.Sprintf("%s%d", timeStr, race.Time)
		distStr = fmt.Sprintf("%s%d", distStr, race.Distance)
	}
	return util.Atoi(timeStr), util.Atoi(distStr)
}

func parseInput(input []string) []Race {
	races := make([]Race, 0)
	times := util.ExtractNumbers(strings.Split(input[0], ":")[1])
	distances := util.ExtractNumbers(strings.Split(input[1], ":")[1])
	for i, t := range times {
		races = append(races, Race{
			Time:     t,
			Distance: distances[i],
		})
	}
	return races
}

func init() {
	challenges.RegisterChallengeFunc(2023, 6, 1, "day06.txt", part1)
	challenges.RegisterChallengeFunc(2023, 6, 2, "day06.txt", part2)
}
