package aoc2022_day22

import (
	"fmt"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func part1(input []string) (string, error) {
	grid, instructions := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", grid)

	var rows = make(map[int]Row)
	var cols = make(map[int]Col)

	var colStarts = make(map[int]int)
	var colEnds = make(map[int]int)
	var colPrevs = make(map[int]int)
	for i := 0; i < len(grid); i++ {
		s := 0
		e := 0
		prev := VOID
		for j := 0; j < len(grid[i]); j++ {
			if i == 0 {
				colPrevs[j] = VOID
			}
			c := grid[i][j]
			switch c {
			case VOID:
				if prev != VOID {
					e = j - 1
				}
				if colPrevs[j] != VOID {
					colEnds[j] = i - 1
				}
			case EMPTY:
				if prev == VOID {
					s = j
				}
				if colPrevs[j] == VOID {
					colStarts[j] = i
				}
			case WALL:
				if prev == VOID {
					s = j
				}
				if colPrevs[j] == VOID {
					colStarts[j] = i
				}
			}
			prev = c
			colPrevs[j] = c
		}
		if prev != VOID {
			e = len(grid[i]) - 1
		}
		rows[i] = Row{s, e}
		for j, c := range colPrevs {
			if c != VOID {
				colEnds[j] = len(grid) - 1
			}
			cols[j] = Col{colStarts[j], colEnds[j]}
		}
	}

	var row = 0
	var col = rows[0].Start
	var facing = R

	for _, inst := range instructions {
		if inst.Type == MOVE {
			switch facing {
			case R:
				thisRow := rows[row]
				for m := 0; m < inst.Value; m++ {
					next := col + 1
					if col == thisRow.End {
						next = thisRow.Start
					}
					if grid[row][next] != WALL {
						col = next
					}
				}
			case L:
				thisRow := rows[row]
				for m := 0; m < inst.Value; m++ {
					next := col - 1
					if col == thisRow.Start {
						next = thisRow.End
					}
					if grid[row][next] != WALL {
						col = next
					}
				}
			case U:
				thisCol := cols[col]
				for m := 0; m < inst.Value; m++ {
					next := row - 1
					if row == thisCol.Start {
						next = thisCol.End
					}
					if grid[next][col] != WALL {
						row = next
					}
				}
			case D:
				thisCol := cols[col]
				for m := 0; m < inst.Value; m++ {
					next := row + 1
					if row == thisCol.End {
						next = thisCol.Start
					}
					if grid[next][col] != WALL {
						row = next
					}
				}
			}
		} else if inst.Type == TURN {
			switch facing {
			case R:
				if inst.Value == R {
					facing = D
				} else {
					facing = U
				}
			case L:
				if inst.Value == R {
					facing = U
				} else {
					facing = D
				}
			case U:
				if inst.Value == R {
					facing = R
				} else {
					facing = L
				}
			case D:
				if inst.Value == R {
					facing = L
				} else {
					facing = R
				}
			}
		}
	}

	var result = (1000 * (row + 1)) + (4 * (col + 1)) + facing
	return fmt.Sprintf("%d", result), nil
}

func part2(input []string) (string, error) {
	grid, instructions := parseInput(input)
	// log.Debugf("Parsed Input:\n%v\n", grid)

	var lines map[string]*CubeLine
	if len(grid) <= 16 {
		lines = setupPart2_example()
	} else {
		lines = setupPart2()
	}

	var row = 0
	var col = lines["row0"].Start
	var facing = R

	log.Debugf("\trow: %d, col: %d, facing: %s\n", row, col, DIR_NAMES[facing])
	for _, inst := range instructions {
		if inst.Type == MOVE {
			log.Debugf("INSTR: %v\n", inst)
			for m := 0; m < inst.Value; m++ {
				if grid[row][col] == VOID {
					log.Debugf("I have strayed off the path! r: %d, c: %d (m: %d) facing: %s\n", row, col, m, DIR_NAMES[facing])
				}
				switch facing {
				case R:
					thisRow := lines[fmt.Sprintf("row%d", row)]
					nextCol := col + 1
					nextRow := row
					shift := false
					if col == thisRow.End {
						nextCol = thisRow.AfterCoord.X
						nextRow = thisRow.AfterCoord.Y
						shift = true
					}
					if grid[nextRow][nextCol] != WALL {
						col = nextCol
						row = nextRow
						if shift {
							facing = thisRow.AfterDir
						}
					}
					// log.Debugf("row: %d, col: %d, nextCol: %d, nextRow: %d, thisRow: %#v, afterCoord: %v\n", row, col, nextCol, nextRow, thisRow, thisRow.AfterLine.AfterCoord)
				case L:
					thisRow := lines[fmt.Sprintf("row%d", row)]
					nextCol := col - 1
					nextRow := row
					shift := false
					if col == thisRow.Start {
						nextCol = thisRow.BeforeCoord.X
						nextRow = thisRow.BeforeCoord.Y
						shift = true
					}
					if grid[nextRow][nextCol] != WALL {
						col = nextCol
						row = nextRow
						if shift {
							facing = thisRow.BeforeDir
						}
					}
				case U:
					thisCol := lines[fmt.Sprintf("col%d", col)]
					nextRow := row - 1
					nextCol := col
					shift := false
					if row == thisCol.Start {
						nextCol = thisCol.BeforeCoord.X
						nextRow = thisCol.BeforeCoord.Y
						shift = true
					}
					// log.Debugf("row: %d, col: %d, nextCol: %d, nextRow: %d, thisCol: %#v, beforeCoord: %v\n", row, col, nextCol, nextRow, thisCol, thisCol.BeforeCoord)
					if grid[nextRow][nextCol] != WALL {
						col = nextCol
						row = nextRow
						if shift {
							facing = thisCol.BeforeDir
						}
					}
				case D:
					thisCol := lines[fmt.Sprintf("col%d", col)]
					nextRow := row + 1
					nextCol := col
					shift := false
					if row == thisCol.End {
						nextCol = thisCol.AfterCoord.X
						nextRow = thisCol.AfterCoord.Y
						shift = true
					}
					if grid[nextRow][nextCol] != WALL {
						col = nextCol
						row = nextRow
						if shift {
							facing = thisCol.AfterDir
						}
					}
				}
			}
		} else if inst.Type == TURN {
			log.Debugf("INSTR: {%s %s}\n", inst.Type, DIR_NAMES[inst.Value])
			switch facing {
			case R:
				if inst.Value == R {
					facing = D
				} else {
					facing = U
				}
			case L:
				if inst.Value == R {
					facing = U
				} else {
					facing = D
				}
			case U:
				if inst.Value == R {
					facing = R
				} else {
					facing = L
				}
			case D:
				if inst.Value == R {
					facing = L
				} else {
					facing = R
				}
			}
		}
		log.Debugf("\trow: %d, col: %d, facing: %s\n", row, col, DIR_NAMES[facing])
	}

	var result = (1000 * (row + 1)) + (4 * (col + 1)) + facing
	return fmt.Sprintf("%d", result), nil
}

func parseInput(input []string) (Grid, []Instruction) {
	grid := make(Grid, len(input)-2)
	for i := range grid {
		grid[i] = make([]int, len(input[0]))
	}
	var instructionStr []rune
	for i, s := range input {
		if s == "" {
			instructionStr = []rune(input[i+1])
			break
		}
		for j, c := range []rune(s) {
			switch c {
			case ' ':
				grid[i][j] = VOID
			case '.':
				grid[i][j] = EMPTY
			case '#':
				grid[i][j] = WALL
			}
		}
	}

	instructions := make([]Instruction, 0)
	for i := 0; i < len(instructionStr); i++ {
		if d, ok := DIRS[instructionStr[i]]; ok {
			instructions = append(instructions, Instruction{TURN, d})
		} else {
			n := DIGITS[instructionStr[i]]
			if i < len(instructionStr)-1 {
				if n2, ok := DIGITS[instructionStr[i+1]]; ok {
					n = n*10 + n2
					i++
				}
			}
			instructions = append(instructions, Instruction{MOVE, n})
		}
	}

	return grid, instructions
}

const (
	VOID int = iota
	EMPTY
	WALL

	R int = 0
	D int = 1
	L int = 2
	U int = 3

	TURN string = "TURN"
	MOVE string = "MOVE"

	ROW int = 0
	COL int = 1
)

var DIR_NAMES = map[int]string{
	R: "R",
	D: "D",
	L: "L",
	U: "U",
}

var (
	DIGITS = map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
	}
	DIRS = map[rune]int{
		'R': R,
		'D': D,
		'L': L,
		'U': U,
	}
)

type Instruction struct {
	Type  string
	Value int
}

type Grid [][]int

type Coord struct {
	X int
	Y int
}

type Row struct {
	Start int
	End   int
}

type Col struct {
	Start int
	End   int
}

func setupPart2() map[string]*CubeLine {
	var lines = make(map[string]*CubeLine)
	// rows
	start := 50
	end := 149
	bx := 149
	by := 0
	bd := R
	ax := 149
	ay := 99
	ad := L
	for i := 0; i < 200; i++ {
		if i < 50 {
			start = 50
			end = 149
			bx = 0
			by = 149 - i
			bd = R
			ax = 99
			ay = 149 - i
			ad = L
		} else if i >= 50 && i < 100 {
			start = 50
			end = 99
			bx = i - 50
			by = 100
			bd = D
			ax = i + 50
			ay = 49
			ad = U
		} else if i >= 100 && i < 150 {
			start = 0
			end = 99
			bx = 50
			by = 149 - i
			bd = R
			ax = 149
			ay = 149 - i
			ad = L
		} else if i >= 150 {
			start = 0
			end = 49
			bx = i - 100
			by = 0
			bd = D
			ax = i - 100
			ay = 149
			ad = U
		}

		r := &CubeLine{
			Index:       i,
			Start:       start,
			End:         end,
			Type:        ROW,
			BeforeCoord: Coord{bx, by},
			AfterCoord:  Coord{ax, ay},
			BeforeDir:   bd,
			AfterDir:    ad,
		}
		lines[r.Key()] = r
	}

	for i := 0; i < 150; i++ {
		if i < 50 {
			start = 100
			end = 199
			bx = 50
			by = i + 50
			bd = R
			ax = i + 100
			ay = 0
			ad = D
		} else if i >= 50 && i < 100 {
			start = 0
			end = 149
			bx = 0
			by = i + 100
			bd = R
			ax = 49
			ay = i + 100
			ad = L
		} else if i >= 100 {
			start = 0
			end = 49
			bx = i - 100
			by = 199
			bd = U
			ax = 99
			ay = i - 50
			ad = L
		}

		c := &CubeLine{
			Index:       i,
			Start:       start,
			End:         end,
			Type:        COL,
			BeforeCoord: Coord{bx, by},
			AfterCoord:  Coord{ax, ay},
			BeforeDir:   bd,
			AfterDir:    ad,
		}
		lines[c.Key()] = c
	}

	// log.Debugln(util.Keys(lines))

	key := ""
	bkey := ""
	akey := ""
	var l *CubeLine
	for i := 0; i < 200; i++ {
		if i < 50 {
			bkey = fmt.Sprintf("row%d", 149-i)
			akey = fmt.Sprintf("row%d", 149-i)
		} else if i >= 50 && i < 100 {
			bkey = fmt.Sprintf("col%d", i-50)
			akey = fmt.Sprintf("col%d", i+50)
		} else if i >= 100 && i < 150 {
			bkey = fmt.Sprintf("row%d", 149-i)
			akey = fmt.Sprintf("row%d", 149-i)
		} else if i >= 150 {
			bkey = fmt.Sprintf("col%d", i-100)
			akey = fmt.Sprintf("col%d", i-100)
		}

		key = fmt.Sprintf("row%d", i)
		l = lines[key]
		// log.Debugf("key: %s | before key: %s | after key: %s\n", key, bkey, akey)
		l.BeforeLine = lines[bkey]
		l.AfterLine = lines[akey]
	}

	for i := 0; i < 150; i++ {
		if i < 50 {
			bkey = fmt.Sprintf("row%d", i+50)
			akey = fmt.Sprintf("col%d", i+100)
		} else if i >= 50 && i < 100 {
			bkey = fmt.Sprintf("row%d", i+100)
			akey = fmt.Sprintf("row%d", i+100)
		} else if i >= 100 {
			bkey = fmt.Sprintf("col%d", i-100)
			akey = fmt.Sprintf("row%d", i-50)
		}

		key = fmt.Sprintf("col%d", i)
		l = lines[key]
		l.BeforeLine = lines[bkey]
		l.AfterLine = lines[akey]
	}
	return lines
}

func setupPart2_example() map[string]*CubeLine {
	var lines = make(map[string]*CubeLine)
	// rows
	var start int
	var end int
	var bc int
	var br int
	var bd int
	var ac int
	var ar int
	var ad int
	for i := 0; i < 12; i++ {
		if i < 4 {
			start = 8
			end = 11
			bc = i + 4
			br = 4
			bd = D
			ac = 15
			ar = 11 - i
			ad = L
		} else if i >= 4 && i < 8 {
			start = 0
			end = 11
			bc = 19 - i
			br = 11
			bd = U
			ac = 19 - i
			ar = 8
			ad = D
		} else if i >= 8 && i < 12 {
			start = 8
			end = 15
			bc = 15 - i
			br = 7
			bd = U
			ac = 11
			ar = 11 - i
			ad = L
		}

		r := &CubeLine{
			Index:       i,
			Start:       start,
			End:         end,
			Type:        ROW,
			BeforeCoord: Coord{bc, br},
			AfterCoord:  Coord{ac, ar},
			BeforeDir:   bd,
			AfterDir:    ad,
		}
		lines[r.Key()] = r
	}

	for i := 0; i < 16; i++ {
		if i < 4 {
			start = 4
			end = 7
			bc = 11 - i
			br = 0
			bd = D
			ac = 11 - i
			ar = 11
			ad = U
		} else if i >= 4 && i < 8 {
			start = 4
			end = 7
			bc = 8
			br = i - 4
			bd = R
			ac = 8
			ar = 15 - i
			ad = R
		} else if i >= 8 {
			start = 0
			end = 11
			bc = 11 - i
			br = 4
			bd = D
			ac = 11 - i
			ar = 7
			ad = U
		} else if i >= 12 {
			start = 8
			end = 11
			bc = 11
			br = 19 - i
			bd = L
			ac = 0
			ar = 19 - i
			ad = R
		}

		c := &CubeLine{
			Index:       i,
			Start:       start,
			End:         end,
			Type:        COL,
			BeforeCoord: Coord{bc, br},
			AfterCoord:  Coord{ac, ar},
			BeforeDir:   bd,
			AfterDir:    ad,
		}
		lines[c.Key()] = c
	}

	// log.Debugln(util.Keys(lines))

	key := ""
	bkey := ""
	akey := ""
	var l *CubeLine
	for i := 0; i < 12; i++ {
		if i < 4 {
			bkey = fmt.Sprintf("col%d", i+4)
			akey = fmt.Sprintf("row%d", 11-i)
		} else if i >= 4 && i < 8 {
			bkey = fmt.Sprintf("col%d", 19-i)
			akey = fmt.Sprintf("col%d", 19-i)
		} else if i >= 8 && i < 12 {
			bkey = fmt.Sprintf("col%d", 15-i)
			akey = fmt.Sprintf("row%d", 11-i)
		}

		key = fmt.Sprintf("row%d", i)
		l = lines[key]
		// log.Debugf("key: %s | before key: %s | after key: %s\n", key, bkey, akey)
		l.BeforeLine = lines[bkey]
		l.AfterLine = lines[akey]
	}

	for i := 0; i < 16; i++ {
		if i < 4 {
			bkey = fmt.Sprintf("col%d", 11-i)
			akey = fmt.Sprintf("col%d", 11-i)
		} else if i >= 4 && i < 8 {
			bkey = fmt.Sprintf("row%d", i-4)
			akey = fmt.Sprintf("row%d", 15-i)
		} else if i >= 8 {
			bkey = fmt.Sprintf("col%d", 11-i)
			akey = fmt.Sprintf("col%d", 11-i)
		} else if i >= 12 {
			bkey = fmt.Sprintf("row%d", 19-i)
			akey = fmt.Sprintf("row%d", 19-i)
		}

		key = fmt.Sprintf("col%d", i)
		l = lines[key]
		l.BeforeLine = lines[bkey]
		l.AfterLine = lines[akey]
	}
	return lines
}

type CubeLine struct {
	Index       int
	Start       int
	End         int
	Type        int // ROW or COL
	AfterCoord  Coord
	BeforeCoord Coord
	BeforeLine  *CubeLine
	AfterLine   *CubeLine
	AfterDir    int // R, D, L, U
	BeforeDir   int
}

func (c CubeLine) Key() string {
	t := "row"
	if c.Type == COL {
		t = "col"
	}
	return fmt.Sprintf("%s%d", t, c.Index)
}

type Line interface {
	Get(int) int
}

func Move[T Row | Col](c *Coord, magnitude int, direction int) {

}

func init() {
	challenges.RegisterChallengeFunc(2022, 22, 1, "day22.txt", part1)
	challenges.RegisterChallengeFunc(2022, 22, 2, "day22.txt", part2)
}
