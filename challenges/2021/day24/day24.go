package aoc2021_day24

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davejhilton/adventofcode/challenges"
	"github.com/davejhilton/adventofcode/log"
	"github.com/davejhilton/adventofcode/util"
)

func part1(input []string) (string, error) {
	instructions := parseInput(input)

	PrintStuffToHelpReverseEngineer(instructions)
	stackOps := getStackInstructions(instructions)
	for i, op := range stackOps {
		if op.Opcode == "PUSH" {
			log.Debugf("PUSH: digits[%02d] + %d\n", i, op.ValueModifier)
		} else {
			log.Debugf("POP:  {popped} - %d SHOULD EQUAL digits[%02d]\n", op.ValueModifier, i)
		}
	}

	modelNum := 96929994293996

	execStack(stackOps, modelNum)

	alu := NewALU(0)
	result := RunProgram(instructions, modelNum, alu)
	if result == 0 {
		return fmt.Sprintf("%d", modelNum), nil
	} else {
		log.Debugf("CRAP. result = %d instead of 0\n", result)
	}

	return "no result", nil
}

func part2(input []string) (string, error) {
	instructions := parseInput(input)

	PrintStuffToHelpReverseEngineer(instructions)
	stackOps := getStackInstructions(instructions)
	for i, op := range stackOps {
		if op.Opcode == "PUSH" {
			log.Debugf("PUSH: digits[%02d] + %d\n", i, op.ValueModifier)
		} else {
			log.Debugf("POP:  {popped} - %d SHOULD EQUAL digits[%02d]\n", op.ValueModifier, i)
		}
	}

	modelNum := 41811761181141

	execStack(stackOps, modelNum)

	alu := NewALU(0)
	result := RunProgram(instructions, modelNum, alu)
	if result == 0 {
		return fmt.Sprintf("%d", modelNum), nil
	} else {
		log.Debugf("CRAP. result = %d instead of 0\n", result)
	}

	return "no result", nil
}

func parseInput(input []string) []Instruction {
	instructions := make([]Instruction, 0, len(input))
	for _, s := range input {
		parts := strings.Split(s, " ")
		op := Instruction{
			Opcode: parts[0],
			Arg1:   parts[1],
		}
		if len(parts) > 2 {
			op.Arg2 = parts[2]
		}
		instructions = append(instructions, op)
	}
	return instructions
}

/*

00 PUSH: digits[00] + 8
	01 PUSH: digits[01] + 16
		02 PUSH: digits[02] + 4
		03 POP:  {popped} - 11 SHOULD EQUAL digits[03]

		04 PUSH: digits[04] + 13
			05 PUSH: digits[05] + 5
				06 PUSH: digits[06] + 0
				07 POP:  {popped} - 5 SHOULD EQUAL digits[07]

				08 PUSH: digits[08] + 7
				09 POP:  {popped} - 0 SHOULD EQUAL digits[09]
			10 POP:  {popped} - 11 SHOULD EQUAL digits[10]
		11 POP:  {popped} - 13 SHOULD EQUAL digits[11]
	12 POP:  {popped} - 13 SHOULD EQUAL digits[12]
13 POP:  {popped} - 11 SHOULD EQUAL digits[13]


[ 0] +  8        [ 0] - 3 = [13]   9   4
[13] - 11        [13] + 3 = [ 0]   6   1

[ 1] + 16        [ 1] + 3 = [12]   6   1
[12] - 13        [12] - 3 = [ 1]   9   4

[ 2] +  4        [ 2] - 7 = [ 3]   9   8
[ 3] - 11        [ 3] + 7 = [ 2]   2   1

[ 4] + 13        [ 4] + 0 = [11]   9   1
[11] - 13        [11] - 0 = [ 4]   9   1

[ 5] +  5        [ 5] - 6 = [10]   9   7
[10] - 11        [10] + 6 = [ 5]   3   1

[ 6] +  0        [ 6] - 5 = [ 7]   9   6
[ 7] -  5        [ 7] + 5 = [ 6]   4   1

[ 8] +  7        [ 8] + 7 = [ 9]   2   1
[ 9] -  0        [ 9] - 7 = [ 8]   9   8

96929994293996
41811761181141

 0 = digit[xxx] +
 1 = digit[xxx] +
 2 = digit[  3] +
 3 = digit[  2] +
 4 = digit[xxx] +
 5 = digit[xxx] +
 6 = digit[  7] - 5
 7 = digit[  6] + 5
 8 = digit[xxx] +
 9 = digit[xxx] +
10 = digit[xxx] +
11 = digit[xxx] +
12 = digit[xxx] +
13 = digit[xxx] +






digit_03 = digit_02 + 7
digit_07 = digit_06 + 5
digit_09 = digit_08 + 7
digit_10 = digit_05 + 6
digit_11 = digit_04 + 0
digit_12 = digit_01 + 3
digit_13 = digit_00 + 3


digit 00 = 9
digit 01 = 9
digit 02 = 9
digit 03 = 2
digit 04 = 9
digit 05 = 9
digit 06 = 9
digit 07 = 4
digit 08 = 9
digit 09 = 2
digit 10 = 3
digit 11 = 9
digit 12 = 6
digit 13 = 6

99929994923966

highest value:

digit 00 = 6
digit 01 = 6
digit 02 = 2
digit 03 = 9
digit 04 = 9
digit 05 = 3
digit 06 = 4
digit 07 = 9
digit 08 = 2
digit 09 = 9
digit 10 = 9
digit 11 = 9
digit 12 = 9
digit 13 = 9


66299349299999
*/

func getStackInstructions(orig []Instruction) []StackInstruction {
	stackOps := make([]StackInstruction, 0)
	chunk := make([]Instruction, 0)
	for i, instr := range orig {
		if instr.Opcode == "inp" || i == len(orig)-1 {
			if len(chunk) > 0 {
				div := util.Atoi(chunk[4].Arg2)
				sub := util.Atoi(chunk[5].Arg2)
				add := util.Atoi(chunk[15].Arg2)

				if div == 1 { // push
					stackOps = append(stackOps, StackInstruction{
						Opcode:        "PUSH",
						ValueModifier: add,
					})
				} else if div == 26 { // pop
					stackOps = append(stackOps, StackInstruction{
						Opcode:        "POP",
						ValueModifier: sub * -1,
					})
				}
				chunk = make([]Instruction, 0)
			}
		}
		chunk = append(chunk, instr)
	}
	return stackOps
}

func execStack(ops []StackInstruction, modelNum int) {
	stack := []int{}
	digitStrings := strings.Split(fmt.Sprintf("%d", modelNum), "")
	digits := make([]int, 0)
	for _, d := range digitStrings {
		digits = append(digits, util.Atoi(d))
	}

	for i, op := range ops {
		digit := digits[i]
		if op.Opcode == "PUSH" {
			stack = append(stack, digit+op.ValueModifier)
			log.Debugf("%2d - PUSHED %d (%d + %d)\n", i, digit+op.ValueModifier, digit, op.ValueModifier)
		} else {
			popped := stack[len(stack)-1]
			val := popped - op.ValueModifier
			log.Debugf("%2d - POPPED %d (%d + %d)\n", i, popped, digit, op.ValueModifier)
			if val != digit {
				log.Debugf("   x NO MATCH: (%d - %d) != %d\n", popped, op.ValueModifier, digit)
			}
			stack = stack[:len(stack)-1]
		}
	}
}

type StackInstruction struct {
	Opcode        string
	ValueModifier int
}

func PrintStuffToHelpReverseEngineer(instructions []Instruction) {
	chunks := make([][]Instruction, 0)
	chunk := make([]Instruction, 0)
	chunkNum := 0
	for i, instr := range instructions {
		if instr.Opcode == "inp" || i == len(instructions)-1 {
			if len(chunk) > 0 {
				log.Debugf("chunk %d\n  - Divides by: %s\n  - Adds: %s\n  - Adds: %s\n\n",
					chunkNum,
					chunk[4].Arg2,
					chunk[5].Arg2,
					chunk[15].Arg2,
				)
				chunks = append(chunks, chunk)
				chunk = make([]Instruction, 0)
				chunkNum++
			}
		}
		chunk = append(chunk, instr)
	}
	_ = chunks
}

func RunProgram(instructions []Instruction, input int, alu *ALU) int {

	if alu == nil {
		alu = NewALU(input)
	} else {
		alu.Reset(input)
	}

	for _, instr := range instructions {
		alu.Execute(instr)
	}

	// log.Debugf("Program finished execution. Register status:\n\tw: %d\n\tx: %d\n\ty: %d\n\tz: %d\n",
	// 	alu.Registers["w"],
	// 	alu.Registers["x"],
	// 	alu.Registers["y"],
	// 	alu.Registers["z"],
	// )

	return alu.Registers["z"]
}

type ALU struct {
	Inputs     []int
	InputIndex int
	Registers  map[string]int
}

func NewALU(input int) *ALU {
	alu := &ALU{
		Registers: make(map[string]int),
	}
	alu.Reset(input)
	return alu
}

func (alu *ALU) Reset(input int) {
	runes := []rune(fmt.Sprintf("%d", input))
	digits := make([]int, 0, len(runes))

	for _, r := range runes {
		n, _ := strconv.Atoi(string(r))
		digits = append(digits, n)
	}

	alu.Inputs = digits
	alu.InputIndex = 0
	alu.Registers["w"] = 0
	alu.Registers["x"] = 0
	alu.Registers["y"] = 0
	alu.Registers["z"] = 0
}

func (alu *ALU) Execute(instr Instruction) {
	a := instr.Arg1
	b := instr.Arg2
	switch instr.Opcode {
	case "inp":
		in := alu.GetInput()
		log.Debugf("INPUT into %s (%d)", a, in)
		alu.SetValue(a, in)
	case "add":
		nA := alu.GetValue(a)
		nB := alu.GetValue(b)
		log.Debugf("ADD %s += %s (%d + %d)", a, b, nA, nB)
		alu.SetValue(a, nA+nB)
	case "mod":
		nA := alu.GetValue(a)
		nB := alu.GetValue(b)
		if nA < 0 || nB <= 0 {
			fmt.Printf("%s%s\n", log.Colorize("BAD?!?", log.Red, 0), instr)
		}
		log.Debugf("MOD %s %%= %s (%d %% %d)", a, b, nA, nB)
		alu.SetValue(a, nA%nB)
	case "mul":
		nA := alu.GetValue(a)
		nB := alu.GetValue(b)
		log.Debugf("MUL %s *= %s (%d * %d)", a, b, nA, nB)
		alu.SetValue(a, nA*nB)
	case "div":
		nA := alu.GetValue(a)
		nB := alu.GetValue(b)
		if nB == 0 {
			fmt.Printf("%s%s\n", log.Colorize("BAD?!?", log.Red, 0), instr)
		}
		log.Debugf("DIV %s /= %s (%d / %d)", a, b, nA, nB)
		alu.SetValue(a, nA/nB)
	case "eql":
		nA := alu.GetValue(a)
		nB := alu.GetValue(b)
		log.Debugf("EQL %s == %s (%d == %d)", a, b, nA, nB)
		if nA == nB {
			alu.SetValue(a, 1)
		} else {
			alu.SetValue(a, 0)
		}
	}
	log.Debugf(" - \t%d\t%d\t%d\t%d\n", alu.Registers["w"], alu.Registers["x"], alu.Registers["y"], alu.Registers["z"])
}

func (alu *ALU) GetInput() int {
	n := alu.Inputs[alu.InputIndex]
	alu.InputIndex += 1
	return n
}

func (alu *ALU) GetValue(v string) int {
	if v == "w" || v == "x" || v == "y" || v == "z" {
		return alu.Registers[v]
	}
	num, _ := strconv.Atoi(v)
	return num
}

func (alu *ALU) SetValue(loc string, val int) {
	alu.Registers[loc] = val
}

type Instruction struct {
	Opcode string
	Arg1   string
	Arg2   string
}

func (instr Instruction) String() string {
	return fmt.Sprintf("\n%s %s %s", instr.Opcode, instr.Arg1, instr.Arg2)
}

func init() {
	challenges.RegisterChallengeFunc(2021, 24, 1, "day24.txt", part1)
	challenges.RegisterChallengeFunc(2021, 24, 2, "day24.txt", part2)
}
