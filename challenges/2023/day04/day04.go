package aoc2023_day4

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/util"
)

type Card struct {
	CardId         int
	WinningNumbers []int
	MyNumbers      []int
}

func part1(input []string) (string, error) {
	cards := parseInput(input)
	result := 0
	for _, card := range cards {
		cardScore := 0
		for _, mine := range card.MyNumbers {
			for _, num := range card.WinningNumbers {
				if num == mine {
					if cardScore == 0 {
						cardScore = 1
					} else {
						cardScore *= 2
					}
					break
				}
			}
		}
		result += cardScore
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	cards := parseInput(input)
	cardCounts := make(map[int]int)
	result := 0
	for i, card := range cards {
		cardCounts[card.CardId] += 1 // count this card, too
		result += cardCounts[card.CardId]
		nMatches := 0
		for _, mine := range card.MyNumbers {
			for _, num := range card.WinningNumbers {
				if num == mine {
					nMatches++
					if i+nMatches < len(cards) {
						cardCounts[cards[i+nMatches].CardId] += cardCounts[card.CardId]
					}
					break
				}
			}
		}
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []Card {
	cards := make([]Card, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, ": ")
		numSets := strings.Split(parts[1], " | ")
		cards = append(cards, Card{
			CardId:         util.ExtractNumbers(parts[0])[0],
			WinningNumbers: util.ExtractNumbers(numSets[0]),
			MyNumbers:      util.ExtractNumbers(numSets[1]),
		})
	}
	return cards
}

func init() {
	challenges.RegisterChallengeFunc(2023, 4, 1, "day04.txt", part1)
	challenges.RegisterChallengeFunc(2023, 4, 2, "day04.txt", part2)
}
