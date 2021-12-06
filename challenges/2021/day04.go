package challenges2021

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
)

func day04_part1(input []string) (string, error) {
	nums, boards := day04_parse(input)
	winningBoard, lastCall := PlayBingo(nums, boards, true)

	log.Printf("\n\n==================\n\n%s\n\n", winningBoard)

	sum := 0
	for i, row := range winningBoard.board {
		for j, n := range row {
			if winningBoard.marks[i][j] == 0 {
				sum += n
			}
		}
	}
	result := sum * lastCall

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

func PlayBingo(calls []int, boards []*day04_board, firstWinner bool) (*day04_board, int) {
	var bingoCount int
	for _, call := range calls {
		fmt.Println("\n\n----------------------")
		fmt.Printf("CALL: %d\n", call)
		fmt.Println("----------------------")
		for i, board := range boards {
			if !board.bingo {
				if MarkAndCheck(board, call) {
					board.bingo = true
					bingoCount++
					fmt.Printf("\n%s\n%s\n", log.Colorize(fmt.Sprintf("BOARD %d - BINGO", i), log.Red, 0), board)
					if firstWinner || bingoCount == len(boards) {
						return board, call
					}
				} else {
					fmt.Printf("\nBOARD %d:\n%s\n", i, board)
				}
			}
		}
	}
	return &day04_board{}, 0
}

func CalcScore(board *day04_board, lastCall int) {
	//
}

func MarkAndCheck(board *day04_board, call int) bool {
	for i, row := range board.board {
		for j, n := range row {
			if n == call {
				board.marks[i][j] = 1
				return CheckBingo(board, i, j)
			}
		}
	}
	return false
}

func CheckBingo(board *day04_board, rowNum int, colNum int) bool {
	colBingo := true
	for i, row := range board.marks {
		if i == rowNum {
			rowBingo := true
			for _, n := range row {
				if n == 0 {
					rowBingo = false
					break
				}
			}
			if rowBingo {
				return true
			}
		}
		if row[colNum] == 0 {
			colBingo = false
		}
	}
	return colBingo
}

func day04_part2(input []string) (string, error) {
	nums, boards := day04_parse(input)
	winningBoard, lastCall := PlayBingo(nums, boards, false)

	log.Printf("\n\n==================\n\n%s\n\n", winningBoard)

	sum := 0
	for i, row := range winningBoard.board {
		for j, n := range row {
			if winningBoard.marks[i][j] == 0 {
				sum += n
			}
		}
	}
	result := sum * lastCall

	log.Debugf("Result: %d\n", result)
	return fmt.Sprintf("%d", result), nil
}

var whitespace = regexp.MustCompile(`\s+`)

func day04_parse(input []string) ([]int, []*day04_board) {
	nums := make([]int, 0)
	strNums := strings.Split(input[0], ",")
	for _, s := range strNums {
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	boards := make([]*day04_board, 0)

	curBoard := NewBoard()
	for i := 2; i < len(input); i++ {
		line := strings.Trim(input[i], " ")
		if len(line) == 0 {
			boards = append(boards, curBoard)
			curBoard = NewBoard()
			continue
		}

		strNums = whitespace.Split(line, -1)
		row := make([]int, 0)
		for _, s := range strNums {
			n, _ := strconv.Atoi(s)
			row = append(row, n)
		}
		curBoard.board = append(curBoard.board, row)
	}
	boards = append(boards, curBoard)

	return nums, boards
}

type day04_board struct {
	board [][]int
	marks [][]int
	bingo bool
}

func NewBoard() *day04_board {
	marks := make([][]int, 5)
	for i := range marks {
		marks[i] = make([]int, 5)
	}
	return &day04_board{
		board: make([][]int, 0),
		marks: marks,
		bingo: false,
	}
}

func (b day04_board) String() string {
	var sb strings.Builder
	for i, row := range b.board {
		if i != 0 {
			sb.WriteString("\n")
		}
		for j, n := range row {
			if j > 0 {
				sb.WriteString(" ")
			}
			if b.marks[i][j] == 1 {
				sb.WriteString(log.Colorize(n, log.Green, 2))
			} else {
				fmt.Fprintf(&sb, "%2d", n)
			}
		}
	}
	return sb.String()
}

func init() {
	challenges.RegisterChallengeFunc(2021, 4, 1, "day04.txt", day04_part1)
	challenges.RegisterChallengeFunc(2021, 4, 2, "day04.txt", day04_part2)
}
