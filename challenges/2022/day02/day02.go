package aoc2022_day2

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type matchup struct {
	TheirMove string
	Response  string
}

func part1(input []string) (string, error) {
	rounds := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", rounds)

	movePoints := map[string]int{
		"X": 1, // Rock
		"Y": 2, // Paper
		"Z": 3, // Scissors
	}
	matchPoints := map[string]int{
		"A X": 3, // Rock      : Rock      ?  Draw
		"A Y": 6, // Rock      : Paper     ?  Win
		"A Z": 0, // Rock      : Scissors  ?  Lose
		"B X": 0, // Paper     : Rock      ?  Lose
		"B Y": 3, // Paper     : Paper     ?  Draw
		"B Z": 6, // Paper     : Scissors  ?  Win
		"C X": 6, // Scissors  : Rock      ?  Win
		"C Y": 0, // Scissors  : Paper     ?  Lose
		"C Z": 3, // Scissors  : Scissors  ?  Draw
	}

	score := 0
	for _, round := range rounds {
		moveScore := movePoints[round.Response]
		matchScore := matchPoints[fmt.Sprintf("%s %s", round.TheirMove, round.Response)]
		score += moveScore + matchScore
	}

	return fmt.Sprintf("%d", score), nil
}

func part2(input []string) (string, error) {
	rounds := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", rounds)

	matchMoves := map[string]string{
		"A X": "C", // Rock      -->  Lose  ?=  Scissors
		"A Y": "A", // Rock      -->  Draw  ?=  Rock
		"A Z": "B", // Rock      -->  Win   ?=  Paper
		"B X": "A", // Paper     -->  Lose  ?=  Rock
		"B Y": "B", // Paper     -->  Draw  ?=  Paper
		"B Z": "C", // Paper     -->  Win   ?=  Scissors
		"C X": "B", // Scissors  -->  Lose  ?=  Rock
		"C Y": "C", // Scissors  -->  Draw  ?=  Scissors
		"C Z": "A", // Scissors  -->  Win   ?=  Paper
	}
	movePoints := map[string]int{
		"A": 1, // Rock
		"B": 2, // Paper
		"C": 3, // Scissors
	}
	matchPoints := map[string]int{
		"X": 0, // Lose
		"Y": 3, // Draw
		"Z": 6, // Win
	}

	score := 0
	for _, round := range rounds {
		myMove := matchMoves[fmt.Sprintf("%s %s", round.TheirMove, round.Response)]
		moveScore := movePoints[myMove]
		matchScore := matchPoints[round.Response]
		score += moveScore + matchScore
	}

	return fmt.Sprintf("%d", score), nil
}

func parseInput(input []string) []matchup {
	rounds := make([]matchup, 0, len(input))
	for _, s := range input {
		moves := strings.Split(s, " ")
		rounds = append(rounds, matchup{moves[0], moves[1]})
	}
	return rounds
}

func init() {
	challenges.RegisterChallengeFunc(2022, 2, 1, "day02.txt", part1)
	challenges.RegisterChallengeFunc(2022, 2, 2, "day02.txt", part2)
}
