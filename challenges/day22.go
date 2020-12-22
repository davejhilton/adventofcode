package challenges

import (
	"fmt"
	"strconv"

	"github.com/davejhilton/adventofcode2020/log"
)

func day22_part1(input []string) (string, error) {
	deck1, deck2 := day22_parse(input)

	log.Debugf("PLAYER 1: %v\n", deck1)
	log.Debugf("PLAYER 2: %v\n", deck2)

	round := 0

	for len(deck1) > 0 && len(deck2) > 0 {
		c1 := deck1[0]
		c2 := deck2[0]
		deck1 = deck1[1:]
		deck2 = deck2[1:]

		if c1 > c2 {
			deck1 = append(deck1, c1, c2)
		} else {
			deck2 = append(deck2, c2, c1)
		}
		round++
	}

	log.Debugf("AFTER %d ROUNDS:\n\n", round)
	log.Debugf("PLAYER 1: %v\n", deck1)
	log.Debugf("PLAYER 2: %v\n", deck2)

	winningDeck := deck1
	if len(deck1) == 0 {
		winningDeck = deck2
	}

	sum := 0
	for i := len(winningDeck) - 1; i >= 0; i-- {
		m := len(winningDeck) - i
		sum += winningDeck[i] * m
	}

	return fmt.Sprintf("%d", sum), nil
}

var gameCount = 0

func day22_part2(input []string) (string, error) {
	deck1, deck2 := day22_parse(input)

	_, points := playGame(deck1, deck2)

	return fmt.Sprintf("%d", points), nil
}

func playGame(deck1 []int, deck2 []int) (int, int) {

	gameCount++
	log.Debugf("GAME %d\n", gameCount)
	log.Debugf("PLAYER 1: %v\n", deck1)
	log.Debugf("PLAYER 2: %v\n", deck2)

	prevRounds := make(map[string]bool)
	round := 0

	for len(deck1) > 0 && len(deck2) > 0 {

		startingConfig := fmt.Sprintf("%v vs %v", deck1, deck2)
		if _, ok := prevRounds[startingConfig]; ok {
			return 1, 1
		} else {
			prevRounds[startingConfig] = true
		}

		c1 := deck1[0]
		c2 := deck2[0]
		deck1 = deck1[1:]
		deck2 = deck2[1:]

		log.Debugf("Game %d, Round %d: P1: %d vs P2: %d\n", gameCount, round, c1, c2)
		if c1 <= len(deck1) && c2 <= len(deck2) {
			log.Debugf("Recursing into subgame!\n")
			d1Copy := make([]int, c1)
			copy(d1Copy, deck1[:c1])
			d2Copy := make([]int, c2)
			copy(d2Copy, deck2[:c2])
			winner, _ := playGame(d1Copy, d2Copy)
			if winner == 1 {
				log.Debugf("Player 1 wins round %d, game %d\n", round, gameCount)
				deck1 = append(deck1, c1, c2)
			} else {
				log.Debugf("Player 2 wins round %d, game %d\n", round, gameCount)
				deck2 = append(deck2, c2, c1)
			}
		} else if c1 > c2 {
			log.Debugf("Player 1 wins round %d, game %d\n", round, gameCount)
			deck1 = append(deck1, c1, c2)
		} else {
			log.Debugf("Player 2 wins round %d, game %d\n", round, gameCount)
			deck2 = append(deck2, c2, c1)
		}
		round++
	}

	winner := 1
	var winningDeck []int
	if len(deck1) != 0 {
		log.Debugf("Player 1 wins game %d (rounds: %d)\n", gameCount, round)
		winningDeck = deck1
	} else {
		log.Debugf("Player 2 wins game %d (rounds: %d)\n", gameCount, round)
		winningDeck = deck2
		winner = 2
	}

	points := 0
	for i := len(winningDeck) - 1; i >= 0; i-- {
		m := len(winningDeck) - i
		points += winningDeck[i] * m
	}
	return winner, points
}

func day22_parse(input []string) ([]int, []int) {

	var p1Cards []int

	cards := make([]int, 0)
	for i := 1; i < len(input); i++ {
		var s = input[i]
		if s == "" {
			i++
			p1Cards = cards
			cards = make([]int, 0)
			continue
		}
		n, _ := strconv.Atoi(s)
		cards = append(cards, n)
	}
	return p1Cards, cards
}

func init() {
	registerChallengeFunc(22, 1, "day22.txt", day22_part1)
	registerChallengeFunc(22, 2, "day22.txt", day22_part2)
}
