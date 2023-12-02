package aoc2023_day2

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

type CubeCount struct {
	Red   int
	Green int
	Blue  int
}

func (c CubeCount) String() string {
	return fmt.Sprintf("red %d, green %d, blue%d", c.Red, c.Green, c.Blue)
}

type Game struct {
	Id         int
	CubeCounts []CubeCount
}

func part1(input []string) (string, error) {
	games := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", games)
	maxCount := CubeCount{
		Red:   12,
		Green: 13,
		Blue:  14,
	}

	result := 0
	for _, game := range games {
		// log.Debugf("Game %d: %v\n", game.Id, game.CubeCounts)
		possible := true
		for _, count := range game.CubeCounts {
			if count.Red > maxCount.Red || count.Green > maxCount.Green || count.Blue > maxCount.Blue {
				possible = false
				break
			}
		}
		if possible {
			log.Debugf("Game %d is possible\n", game.Id)
			result += game.Id
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	games := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", games)

	result := 0
	for _, game := range games {
		log.Debugf("Game %d: %v\n", game.Id, game.CubeCounts)
		mins := CubeCount{}
		for _, count := range game.CubeCounts {
			if count.Red > mins.Red {
				mins.Red = count.Red
			}
			if count.Green > mins.Green {
				mins.Green = count.Green
			}
			if count.Blue > mins.Blue {
				mins.Blue = count.Blue
			}
		}
		result += (mins.Red * mins.Green * mins.Blue)
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []Game {
	games := make([]Game, 0, len(input))
	for _, line := range input {
		game := Game{
			CubeCounts: make([]CubeCount, 0),
		}
		parts := strings.Split(line, ":")
		gameStr, countsStr := parts[0], parts[1]
		game.Id = util.Atoi(strings.Replace(gameStr, "Game ", "", 1))
		sets := strings.Split(countsStr, ";")
		for _, set := range sets {
			count := CubeCount{}
			cubeStr := strings.Split(set, ",")
			for _, cube := range cubeStr {
				cube = strings.TrimSpace(cube)
				cubeParts := strings.Split(cube, " ")
				color := cubeParts[1]
				num := util.Atoi(cubeParts[0])
				switch color {
				case "red":
					count.Red = num
				case "green":
					count.Green = num
				case "blue":
					count.Blue = num
				}
			}
			game.CubeCounts = append(game.CubeCounts, count)
		}
		games = append(games, game)
	}
	return games
}

func init() {
	challenges.RegisterChallengeFunc(2023, 2, 1, "day02.txt", part1)
	challenges.RegisterChallengeFunc(2023, 2, 2, "day02.txt", part2)
}
