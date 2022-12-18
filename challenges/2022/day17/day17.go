package aoc2022_day17

import (
	"fmt"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

type Direction int
type RockType int
type RowType int

func part1(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", *grid)

	for grid.RocksAbsorbed < 2022 {
		grid.DropNewRock()
		for {
			grid.AirJetPush()
			if !grid.MoveRockDown() {
				break
			}
		}
		grid.AbsorbRock()
	}

	var result = grid.Height
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid := parseInput(input)
	log.Debugf("Parsed Input:\n%v\n", *grid)

	var pattern = new(Pattern)
	var found bool
	for !found && grid.RocksAbsorbed < 10000 {
		grid.DropNewRock()
		for {
			grid.AirJetPush()
			if !grid.MoveRockDown() {
				break
			}
		}
		grid.AbsorbRock()

		found = pattern.CheckForPattern(grid.RocksAbsorbed, grid.Height)
	}

	if pattern.Found {
		log.Debugf(
			"FOUND A PATTERN!\nstarts at: %d, length: %d, pattern:\n%s\n\n%s\n\n",
			pattern.StartsAt,
			pattern.Length,
			string(pattern.chars[pattern.StartsAt:pattern.StartsAt+pattern.Length]),
			string(pattern.chars[pattern.StartsAt+pattern.Length:pattern.StartsAt+pattern.Length+pattern.Length]),
		)
	}

	var result = pattern.Extrapolate(1_000_000_000_000)
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) *Grid {
	dirs := make([]Direction, 0, len(input[0]))
	charMap := map[rune]Direction{
		'<': LEFT,
		'>': RIGHT,
	}
	for _, c := range input[0] {
		dirs = append(dirs, charMap[c])
	}
	return &Grid{
		Height:          0,
		RockOffset:      0,
		Rock:            nil,
		TopRow:          nil,
		AirJets:         dirs,
		NextAirJetIdx:   0,
		NextRockTypeIdx: 0,
		RocksAbsorbed:   0,
	}
}

const (
	GRID_WIDTH int = 7

	LEFT Direction = iota
	RIGHT

	ROCK_DASH RockType = iota
	ROCK_PLUS
	ROCK_J
	ROCK_I
	ROCK_SQUARE

	ROCK RowType = iota
	GRID
)

var (
	ROCK_TYPES = [5]RockType{
		ROCK_DASH,
		ROCK_PLUS,
		ROCK_J,
		ROCK_I,
		ROCK_SQUARE,
	}

	ROCK_CHARS = map[int]rune{
		0: '.',
		1: '@',
	}
	GRID_CHARS = map[int]rune{
		0: '.',
		1: '#',
	}

	prevHeight  int
	heightDiffs = make([]int, 0)

	separator = "=========="
)

type Row struct {
	R        []int
	RowAbove *Row
	RowBelow *Row
	Type     RowType
}

func (r Row) String() string {
	var b strings.Builder
	for i := 0; i < len(r.R); i++ {
		if r.Type == ROCK {
			b.WriteRune(ROCK_CHARS[r.R[i]])
		} else {
			b.WriteRune(GRID_CHARS[r.R[i]])
		}
	}
	return b.String()
}

type Rock struct {
	Type       RockType
	Height     int
	Width      int
	LeftOffset int
	BottomRow  *Row
}

func (r Rock) TopRow() *Row {
	var rockTop = r.BottomRow
	for rockTop.RowAbove != nil {
		rockTop = rockTop.RowAbove
	}
	return rockTop
}

type Grid struct {
	TopRow          *Row
	Height          int
	Rock            *Rock
	RockOffset      int // 0 means the rock is directly on top of the grid's highest point. -1 means it's overlapping.
	AirJets         []Direction
	NextAirJetIdx   int
	NextRockTypeIdx int
	RocksAbsorbed   int
}

func (g *Grid) AirJetPush() {

	if g.NextAirJetIdx == 15 {
		log.Debug("JUST DID THE WEIRD DROP...\n")
		log.Debugf("RockOffset: %d, LeftOffset: %d\n\n", g.RockOffset, g.Rock.LeftOffset)
	}

	var r *Rock = g.Rock
	var d Direction = g.AirJets[g.NextAirJetIdx]
	dirName := "RIGHT"
	if d == LEFT {
		dirName = "LEFT"
	}

	// if g.NextAirJetIdx+1 >= len(g.AirJets) {
	// 	// log.Printf("Repeating air jet sequence! Height: %d, Rocks: %d\n", g.Height, g.RocksAbsorbed)
	// 	resets = append(resets, []int{g.Height, g.RocksAbsorbed, g.RockOffset})
	// }
	g.NextAirJetIdx = (g.NextAirJetIdx + 1) % len(g.AirJets)

	if (d == LEFT && r.LeftOffset == 0) || (d == RIGHT && r.LeftOffset+r.Width == GRID_WIDTH) {
		log.Debugf("TRIED MOVING ROCK %s, BUT IT COLLIDED WITH THE WALL\n\n", dirName)
		return
	}

	if g.RockOffset >= 0 {
		// rock is still completely above the grid. just push it.
		if d == LEFT {
			r.LeftOffset -= 1
		} else {
			r.LeftOffset += 1
		}
	} else {
		gridRow := g.TopRow
		rockRow := r.TopRow()
		gridRowIdx := g.RockOffset + r.Height - 1
		if gridRowIdx <= -1 {
			for i := -1; i > gridRowIdx; i-- {
				gridRow = gridRow.RowBelow
			}
		} else {
			for i := gridRowIdx; i >= 0; i-- {
				rockRow = rockRow.RowBelow
			}
		}
		log.Debugf("rock Offset: %d, gridRowIdx: %d\n\n", g.RockOffset, gridRowIdx)

		// gridRow and rockRow are both aligned at their uppermost overlapping row

		if d == LEFT {
			for rockRow != nil {
				rockLeft := 0
				for rockRow.R[rockLeft] != 1 {
					rockLeft++
				}
				if gridRow.R[rockLeft+r.LeftOffset-1] == 1 {
					log.Debug("TRIED MOVING ROCK LEFT, BUT IT COLLIDED WITH THE GRID\n\n")
					return // collided, can't move
				}
				rockRow = rockRow.RowBelow
				gridRow = gridRow.RowBelow
			}
			r.LeftOffset -= 1 // successfully moved
		} else {
			for rockRow != nil {
				rockRight := r.Width - 1
				for rockRow.R[rockRight] != 1 {
					rockRight--
				}
				log.Debugf("rock right: %d, grid check: %d\n", rockRight, rockRight+r.LeftOffset+1)
				log.Debugf("grid val: %d\n\n", gridRow.R[rockRight+r.LeftOffset+1])
				if gridRow.R[rockRight+r.LeftOffset+1] == 1 {
					log.Debug("TRIED MOVING ROCK RIGHT, BUT IT COLLIDED WITH THE GRID\n\n")
					return // collided, can't move
				}
				rockRow = rockRow.RowBelow
				gridRow = gridRow.RowBelow
			}
			r.LeftOffset += 1 // successfully moved
		}
	}
	log.Debugf("%s\nPUSHED ROCK %s:\n%s\n%s\n\n", separator, dirName, *g, separator)
}

func (g *Grid) AbsorbRock() {
	r := g.Rock

	rockRow := r.BottomRow

	// starting from the bottom of the rock, merge any overlapping rock/grid rows
	if g.RockOffset < 0 {

		// first, find the grid row that's aligned with the rock's bottom row
		gridRow := g.TopRow
		for i := -1; i > g.RockOffset; i-- {
			gridRow = gridRow.RowBelow
		}

		// merge rock rows with grid rows wherever they overlap, moving upward
		// until no more rows overlap
		for rockRow != nil && gridRow != nil {
			for i := 0; i < r.Width; i++ {
				if rockRow.R[i] == 1 {
					gridRow.R[i+r.LeftOffset] = 1
				}
			}
			rockRow = rockRow.RowAbove
			gridRow = gridRow.RowAbove
		}
	}

	// now, handle any rock rows that stick up above the existing grid's top
	for rockRow != nil {
		newRow := &Row{R: make([]int, GRID_WIDTH), Type: GRID}
		for i := 0; i < r.Width; i++ {
			if rockRow.R[i] == 1 {
				newRow.R[i+r.LeftOffset] = 1
			}
		}
		if g.TopRow != nil {
			g.TopRow.RowAbove = newRow
			newRow.RowBelow = g.TopRow
		}
		g.TopRow = newRow
		g.Height++

		rockRow = rockRow.RowAbove
	}

	g.Rock = nil
	g.RockOffset = (-1 * g.Height) - 1 // reset this to below the grid (invalid)
	g.RocksAbsorbed += 1

	heightDiffs = append(heightDiffs, g.Height-prevHeight)
	pattern.WriteString(fmt.Sprintf("%d", g.Height-prevHeight))
	prevHeight = g.Height

	log.Debugf("%s\nABSORBED ROCK:\n%s\n%s\n\n", separator, *g, separator)
}

var pattern strings.Builder

func (g *Grid) MoveRockDown() (movedWithoutColliding bool) {
	r := g.Rock

	if g.RockOffset > 0 {
		// rock is still completely above the grid
		g.RockOffset--
	} else if g.Height == 0 {
		log.Debug("TRIED MOVING ROCK DOWN, BUT IT COLLIDED WITH THE FLOOR\n\n")
		return false // this is the first rock, and it's colliding with the "floor"
	} else {
		gridRow := g.TopRow
		// move down until `gridRow` is the grid row that's JUST BELOW the rock's current bottom
		for i := -1; i >= g.RockOffset; i-- {
			gridRow = gridRow.RowBelow
		}

		if gridRow == nil {
			log.Debug("TRIED MOVING ROCK DOWN, BUT IT COLLIDED WITH THE FLOOR\n\n")
			return false
		}

		// check if the bottom row of the rock would collide with the grid row below it
		for i, v := range r.BottomRow.R {
			if v == 1 && gridRow.R[i+r.LeftOffset] == 1 {
				log.Debug("TRIED MOVING ROCK DOWN, BUT IT COLLIDED WITH THE GRID\n\n")
				return false // collided!
			}
		}

		// the '+' shaped rock has an overhang that can collide too - so check that row too!
		if r.Type == ROCK_PLUS && gridRow.RowAbove != nil {
			gridRow = gridRow.RowAbove
			for i, v := range r.BottomRow.RowAbove.R {
				if v == 1 && gridRow.R[i+r.LeftOffset] == 1 {
					log.Debug("TRIED MOVING ROCK DOWN, BUT IT COLLIDED WITH THE GRID\n\n")
					return false // collided!
				}
			}
		}
		g.RockOffset -= 1
	}

	log.Debugf("%s\nMOVED ROCK DOWN:\n%s\n%s\n\n", separator, *g, separator)
	return true
}

func (g *Grid) DropNewRock() {

	var (
		rType  = ROCK_TYPES[g.NextRockTypeIdx]
		height int
		width  int
		bottom *Row
	)
	switch rType {
	case ROCK_DASH:
		bottom = &Row{R: []int{1, 1, 1, 1}, Type: ROCK}
		height = 1
		width = 4
	case ROCK_PLUS:
		top := &Row{R: []int{0, 1, 0}, Type: ROCK}
		middle := &Row{R: []int{1, 1, 1}, Type: ROCK, RowAbove: top}
		bottom = &Row{R: []int{0, 1, 0}, Type: ROCK, RowAbove: middle}
		top.RowBelow = middle
		middle.RowBelow = bottom
		height = 3
		width = 3
	case ROCK_J:
		top := &Row{R: []int{0, 0, 1}, Type: ROCK}
		middle := &Row{R: []int{0, 0, 1}, Type: ROCK, RowAbove: top}
		bottom = &Row{R: []int{1, 1, 1}, Type: ROCK, RowAbove: middle}
		top.RowBelow = middle
		middle.RowBelow = bottom
		height = 3
		width = 3
	case ROCK_I:
		top := &Row{R: []int{1}}
		middle1 := &Row{R: []int{1}, Type: ROCK, RowAbove: top}
		middle2 := &Row{R: []int{1}, Type: ROCK, RowAbove: middle1}
		bottom = &Row{R: []int{1}, Type: ROCK, RowAbove: middle2}
		top.RowBelow = middle1
		middle1.RowBelow = middle2
		middle2.RowBelow = bottom
		height = 4
		width = 1
	case ROCK_SQUARE:
		top := &Row{R: []int{1, 1}, Type: ROCK}
		bottom = &Row{R: []int{1, 1}, Type: ROCK, RowAbove: top}
		top.RowBelow = bottom
		height = 2
		width = 2
	}

	g.RockOffset = 3
	g.NextRockTypeIdx = (g.NextRockTypeIdx + 1) % len(ROCK_TYPES)
	g.Rock = &Rock{
		Type:       rType,
		BottomRow:  bottom,
		LeftOffset: 2,
		Height:     height,
		Width:      width,
	}
	log.Debugf("%s\nSTARTED DROPPING NEW ROCK:\n%s\n%s\n\n", separator, *g, separator)
}

func (g Grid) String() string {
	stringRows := make([]string, 0)

	var r = g.Rock

	var rockRow, gridRow *Row

	if r != nil && g.RockOffset+r.Height > 0 {
		rockRow = r.TopRow()

		var b strings.Builder
		for i := 0; i < g.RockOffset+r.Height; i++ {
			if rockRow != nil {
				b.WriteString(fmt.Sprintf(
					"|%s%s%s|",
					strings.Repeat(".", r.LeftOffset),
					*rockRow,
					strings.Repeat(".", GRID_WIDTH-(r.LeftOffset+r.Width)),
				))
				stringRows = append(stringRows, b.String())
				b.Reset()
				rockRow = rockRow.RowBelow
			} else {
				stringRows = append(stringRows, "|.......|")
			}
		}
	}

	gridRow = g.TopRow
	rockOffsetIdx := -1 // for rock offset, -1 is the top of the grid, since 0 is the row above it

	for gridRow != nil && len(stringRows) < 20 { // only print the top 20 rows
		if r != nil && rockRow == nil && rockOffsetIdx == g.RockOffset+r.Height {
			rockRow = r.TopRow()
		}
		if rockRow != nil { // if g.RockOffset <= rockOffsetIdx && g.RockOffset+r.Height >= rockOffsetIdx
			// we have both a rock and part of the grid in this row

			chars := make([]rune, GRID_WIDTH)
			for i, c := range gridRow.R {
				if c == 1 {
					chars[i] = GRID_CHARS[1]
				} else {
					if i >= r.LeftOffset && i < r.LeftOffset+r.Width && rockRow.R[i-r.LeftOffset] == 1 {
						chars[i] = ROCK_CHARS[1]
					} else {
						chars[i] = '.'
					}
				}
			}
			stringRows = append(stringRows, fmt.Sprintf("|%s|", string(chars)))
		} else {
			// just the grid on this row, no rock.
			stringRows = append(stringRows, fmt.Sprintf("|%s|", *gridRow))
		}
		rockOffsetIdx--
		gridRow = gridRow.RowBelow
		if rockRow != nil {
			rockRow = rockRow.RowBelow
		}
	}

	stringRows = append(stringRows, "+-------+")
	return strings.Join(stringRows, "\n")
}

type Pattern struct {
	HeightDeltas []int
	CurHeight    int

	Found    bool
	StartsAt int
	Length   int

	chars []rune
}

func (p *Pattern) CheckForPattern(nRocks, height int) bool {
	dh := height - p.CurHeight
	p.HeightDeltas = append(p.HeightDeltas, dh)
	p.chars = append(p.chars, rune(dh+48)) // 48 is the ascii table offset for 0
	p.CurHeight = height

	l := len(p.chars)
	if l >= 500 && l%100 == 0 {
		str := string(p.chars)
		chunkIdx := l - 30
		testChunk := str[chunkIdx:]
		if prevIdx := strings.LastIndex(str[:chunkIdx], testChunk); prevIdx != -1 {
			// we've seen this (piece of the) pattern before!
			// but we can't find the full pattern unless this is at least the 3rd time this chunk repeat
			if firstIdx := strings.LastIndex(str[:prevIdx], testChunk); firstIdx != -1 {
				fullPattern := str[firstIdx:prevIdx]
				if str[prevIdx:chunkIdx] == fullPattern {
					// Tada! we've found two consecutive occurrences of the same pattern, and confirmed the start of a third
					for firstIdx >= 1 && str[firstIdx-1] == str[prevIdx-1] {
						// for the sake of correctness... let's rewind until we find the actual start of the pattern
						firstIdx--
						prevIdx--
					}

					p.Found = true
					p.StartsAt = firstIdx
					p.Length = len(fullPattern)
					return true
				}
			}
		}
	}

	return false
}

func (p *Pattern) Extrapolate(iterations int) int {
	prePatternHeight := 0
	for i := 0; i < p.StartsAt; i++ {
		prePatternHeight += p.HeightDeltas[i]
	}

	patternHeightGain := 0
	for i := 0; i < p.Length; i++ {
		patternHeightGain += p.HeightDeltas[p.StartsAt+i]
	}

	nRepeats := (iterations - p.StartsAt) / p.Length // INTEGER DIVISION (truncation is what we want)
	tailLength := (iterations - p.StartsAt) % p.Length

	tailHeightGain := 0
	for i := 0; i < tailLength; i++ {
		tailHeightGain += p.HeightDeltas[p.StartsAt+i]
	}

	return prePatternHeight + (patternHeightGain * nRepeats) + tailHeightGain
}

func init() {
	challenges.RegisterChallengeFunc(2022, 17, 1, "day17.txt", part1)
	challenges.RegisterChallengeFunc(2022, 17, 2, "day17.txt", part2)
}
