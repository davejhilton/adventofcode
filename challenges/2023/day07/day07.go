package aoc2023_day7

import (
	"fmt"
	"sort"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

var (
	Part1Values = map[string]string{
		"2": "A",
		"3": "B",
		"4": "C",
		"5": "D",
		"6": "E",
		"7": "F",
		"8": "G",
		"9": "H",
		"T": "I",
		"J": "J",
		"Q": "K",
		"K": "L",
		"A": "M",
	}
	Part2Values = map[string]string{
		"J": "A",
		"2": "B",
		"3": "C",
		"4": "D",
		"5": "E",
		"6": "F",
		"7": "G",
		"8": "H",
		"9": "I",
		"T": "J",
		"Q": "K",
		"K": "L",
		"A": "M",
	}

	StrTypeValues = map[string]string{
		"Five of a Kind":  "Z",
		"Four of a Kind":  "Y",
		"Full House":      "X",
		"Three of a Kind": "W",
		"Two Pair":        "V",
		"One Pair":        "U",
		"High Card":       "T",
	}
)

type Hand struct {
	Hand     string
	Cards    []string
	Bid      int
	Type     string
	Strength string
}

func (h *Hand) SetStrength(part int) string {
	if h.Type == "" {
		if part == 1 {
			h.SetTypePart1()
		} else {
			h.SetTypePart2()
		}
	}
	strength := StrTypeValues[h.Type]
	for _, card := range h.Cards {
		if part == 1 {
			strength = fmt.Sprintf("%s%s", strength, Part1Values[card])
		} else {
			strength = fmt.Sprintf("%s%s", strength, Part2Values[card])
		}
	}
	h.Strength = strength
	return strength
}

func (h *Hand) SetTypePart1() string {
	labels := make(map[string]int)
	for _, card := range h.Cards {
		labels[card] += 1
	}
	counts := make(map[int]int)
	for _, v := range labels {
		counts[v] += 1
	}
	if counts[5] == 1 {
		h.Type = "Five of a Kind"
	} else if counts[4] == 1 {
		h.Type = "Four of a Kind"
	} else if counts[3] == 1 && counts[2] == 1 {
		h.Type = "Full House"
	} else if counts[3] == 1 {
		h.Type = "Three of a Kind"
	} else if counts[2] == 2 {
		h.Type = "Two Pair"
	} else if counts[2] == 1 {
		h.Type = "One Pair"
	} else {
		h.Type = "High Card"
	}
	return h.Type
}

func (h *Hand) SetTypePart2() string {
	labels := make(map[string]int)
	jokers := 0
	for _, card := range h.Cards {
		if card == "J" {
			jokers++
		} else {
			labels[card] += 1
		}
	}
	counts := make(map[int]int)
	for _, v := range labels {
		counts[v] += 1
	}
	if counts[5] == 1 {
		h.Type = "Five of a Kind"
		log.Debugf("[%s]: counts[5] == 1. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else if counts[4] == 1 {
		if jokers == 1 {
			h.Type = "Five of a Kind"
		} else {
			h.Type = "Four of a Kind"
		}
		log.Debugf("[%s]: counts[4] == 1. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else if counts[3] == 1 && counts[2] == 1 {
		h.Type = "Full House"
		log.Debugf("[%s]: counts[x] = FH. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else if counts[3] == 1 {
		if jokers == 2 {
			h.Type = "Five of a Kind"
		} else if jokers == 1 {
			h.Type = "Four of a Kind"
		} else {
			h.Type = "Three of a Kind"
		}
		log.Debugf("[%s]: counts[3] == 1. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else if counts[2] == 2 {
		if jokers == 1 {
			h.Type = "Full House"
		} else {
			h.Type = "Two Pair"
		}
		log.Debugf("[%s]: counts[2] == 2. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else if counts[2] == 1 {
		if jokers == 3 {
			h.Type = "Five of a Kind"
		} else if jokers == 2 {
			h.Type = "Four of a Kind"
		} else if jokers == 1 {
			h.Type = "Three of a Kind"
		} else {
			h.Type = "One Pair"
		}
		log.Debugf("[%s]: counts[2] == 1. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	} else {
		if jokers >= 4 {
			h.Type = "Five of a Kind"
		} else if jokers == 3 {
			h.Type = "Four of a Kind"
		} else if jokers == 2 {
			h.Type = "Three of a Kind"
		} else if jokers == 1 {
			h.Type = "One Pair"
		} else {
			h.Type = "High Card"
		}
		log.Debugf("[%s]: counts[x] == x. jokers: %d ==> %s\n", h.Hand, jokers, h.Type)
	}
	return h.Type
}

func part1(input []string) (string, error) {
	hands := parseInput(input)

	for _, hand := range hands {
		hand.SetStrength(1)
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Strength < hands[j].Strength
	})

	var result int
	for i, hand := range hands {
		log.Debugf("%d: %s (%s)\n", i+1, hand.Hand, hand.Strength)
		result += (i + 1) * hand.Bid
	}

	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	hands := parseInput(input)

	for _, hand := range hands {
		hand.SetStrength(2)
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Strength < hands[j].Strength
	})

	var result int
	for i, hand := range hands {
		log.Debugf("%d: %s - %s (%s)\n", i+1, hand.Hand, hand.Type, hand.Strength)
		result += (i + 1) * hand.Bid
	}

	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) []*Hand {
	hands := make([]*Hand, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " ")
		hands = append(hands, &Hand{
			Hand:  parts[0],
			Cards: strings.Split(parts[0], ""),
			Bid:   util.Atoi(parts[1]),
		})
	}
	return hands
}

func init() {
	challenges.RegisterChallengeFunc(2023, 7, 1, "day07.txt", part1)
	challenges.RegisterChallengeFunc(2023, 7, 2, "day07.txt", part2)
}
