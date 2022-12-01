package aoc2021_day21

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	game := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", game)

	turn := 0
	for !game.HasWinner() {
		rolls, roll := game.Die.Roll(3)
		space := game.Positions[turn].Advance(roll)
		game.Scores[turn] += space
		log.Debugf("Player %d rolls %d+%d+%d (%d) and moves to space %d - total score: %d\n",
			turn+1, rolls[0], rolls[1], rolls[2], roll, space, game.Scores[turn])
		turn = (turn + 1) % 2
	}
	losingScore := game.GetLosingScore()
	rollCount := game.Die.Count()
	log.Debugln("-------------------")
	log.Debugln(game)
	result := losingScore * rollCount
	return fmt.Sprintf("%d", result), nil
}

/*
1+1+1
1+1+2    1+2+1    2+1+1
1+2+2    2+1+2    2+2+1    3+1+1    1+1+3    1+3+1
1+2+3    1+3+2    2+1+3    2+3+1    3+1+2    3+2+1    2+2+2
2+2+3    2+3+2    3+2+2    3+3+1    3+1+3    1+3+3
3+3+2    3+2+3    2+3+3
3+3+3
*/

var rollOccurrences = map[int8]int{
	3: 1,
	4: 3,
	5: 6,
	6: 7,
	7: 6,
	8: 3,
	9: 1,
}

func part2(input []string) (string, error) {
	game := parseInput(input)

	p1_wins := 0
	p2_wins := 0
	pos := [2]int8{int8(game.Positions[0].Value()), int8(game.Positions[1].Value())}

	for threeRollSum, occurrences := range rollOccurrences {
		splitTheUniverse(0, threeRollSum, occurrences, pos, [2]int16{0, 0}, [2]*int{&p1_wins, &p2_wins})
	}
	result := util.Max(p1_wins, p2_wins)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *Game {
	p1 := 0
	p2 := 0
	for i, s := range input {
		num := &p1
		if i == 1 {
			num = &p2
		}
		idx := 0
		fmt.Sscanf(s, "Player %d starting position: %d", &idx, num)
	}
	return &Game{
		Scores:    make(map[int]int),
		Positions: map[int]*Position{0: NewPosition(p1), 1: NewPosition(p2)},
		Die:       NewDie(),
	}
}

func splitTheUniverse(player int8, threeRollSum int8, cumulativeCount int, pos [2]int8, scores [2]int16, wins [2]*int) {
	newPos := pos[player] + threeRollSum
	for newPos > 10 {
		newPos -= 10
	}
	newScore := scores[player] + int16(newPos)
	if newScore >= 21 {
		// log.Debugf("Player %d wins %d more times\n", player+1, cumulativeCount)
		*(wins[player]) += cumulativeCount
	} else {
		otherPlayer := (player + 1) % 2
		posCopy := [2]int8{pos[0], pos[1]}
		posCopy[player] = newPos
		scoresCopy := [2]int16{scores[0], scores[1]}
		scoresCopy[player] = newScore
		for nextRoll, occurrences := range rollOccurrences {
			splitTheUniverse(otherPlayer, nextRoll, cumulativeCount*occurrences, posCopy, scoresCopy, wins)
		}
	}
}

func NewDie() *Die {
	return &Die{Val: 1}
}

type Game struct {
	Scores    map[int]int
	Positions map[int]*Position
	Die       *Die
}

func (g *Game) HasWinner() bool {
	return g.Scores[0] >= 1000 || g.Scores[1] >= 1000
}

func (g *Game) GetLosingScore() int {
	if g.Scores[0] >= 1000 {
		return g.Scores[1]
	} else if g.Scores[1] >= 1000 {
		return g.Scores[0]
	} else {
		return -1
	}
}

func (g *Game) String() string {
	return fmt.Sprintf(`GAME STATE:
--------------------
Player 1 score: %4d
Player 2 score: %4d

Player 1 position: %2d
Player 2 position: %2d

Die Value:   %3d
Roll Count: %4d
`, g.Scores[0], g.Scores[1], g.Positions[0].Value(), g.Positions[1].Value(), g.Die.Val, g.Die.Count())
}

type Die struct {
	Val       int
	rollCount int
}

type Position struct {
	val int
}

func NewPosition(start int) *Position {
	return &Position{
		val: start,
	}
}

func (p *Position) Advance(n int) int {
	for i := 0; i < n; i++ {
		p.val++
		if p.val > 10 {
			p.val -= 10
		}
	}
	return p.val
}

func (p *Position) Value() int {
	return p.val
}

func (d *Die) Roll(n int) ([]int, int) {
	if d.Val == 0 {
		d.Val = 1
	}
	sum := 0
	rolls := make([]int, 0, n)
	for i := 0; i < n; i++ {
		sum += d.Val
		rolls = append(rolls, d.Val)
		d.rollCount++
		d.Val += 1
		if d.Val > 100 {
			d.Val -= 100
		}
	}
	return rolls, sum
}

func (d *Die) Count() int {
	return d.rollCount
}

func (d *Die) Reset() {
	d.rollCount = 0
	d.Val = 1
}

func init() {
	challenges.RegisterChallengeFunc(2021, 21, 1, "day21.txt", part1)
	challenges.RegisterChallengeFunc(2021, 21, 2, "day21.txt", part2)
}
